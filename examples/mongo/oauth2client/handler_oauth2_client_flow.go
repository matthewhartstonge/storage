package oauth2client

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
)

func ClientEndpoint(c clientcredentials.Config) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte("<h1>Client Credentials Grant</h1>"))
		token, err := c.Token(context.Background())
		if err != nil {
			_, _ = fmt.Fprintf(rw, `<p>I tried to get a token but received an error: %s</p>`, err.Error())
			return
		}
		_, _ = fmt.Fprintf(rw, `<p>Awesome, you just received an access token!<br><br>%s<br><br><strong>more info:</strong><br><br>%v</p>`, token.AccessToken, token)
		_, _ = rw.Write([]byte(`<p><a href="/">Go back</a></p>`))
	}
}
