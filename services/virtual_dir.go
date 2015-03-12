package services

import (
	"application/models"
)

type VirtualDirService struct {
	*Repository
}

func NewVirtualDirService(repository *Repository) *VirtualDirService {
	repository.DbContext.AddTableWithName(models.DtoVirtualDir{}, repository.Table).SetKeys(false, "token")
	return &VirtualDirService{
		repository,
	}
}

func (virtualdirservice *VirtualDirService) Get(token string) (virtualdir *models.DtoVirtualDir, err error) {
	virtualdir = new(models.DtoVirtualDir)
	err = virtualdirservice.DbContext.SelectOne(virtualdir, "select * from "+virtualdirservice.Table+" where token = ?", token)
	if err != nil {
		log.Error("Error during getting virtual dir object from database %v with value %v", err, token)
		return nil, err
	}

	return virtualdir, nil
}

func (virtualdirservice *VirtualDirService) Create(virtualdir *models.DtoVirtualDir) (err error) {
	err = virtualdirservice.DbContext.Insert(virtualdir)
	if err != nil {
		log.Error("Error during creating virtual dir object in database %v", err)
		return err
	}

	return nil
}
