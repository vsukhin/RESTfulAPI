package services

import (
	"application/models"
)

type PaymentRepository interface {
	Exists(unit_id int64) (found bool, err error)
	Get(unit_id int64) (payment *models.DtoPayment, err error)
	Create(payment *models.DtoPayment) (err error)
	Update(payment *models.DtoPayment) (err error)
	Save(payment *models.DtoPayment) (err error)
}

type PaymentService struct {
	*Repository
}

func NewPaymentService(repository *Repository) *PaymentService {
	repository.DbContext.AddTableWithName(models.DtoPayment{}, repository.Table).SetKeys(false, "unit_id")
	return &PaymentService{Repository: repository}
}

func (paymentservice *PaymentService) Exists(unit_id int64) (found bool, err error) {
	var count int64
	count, err = paymentservice.DbContext.SelectInt("select count(*) from "+paymentservice.Table+
		" where unit_id = ?", unit_id)
	if err != nil {
		log.Error("Error during checking payment object in database %v with value %v", err, unit_id)
		return false, err
	}

	return count != 0, nil
}

func (paymentservice *PaymentService) Get(unit_id int64) (payment *models.DtoPayment, err error) {
	payment = new(models.DtoPayment)
	err = paymentservice.DbContext.SelectOne(payment, "select * from "+paymentservice.Table+" where unit_id = ?", unit_id)
	if err != nil {
		log.Error("Error during getting payment object from database %v with value %v", err, unit_id)
		return nil, err
	}

	return payment, nil
}

func (paymentservice *PaymentService) Create(payment *models.DtoPayment) (err error) {
	err = paymentservice.DbContext.Insert(payment)
	if err != nil {
		log.Error("Error during creating payment object in database %v", err)
		return err
	}

	return nil
}

func (paymentservice *PaymentService) Update(payment *models.DtoPayment) (err error) {
	_, err = paymentservice.DbContext.Update(payment)
	if err != nil {
		log.Error("Error during updating payment object in database %v with value %v", err, payment.Unit_ID)
		return err
	}

	return nil
}

func (paymentservice *PaymentService) Save(payment *models.DtoPayment) (err error) {
	count, err := paymentservice.DbContext.SelectInt("select count(*) from "+paymentservice.Table+
		" where unit_id = ?", payment.Unit_ID)
	if err != nil {
		log.Error("Error during saving payment object in database %v with value %v", err, payment.Unit_ID)
		return err
	}
	if count == 0 {
		err = paymentservice.Create(payment)
	} else {
		err = paymentservice.Update(payment)
	}

	return err
}
