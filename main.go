package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/k3a/html2text"
	"github.com/mattn/go-mastodon"
)

func main() {
	tootURL, err := url.Parse(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	c := mastodon.NewClient(&mastodon.Config{
		Server:   fmt.Sprintf("%s://%s:%s", tootURL.Scheme, tootURL.Hostname(), tootURL.Port()),
		ClientID: "plaintoot",
	})

	parts := strings.Split(tootURL.Path, "/")

	status, err := c.GetStatus(context.Background(), mastodon.ID(parts[len(parts)-1]))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n\n-- %s@%s\n", html2text.HTML2Text(status.Content), status.Account.Username, tootURL.Hostname())
}
