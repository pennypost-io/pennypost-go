// Generated from api/openapi.json (sha cc251884d96d) by sdk/generate.ts — do not edit by hand.
// PennyPost: the affordable email API. https://pennypost.io/docs
package pennypost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const defaultBase = "https://api.pennypost.io"

// SendEmailRequest is the body for SendEmail. To is always a slice (1-50).
type SendEmailRequest struct {
	From     string            `json:"from"`
	To       []string          `json:"to"`
	Subject  string            `json:"subject"`
	HTML     string            `json:"html,omitempty"`
	Text     string            `json:"text,omitempty"`
	ReplyTo  string            `json:"reply_to,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

type RemoveResult struct {
	Removed bool `json:"removed"`
}

// Error is the API's typed error contract.
type Error struct {
	Status    int
	Type      string `json:"type"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Param     string `json:"param"`
	Retryable bool   `json:"retryable"`
}

func (e *Error) Error() string { return fmt.Sprintf("pennypost: %s (%s)", e.Message, e.Code) }

type Client struct {
	key  string
	base string
	hc   *http.Client
}

func New(apiKey string, opts ...func(*Client)) *Client {
	c := &Client{key: apiKey, base: defaultBase, hc: http.DefaultClient}
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithBaseURL overrides the API base (tests, staging).
func WithBaseURL(u string) func(*Client) {
	return func(c *Client) { c.base = strings.TrimRight(u, "/") }
}

func (c *Client) request(method, path string, body any, query map[string]string, idem string, out any) error {
	u := c.base + path
	if len(query) > 0 {
		q := url.Values{}
		for k, v := range query {
			if v != "" {
				q.Set(k, v)
			}
		}
		if enc := q.Encode(); enc != "" {
			u += "?" + enc
		}
	}
	var rd *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rd = bytes.NewReader(b)
	} else {
		rd = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(method, u, rd)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	res, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		var wrap struct {
			Error Error `json:"error"`
		}
		_ = json.NewDecoder(res.Body).Decode(&wrap)
		e := wrap.Error
		e.Status = res.StatusCode
		if e.Message == "" {
			e.Message = fmt.Sprintf("HTTP %d", res.StatusCode)
			e.Retryable = res.StatusCode >= 500
		}
		return &e
	}
	return json.NewDecoder(res.Body).Decode(out)
}

// SendEmail calls POST /v1/emails.
func (c *Client) SendEmail(req *SendEmailRequest, idempotencyKey string) (*SendEmailResponse, error) {
	out := new(SendEmailResponse)
	if err := c.request("POST", "/v1/emails", req, nil, idempotencyKey, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListEmails calls GET /v1/emails.
func (c *Client) ListEmails(params map[string]string) (*EmailPage, error) {
	out := new(EmailPage)
	if err := c.request("GET", "/v1/emails", nil, params, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetEmail calls GET /v1/emails/{id}.
func (c *Client) GetEmail(id string) (*Email, error) {
	out := new(Email)
	if err := c.request("GET", "/v1/emails/" + url.PathEscape(id), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListSuppressions calls GET /v1/suppressions.
func (c *Client) ListSuppressions(params map[string]string) (*SuppressionPage, error) {
	out := new(SuppressionPage)
	if err := c.request("GET", "/v1/suppressions", nil, params, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// AddSuppression calls POST /v1/suppressions.
func (c *Client) AddSuppression(email string) (*Suppression, error) {
	out := new(Suppression)
	if err := c.request("POST", "/v1/suppressions", map[string]string{"email": email}, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveSuppression calls DELETE /v1/suppressions/{email}.
func (c *Client) RemoveSuppression(email string) (*RemoveResult, error) {
	out := new(RemoveResult)
	if err := c.request("DELETE", "/v1/suppressions/" + url.PathEscape(email), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListWebhooks calls GET /v1/webhooks.
func (c *Client) ListWebhooks() (*WebhookList, error) {
	out := new(WebhookList)
	if err := c.request("GET", "/v1/webhooks", nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateWebhook calls POST /v1/webhooks.
func (c *Client) CreateWebhook(req *CreateWebhookRequest) (*WebhookEndpoint, error) {
	out := new(WebhookEndpoint)
	if err := c.request("POST", "/v1/webhooks", req, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteWebhook calls DELETE /v1/webhooks/{id}.
func (c *Client) DeleteWebhook(id string) (*RemoveResult, error) {
	out := new(RemoveResult)
	if err := c.request("DELETE", "/v1/webhooks/" + url.PathEscape(id), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// TestWebhook calls POST /v1/webhooks/{id}/test.
func (c *Client) TestWebhook(id string) (*WebhookTestResult, error) {
	out := new(WebhookTestResult)
	if err := c.request("POST", "/v1/webhooks/" + url.PathEscape(id) + "/test", nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAccount calls GET /v1/account.
func (c *Client) GetAccount() (*Account, error) {
	out := new(Account)
	if err := c.request("GET", "/v1/account", nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListKeys calls GET /v1/keys.
func (c *Client) ListKeys() (*ApiKeyList, error) {
	out := new(ApiKeyList)
	if err := c.request("GET", "/v1/keys", nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateKey calls POST /v1/keys.
func (c *Client) CreateKey(req *CreateKeyRequest) (*ApiKeySummary, error) {
	out := new(ApiKeySummary)
	if err := c.request("POST", "/v1/keys", req, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}

// RevokeKey calls DELETE /v1/keys/{id}.
func (c *Client) RevokeKey(id string) (*RemoveResult, error) {
	out := new(RemoveResult)
	if err := c.request("DELETE", "/v1/keys/" + url.PathEscape(id), nil, nil, "", out); err != nil {
		return nil, err
	}
	return out, nil
}
