package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// API creates an HTTP server for serving API requests
type API struct {
	*httprouter.Router

	auth     *auth
	sessions *sessions
}

// New creates a new API
func New() (*API, error) {
	router := httprouter.New()

	auth, err := newAuth()
	if err != nil {
		return nil, err
	}

	sessions := newSessions()

	api := &API{
		Router:   router,
		auth:     auth,
		sessions: sessions,
	}

	router.GET("/login", api.login)
	router.GET("/callback", api.callback)
	router.GET("/user", api.user)
	router.GET("/logout", api.logout)

	router.GET("/playlists", api.playlists)

	return api, nil
}

func (a *API) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	url := a.auth.RedirectURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *API) callback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token, err := a.auth.Exchange(r)

	noauth := "/#/noauth"
	if err != nil {
		http.Redirect(w, r, noauth, http.StatusMovedPermanently)
		return
	}

	client := a.auth.Client(token)
	user, err := client.CurrentUser()
	if err != nil {
		http.Redirect(w, r, noauth, http.StatusMovedPermanently)
		return
	}

	if err := a.sessions.New(token, user, w, r); err != nil {
		http.Redirect(w, r, noauth, http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, "/#/", http.StatusMovedPermanently)
}

func (a *API) user(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := a.sessions.User(w, r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	data, _ := json.Marshal(user)

	w.Write(data)
}

func (a *API) logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := a.sessions.Flush(w, r); err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func (a *API) playlists(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token, err := a.sessions.Token(w, r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user, err := a.sessions.User(w, r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	client := a.auth.Client(token)

	playlists, err := client.GetPlaylistsForUser(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(playlists)

	w.Write(data)
}
