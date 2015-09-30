package models

import (
	"testing"
	"time"
)

func TestNewApiFile(t *testing.T) {
	var id int64 = 1
	var apiFile *ApiFile

	apiFile = NewApiFile(id)
	if apiFile.ID != id {
		t.Error("ID field is not properly initialized")
	}
}

func TestNewApiImage(t *testing.T) {
	var id int64 = 1
	var apiImage *ApiImage

	apiImage = NewApiImage(id)
	if apiImage.ID != id {
		t.Error("ID field is not properly initialized")
	}
}

func TestNewDtoFile(t *testing.T) {
	var id int64 = 1
	var name = "facility"
	var path = "/some/where/in"
	var created = time.Now()
	var permanent = true
	var export_ready = true
	var export_percentage byte = 50
	var export_object_id int64 = 1
	var export_error = true
	var export_errordescription = "description"
	var filedata = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	var dtoFile *DtoFile

	dtoFile = NewDtoFile(id, name, path, created, permanent, export_ready, export_percentage, export_object_id,
		export_error, export_errordescription, filedata)
	if dtoFile.ID != id {
		t.Error("ID field is not properly initialized")
	}
	if dtoFile.Name != name {
		t.Error("Name field is not properly initialized")
	}
	if dtoFile.Path != path {
		t.Error("Description field is not properly initialized")
	}
	if dtoFile.Created != created {
		t.Error("Created field is not properly initialized")
	}
	if dtoFile.Permanent != permanent {
		t.Error("Permanent field is not properly initialized")
	}
	if dtoFile.Export_Ready != export_ready {
		t.Error("Ready field is not properly initialized")
	}
	if dtoFile.Export_Percentage != export_percentage {
		t.Error("Percentage field is not properly initialized")
	}
	if dtoFile.Export_Object_ID != export_object_id {
		t.Error("Object_ID field is not properly initialized")
	}
	if dtoFile.Export_Error != export_error {
		t.Error("Error field is not properly initialized")
	}
	if dtoFile.Export_ErrorDescription != export_errordescription {
		t.Error("Error description field is not properly initialized")
	}
	if len(dtoFile.FileData) != len(filedata) {
		t.Error("FileData field is not properly initialized")
	} else {
		for i := 0; i < len(dtoFile.FileData); i++ {
			if dtoFile.FileData[i] != filedata[i] {
				t.Error("FileData field is not properly initialized")
				break
			}
		}
	}
}
