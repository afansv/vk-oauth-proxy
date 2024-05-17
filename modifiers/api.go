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
	apiUsersGetPath           = "/method/users.get"
	apiUsersGetUserIDFieldKey = "id"

	apiResponseFieldKey = "response"
)

type APIResponseModifier struct {
	userEmailStore *store.UserEmail
}

func NewAPIResponseModifier(userEmailStore *store.UserEmail) *APIResponseModifier {
	return &APIResponseModifier{userEmailStore: userEmailStore}
}

func (modifier APIResponseModifier) Modify(resp *http.Response) error {
	if resp.Request.URL.Path != apiUsersGetPath || resp.StatusCode != http.StatusOK {
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

	// unpack response
	responseObj, responseObjOK := respMap[apiResponseFieldKey]
	if !responseObjOK {
		return nil
	}

	// find user id
	responseObjArr, responseObjArrOK := responseObj.([]interface{})
	if !responseObjArrOK {
		return nil
	}
	if len(responseObjArr) == 0 {
		return nil
	}
	responseObjMap, responseObjMapOK := responseObjArr[0].(map[string]interface{})
	if !responseObjMapOK {
		return nil
	}
	userIDValue, userIDValueOK := responseObjMap[apiUsersGetUserIDFieldKey]
	if !userIDValueOK {
		return nil
	}
	userID, userIDOK := userIDValue.(float64) // see https://pkg.go.dev/encoding/json#Unmarshal
	if !userIDOK {
		return nil
	}

	// find email by user id
	email, emailFound := modifier.userEmailStore.Get(int(userID))
	if emailFound {
		// inject email to response
		responseObjMap[emailFieldKey] = email
	}

	modifiedBodyData, err := json.Marshal(responseObjMap)
	if err != nil {
		return fmt.Errorf("marshal modified response: %w", err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(modifiedBodyData))

	return nil
}
