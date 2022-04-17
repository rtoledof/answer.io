package handler

import (
	"answer.io/pkg/model"
)

type QuestionManager interface {
	New(key model.Key, value model.Value) (*model.Question, error)
	Update(key model.Key, value model.Value) error
	Delete(key model.Key) error
	Get(key model.Key) (model.Question, error)
	List() ([]model.Question, error)
	History(model.Key) ([]model.Event, error)
}
