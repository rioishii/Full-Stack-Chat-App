package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	query := r.URL.Query().Get("url")
	if len(query) < 1 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	bodyStream, err := fetchHTML(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	summary, err := extractSummary(query, bodyStream)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer bodyStream.Close()
	json.NewEncoder(w).Encode(summary)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("response status code was %v", resp.StatusCode)
	}
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, fmt.Errorf("response content type was %s not text/html", ctype)
	}

	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	tokenizer := html.NewTokenizer(htmlStream)
	ps := PageSummary{}
	images := []*PreviewImage{}
	counter := -1
	ogTitleFound := false
	ogDescriptionFound := false

	for {
		//get the next token type
		tokenType := tokenizer.Next()

		//if it's an error token, we either reached
		//the end of the file, or the HTML was malformed
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			return nil, fmt.Errorf("error tokenizing HTML: %v", tokenizer.Err())
		}

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if token.Data == "meta" {
				ogType, ok := extractMetaProperty(token, "og:type")
				if ok {
					ps.Type = ogType
				}
				ogURL, ok := extractMetaProperty(token, "og:url")
				if ok {
					ps.URL = ogURL
				}
				ogTitle, ok := extractMetaProperty(token, "og:title")
				if ok {
					ogTitleFound = true
					ps.Title = ogTitle
				}
				ogSiteName, ok := extractMetaProperty(token, "og:site_name")
				if ok {
					ps.SiteName = ogSiteName
				}
				ogDescription, ok := extractMetaProperty(token, "og:description")
				if ok {
					ogDescriptionFound = true
					ps.Description = ogDescription
				}
				description, ok := extractMetaName(token, "description")
				if ok {
					if !ogDescriptionFound {
						ps.Description = description
					}
				}
				author, ok := extractMetaName(token, "author")
				if ok {
					ps.Author = author
				}
				keywords, ok := extractMetaName(token, "keywords")
				if ok {
					s := strings.Split(keywords, ",")
					for i := range s {
						s[i] = strings.TrimSpace(s[i])
					}
					ps.Keywords = s
				}
				image, ok := extractMetaProperty(token, "og:image")
				if ok {
					counter++
					pi := PreviewImage{}
					u, _ := url.Parse(image)
					if u.IsAbs() {
						pi.URL = u.String()
					} else {
						base, _ := url.Parse(pageURL)
						pi.URL = base.ResolveReference(u).String()
					}
					images = append(images, &pi)
				}
				imageType, ok := extractMetaProperty(token, "og:image:type")
				if ok {
					pi := images[counter]
					pi.Type = imageType
				}
				imageSecureURL, ok := extractMetaProperty(token, "og:image:secure_url")
				if ok {
					pi := images[counter]
					pi.SecureURL = imageSecureURL
				}
				imageWidth, ok := extractMetaProperty(token, "og:image:width")
				if ok {
					width, _ := strconv.Atoi(imageWidth)
					pi := images[counter]
					pi.Width = width
				}
				imageHeight, ok := extractMetaProperty(token, "og:image:height")
				if ok {
					height, _ := strconv.Atoi(imageHeight)
					pi := images[counter]
					pi.Height = height
				}
				imageAlt, ok := extractMetaProperty(token, "og:image:alt")
				if ok {
					pi := images[counter]
					pi.Alt = imageAlt
				}
			}
			if token.Data == "title" {
				if !ogTitleFound {
					tokenType = tokenizer.Next()
					if tokenType == html.TextToken {
						ps.Title = tokenizer.Token().Data
					}
				}
			}
			if token.Data == "link" {
				icon := extractLinkAttribute(token, pageURL)
				ps.Icon = icon
			}
		}
	}

	if len(images) > 0 {
		ps.Images = images
	}

	return &ps, nil
}

func extractMetaProperty(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}

func extractMetaName(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "name" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}

func extractLinkAttribute(t html.Token, pageURL string) *PreviewImage {
	icon := PreviewImage{}
	for _, attr := range t.Attr {
		if attr.Key == "href" {
			u, _ := url.Parse(attr.Val)
			if u.IsAbs() {
				icon.URL = u.String()
			} else {
				base, _ := url.Parse(pageURL)
				icon.URL = base.ResolveReference(u).String()
			}
		}
		if attr.Key == "sizes" {
			if attr.Val != "any" {
				s := strings.Split(attr.Val, "x")
				height, _ := strconv.Atoi(s[0])
				icon.Height = height
				width, _ := strconv.Atoi(s[1])
				icon.Width = width
			}
		}
		if attr.Key == "type" {
			icon.Type = attr.Val
		}
	}

	return &icon
}
