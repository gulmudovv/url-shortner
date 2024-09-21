package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gulmudovv/url-shortener/internal/http-server/handlers/url/save"
	"github.com/gulmudovv/url-shortener/internal/lib/api"
	"github.com/gulmudovv/url-shortener/internal/lib/api/random"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

const (
	host = "localhost:8080"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   faker.Internet().Url(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("superuser", "secret").
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}

//nolint:funlen
func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   faker.Internet().Url(),
			alias: random.NewRandomString(10),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: random.NewRandomString(10),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   faker.Internet().Url(),
			alias: "",
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("superuser", "secret").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect

			testRedirect(t, alias, tc.url)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}
