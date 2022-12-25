package plaintoot

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/k3a/html2text"
	"github.com/mattn/go-mastodon"
)

func NewRepository(ctx context.Context) (Repository, error) {
	return &DefaultRepository{}, nil
}

type DefaultRepository struct {
}

func (r *DefaultRepository) Lookup(tootURL *url.URL) (*Post, error) {
	c := mastodon.NewClient(&mastodon.Config{
		Server:   fmt.Sprintf("%s://%s:%s", tootURL.Scheme, tootURL.Hostname(), tootURL.Port()),
		ClientID: "plaintoot",
	})

	parts := strings.Split(tootURL.Path, "/")

	status, err := c.GetStatus(context.Background(), mastodon.ID(parts[len(parts)-1]))

	if err != nil {
		log.Fatal(err)
	}

	return &Post{
		Text: html2text.HTML2Text(status.Content),
		User: status.Account.Username,
		Host: tootURL.Hostname(),
	}, nil
}
