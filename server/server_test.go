package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/uhlig-it/plaintoot/plaintoot"
	ptServer "github.com/uhlig-it/plaintoot/server"
)

var _ = Describe("http-server", func() {
	var (
		request  *http.Request
		response *httptest.ResponseRecorder
		context  echo.Context
		e        *echo.Echo
		server   ptServer.Server
	)

	BeforeEach(func() {
		e = echo.New()
		response = httptest.NewRecorder()

		r, err := plaintoot.NewMockRepository("https://example.com/@foo@example.net/42", "Hello World")
		Expect(err).ToNot(HaveOccurred())

		server = *ptServer.NewServer(r)
	})

	JustBeforeEach(func() {
		context = e.NewContext(request, response)
	})

	Context("GET /", func() {
		BeforeEach(func() {
			request = httptest.NewRequest(http.MethodGet, "/", nil)
		})

		JustBeforeEach(func() {
			server.Root(context)
		})

		It("succeeds", func() {
			Expect(response.Code).To(Equal(200))
		})

		Context("with blurb", func() {
			BeforeEach(func() {
				server = *server.WithBlurb("Fiat lux")
			})

			It("responds with the blurb", func() {
				Expect(response.Body).To(ContainSubstring("Fiat lux"))
			})
		})
	})

	Context("GET /version", func() {
		BeforeEach(func() {
			request = httptest.NewRequest(http.MethodGet, "/version", nil)
		})

		JustBeforeEach(func() {
			server.Version(context)
		})

		It("succeeds", func() {
			Expect(response.Code).To(Equal(200))
		})

		Context("the response body", func() {
			var body *bytes.Buffer

			JustBeforeEach(func() {
				body = response.Body
			})

			It("has the version number", func() {
				Expect(body).To(MatchRegexp("^plaintoot v"))
			})

			It("has the build date", func() {
				Expect(body).To(ContainSubstring("built on"))
			})
		})
	})

	Context("POST / with query param", func() {
		BeforeEach(func() {
			request = httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(url.Values{
					"url": []string{"https://example.com/@foo@example.net/42"},
				}.Encode(),
				),
			)
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		})

		JustBeforeEach(func() {
			server.Lookup(context)
		})

		It("succeeds", func() {
			Expect(response.Code).To(Equal(200))
		})
	})

	Context("GET /liveness", func() {
		BeforeEach(func() {
			request = httptest.NewRequest(http.MethodGet, "/liveness", nil)
		})

		JustBeforeEach(func() {
			server.Liveness(context)
		})

		It("succeeds", func() {
			Expect(response.Code).To(Equal(200))
		})
	})

	// TODO The mock repo doesn't know about the default lookup that we do for readiness
	PContext("GET /readiness", func() {
		BeforeEach(func() {
			request = httptest.NewRequest(http.MethodGet, "/readiness", nil)
		})

		JustBeforeEach(func() {
			server.Readiness(context)
		})

		It("succeeds", func() {
			Expect(response.Code).To(Equal(200))
		})
	})
})
