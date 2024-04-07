# go-idp

## トークンの払い出し
1. `http://localhost:3846/oauth2/auth?client_id=my-client&redirect_uri=http://localhost:3846/callback&state=64aa6f2d-52d1-ec96-04b7-832f8720e7a7&response_type=code`にアクセス

2. usernameに「peter」と入力

3. callbackで404が返ってくるがそのクエリパラメータから認可コードを取得

4. トークンリクエストを行う
```shell
curl --location 'http://localhost:3846/oauth2/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--header 'Authorization: Basic bXktY2xpZW50OmZvb2Jhcg==' \
--data-urlencode 'grant_type=authorization_code' \
--data-urlencode 'redirect_uri=http://localhost:3846/callback' \
--data-urlencode 'code={取得した認可コード}'
```