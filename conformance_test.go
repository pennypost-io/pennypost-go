package pennypost

// Conformance runner: replays ../conformance/vectors.json against the memory API.
// Start it first: bash sdk/serve-test.sh start

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

const base = "http://127.0.0.1:8799"

var key = "pp_live_" + strings.Repeat("a", 48)

type vec struct {
	ID          string            `json:"id"`
	Op          string            `json:"op"`
	Req         json.RawMessage   `json:"req"`
	RawReq      json.RawMessage   `json:"raw_req"`
	Params      map[string]any    `json:"params"`
	After       string            `json:"after"`
	Auth        string            `json:"auth"`
	Idem        string            `json:"idempotency_key"`
	Repeat      int               `json:"repeat"`
	Expect      map[string]any    `json:"expect"`
	ExpectError map[string]any    `json:"expect_error"`
}

func strParams(m map[string]any) map[string]string {
	out := map[string]string{}
	for k, v := range m {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func toMap(t *testing.T, v any) map[string]any {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func runOp(t *testing.T, c *Client, v vec, ctx map[string]map[string]any) (map[string]any, error) {
	switch v.Op {
	case "sendEmail":
		raw := v.Req
		if raw == nil {
			raw = v.RawReq
		}
		if v.RawReq != nil {
			// invalid-shape vector: bypass the typed struct on purpose
			var e *Error
			err := c.request("POST", "/v1/emails", json.RawMessage(raw), nil, v.Idem, &map[string]any{})
			if err == nil {
				return map[string]any{}, nil
			}
			if es, ok := err.(*Error); ok {
				e = es
			}
			return nil, e
		}
		var req SendEmailRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			t.Fatal(err)
		}
		r, err := c.SendEmail(&req, v.Idem)
		if err != nil {
			return nil, err
		}
		return toMap(t, r), nil
	case "getEmail":
		id, _ := v.Params["id"].(string)
		if id == "" {
			acc := ctx[v.After]["accepted"].([]any)
			id = acc[0].(map[string]any)["id"].(string)
		}
		r, err := c.GetEmail(id)
		if err != nil {
			return nil, err
		}
		return toMap(t, r), nil
	case "listEmails":
		r, err := c.ListEmails(strParams(v.Params))
		if err != nil {
			return nil, err
		}
		return toMap(t, r), nil
	case "listSuppressions":
		r, err := c.ListSuppressions(strParams(v.Params))
		if err != nil {
			return nil, err
		}
		return toMap(t, r), nil
	case "addSuppression":
		var body map[string]string
		_ = json.Unmarshal(v.Req, &body)
		r, err := c.AddSuppression(body["email"])
		if err != nil {
			return nil, err
		}
		return toMap(t, r), nil
	case "removeSuppression":
		r, err := c.RemoveSuppression(v.Params["email"].(string))
		if err != nil {
			return nil, err
		}
		return toMap(t, r), nil
	}
	t.Fatalf("unknown op %s", v.Op)
	return nil, nil
}

func TestConformance(t *testing.T) {
	raw, err := os.ReadFile("../conformance/vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Vectors []vec `json:"vectors"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	ctx := map[string]map[string]any{}
	for _, v := range doc.Vectors {
		c := New(key, WithBaseURL(base))
		if v.Auth != "" {
			c = New(v.Auth, WithBaseURL(base))
		}
		if v.ExpectError != nil {
			_, err := runOp(t, c, v, ctx)
			apiErr, ok := err.(*Error)
			if !ok {
				t.Fatalf("%s: expected API error, got %v", v.ID, err)
			}
			if s, ok := v.ExpectError["status"].(float64); ok && apiErr.Status != int(s) {
				t.Fatalf("%s: status %d != %v", v.ID, apiErr.Status, s)
			}
			if code, ok := v.ExpectError["code"].(string); ok && apiErr.Code != code {
				t.Fatalf("%s: code %s", v.ID, apiErr.Code)
			}
			if p, ok := v.ExpectError["param"].(string); ok && apiErr.Param != p {
				t.Fatalf("%s: param %s", v.ID, apiErr.Param)
			}
			fmt.Println("  ok", v.ID)
			continue
		}
		reps := v.Repeat
		if reps == 0 {
			reps = 1
		}
		var results []map[string]any
		for i := 0; i < reps; i++ {
			r, err := runOp(t, c, v, ctx)
			if err != nil {
				t.Fatalf("%s: %v", v.ID, err)
			}
			results = append(results, r)
		}
		r := results[0]
		ctx[v.ID] = r
		e := v.Expect
		arr := func(k string) []any { a, _ := r[k].([]any); return a }
		if n, ok := e["accepted_len"].(float64); ok && len(arr("accepted")) != int(n) {
			t.Fatalf("%s: accepted %v", v.ID, r)
		}
		if n, ok := e["suppressed_len"].(float64); ok && len(arr("suppressed")) != int(n) {
			t.Fatalf("%s: suppressed", v.ID)
		}
		if p, ok := e["id_prefix"].(string); ok {
			id := arr("accepted")[0].(map[string]any)["id"].(string)
			if !strings.HasPrefix(id, p) {
				t.Fatalf("%s: id prefix", v.ID)
			}
		}
		if b, _ := e["same_id_across_repeats"].(bool); b {
			id0 := results[0]["accepted"].([]any)[0].(map[string]any)["id"]
			id1 := results[1]["accepted"].([]any)[0].(map[string]any)["id"]
			if id0 != id1 {
				t.Fatalf("%s: idempotency", v.ID)
			}
		}
		if s, ok := e["subject"].(string); ok && r["subject"] != s {
			t.Fatalf("%s: subject", v.ID)
		}
		if b, _ := e["has_events_array"].(bool); b {
			if _, ok := r["events"].([]any); !ok {
				t.Fatalf("%s: events", v.ID)
			}
		}
		if n, ok := e["data_len"].(float64); ok && len(arr("data")) != int(n) {
			t.Fatalf("%s: data len %v", v.ID, r)
		}
		if s, ok := e["first_to"].(string); ok && arr("data")[0].(map[string]any)["to"] != s {
			t.Fatalf("%s: first to", v.ID)
		}
		if s, ok := e["reason"].(string); ok && r["reason"] != s {
			t.Fatalf("%s: reason", v.ID)
		}
		if s, ok := e["suppressed_reason"].(string); ok && arr("suppressed")[0].(map[string]any)["reason"] != s {
			t.Fatalf("%s: sup reason", v.ID)
		}
		if want, ok := e["contains_email"].(string); ok {
			found := false
			for _, s := range arr("data") {
				if s.(map[string]any)["email"] == want {
					found = true
				}
			}
			if !found {
				t.Fatalf("%s: contains", v.ID)
			}
		}
		if b, ok := e["removed"].(bool); ok && r["removed"] != b {
			t.Fatalf("%s: removed", v.ID)
		}
		fmt.Println("  ok", v.ID)
	}
}
