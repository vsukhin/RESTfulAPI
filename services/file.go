package services

import (
	"application/config"
	"application/models"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type FileRepository interface {
	Get(fileid int64) (file *models.DtoFile, err error)
	GetBriefly(fileid int64) (file *models.DtoFile, err error)
	GetExpired(timeout time.Duration) (files *[]models.DtoFile, err error)
	FindByType(filetype string) (file *models.DtoFile, err error)
	Create(file *models.DtoFile, data *models.ViewFile) (err error)
	Update(file *models.DtoFile) (err error)
	Delete(file *models.DtoFile) (err error)
}

type FileService struct {
	*Repository
}

func NewFileService(repository *Repository) *FileService {
	repository.DbContext.AddTableWithName(models.DtoFile{}, repository.Table).SetKeys(true, "id")
	return &FileService{
		repository,
	}
}

func (fileservice *FileService) Get(fileid int64) (file *models.DtoFile, err error) {
	file = new(models.DtoFile)
	err = fileservice.DbContext.SelectOne(file, "select * from "+fileservice.Table+" where id = ?", fileid)
	if err != nil {
		log.Error("Error during getting file object from database %v with value %v", err, fileid)
		return nil, err
	}

	file.FileData, err = ioutil.ReadFile(filepath.Join(config.Configuration.FileStorage, file.Path, fmt.Sprintf("%08d", file.ID)))
	if err != nil {
		log.Error("Error during getting file in filesystem %v with value %v", err, file.ID)
		return nil, err
	}

	return file, nil
}

func (fileservice *FileService) GetBriefly(fileid int64) (file *models.DtoFile, err error) {
	file = new(models.DtoFile)
	err = fileservice.DbContext.SelectOne(file, "select * from "+fileservice.Table+" where id = ?", fileid)
	if err != nil {
		log.Error("Error during getting briefly file object from database %v with value %v", err, fileid)
		return nil, err
	}

	return file, nil
}

func (fileservice *FileService) GetExpired(timeout time.Duration) (files *[]models.DtoFile, err error) {
	files = new([]models.DtoFile)
	_, err = fileservice.DbContext.Select(files, "select * from "+fileservice.Table+" where created < ?"+
		" and id not in (select file_id from input_files union select file_id from appendices union select file_id from contracts"+
		" union select file_id from documents union select picDisable_id from services union select picNormal_id from services"+
		" union select picOver_id from services union select picSoon_id from services)", time.Now().Add(-timeout))
	if err != nil {
		log.Error("Error during getting file object from database %v with value %v", err, timeout)
		return nil, err
	}

	return files, nil
}

func (fileservice *FileService) FindByType(filetype string) (file *models.DtoFile, err error) {
	file = new(models.DtoFile)
	err = fileservice.DbContext.SelectOne(file, "select * from "+fileservice.Table+" where name = ?", filetype)
	if err != nil {
		log.Error("Error during getting file object from database %v with value %v", err, filetype)
		return nil, err
	}

	return file, nil
}

func (fileservice *FileService) Create(file *models.DtoFile, data *models.ViewFile) (err error) {
	var srcfile multipart.File
	var dstfile *os.File

	err = fileservice.DbContext.Insert(file)
	if err != nil {
		log.Error("Error during creating file object in database %v", err)
		return err
	}

	fullpath := filepath.Join(config.Configuration.FileStorage, file.Path)
	if _, err = os.Stat(fullpath); os.IsNotExist(err) {
		err = os.MkdirAll(fullpath, 0777)
		if err != nil {
			log.Error("Error during directory creatiion %v with value %v", err, fullpath)
			return err
		}
	}

	dstfile, err = os.Create(filepath.Join(fullpath, fmt.Sprintf("%08d", file.ID)))
	defer dstfile.Close()
	if err != nil {
		log.Error("Error during creating file in filesystem %v with value %v", err, file.ID)
		return err
	}

	if data != nil {
		srcfile, err = data.FileData.Open()
		defer srcfile.Close()
		if err != nil {
			log.Error("Error during creating file in filesystem %v", err)
			return err
		}

		if _, err = io.Copy(dstfile, srcfile); err != nil {
			log.Error("Error during creating file in filesystem %v", err)
			return err
		}
	}

	return nil
}

func (fileservice *FileService) Update(file *models.DtoFile) (err error) {
	_, err = fileservice.DbContext.Update(file)
	if err != nil {
		log.Error("Error during updating file object in database %v with value %v", err, file.ID)
		return err
	}

	return nil
}

func (fileservice *FileService) Delete(file *models.DtoFile) (err error) {
	_, err = fileservice.DbContext.Exec("delete from "+fileservice.Table+" where id = ?", file.ID)
	if err != nil {
		log.Error("Error during deleting file object in database %v with value %v", err, file.ID)
		return err
	}

	err = os.Remove(filepath.Join(config.Configuration.FileStorage, file.Path, fmt.Sprintf("%08d", file.ID)))
	if err != nil {
		log.Error("Error during deleting file in filesystem %v with value %v", err, file.ID)
		return err
	}

	return nil
}
