package plaintoot

import (
	"errors"
	"net/url"
)

type MockRepository struct {
	uri  *url.URL
	text string
}

func NewMockRepository(u, t string) (*MockRepository, error) {
	uri, err := url.Parse(u)

	if err != nil {
		return nil, err
	}

	return &MockRepository{uri: uri, text: t}, nil
}

func (r *MockRepository) Lookup(uri *url.URL) (*Post, error) {
	if uri.String() == r.uri.String() {
		return &Post{Text: "example.com", User: "foo", Host: "example.net"}, nil
	} else {
		return nil, errors.New("no post found with that URL")
	}
}
