package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type InputFtpRepository interface {
	Exists(order_id int64) (found bool, err error)
	Get(order_id int64) (inputftp *models.DtoInputFtp, err error)
	Create(dtoinputftp *models.DtoInputFtp, trans *gorp.Transaction) (err error)
	Delete(order_id int64, trans *gorp.Transaction) (err error)
}

type InputFtpService struct {
	*Repository
}

func NewInputFtpService(repository *Repository) *InputFtpService {
	repository.DbContext.AddTableWithName(models.DtoInputFtp{}, repository.Table).SetKeys(false, "order_id")
	return &InputFtpService{Repository: repository}
}

func (inputftpservice *InputFtpService) Exists(order_id int64) (found bool, err error) {
	var count int64
	count, err = inputftpservice.DbContext.SelectInt("select count(*) from "+inputftpservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting input ftp object from database %v with value %v", err, order_id)
		return false, err
	}

	return count != 0, nil
}

func (inputftpservice *InputFtpService) Get(order_id int64) (inputftp *models.DtoInputFtp, err error) {
	inputftp = new(models.DtoInputFtp)
	err = inputftpservice.DbContext.SelectOne(inputftp, "select * from "+inputftpservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting input ftp object from database %v with value %v", err, order_id)
		return nil, err
	}

	return inputftp, nil
}

func (inputftpservice *InputFtpService) Create(dtoinputftp *models.DtoInputFtp, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoinputftp)
	} else {
		err = inputftpservice.DbContext.Insert(dtoinputftp)
	}
	if err != nil {
		log.Error("Error during creating input ftp object in database %v", err)
		return err
	}

	return nil
}

func (inputftpservice *InputFtpService) Delete(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+inputftpservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = inputftpservice.DbContext.Exec("delete from "+inputftpservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting input ftp object for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
