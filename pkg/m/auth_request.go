package m

import (
	"encoding/json"

	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

const (
	// CustomScope is an example for how to use custom scopes in this library
	//(in this scenario, when requested, it will return a custom claim)
	CustomScope = "custom_scope"

	// CustomClaim is an example for how to return custom claims with this library
	CustomClaim = "custom_claim"

	// CustomScopeImpersonatePrefix is an example scope prefix for passing user id to impersonate using token exchage
	CustomScopeImpersonatePrefix = "custom_scope:impersonate:"
)

type AuthRequest struct {
	ID           string
	CreationDate time.Time
	AuthReq      oidc.AuthRequest
	UserID       uuid.UUID
	IsDone       bool
	AuthTime     time.Time
}

func (a *AuthRequest) GetID() string {
	return a.ID
}

func (a *AuthRequest) GetACR() string {
	return "" // we won't handle acr in this example
}

func (a *AuthRequest) GetAMR() []string {
	// this example only uses password for authentication
	if a.IsDone {
		return []string{"pwd"}
	}
	return nil
}

func (a *AuthRequest) GetAudience() []string {
	return []string{a.AuthReq.ClientID} // always just use the client_id as audience
}

func (a *AuthRequest) GetAuthTime() time.Time {
	return a.AuthTime
}

func (a *AuthRequest) GetClientID() string {
	return a.AuthReq.ClientID
}

func (a *AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(&OIDCCodeChallenge{
		Challenge: a.AuthReq.CodeChallenge,
		Method:    string(a.AuthReq.CodeChallengeMethod),
	})
}

func (a *AuthRequest) GetNonce() string {
	return a.AuthReq.Nonce
}

func (a *AuthRequest) GetRedirectURI() string {
	return a.AuthReq.RedirectURI
}

func (a *AuthRequest) GetResponseType() oidc.ResponseType {
	return a.AuthReq.ResponseType
}

func (a *AuthRequest) GetResponseMode() oidc.ResponseMode {
	return "" // we won't handle response mode in this example
}

func (a *AuthRequest) GetScopes() []string {
	return a.AuthReq.Scopes
}

func (a *AuthRequest) GetState() string {
	return a.AuthReq.State
}

func (a *AuthRequest) GetSubject() string {
	return a.UserID.String()
}

func (a *AuthRequest) Done() bool {
	return a.IsDone
}

func (a *AuthRequest) Content() string {
	// json a.AuthReq
	js, err := json.Marshal(a.AuthReq)
	if err != nil {
		logrus.Fatal(err)
	}
	return string(js)
}

func (a *AuthRequest) SetContent(ct string) error {
	var arq oidc.AuthRequest
	err := json.Unmarshal([]byte(ct), &arq)
	if err != nil {
		return err
	}
	a.AuthReq = arq
	return nil
}
