// Generated from api/openapi.json (sha 7c58f5917c4b) by sdk/generate.ts — do not edit by hand.
// PennyPost: the affordable email API. https://pennypost.io/docs
package pennypost

const SpecSHA = "7c58f5917c4b"

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
