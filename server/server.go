package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uhlig-it/plaintoot/plaintoot"
)

type Server struct {
	startedAt    time.Time
	maxUptime    time.Duration
	startupDelay time.Duration
	blurb        string
	repository   plaintoot.Repository
}

func NewServer(r plaintoot.Repository) *Server {
	return &Server{repository: r}
}

func (s *Server) WithMaxUptime(maxUpTime time.Duration) *Server {
	s.maxUptime = maxUpTime
	return s
}

// WithStartupDelay simulates the time between creating this process and the HTTP server becoming available
// The kubelet uses startup probes to know when a container application has started.
func (s *Server) WithStartupDelay(startupDelay time.Duration) *Server {
	s.startupDelay = startupDelay
	return s
}

func (s *Server) WithBlurb(blurb string) *Server {
	s.blurb = blurb
	return s
}

func (s *Server) Start(addr string) error {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Response().Header().Set("Server", plaintoot.VersionStringShort())
			return next(ctx)
		}
	})

	e.GET("/", s.Root)
	e.POST("/", s.Lookup)
	e.GET("/liveness", s.Liveness)
	e.GET("/readiness", s.Readiness)
	e.GET("/version", s.Version)

	log.Printf("Starting server %s on port %s; pid is %d", plaintoot.VersionStringShort(), addr, os.Getpid())

	if s.maxUptime > 0 {
		log.Printf("Maximum allowed uptime set to %v; afterwards /liveness will report an error", s.maxUptime)
	}

	if s.startupDelay > 0 {
		log.Printf("Simulating startup time of %v", s.startupDelay)
	}

	time.Sleep(s.startupDelay)

	s.startedAt = time.Now()
	log.Printf("Server started at %v", s.startedAt)

	return e.Start(addr)
}

func (s *Server) Root(ctx echo.Context) error {
	if acceptsAnyOfContentType(ctx, echo.MIMETextHTML, echo.MIMETextHTMLCharsetUTF8) {
		return ctx.HTML(http.StatusOK, `
<!DOCTYPE html lang=en>
<head><title>plaintoot</title></head>
<body>
<p>Provides a plaintext representation of a Mastodon post ("toot").</p>
<form method=post>
	<p><label for=url>Enter a URL:</label></p>
	<p>
		<input type=text id=url name=url required minlength=12 size=80 />
		<input type="submit" />
	</p>
</form>
<p>
If you prefer the command line, use this example:
<pre>
curl http://localhost:8080 -d url=https://example.com/@someone@example.net/1234567890
</pre>
and replace <code>http://localhost:8080</code> with the real hostname and port.
</p>
</body>
		`)
	}

	return ctx.String(http.StatusOK, s.blurb)
}

func acceptsAnyOfContentType(ctx echo.Context, mimetypes ...string) bool {
	for _, acceptedHeader := range ctx.Request().Header["Accept"] {
		acceptedContentTypes := strings.Split(acceptedHeader, ",")

		for _, mimetype := range mimetypes {
			if slices.Contains(acceptedContentTypes, mimetype) {
				return true
			}
		}
	}

	return false
}

func (s *Server) Lookup(ctx echo.Context) error {
	url, err := url.Parse(ctx.FormValue("url"))

	if err != nil {
		return ctx.String(http.StatusUnprocessableEntity, "not a valid URL")
	}

	post, err := s.repository.Lookup(url)

	if err != nil {
		return ctx.String(http.StatusNotFound, "no such post")
	}

	return ctx.String(http.StatusOK, post.String())
}

// The kubelet uses liveness probes to know when to restart a container
func (s *Server) Liveness(ctx echo.Context) error {
	ctx.Response().Header().Add("X-Uptime", time.Since(s.startedAt).Round(time.Second).String())

	if s.maxUptime == 0 || time.Since(s.startedAt) < s.maxUptime {
		return ctx.String(http.StatusOK, "OK")
	} else {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error: Maximum uptime of %v reached", s.maxUptime))
	}
}

// The kubelet uses readiness probes to know when a container is ready to start accepting traffic
func (s *Server) Readiness(ctx echo.Context) error {
	u, err := url.Parse("https://chaos.social/@nixCraft@mastodon.social/111108182085516402")

	if err != nil {
		return ctx.String(http.StatusUnprocessableEntity, "could not parse default URL")
	}

	_, err = s.repository.Lookup(u)

	if err == nil {
		return ctx.String(http.StatusOK, "OK")
	} else {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", s.maxUptime))
	}
}

func (s *Server) Version(ctx echo.Context) error {
	return ctx.String(http.StatusOK, plaintoot.VersionString())
}
