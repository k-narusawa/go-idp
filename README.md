# go-idp

## トークンの払い出し(authorization_code)

1. go のアプリケーションを起動

```shell
go run server.go
```

2. sp フォルダに移動

```shell
cd sp
npm install
```

3. next を起動

```shell
npm run dev
```

4. `localhost:3000`にアクセスしてログインを行う

## トークンの払い出し(client_credentials)

1. 以下のコマンドを実行

```shell
curl --location 'http://localhost:3846/oauth2/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=client_credentials' \
--data-urlencode 'client_id=my-client' \
--data-urlencode 'client_secret=foobar'
```

2. 返却されたレスポンスを利用して有効性確認

```shell
curl --location 'http://localhost:3846/oauth2/introspect' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'Authorization: Basic bXktY2xpZW50OmZvb2Jhcg==' \
--data-urlencode 'token=ory_at_2IoxjMxaj_NLpDSCjXmdNKcJEJAv4GWUrJaltyqHhao.-3jlEpxe9p3fXXdsoA2t7DrJCTxn9tjc_orUzszfmf4'
```
