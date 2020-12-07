package main

import (
	"github.com/okzk/go-ciba"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"log"
	"os"
	"strings"
)

func main() {
	client := ciba.NewClient(
		os.Getenv("CIBA_ISSUER"),
		os.Getenv("CIBA_AUTHN_ENDPOINT"),
		os.Getenv("CIBA_TOKEN_ENDPOINT"),
		os.Getenv("CIBA_SCOPE"),
		os.Getenv("CIBA_CLIENT_ID"),
		os.Getenv("CIBA_CLIENT_SECRET"),
	)
	sep := os.Getenv("USERNAME_SEPARATOR")

	accessRequestHandler := func(w radius.ResponseWriter, r *radius.Request) {
		if r.Code != radius.CodeAccessRequest {
			return
		}
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
		if username == "" || password == "" {
			w.Write(r.Response(radius.CodeAccessReject))
			return
		}

		token, err := client.Authenticate(r.Context(), ciba.LoginHint(username), ciba.UserCode(password))
		if err != nil {
			log.Printf("[INFO] authn failed. user: %s, error: %v", username, err)
			w.Write(r.Response(radius.CodeAccessReject))
		} else if token.Claims()["sub"] != username {
			log.Printf("[INFO] authn failed. user: %s, returned_sub: %s", username, token.Claims()["sub"])
			w.Write(r.Response(radius.CodeAccessReject))
		} else {
			log.Printf("[INFO] authn success. user: %s", username)
			w.Write(r.Response(radius.CodeAccessAccept))
		}
	}

	accountingRequestHandler := func(w radius.ResponseWriter, r *radius.Request) {
		if r.Code != radius.CodeAccountingRequest {
			return
		}
		w.Write(r.Response(radius.CodeAccountingResponse))
	}

	secret := []byte(os.Getenv("RADIUS_SECRET"))

	go func() {
		server := radius.PacketServer{
			Addr:         ":1813",
			Handler:      radius.HandlerFunc(accountingRequestHandler),
			SecretSource: radius.StaticSecretSource(secret),
		}
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(accessRequestHandler),
		SecretSource: radius.StaticSecretSource(secret),
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
