package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type MessageRepository interface {
	Get(id int64) (message *models.DtoMessage, err error)
	IsReadByUser(user_id, id int64) (read bool, err error)
	IsViewed(id int64) (read bool, err error)
	IsLastForUser(user_id int64, id int64) (last bool, err error)
	GetByOrder(order_id int64, user_id int64, filter string) (messages *[]models.ApiLongMessage, err error)
	GetMetaByOrder(order_id int64, user_id int64) (message *models.ApiMetaMessage, err error)
	SetReadByUserForOrder(user_id int64, order_id int64) (err error)
	SetReadByUser(user_id int64, message_id int64) (err error)
	Create(message *models.DtoMessage, inTrans bool) (err error)
	Update(message *models.DtoMessage) (err error)
	Delete(message *models.DtoMessage, inTrans bool) (err error)
	DeleteByUser(user_id int64, trans *gorp.Transaction) (err error)
}

type MessageService struct {
	*Repository
}

func NewMessageService(repository *Repository) *MessageService {
	repository.DbContext.AddTableWithName(models.DtoMessage{}, repository.Table).SetKeys(true, "id")
	return &MessageService{
		repository,
	}
}

func (messageservice *MessageService) Get(id int64) (message *models.DtoMessage, err error) {
	message = new(models.DtoMessage)
	err = messageservice.DbContext.SelectOne(message, "select * from "+messageservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting message object from database %v with value %v", err, id)
		return nil, err
	}

	return message, nil
}

func (messageservice *MessageService) IsReadByUser(user_id, id int64) (read bool, err error) {
	count, err := messageservice.DbContext.SelectInt("select count(*) from user_messages where user_id = ? and message_id = ?", user_id, id)
	if err != nil {
		log.Error("Error during getting message object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (messageservice *MessageService) IsViewed(id int64) (read bool, err error) {
	count, err := messageservice.DbContext.SelectInt("select count(*) from "+messageservice.Table+
		" m inner join user_messages u on m.id = u.message_id where m.user_id != u.user_id and m.id = ?", id)
	if err != nil {
		log.Error("Error during getting message object from database %v with value %v", err, id)
		return false, err
	}

	return count != 0, nil
}

func (messageservice *MessageService) IsLastForUser(user_id int64, id int64) (last bool, err error) {
	count, err := messageservice.DbContext.SelectInt("select count(*) from "+messageservice.Table+" where id = ?"+
		" and id = (select id from "+messageservice.Table+" where user_id = ? order by created desc limit 1)", id, user_id)
	if err != nil {
		log.Error("Error during getting message object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (messageservice *MessageService) GetByOrder(order_id int64, user_id int64, filter string) (messages *[]models.ApiLongMessage, err error) {
	messages = new([]models.ApiLongMessage)
	_, err = messageservice.DbContext.Select(messages,
		"select m.id, m.created, m.user_id as userId, m.content as message, m.user_id = ? as isMine,"+
			"coalesce((select u.user_id from user_messages u where u.message_id = m.id and u.user_id = ?), 0) <> ? as new from "+messageservice.Table+
			" m where m.order_id = ?"+filter, user_id, user_id, user_id, order_id)
	if err != nil {
		log.Error("Error during getting all message object from database %v with value %v, %v", err, order_id, user_id)
		return nil, err
	}

	return messages, nil
}

func (messageservice *MessageService) GetMetaByOrder(order_id int64, user_id int64) (message *models.ApiMetaMessage, err error) {
	message = new(models.ApiMetaMessage)
	message.NumOfAll, err = messageservice.DbContext.SelectInt(
		"select count(*) from "+messageservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting meta message object from database %v with value %v", err, order_id)
		return nil, err
	}
	message.NumOfNew, err = messageservice.DbContext.SelectInt(
		"select count(*) from "+messageservice.Table+" where order_id = ? and id not in (select message_id from user_messages where user_id = ?)",
		order_id, user_id)
	if err != nil {
		log.Error("Error during getting meta message object from database %v with value %v, %v", err, order_id, user_id)
		return nil, err
	}

	return message, nil
}

func (messageservice *MessageService) SetReadByUserForOrder(user_id int64, order_id int64) (err error) {
	_, err = messageservice.DbContext.Exec("insert into user_messages (message_id, user_id) select ?, id from "+
		messageservice.Table+" where order_id = ? and id not in (select message_id from user_messages where user_id = ?)",
		user_id, order_id, user_id)
	if err != nil {
		log.Error("Error during setting message objects for user object in database %v with value %v, %v", err, order_id, user_id)
		return err
	}

	return nil
}

func (messageservice *MessageService) SetReadByUser(user_id int64, message_id int64) (err error) {
	_, err = messageservice.DbContext.Exec("insert into user_messages (message_id, user_id) select ?, ? from "+
		messageservice.Table+" where id = ? and id not in (select message_id from user_messages where user_id = ?)",
		message_id, user_id, message_id, user_id)
	if err != nil {
		log.Error("Error during setting message object for user object in database %v with value %v, %v", err, message_id, user_id)
		return err
	}

	return nil
}

func (messageservice *MessageService) Create(message *models.DtoMessage, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = messageservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating message object in database %v", err)
			return err
		}
	}

	if inTrans {
		err = trans.Insert(message)
	} else {
		err = messageservice.DbContext.Insert(message)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating message object in database %v", err)
		return err
	}

	if inTrans {
		_, err = trans.Exec("insert into user_messages (message_id, user_id) values (?, ?)", message.ID, message.User_ID)
	} else {
		_, err = messageservice.DbContext.Exec("insert into user_messages (message_id, user_id) values (?, ?)", message.ID, message.User_ID)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating message object in database %v with value %v, %v", err, message.ID, message.User_ID)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating message object in database %v", err)
			return err
		}
	}

	return nil
}

func (messageservice *MessageService) Update(message *models.DtoMessage) (err error) {
	_, err = messageservice.DbContext.Update(message)
	if err != nil {
		log.Error("Error during updating message object in database %v with value %v", err, message.ID)
		return err
	}

	return nil
}

func (messageservice *MessageService) Delete(message *models.DtoMessage, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = messageservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during deleting message object in database %v", err)
			return err
		}
	}

	if inTrans {
		_, err = trans.Exec("delete from user_messages where message_id = ? and user_id = ?", message.ID, message.User_ID)
	} else {
		_, err = messageservice.DbContext.Exec("delete from user_messages where message_id = ? and user_id = ?", message.ID, message.User_ID)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting message object in database %v with value %v, %v", err, message.ID, message.User_ID)
		return err
	}

	if inTrans {
		_, err = trans.Exec("delete from "+messageservice.Table+" where id = ?", message.ID)
	} else {
		_, err = messageservice.DbContext.Exec("delete from "+messageservice.Table+" where id = ?", message.ID)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting message object in database %v with value %v", err, message.ID)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during deleting message object in database %v", err)
			return err
		}
	}

	return nil
}

func (messageservice *MessageService) DeleteByUser(user_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from user_messages where user_id = ?", user_id)
	} else {
		_, err = messageservice.DbContext.Exec("delete from user_messages where user_id = ?", user_id)
	}
	if err != nil {
		log.Error("Error during deleting message object in database %v with value %v", err, user_id)
		return err
	}

	if trans != nil {
		_, err = trans.Exec("delete from "+messageservice.Table+" where user_id = ?", user_id)
	} else {
		_, err = messageservice.DbContext.Exec("delete from "+messageservice.Table+" where user_id = ?", user_id)
	}
	if err != nil {
		log.Error("Error during deleting message object in database %v with value %v", err, user_id)
		return err
	}

	return nil
}
