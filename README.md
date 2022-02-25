# oidc-radius

A RADIUS server implementation with OpenID CIBA flow.

This server uses POLL mode in CIBA flow.

## How to run

```sh
docker run -d -p 1812:1812/udp -p 1813:1813/udp \
-e RADIUS_SECRET="..." \
-e CIBA_ISSUER="https://example.com" \
-e CIBA_AUTHN_ENDPOINT="https://example.com/backchannel/authn" \
-e CIBA_TOKEN_ENDPOINT="https://example.com/token" \
-e CIBA_CLIENT_ID="..." \
-e CIBA_CLIENT_SECRET="..." \
okzk/oidc-radius
```

## User-Name and User-Password

This server uses `User-Name` as `login_hint`, and `User-Password` as `user_code` in CIBA flow.


## Environment Variables

### `RADIUS_SECRET`
The secret used for authorizing and decrypting RADIUS packets.
**REQUIRED**.

### `CIBA_ISSUER`

The value of `issuer` defined in OpenID Connection.
**REQUIRED**.

### `CIBA_AUTHN_ENDPOINT`

The value of `backchannel_authentication_endpoint` defined in OpenID Connection.
**REQUIRED**.


### `CIBA_TOKEN_ENDPOINT`

The value of `scope` defined in OpenID Connection.
**REQUIRED**.

### `CIBA_SCOPE`
The value of `token_endpoint` defined in OpenID Connection.
Default is `openid`

### `CIBA_CLIENT_ID`

The value of `client_id` defined in OpenID Connection.
**REQUIRED**.

### `CIBA_CLIENT_SECRET`

The value of `client_secret` defined in OpenID Connection.
**REQUIRED**.


### `USERNAME_SEPARATOR`

If not empty, RADIUS User-Name is splitted into `login_hint` and `user_code` by this value.
Default is empty.

This environment value is useful when RADIUS clients not support PAP.

## License
MIT
