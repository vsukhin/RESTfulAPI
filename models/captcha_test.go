package models

import (
	"testing"
	"time"
)

func TestNewApiCaptcha(t *testing.T) {
	var hash = "1234567890"
	var image = "ABCDEF"
	var apiCaptcha *ApiCaptcha

	apiCaptcha = NewApiCaptcha(hash, image)
	if apiCaptcha.Hash != hash {
		t.Error("Hash field is not properly initialized")
	}
	if apiCaptcha.Image != image {
		t.Error("Image field is not properly initialized")
	}
}

func TestNewDtoCaptcha(t *testing.T) {
	var hash = "1234567890"
	var image = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	var value = "123456"
	var created = time.Now()
	var inuse = true
	var dtoCaptcha *DtoCaptcha

	dtoCaptcha = NewDtoCaptcha(hash, image, value, created, inuse)
	if dtoCaptcha.Hash != hash {
		t.Error("Hash field is not properly initialized")
	}
	if len(dtoCaptcha.Image) != len(image) {
		t.Error("Image field is not properly initialized")
	} else {
		for i := 0; i < len(dtoCaptcha.Image); i++ {
			if dtoCaptcha.Image[i] != image[i] {
				t.Error("Image field is not properly initialized")
				break
			}
		}
	}
	if dtoCaptcha.Value != value {
		t.Error("Value field is not properly initialized")
	}
	if dtoCaptcha.Created != created {
		t.Error("Created field is not properly initialized")
	}
	if dtoCaptcha.InUse != inuse {
		t.Error("InUse field is not properly initialized")
	}
}
