package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	sid, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, ErrNoSessionID
	}
	store.Save(sid, sessionState)
	bearer := "Bearer " + sid.String()
	w.Header().Add("Authorization", bearer)

	return sid, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) == 2 {
		reqToken = strings.TrimSpace(splitToken[1])
	} else {
		reqToken = r.URL.Query().Get("auth")
		if reqToken != "" {
			splitToken := strings.Split(reqToken, "Bearer")
			if len(splitToken) != 2 {
				return InvalidSessionID, ErrInvalidScheme
			}
			reqToken = strings.TrimSpace(splitToken[1])
		} else {
			return InvalidSessionID, ErrInvalidScheme
		}
	}

	sid, err := ValidateID(reqToken, signingKey)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}
	return sid, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}
	err = store.Get(sid, sessionState)
	if err != nil {
		return InvalidSessionID, ErrStateNotFound
	}
	return sid, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	sid, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}
	err = store.Delete(sid)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}
	return sid, nil
}
