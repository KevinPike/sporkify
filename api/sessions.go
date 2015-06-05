package api

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/astaxie/beego/session"
	"github.com/zmb3/spotify"
)

type sessions struct {
	sessions *session.Manager
}

func newSessions() *sessions {
	cookie := fmt.Sprintf(`{"cookieName":"sporkify","gclifetime":3600}`)
	s, _ := session.NewManager("memory", cookie)
	go s.GC()

	return &sessions{
		sessions: s,
	}
}

func (s *sessions) Token(w http.ResponseWriter, r *http.Request) (*oauth2.Token, error) {
	session, err := s.sessions.SessionStart(w, r)
	if err != nil {
		return nil, err
	}
	defer session.SessionRelease(w)

	rawToken := session.Get("token")
	token, ok := rawToken.(*oauth2.Token)
	if !ok {
		return nil, errors.New("no token")
	}

	return token, nil
}

func (s *sessions) User(w http.ResponseWriter, r *http.Request) (user *spotify.PrivateUser, err error) {
	session, err := s.sessions.SessionStart(w, r)
	if err != nil {
		return
	}
	defer session.SessionRelease(w)

	rawUser := session.Get("user")
	user, ok := rawUser.(*spotify.PrivateUser)
	if !ok {
		err = errors.New("no user")
		return
	}

	return
}

func (s *sessions) New(tok *oauth2.Token, user *spotify.PrivateUser, w http.ResponseWriter, r *http.Request) error {
	session, _ := s.sessions.SessionStart(w, r)
	defer session.SessionRelease(w)

	session.Set("token", tok)
	session.Set("user", user)

	return nil
}

func (s *sessions) Flush(w http.ResponseWriter, r *http.Request) error {
	session, err := s.sessions.SessionStart(w, r)
	if err != nil {
		return err
	}

	defer session.SessionRelease(w)

	return session.Flush()
}
