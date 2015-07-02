package services

import (
	"application/config"
	"application/models"
	"crypto/tls"
	"github.com/coopernurse/gorp"
	"net"
	"net/smtp"
	"strconv"
)

type EmailRepository interface {
	SendEmail(email string, subject string, body string, headers string, from string) (err error)
	SendHTML(email string, subject string, body string, headers string, from string) (err error)
	Exists(email string) (found bool, err error)
	FindByCode(code string) (email *models.DtoEmail, err error)
	Get(email string) (dtoemail *models.DtoEmail, err error)
	GetByUser(userid int64) (emails *[]models.DtoEmail, err error)
	Create(email *models.DtoEmail, trans *gorp.Transaction) (err error)
	Update(email *models.DtoEmail, trans *gorp.Transaction) (err error)
	Delete(email string, trans *gorp.Transaction) (err error)
	DeleteByUser(userid int64, trans *gorp.Transaction) (err error)
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

/*
 * server at addr, switches to TLS if possible,
 * authenticates with mechanism a if possible, and then sends an email from
 * address from, to addresses to, with message msg.
 */
func SendEmailTLS(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	host, _, _ := net.SplitHostPort(addr)
	config := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	if err = c.Auth(a); err != nil {
		return err
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func (emailservice *EmailService) SendEmail(email string, subject string, body string, headers string, from string) (err error) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		config.Configuration.Mail.Login,
		config.Configuration.Mail.Password,
		config.Configuration.Mail.Host,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	data := ""
	if headers != "" {
		data += headers + "\r\n"
	}
	if from == "" {
		from = config.Configuration.Mail.Sender
	}
	data += "From: " + from + "\r\nTo: " + email + "\r\nSubject: " + subject + "\r\n\r\n" + body

	err = SendEmailTLS(config.Configuration.Mail.Host+":"+strconv.Itoa(config.Configuration.Mail.Port),
		auth,
		from,
		[]string{email},
		[]byte(data))
	if err != nil {
		log.Error("Error during sending email %v with value %v", err, email)
	}

	return err
}

func (emailservice *EmailService) SendHTML(email string, subject string, body string, headers string, from string) (err error) {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		config.Configuration.Mail.Login,
		config.Configuration.Mail.Password,
		config.Configuration.Mail.Host,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	data := ""
	if headers != "" {
		data += headers + "\r\n"
	}
	if from == "" {
		from = config.Configuration.Mail.Sender
	}
	data += "MIME-Version: 1.0\r\nContent-Type: text/html; charset=utf-8\r\nFrom: " + from + "\r\nTo: " + email + "\r\nSubject: " + subject + "\r\n\r\n" + body

	err = SendEmailTLS(config.Configuration.Mail.Host+":"+strconv.Itoa(config.Configuration.Mail.Port),
		auth,
		from,
		[]string{email},
		[]byte(data))
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

func (emailservice *EmailService) Create(email *models.DtoEmail, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(email)
	} else {
		err = emailservice.DbContext.Insert(email)
	}
	if err != nil {
		log.Error("Error during creating email object in database %v", err)
		return err
	}

	return nil
}

func (emailservice *EmailService) Update(email *models.DtoEmail, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Update(email)
	} else {
		_, err = emailservice.DbContext.Update(email)
	}
	if err != nil {
		log.Error("Error during updating email object in database %v with value %v", err, email.Email)
		return err
	}

	return nil
}

func (emailservice *EmailService) Delete(email string, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+emailservice.Table+" where email = ?", email)
	} else {
		_, err = emailservice.DbContext.Exec("delete from "+emailservice.Table+" where email = ?", email)
	}
	if err != nil {
		log.Error("Error during deleting email object in database %v with value %v", err, email)
		return err
	}

	return nil
}

func (emailservice *EmailService) DeleteByUser(userid int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+emailservice.Table+" where user_id = ?", userid)
	} else {
		_, err = emailservice.DbContext.Exec("delete from "+emailservice.Table+" where user_id = ?", userid)
	}
	if err != nil {
		log.Error("Error during deleting email object in database %v with value %v", err, userid)
		return err
	}

	return nil
}
