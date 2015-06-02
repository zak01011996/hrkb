package notify

import (
	"errors"
	"time"

	M "hrkb/models"
)

const (
	ErrLimit = "Limit zero or negative number"
)

type MailStore struct {
	db    *M.DM
	retry int
}

func (s *MailStore) GetFailed(mails interface{}, limit int) error {

	if limit <= 0 {
		return errors.New(ErrLimit)
	}

	err := s.db.FindAll(&M.Mail{}, mails, M.Sf{}, M.Where{And: M.W{"Status": false, "Try<": s.retry}}, M.NewParams(M.Params{Limit: limit}))

	if err != nil {
		return err
	}

	return nil
}

func (s *MailStore) Update(m *M.Mail) error {
	m.Updated = time.Now()

	if _, err := s.db.Update(m); err != nil {
		return err
	}
	return nil
}

func (s *MailStore) Insert(m *M.Mail) error {
	m.Created = time.Now()
	m.Active = true
	m.Status = true

	if _, err := s.db.Insert(m); err != nil {
		return err
	}

	return nil
}
