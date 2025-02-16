package mailgun

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addValidationRoutes(r chi.Router) {
	r.Get("/v3/address/validate", ms.validateEmail)
	r.Get("/v3/address/parse", ms.parseEmail)
	r.Get("/v3/address/private/validate", ms.validateEmail)
	r.Get("/v3/address/private/parse", ms.parseEmail)
	r.Get("/v4/address/validate", ms.validateEmailV4)
}

func (ms *mockServer) validateEmailV4(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("address") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'address' parameter is required"})
		return
	}

	var results v4EmailValidationResp
	results.Risk = "unknown"
	parts, err := mail.ParseAddress(r.FormValue("address"))
	if err == nil {
		results.Risk = "low"
		results.Parts.Domain = strings.Split(parts.Address, "@")[1]
		results.Parts.LocalPart = strings.Split(parts.Address, "@")[0]
		results.Parts.DisplayName = parts.Name
	}
	results.Reason = []string{"no-reason"}
	results.Result = "deliverable"
	results.Engagement = &EngagementData{
		Engaging: false,
		Behavior: "disengaged",
		IsBot:    false,
	}
	toJSON(w, results)
}

func (ms *mockServer) validateEmail(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("address") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'address' parameter is required"})
		return
	}

	var results EmailVerification
	parts, err := mail.ParseAddress(r.FormValue("address"))
	if err == nil {
		results.IsValid = true
		results.Parts.Domain = strings.Split(parts.Address, "@")[1]
		results.Parts.LocalPart = strings.Split(parts.Address, "@")[0]
		results.Parts.DisplayName = parts.Name
	}
	results.Reason = "no-reason"
	results.Risk = "unknown"
	toJSON(w, results)
}

func (ms *mockServer) parseEmail(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("addresses") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'addresses' parameter is required"})
		return
	}

	addresses := strings.Split(r.FormValue("addresses"), ",")

	var results addressParseResult
	for _, address := range addresses {
		_, err := mail.ParseAddress(address)
		if err != nil {
			results.Unparseable = append(results.Unparseable, address)
		} else {
			results.Parsed = append(results.Parsed, address)
		}
	}
	toJSON(w, results)
}
