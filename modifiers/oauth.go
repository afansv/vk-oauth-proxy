package modifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/afansv/vk-oauth-proxy/store"
)

const (
	oauthAccessTokenPath = "/access_token"

	oauthTokenTypeFieldKey   = "token_type"
	oauthTokenTypeFieldValue = "Bearer"

	oauthUserIDFieldKey = "user_id"
)

type OAuthResponseModifier struct {
	userEmailStore *store.UserEmail
}

func NewOAuthResponseModifier(
	userEmailStore *store.UserEmail,
) *OAuthResponseModifier {
	return &OAuthResponseModifier{userEmailStore: userEmailStore}
}

func (modifier OAuthResponseModifier) Modify(resp *http.Response) error {
	if resp.Request.URL.Path != oauthAccessTokenPath || resp.StatusCode != http.StatusOK {
		return nil
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	var upstreamResp interface{}
	if err := json.Unmarshal(rawBody, &upstreamResp); err != nil {
		return fmt.Errorf("unmarshal upstream response: %w", err)
	}

	respMap, ok := upstreamResp.(map[string]interface{})
	if !ok {
		return nil
	}

	// store user email in ttl cache if possible
	// ...
	// Locality of Behaviour (LoB) :)
	storeUserEmailIfPossible := func() {
		emailValue, emailValueOK := respMap[emailFieldKey]
		userIDValue, userIDValueOK := respMap[oauthUserIDFieldKey]
		if !emailValueOK || !userIDValueOK {
			return
		}

		email, emailOK := emailValue.(string)
		userID, userIDOK := userIDValue.(float64) // see https://pkg.go.dev/encoding/json#Unmarshal
		if emailOK && userIDOK {
			modifier.userEmailStore.Set(int(userID), email)
		}
	}

	storeUserEmailIfPossible()
	respMap[oauthTokenTypeFieldKey] = oauthTokenTypeFieldValue

	modifiedBodyData, err := json.Marshal(respMap)
	if err != nil {
		return fmt.Errorf("marshal modified response: %w", err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(modifiedBodyData))

	return nil
}
