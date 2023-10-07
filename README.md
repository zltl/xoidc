# xoidc
oidc implement

client:
```bash
CLIENT_ID=MA== CLIENT_SECRET=123456 ISSUER=http://localhost:9998/ SCOPES="openid profile" PORT=9999 go run github.com/zitadel/oidc/v2/example/client/app
```

then open http://127.0.0.1:9999/login

