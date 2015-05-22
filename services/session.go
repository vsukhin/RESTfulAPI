package services

import (
	"application/config"
	"application/models"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"time"
)

type SessionRepository interface {
	GenerateToken(length int) (token string, err error)
	GetAndSaveSession(request *http.Request, r render.Render, params martini.Params,
		updateSession bool, takeFromURI bool, quietMode bool) (session *models.DtoSession, token string, err error)
	Get(token string) (session *models.DtoSession, err error)
	Create(session *models.DtoSession, inTrans bool) (err error)
	Update(session *models.DtoSession, briefly bool, inTrans bool) (err error)
	Delete(token string, inTrans bool) (err error)
	DeleteByUser(userid int64, trans *gorp.Transaction) (err error)
}

type SessionService struct {
	GroupRepository GroupRepository
	*Repository
}

const (
	ACCESS_TOKEN_HEADER_NAME = "X-Access-Token"
	ACCESS_TOKEN_QUERY_NAME  = "access-token"
	ACCESS_TOKEN_COOKIE_NAME = "Access-Token"
	ACCESS_TOKEN_PARAM_NAME  = "token"
	ACCESS_TOKEN_LENGTH      = 255
)

func NewSessionService(repository *Repository) *SessionService {
	repository.DbContext.AddTableWithName(models.DtoSession{}, repository.Table).SetKeys(false, "token")
	return &SessionService{Repository: repository}
}

func (sessionservice *SessionService) GenerateToken(length int) (token string, err error) {
	tokenRaw := make([]byte, length)
	if _, err = rand.Read(tokenRaw); nil != err {
		log.Error("Error during token generation %v", err)
		return "", err
	}

	return base64.URLEncoding.EncodeToString(tokenRaw), nil
}

func (sessionservice *SessionService) GetAndSaveSession(request *http.Request, r render.Render,
	params martini.Params, updateSession bool, takeFromURI bool, quietMode bool) (session *models.DtoSession, token string, err error) {
	session = new(models.DtoSession)

	if takeFromURI {
		token = params[ACCESS_TOKEN_PARAM_NAME]
	}

	if token == "" {
		token = request.URL.Query().Get(ACCESS_TOKEN_QUERY_NAME)
	}

	if token == "" {
		token = request.Header.Get(ACCESS_TOKEN_HEADER_NAME)
	}

	if token == "" {
		cookie, errCookie := request.Cookie(ACCESS_TOKEN_COOKIE_NAME)
		if errCookie == nil {
			if (cookie.Expires.Sub(time.Now()) > 0) && (cookie.Domain == config.Configuration.Server.Host) {
				token = cookie.Value
			}
		}
	}

	if token == "" || len(token) > ACCESS_TOKEN_LENGTH {
		if !quietMode {
			log.Error("Can't find session token %v", token)
		}
		return nil, "", errors.New("Missing token")
	}

	session, err = sessionservice.Get(token)
	if err != nil {
		if !quietMode {
			log.Error("Can't find session object in database %v with token %v", err, token)
		}
		return nil, token, err
	}

	if time.Now().Sub(session.LastActivity) > config.Configuration.Server.SessionTimeout {
		if !quietMode {
			log.Error("Session has been expired %v with value %v", session.LastActivity, token)
		}
		return nil, token, errors.New("Expired session")
	}

	if updateSession {
		session.LastActivity = time.Now()
		err = sessionservice.Update(session, true, false)
		if err != nil {
			if !quietMode {
				log.Error("Can't update session object in database %v with value %v", err, token)
			}
			return nil, token, err
		}
	}

	return session, token, nil
}

func (sessionservice *SessionService) Get(token string) (session *models.DtoSession, err error) {
	session = new(models.DtoSession)
	err = sessionservice.DbContext.SelectOne(session, "select * from "+sessionservice.Table+" where token = ?", token)
	if err != nil {
		log.Error("Error during getting session object from database %v with value %v", err, token)
		return nil, err
	}

	var roles *[]models.UserRole
	roles, err = sessionservice.GroupRepository.GetBySession(token)
	if err != nil {
		log.Error("Error during getting session object from database %v with value %v", err, token)
		return nil, err
	}
	session.Roles = *roles

	return session, nil
}

func (sessionservice *SessionService) Create(session *models.DtoSession, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = sessionservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating session object in database %v", err)
			return err
		}
	}

	if inTrans {
		err = trans.Insert(session)
	} else {
		err = sessionservice.DbContext.Insert(session)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating session object in database %v", err)
		return err
	}

	err = sessionservice.GroupRepository.SetBySession(session.AccessToken, &session.Roles, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating session object in database %v with value %v", err, session.AccessToken)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating session object in database %v", err)
			return err
		}
	}

	return nil
}

func (sessionservice *SessionService) Update(session *models.DtoSession, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = sessionservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating session object in database %v", err)
			return err
		}
	}

	if inTrans {
		_, err = trans.Update(session)
	} else {
		_, err = sessionservice.DbContext.Update(session)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating session object in database %v with value %v", err, session.AccessToken)
		return err
	}

	if !briefly {
		err = sessionservice.GroupRepository.SetBySession(session.AccessToken, &session.Roles, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during updating session object in database %v with value %v", err, session.AccessToken)
			return err
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating session object in database %v", err)
			return err
		}
	}

	return nil
}

func (sessionservice *SessionService) Delete(token string, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = sessionservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during deleting session object in database %v", err)
			return err
		}
	}

	err = sessionservice.GroupRepository.SetBySession(token, &[]models.UserRole{}, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting session object in database %v with value %v", err, token)
		return err
	}

	if inTrans {
		_, err = trans.Exec("delete from "+sessionservice.Table+" where token = ?", token)
	} else {
		_, err = sessionservice.DbContext.Exec("delete from "+sessionservice.Table+" where token = ?", token)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting session object in database %v with value %v", err, token)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during deleting session object in database %v", err)
			return err
		}
	}

	return nil
}

func (sessionservice *SessionService) DeleteByUser(userid int64, trans *gorp.Transaction) (err error) {
	err = sessionservice.GroupRepository.DeleteByUser(userid, trans)
	if err != nil {
		log.Error("Error during deleting session object for user in database %v with value %v", err, userid)
		return err
	}

	if trans != nil {
		_, err = trans.Exec("delete from "+sessionservice.Table+" where user_id = ?", userid)
	} else {
		_, err = sessionservice.DbContext.Exec("delete from "+sessionservice.Table+" where user_id = ?", userid)
	}
	if err != nil {
		log.Error("Error during deleting session object for user in database %v with value %v", err, userid)
		return err
	}

	return nil
}
