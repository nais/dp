package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nais/dp/backend/config"
)

const AzureGraphMemberOfEndpoint = "https://graph.microsoft.com/me/memberOf"

type CacheEntry struct {
	groups  []string
	updated time.Time
}

type AzureGroups struct {
	Cache  map[string]CacheEntry
	Client *http.Client
	Config config.Config
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type MemberOfResponse struct {
	Groups []MemberOfGroup `json:"value"`
}

type MemberOfGroup struct {
	Id string `json:"id"`
}

func (a *AzureGroups) GetGroupsForUser(ctx context.Context, token, email string) ([]string, error) {
	entry, found := a.Cache[email]
	if found && entry.updated.Add(1*time.Hour).Before(time.Now()) {
		return entry.groups, nil
	}
	bearerToken, err := a.getBearerTokenOnBehalfOfUser(ctx, token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, AzureGraphMemberOfEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bearerToken))
	response, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	var memberOfResponse MemberOfResponse
	if err := json.NewDecoder(response.Body).Decode(&memberOfResponse); err != nil {
		return nil, err
	}

	var groups []string
	for _, entry := range memberOfResponse.Groups {
		groups = append(groups, entry.Id)
	}

	a.Cache[email] = CacheEntry{
		groups:  groups,
		updated: time.Now(),
	}

	return groups, nil
}

func (a *AzureGroups) getBearerTokenOnBehalfOfUser(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.Config.OAuth2TokenURL, nil)
	if err != nil {
		return "", err
	}
	req.Form.Add("client_id", a.Config.OAuth2ClientID)
	req.Form.Add("client_secret", a.Config.OAuth2ClientSecret)
	req.Form.Add("scope", "https://graph.microsoft.com/.default")
	req.Form.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	req.Form.Add("requested_token_use", "on_behalf_of")
	req.Form.Add("assertion", token)

	response, err := a.Client.Do(req)
	if err != nil {
		return "", err
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}
	return tokenResponse.AccessToken, nil
}
