package oauth2client

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

func OwnerHandler(c oauth2.Config) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte("<h1>Resource Owner Password Credentials Grant</h1>"))
		_ = req.ParseForm()
		if req.Form.Get("username") == "" || req.Form.Get("password") == "" {
			_, _ = rw.Write([]byte(`<form method="post">
			<ul>
				<li>
					<input type="text" name="username" placeholder="username"/> <small>try "peter"</small>
				</li>
				<li>
					<input type="password" name="password" placeholder="password"/> <small>try "secret"</small><br>
				</li>
				<li>
					<input type="submit" />
				</li>
			</ul>
		</form>`))
			_, _ = rw.Write([]byte(`<p><a href="/">Go back</a></p>`))
			return
		}

		token, err := c.PasswordCredentialsToken(context.Background(), req.Form.Get("username"), req.Form.Get("password"))
		if err != nil {
			_, _ = fmt.Fprintf(rw, `<p>I tried to get a token but received an error: %s</p>`, err.Error())
			_, _ = rw.Write([]byte(`<p><a href="/">Go back</a></p>`))
			return
		}
		_, _ = fmt.Fprintf(rw, `<p>Awesome, you just received an access token!<br><br>%s<br><br><strong>more info:</strong><br><br>%v</p>`, token.AccessToken, token)
		_, _ = rw.Write([]byte(`<p><a href="/">Go back</a></p>`))
	}
}
