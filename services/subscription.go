package services

import (
	"application/models"
)

type SubscriptionRepository interface {
	Exists(email string) (found bool, err error)
	FindBySubscrCode(code string) (dtosubscription *models.DtoSubscription, err error)
	FindByUnsubscrCode(code string) (dtosubscription *models.DtoSubscription, err error)
	Get(email string) (dtosubscription *models.DtoSubscription, err error)
	Create(dtosubscription *models.DtoSubscription) (err error)
	Update(dtosubscription *models.DtoSubscription) (err error)
	Delete(email string) (err error)
}

type SubscriptionService struct {
	*Repository
}

func NewSubscriptionService(repository *Repository) *SubscriptionService {
	repository.DbContext.AddTableWithName(models.DtoSubscription{}, repository.Table).SetKeys(false, "email")
	return &SubscriptionService{
		repository,
	}
}

func (subscriptionservice *SubscriptionService) Exists(email string) (found bool, err error) {
	var count int64
	count, err = subscriptionservice.DbContext.SelectInt("select count(*) from "+subscriptionservice.Table+" where email = ?", email)
	if err != nil {
		log.Error("Error during getting subscription object from database %v with value %v", err, email)
		return false, err
	}

	return count != 0, nil
}

func (subscriptionservice *SubscriptionService) FindBySubscrCode(code string) (dtosubscription *models.DtoSubscription, err error) {
	dtosubscription = new(models.DtoSubscription)
	err = subscriptionservice.DbContext.SelectOne(dtosubscription, "select * from "+subscriptionservice.Table+" where subscr_code = ?", code)
	if err != nil {
		log.Error("Error during finding subscription object in database %v with value %v", err, code)
		return nil, err
	}

	return dtosubscription, nil
}

func (subscriptionservice *SubscriptionService) FindByUnsubscrCode(code string) (dtosubscription *models.DtoSubscription, err error) {
	dtosubscription = new(models.DtoSubscription)
	err = subscriptionservice.DbContext.SelectOne(dtosubscription, "select * from "+subscriptionservice.Table+" where unsubscr_code = ?", code)
	if err != nil {
		log.Error("Error during finding subscription object in database %v with value %v", err, code)
		return nil, err
	}

	return dtosubscription, nil
}

func (subscriptionservice *SubscriptionService) Get(email string) (dtosubscription *models.DtoSubscription, err error) {
	dtosubscription = new(models.DtoSubscription)
	err = subscriptionservice.DbContext.SelectOne(dtosubscription, "select * from "+subscriptionservice.Table+" where email = ?", email)
	if err != nil {
		log.Error("Error during getting subscription object from database %v with value %v", err, email)
		return nil, err
	}

	return dtosubscription, nil
}

func (subscriptionservice *SubscriptionService) Create(dtosubscription *models.DtoSubscription) (err error) {
	err = subscriptionservice.DbContext.Insert(dtosubscription)
	if err != nil {
		log.Error("Error during creating subscription object in database %v", err)
		return err
	}

	return nil
}

func (subscriptionservice *SubscriptionService) Update(dtosubscription *models.DtoSubscription) (err error) {
	_, err = subscriptionservice.DbContext.Update(dtosubscription)
	if err != nil {
		log.Error("Error during updating subscription object in database %v with value %v", err, dtosubscription.Email)
		return err
	}

	return nil
}

func (subscriptionservice *SubscriptionService) Delete(email string) (err error) {
	_, err = subscriptionservice.DbContext.Exec("delete from "+subscriptionservice.Table+" where email = ?", email)
	if err != nil {
		log.Error("Error during deleting subscription object in database %v with value %v", err, email)
		return err
	}

	return nil
}
