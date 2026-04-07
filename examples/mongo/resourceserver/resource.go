package resourceserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2/clientcredentials"
)

func ProtectedEndpoint(c clientcredentials.Config) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		resp, err := c.Client(context.Background()).PostForm(strings.ReplaceAll(c.TokenURL, "token", "introspect"), url.Values{"token": []string{req.URL.Query().Get("token")}, "scope": []string{req.URL.Query().Get("scope")}})
		if err != nil {
			_, _ = fmt.Fprintf(rw, "<h1>An error occurred!</h1><p>Could not perform introspection request: %v</p>", err)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		introspection := struct {
			Active bool `json:"active"`
		}{}
		out, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(out, &introspection); err != nil {
			_, _ = fmt.Fprintf(rw, "<h1>An error occurred!</h1>%s\n%s", err.Error(), out)
			return
		}

		if !introspection.Active {
			_, _ = fmt.Fprint(rw, `<h1>Request could not be authorized.</h1>
<a href="/">return</a>`)
			return
		}

		_, _ = fmt.Fprintf(rw, `<h1>Request authorized!</h1>
<code>%s</code><br>
<hr>
<a href="/">return</a>
`,
			out,
		)
	}
}
