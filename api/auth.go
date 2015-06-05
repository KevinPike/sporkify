package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/zmb3/spotify"

	"golang.org/x/oauth2"
)

type auth struct {
	state string
	sa    spotify.Authenticator
}

func newAuth() (*auth, error) {
	scopes := []string{"playlist-modify-public", "playlist-modify-private"}
	redirectURL := os.Getenv("SPOTIFY_REDIRECT")
	if redirectURL == "" {
		return nil, errors.New("define SPOTIFY_REDIRECT")
	}

	sa := spotify.NewAuthenticator(redirectURL, scopes...)
	return &auth{
		state: "random",
		sa:    sa,
	}, nil
}

func (a *auth) RedirectURL() string {
	return a.sa.AuthURL(a.state)
}

func (a *auth) Exchange(r *http.Request) (*oauth2.Token, error) {
	return a.sa.Token(a.state, r)
}

func (a *auth) Client(tok *oauth2.Token) spotify.Client {
	return a.sa.NewClient(tok)
}
