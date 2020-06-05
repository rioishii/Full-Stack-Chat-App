package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/sessions"
)

//UsersHandler handles requests for the "users" resource
func (ctx *HandlerCtx) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		nu := users.NewUser{}
		err := json.NewDecoder(r.Body).Decode(&nu)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := nu.ToUser()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userWithID, err := ctx.UserStore.Insert(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		firstname := strings.ToLower(userWithID.FirstName)
		fsplit := strings.Split(firstname, " ")
		for i := range fsplit {
			fsplit[i] = strings.TrimSpace(fsplit[i])
			ctx.Trie.Add(fsplit[i], userWithID.ID)
		}
		lastname := strings.ToLower(userWithID.LastName)
		lsplit := strings.Split(lastname, " ")
		for i := range lsplit {
			lsplit[i] = strings.TrimSpace(lsplit[i])
			ctx.Trie.Add(lsplit[i], userWithID.ID)
		}
		ctx.Trie.Add(strings.ToLower(userWithID.UserName), userWithID.ID)

		sessionState := &SessionState{time.Now(), userWithID}
		sid, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sessionState, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(sid) == 0 {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userWithID)
	} else if r.Method == http.MethodGet {
		_, err := sessions.GetSessionID(r, ctx.SigningKey)
		if err != nil {
			http.Error(w, "user is not authenticated", http.StatusUnauthorized)
			return
		}
		query := r.URL.Query().Get("q")
		if len(query) < 1 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		userIDs := ctx.Trie.Find(query, 20)
		users := []*users.User{}
		for _, id := range userIDs {
			user, err := ctx.UserStore.GetByID(id)
			if err != nil {
				http.Error(w, "no user found with given ID", http.StatusNotFound)
				return
			}
			users = append(users, user)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	} else {
		http.Error(w, "http method must be GET or POST", http.StatusMethodNotAllowed)
		return
	}
}

//SpecificUserHandler handles requests for a specific user
func (ctx *HandlerCtx) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	_, err := sessions.GetSessionID(r, ctx.SigningKey)
	if err != nil {
		http.Error(w, "user is not authenticated", http.StatusUnauthorized)
		return
	}
	stringID := path.Base(r.URL.Path)
	sessionState := &SessionState{}
	_, err = sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, sessionState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id int64
	if stringID == "me" {
		id = sessionState.User.ID
	} else {
		id, _ = strconv.ParseInt(stringID, 10, 64)
	}

	if stringID != "me" && id != sessionState.User.ID {
		http.Error(w, "user ID does not match the currently authenticated user", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodGet {
		stringID := path.Base(r.URL.Path)
		id, _ := strconv.ParseInt(stringID, 10, 64)
		user, err := ctx.UserStore.GetByID(id)
		if err != nil {
			http.Error(w, "no user found with given ID", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else if r.Method == http.MethodPatch {
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		update := users.Updates{}
		err = json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		oldUser := sessionState.User
		err = sessionState.User.ApplyUpdates(&update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := ctx.UserStore.Update(id, &update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		oldfirstname := strings.ToLower(oldUser.FirstName)
		ofsplit := strings.Split(oldfirstname, " ")
		for i := range ofsplit {
			ofsplit[i] = strings.TrimSpace(ofsplit[i])
			ctx.Trie.Remove(ofsplit[i], oldUser.ID)
		}
		oldlastname := strings.ToLower(oldUser.LastName)
		olsplit := strings.Split(oldlastname, " ")
		for i := range olsplit {
			olsplit[i] = strings.TrimSpace(olsplit[i])
			ctx.Trie.Remove(olsplit[i], oldUser.ID)
		}

		newfirstname := strings.ToLower(user.FirstName)
		nfsplit := strings.Split(newfirstname, " ")
		for i := range nfsplit {
			nfsplit[i] = strings.TrimSpace(nfsplit[i])
			ctx.Trie.Add(nfsplit[i], user.ID)
		}
		newlastname := strings.ToLower(user.LastName)
		nlsplit := strings.Split(newlastname, " ")
		for i := range nlsplit {
			nlsplit[i] = strings.TrimSpace(nlsplit[i])
			ctx.Trie.Add(nlsplit[i], user.ID)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "http method must be GET or POST", http.StatusMethodNotAllowed)
		return
	}
}

//SessionsHandler handles requests for the "sessions" resource, and allows clients
//to begin a new session using an existing user's credentials
func (ctx *HandlerCtx) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-type")
		if contentType != "application/json" {
			http.Error(w, "request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		cred := users.Credentials{}
		err := json.NewDecoder(r.Body).Decode(&cred)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := ctx.UserStore.GetByEmail(cred.Email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		err = user.Authenticate(cred.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		sessionState := &SessionState{time.Now(), user}
		sid, err := sessions.BeginSession(ctx.SigningKey, ctx.SessionStore, sessionState, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(sid) == 0 {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "http method must be POST", http.StatusMethodNotAllowed)
		return
	}
}

//SpecificSessionHandler handles requests related to a specific authenticated session
func (ctx *HandlerCtx) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	_, err := sessions.GetSessionID(r, ctx.SigningKey)
	if err != nil {
		http.Error(w, "user is not authenticated", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodDelete {
		baseURL := path.Base(r.URL.Path)
		if baseURL != "mine" {
			http.Error(w, "status forbidden", http.StatusForbidden)
			return
		}
		_, err := sessions.EndSession(r, ctx.SigningKey, ctx.SessionStore)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "http method must be DELETE", http.StatusMethodNotAllowed)
		return
	}
}
