# go-idp

## クライアント登録

```shell
go run server.go

curl --location 'http://localhost:3846/admin/clients' \
--header 'Content-Type: application/json' \
--data '{
  "client_id": "go-idp",
  "client_secret": "~Ep52-Sp%iQtcEHpSLQ5,LT-,9*HMNfg,7WP",
  "redirect_uris": [
    "http://localhost:3000/api/auth/callback/go-idp"
  ],
  "grant_types": [
    "authorization_code",
    "refresh_token",
    "client_credentials"
  ],
  "response_types": [
    "code"
  ],
  "scopes": [
    "openid",
    "offline"
  ],
  "audience": "go-idp",
  "public": false
}'
```

## トークンの払い出し(authorization_code)

```shell
go run server.go

cd sp
npm install
npm run dev
```

デモサイトにアクセス.
http://localhost:3000
