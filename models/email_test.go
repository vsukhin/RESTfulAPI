package models

import (
	"testing"
	"time"
)

func TestViewApiEmail(t *testing.T) {
	var email = "test@email.com"
	var primary = true
	var confirmed = true
	var subscription = true
	var language = "eng"
	var viewApiEmail *ViewApiEmail

	viewApiEmail = NewViewApiEmail(email, primary, confirmed, subscription, language)
	if viewApiEmail.Email != email {
		t.Error("Email field is not properly initialized")
	}
	if viewApiEmail.Primary != primary {
		t.Error("Primary field is not properly initialized")
	}
	if viewApiEmail.Confirmed != confirmed {
		t.Error("Confirmed field is not properly initialized")
	}
	if viewApiEmail.Subscription != subscription {
		t.Error("Subscription field is not properly initialized")
	}
	if viewApiEmail.Language != language {
		t.Error("Language field is not properly initialized")
	}
}

func TestNewDtoEmail(t *testing.T) {
	var email = "test@email.com"
	var userid int64 = 1
	var created = time.Now()
	var primary = true
	var confirmed = true
	var subscription = true
	var code = "12345"
	var language = "eng"
	var exists = true
	var dtoEmail *DtoEmail

	dtoEmail = NewDtoEmail(email, userid, created, primary, confirmed, subscription, code, language, exists)
	if dtoEmail.Email != email {
		t.Error("Email field is not properly initialized")
	}
	if dtoEmail.UserID != userid {
		t.Error("User id field is not properly initialized")
	}
	if dtoEmail.Created != created {
		t.Error("Created field is not properly initialized")
	}
	if dtoEmail.Primary != primary {
		t.Error("Primary field is not properly initialized")
	}
	if dtoEmail.Confirmed != confirmed {
		t.Error("Confirmed field is not properly initialized")
	}
	if dtoEmail.Subscription != subscription {
		t.Error("Subscription field is not properly initialized")
	}
	if dtoEmail.Code != code {
		t.Error("Code field is not properly initialized")
	}
	if dtoEmail.Language != language {
		t.Error("Language field is not properly initialized")
	}
	if dtoEmail.Exists != exists {
		t.Error("Exists field is not properly initialized")
	}
}
