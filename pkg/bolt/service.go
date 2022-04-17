package bolt

import (
	"answer.io/pkg/derrors"
	"answer.io/pkg/model"
	"answer.io/pkg/utils"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	questionBucket        = []byte("questions")
	deletedQuestionBucket = []byte("deleted_questions")
)

type service struct {
	db *bolt.DB
}

func NewService(db *bolt.DB) (*service, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(deletedQuestionBucket); err != nil {
			return err
		}
		_, err := tx.CreateBucketIfNotExists(questionBucket)
		return err
	})
	return &service{db: db}, err
}

func (s *service) New(key model.Key, value model.Value) (_ *model.Question, err error) {
	defer derrors.WrapStack(&err, "bolt.service.New")
	q := model.New(utils.NextID(), key, value)
	var data bytes.Buffer
	if err := gob.NewEncoder(&data).Encode(q); err != nil {
		return nil, err
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		qBucket := tx.Bucket(questionBucket)
		dBucket := tx.Bucket(deletedQuestionBucket)
		if qBucket == nil || dBucket == nil {
			return errors.New("bucket doesn't exist")
		}
		d := dBucket.Get([]byte(key))
		if data := qBucket.Get([]byte(key)); len(data) > 0 && len(d) == 0 {
			return errors.New("key already exist")
		}
		if err := qBucket.Put([]byte(key), data.Bytes()); err != nil {
			return err
		}
		return dBucket.Delete([]byte(key))
	})
	if err != nil {
		return nil, err
	}
	return q, err
}

func (s *service) Get(key model.Key) (model.Question, error) {
	var q model.Question
	err := s.db.View(func(tx *bolt.Tx) error {
		qBucket := tx.Bucket(questionBucket)
		dBucket := tx.Bucket(deletedQuestionBucket)
		if qBucket == nil || deletedQuestionBucket == nil {
			return errors.New("bucket doesn't exist")
		}
		if data := dBucket.Get([]byte(key)); len(data) > 0 {
			return errors.New("question deleted")
		}
		data := qBucket.Get([]byte(key))
		if len(data) == 0 {
			return fmt.Errorf("question not found")
		}
		if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&q); err != nil {
			return fmt.Errorf("service.Get: %w", err)
		}
		return nil
	})
	return q, err
}

func (s *service) Update(key model.Key, value model.Value) error {
	q, err := s.Get(key)
	if err != nil {
		return err
	}
	if err := q.Update(value); err != nil {
		return err
	}
	var data bytes.Buffer
	if err := gob.NewEncoder(&data).Encode(q); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		qBucket := tx.Bucket(questionBucket)
		if qBucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return qBucket.Put([]byte(key), data.Bytes())
	})
}

func (s *service) Delete(key model.Key) error {
	q, err := s.Get(key)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		dBucket := tx.Bucket(deletedQuestionBucket)
		if dBucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return dBucket.Put([]byte(q.Key), q.Id)
	})
}

func (s *service) List() ([]model.Question, error) {
	var l []model.Question
	return l, s.db.View(func(tx *bolt.Tx) error {
		qBucket := tx.Bucket(questionBucket)
		dBucket := tx.Bucket(deletedQuestionBucket)
		if qBucket == nil || dBucket == nil {
			return fmt.Errorf("bucket not found")
		}
		cursor := qBucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var q model.Question
			if data := dBucket.Get(k); len(data) > 0 {
				continue
			}
			if err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(&q); err != nil {
				return err
			}
			l = append(l, q)
		}
		return nil
	})
}

func (s *service) History(key model.Key) (_ []model.Event, err error) {
	defer derrors.WrapStack(&err, "bolt.service.History")
	var list []model.Event
	err = s.db.View(func(tx *bolt.Tx) error {
		qhBucket := tx.Bucket([]byte(key))
		if qhBucket == nil {
			return fmt.Errorf("bucket doesn't exist")
		}
		c := qhBucket.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var e model.Event
			if err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(e); err != nil {
				return err
			}
			list = append(list, e)
		}
		return nil
	})
	return list, err
}
