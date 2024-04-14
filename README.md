# go-idp

## トークンの払い出し(authorization_code)
1. `http://localhost:3846/oauth2/auth?client_id=my-client&redirect_uri=http://localhost:3846/callback&state=64aa6f2d-52d1-ec96-04b7-832f8720e7a7&response_type=code`にアクセス

2. フォームに以下を入力
```
username: test@example.com
password: password
```

3. トークンが返却される

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