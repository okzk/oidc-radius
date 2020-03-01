package main

import (
	"context"
	"github.com/okzk/go-ciba"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"os"
	"strings"
)

func main() {
	client := ciba.NewClient(
		os.Getenv("CIBA_ISSUER"),
		os.Getenv("CIBA_AUTHN_ENDBPOINT"),
		os.Getenv("CIBA_TOKEN_ENDBPOINT"),
		os.Getenv("CIBA_SCOPE"),
		os.Getenv("CIBA_CLIENT_ID"),
		os.Getenv("CIBA_CLIENT_SECRET"),
	)
	sep := os.Getenv("USERNAME_SEPARATOR")

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password := ""
		if sep != "" {
			parts := strings.SplitN(username, sep, 2)
			if len(parts) == 2 {
				username = parts[0]
				password = parts[1]
			}
		}
		if password == "" {
			password = rfc2865.UserPassword_GetString(r.Packet)
		}

		code := radius.CodeAccessReject
		token, err := client.Authenticate(context.Background(), ciba.LoginHint(username), ciba.UserCode(password))
		if err == nil && token.Claims()["sub"] == username {
			code = radius.CodeAccessAccept
		}
		w.Write(r.Response(code))
	}

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
