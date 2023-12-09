package plaintoot

import (
	"errors"
	"net/url"
)

type MockRepository map[string]string

func (r MockRepository) Lookup(uri *url.URL) (*Post, error) {
	text, found := r[uri.String()]

	if found {
		return &Post{Text: text, User: uri.User.String(), Host: uri.Host}, nil
	} else {
		return nil, errors.New("no post found with that URL")
	}
}
