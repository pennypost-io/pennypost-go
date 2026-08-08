// Generated from api/openapi.json (sha cc251884d96d) by sdk/generate.ts — do not edit by hand.
// PennyPost: the affordable email API. https://pennypost.io/docs
package pennypost

const SpecSHA = "cc251884d96d"

type Recipient struct {
	To string `json:"to"`
	Id *string `json:"id,omitempty"`
	Reason *string `json:"reason,omitempty"`
}

type SendEmailResponse struct {
	Accepted []Recipient `json:"accepted"`
	Suppressed []Recipient `json:"suppressed"`
	Quarantined []Recipient `json:"quarantined"`
	Failed []Recipient `json:"failed"`
}

type EmailEvent struct {
	Type string `json:"type"`
	Code *string `json:"code,omitempty"`
	Reason *string `json:"reason,omitempty"`
	At string `json:"at"`
}

type Email struct {
	Id string `json:"id"`
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Status string `json:"status"`
	Mode *string `json:"mode,omitempty"`
	Tags []string `json:"tags"`
	Metadata *map[string]string `json:"metadata,omitempty"`
	CreatedAt string `json:"created_at"`
	Events []EmailEvent `json:"events"`
}

type EmailPage struct {
	Data []Email `json:"data"`
	HasMore bool `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

type Suppression struct {
	Email string `json:"email"`
	Reason string `json:"reason"`
	SourceEmailId *string `json:"source_email_id,omitempty"`
	At *string `json:"at,omitempty"`
}

type SuppressionPage struct {
	Data []Suppression `json:"data"`
	HasMore bool `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

type CreateWebhookRequest struct {
	Url string `json:"url"`
	Events []string `json:"events"`
}

type WebhookEndpoint struct {
	Id string `json:"id"`
	Url string `json:"url"`
	Events []string `json:"events"`
	Secret *string `json:"secret,omitempty"`
	Status string `json:"status"`
	ConsecutiveFailures *int `json:"consecutive_failures,omitempty"`
	LastSuccessAt *string `json:"last_success_at,omitempty"`
	LastFailureAt *string `json:"last_failure_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

type WebhookList struct {
	Data []WebhookEndpoint `json:"data"`
}

type WebhookTestResult struct {
	Delivered bool `json:"delivered"`
	EndpointStatus string `json:"endpoint_status"`
}

type Account struct {
	Id string `json:"id"`
	Name *string `json:"name,omitempty"`
	ContactEmail *string `json:"contact_email,omitempty"`
	Plan string `json:"plan"`
	Status string `json:"status"`
	DailyCap int `json:"daily_cap"`
	MonthToDateSent *int `json:"month_to_date_sent,omitempty"`
	CardOnFile *bool `json:"card_on_file,omitempty"`
	Enforcement map[string]string `json:"enforcement"`
	CreatedAt *string `json:"created_at,omitempty"`
}

type CreateKeyRequest struct {
	Name *string `json:"name,omitempty"`
	Mode *string `json:"mode,omitempty"`
}

type ApiKeySummary struct {
	Id string `json:"id"`
	Name *string `json:"name,omitempty"`
	Prefix string `json:"prefix"`
	Mode string `json:"mode"`
	Key *string `json:"key,omitempty"`
	LastUsedAt *string `json:"last_used_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

type ApiKeyList struct {
	Data []ApiKeySummary `json:"data"`
}
