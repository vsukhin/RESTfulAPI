package services

import (
	"application/config"
	"application/models"
	"net/smtp"
	"strconv"
)

type EmailRepository interface {
	SendEmail(email string, subject string, body string) (err error)
	Exists(email string) (found bool, err error)
	FindByCode(code string) (email *models.DtoEmail, err error)
	Get(email string) (dtoemail *models.DtoEmail, err error)
	GetByUser(userid int64) (emails *[]models.DtoEmail, err error)
	Create(email *models.DtoEmail) (err error)
	Update(email *models.DtoEmail) (err error)
	Delete(email string) (err error)
	DeleteByUser(userid int64) (err error)
}

type EmailService struct {
	*Repository
}

func NewEmailService(repository *Repository) *EmailService {
	repository.DbContext.AddTableWithName(models.DtoEmail{}, repository.Table).SetKeys(false, "email")
	return &EmailService{
		repository,
	}
}

func (emailservice *EmailService) SendEmail(email string, subject string, body string) (err error) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		config.Configuration.Mail.Login,
		config.Configuration.Mail.Password,
		config.Configuration.Mail.Host,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	data := "To: " + email + "\r\n From: " + config.Configuration.Mail.Sender + "\r\nSubject: " + subject + "\r\n\r\n" + body

	err = smtp.SendMail(
		config.Configuration.Mail.Host+":"+strconv.Itoa(config.Configuration.Mail.Port),
		auth,
		config.Configuration.Mail.Sender,
		[]string{email},
		[]byte(data),
	)
	if err != nil {
		log.Error("Error during sending email %v with value %v", err, email)
	}

	return err
}

func (emailservice *EmailService) Exists(email string) (found bool, err error) {
	var count int64
	count, err = emailservice.DbContext.SelectInt("select count(*) from "+emailservice.Table+" where email = ?", email)
	if err != nil {
		log.Error("Error during getting email object from database %v with value %v", err, email)
		return false, err
	}

	return count != 0, nil
}

func (emailservice *EmailService) FindByCode(code string) (email *models.DtoEmail, err error) {
	email = new(models.DtoEmail)
	err = emailservice.DbContext.SelectOne(email, "select * from "+emailservice.Table+" where code = ?", code)
	if err != nil {
		log.Error("Error during finding email object in database %v with value %v", err, code)
		return nil, err
	}

	return email, nil
}

func (emailservice *EmailService) Get(email string) (dtoemail *models.DtoEmail, err error) {
	dtoemail = new(models.DtoEmail)
	err = emailservice.DbContext.SelectOne(dtoemail, "select * from "+emailservice.Table+" where email = ?", email)
	if err != nil {
		log.Error("Error during getting email object from database %v with value %v", err, email)
		return nil, err
	}

	return dtoemail, nil
}

func (emailservice *EmailService) GetByUser(userid int64) (emails *[]models.DtoEmail, err error) {
	emails = new([]models.DtoEmail)
	_, err = emailservice.DbContext.Select(emails, "select * from "+emailservice.Table+" where user_id = ?", userid)
	if err != nil {
		log.Error("Error during getting email object from database %v with value %v", err, userid)
		return nil, err
	}

	return emails, nil
}

func (emailservice *EmailService) Create(email *models.DtoEmail) (err error) {
	err = emailservice.DbContext.Insert(email)
	if err != nil {
		log.Error("Error during creating email object in database %v", err)
		return err
	}

	return nil
}

func (emailservice *EmailService) Update(email *models.DtoEmail) (err error) {
	_, err = emailservice.DbContext.Update(email)
	if err != nil {
		log.Error("Error during updating email object in database %v with value %v", err, email.Email)
		return err
	}

	return nil
}

func (emailservice *EmailService) Delete(email string) (err error) {
	_, err = emailservice.DbContext.Exec("delete from "+emailservice.Table+" where email = ?", email)
	if err != nil {
		log.Error("Error during deleting email object in database %v with value %v", err, email)
		return err
	}

	return nil
}

func (emailservice *EmailService) DeleteByUser(userid int64) (err error) {
	_, err = emailservice.DbContext.Exec("delete from "+emailservice.Table+" where user_id = ?", userid)
	if err != nil {
		log.Error("Error during deleting email object in database %v with value %v", err, userid)
		return err
	}

	return nil
}
