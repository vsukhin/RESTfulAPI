package models

import (
	"testing"
	"time"
)

func TestApiSession(t *testing.T) {
	var timeout = time.Now()
	var accesstoken = "1234567890"
	var аpiSession *ApiSession

	аpiSession = NewApiSession(timeout, accesstoken)
	if аpiSession.Timeout != timeout {
		t.Error("Timeout field is not properly initialized")
	}
	if аpiSession.AccessToken != accesstoken {
		t.Error("AccessToken field is not properly initialized")
	}
}

func TestNewDtoSession(t *testing.T) {
	var accesstoken = "1234567890"
	var userid int64 = 1
	var roles = []UserRole{USER_ROLE_ADMINISTRATOR, USER_ROLE_DEVELOPER}
	var lastactivity = time.Now()
	var language = "eng"
	var dtoSession *DtoSession

	dtoSession = NewDtoSession(accesstoken, userid, roles, lastactivity, language)
	if dtoSession.AccessToken != accesstoken {
		t.Error("AccessToken field is not properly initialized")
	}
	if dtoSession.UserID != userid {
		t.Error("UserID field is not properly initialized")
	}
	if len(dtoSession.Roles) != len(roles) {
		t.Error("Roles field is not properly initialized")
	} else {
		for i := 0; i < len(dtoSession.Roles); i++ {
			if dtoSession.Roles[i] != roles[i] {
				t.Error("Roles field is not properly initialized")
				break
			}
		}
	}
	if dtoSession.LastActivity != lastactivity {
		t.Error("LastActivity field is not properly initialized")
	}
	if dtoSession.Language != language {
		t.Error("Language field is not properly initialized")
	}
}
