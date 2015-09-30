package workflows

import (
	"application/config"
	"application/services"
	"time"
)

type FileWorkflow struct {
	FileRepository services.FileRepository
}

func NewFileWorkflow(filerepository services.FileRepository) *FileWorkflow {
	return &FileWorkflow{
		FileRepository: filerepository,
	}
}

func (fileworkflow *FileWorkflow) ClearExpired() {
	for {
		files, err := fileworkflow.FileRepository.GetExpired(config.Configuration.FileTimeout)
		if err == nil {
			for _, file := range *files {
				if !file.Permanent {
					err = fileworkflow.FileRepository.Delete(&file)
				}
			}
		}
		time.Sleep(time.Minute)
	}
}
