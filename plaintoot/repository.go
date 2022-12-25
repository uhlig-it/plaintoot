package plaintoot

import (
	"fmt"
	"net/url"
)

type Post struct {
	Text string
	User string
	Host string
}

func (p *Post) String() string {
	return fmt.Sprintf("%s\n\n-- %s@%s\n", p.Text, p.User, p.Host)

}

type Repository interface {
	// returns the Post identified by the given URL
	Lookup(uri *url.URL) (*Post, error)
}
