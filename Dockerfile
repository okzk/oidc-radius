FROM golang AS build-env-golang
ENV CGO_ENABLED=0
RUN go get github.com/okzk/oidc-radius


FROM alpine

RUN apk add --no-cache ca-certificates
COPY --from=build-env-golang /go/bin/oidc-radius /usr/local/bin/

ENV RADIUS_SECRET= \
  CIBA_ISSUER= \
  CIBA_AUTHN_ENDBPOINT= \
  CIBA_TOKEN_ENDBPOINT= \
  CIBA_SCOPE=openid \
  CIBA_CLIENT_ID= \
  CIBA_CLIENT_SECRET= \
  USERNAME_SEPARATOR=

EXPOSE 1812/udp 1813/udp

CMD ["/usr/local/bin/oidc-radius"]
