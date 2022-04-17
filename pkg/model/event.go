package model

import (
	"answer.io/pkg/utils"
	"encoding/gob"
)

func init() {
	gob.Register(QuestionAdded{})
	gob.Register(QuestionUpdate{})
	gob.Register(QuestionDelete{})
}

var _ Event = &QuestionAdded{}

type QuestionAdded struct {
	ID    utils.ID `json:"id"`
	Key   Key      `json:"key"`
	Value Value    `json:"value"`
}

func (q QuestionAdded) IsEvent()       {}
func (q QuestionAdded) String() string { return "add" }
func (q QuestionAdded) Data() Data {
	return Data{
		Key:   string(q.Key),
		Value: string(q.Value),
	}
}

type QuestionUpdate struct {
	Key      Key   `json:"id"`
	NewValue Value `json:"value"`
}

func (q QuestionUpdate) IsEvent()       {}
func (q QuestionUpdate) String() string { return "update" }
func (q QuestionUpdate) Data() Data {
	return Data{
		Key:   string(q.Key),
		Value: string(q.NewValue),
	}
}

type QuestionDelete struct {
	Key Key `json:"key"`
}

func (q QuestionDelete) IsEvent()       {}
func (q QuestionDelete) String() string { return "delete" }
func (q QuestionDelete) Data() Data {
	return Data{
		Key: string(q.Key),
	}
}
