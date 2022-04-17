package model

import (
	"fmt"

	"answer.io/pkg/utils"
)

type Key string
type Value string

type Question struct {
	Id      utils.ID `json:"id"`
	Key     Key      `json:"key"`
	Value   Value    `json:"value"`
	Deleted bool     `json:"deleted"`
	History []Event  `json:"history"`
	Version int      `json:"version"`
}

func NewFromEvents(events []Event) *Question {
	q := &Question{}

	for _, event := range events {
		q.On(event, false)
	}
	q.History = events
	return q
}

func New(id utils.ID, key Key, value Value) *Question {
	q := &Question{
		Id:    id,
		Key:   key,
		Value: value,
	}

	q.raise(QuestionAdded{
		ID:    id,
		Key:   key,
		Value: value,
	})

	return q
}

func (q *Question) Update(value Value) error {
	if q.Deleted {
		return fmt.Errorf("question deleted")
	}

	q.raise(QuestionUpdate{
		Key:      q.Key,
		NewValue: value,
	})
	q.Value = value

	return nil
}

func (q *Question) Delete() error {
	if q.Deleted {
		return fmt.Errorf("question deleted")
	}
	q.Deleted = true
	q.raise(QuestionDelete{Key: q.Key})
	return nil
}

func (q *Question) On(ev Event, new bool) {
	switch e := ev.(type) {
	case QuestionAdded:
		q.Id = e.ID
		q.Key = e.Key
		q.Value = e.Value
	case QuestionUpdate:
		q.Value = e.NewValue
		new = false
	case QuestionDelete:
		q.Deleted = true
		new = false
	}
	if !new {
		q.Version++
	}
}

func (q *Question) raise(event Event) {
	q.History = append(q.History, event)
	q.On(event, true)
}

func (q *Question) Events() []Event {
	return q.History
}
