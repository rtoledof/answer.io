package model

import (
	"testing"

	"answer.io/pkg/utils"
	"github.com/google/go-cmp/cmp"
)

func TestQuestionEvents(t *testing.T) {
	utils.Generator = func() string {
		return "test_id_generator"
	}
	var testCases = []struct {
		name string
		in   Question
		want []Event
	}{
		{
			name: "list of events",
			in: Question{
				Id:    []byte{},
				Key:   Key("new_key"),
				Value: Value("new value"),
				History: []Event{
					&QuestionAdded{
						ID:    utils.NextID(),
						Key:   Key("new_key"),
						Value: Value("new value"),
					},
					&QuestionUpdate{
						Key:      Key("new_key"),
						NewValue: Value("new value"),
					},
					&QuestionUpdate{
						Key:      Key("new_key"),
						NewValue: Value("other value"),
					},
					&QuestionDelete{
						Key: Key("new_key"),
					},
				},
			},
			want: []Event{
				&QuestionAdded{
					ID:    utils.NextID(),
					Key:   Key("new_key"),
					Value: Value("new value"),
				},
				&QuestionUpdate{
					Key:      Key("new_key"),
					NewValue: Value("new value"),
				},
				&QuestionUpdate{
					Key:      Key("new_key"),
					NewValue: Value("other value"),
				},
				&QuestionDelete{
					Key: Key("new_key"),
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.in.Events()); diff != "" {
				t.Errorf("unexpected question mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewFromEvents(t *testing.T) {
	utils.Generator = func() string {
		return "test_id_generator"
	}
	var testCases = []struct {
		name string
		in   []Event
		want *Question
	}{
		{
			name: "create questions from event",
			in: []Event{
				&QuestionAdded{
					ID:    utils.NextID(),
					Key:   Key("new_key"),
					Value: Value("new value"),
				},
				&QuestionUpdate{
					Key:      Key("new_key"),
					NewValue: Value("new value"),
				},
				&QuestionUpdate{
					Key:      Key("new_key"),
					NewValue: Value("other value"),
				},
				&QuestionDelete{
					Key: Key("new_key"),
				},
			},
			want: &Question{
				Id:      utils.NextID(),
				Key:     Key("new_key"),
				Value:   Value("other value"),
				Version: 4,
				Deleted: true,
				History: []Event{
					&QuestionAdded{
						ID:    utils.NextID(),
						Key:   Key("new_key"),
						Value: Value("new value"),
					},
					&QuestionUpdate{
						Key:      Key("new_key"),
						NewValue: Value("new value"),
					},
					&QuestionUpdate{
						Key:      Key("new_key"),
						NewValue: Value("other value"),
					},
					&QuestionDelete{
						Key: Key("new_key"),
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFromEvents(tt.in)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected question mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
