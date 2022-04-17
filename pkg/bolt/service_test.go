package bolt

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"answer.io/pkg/model"
	"answer.io/pkg/utils"

	"github.com/google/go-cmp/cmp"
	bbolt "go.etcd.io/bbolt"
)

func checkAsserts(t testing.TB, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot  '%v' \nwant '%v'", got, want)
	}
}

func checkError(t testing.TB, got, want error) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("\ngot  '%s' \nwant '%s'", got, want)
	}
}

// MustOpenDB is a function to create a open db
func mustOpenDB(t testing.TB) (*bbolt.DB, func(t testing.TB)) {
	t.Helper()
	name := make([]byte, 8)
	rand.Read(name)
	strname := base64.RawURLEncoding.EncodeToString(name)
	f, err := ioutil.TempFile("", strname)
	checkError(t, err, nil)
	db, err := utils.Open(f.Name())
	checkError(t, err, nil)
	fn := func(t testing.TB) {
		t.Helper()
		db.Close()
		os.Remove(f.Name())
	}
	return db, fn
}

func TestServiceNew(t *testing.T) {
	db, clean := mustOpenDB(t)
	defer clean(t)
	utils.Generator = func() string {
		return "test_id_generator"
	}
	var testCases = []struct {
		name string
		in   struct {
			Key   model.Key
			value model.Value
		}
		want    *model.Question
		wantErr bool
	}{
		{
			name: "success create new question",
			in: struct {
				Key   model.Key
				value model.Value
			}{
				"new_key",
				"new_value",
			},
			want: &model.Question{
				Id:    utils.NextID(),
				Key:   model.Key("new_key"),
				Value: model.Value("new_value"),
				History: []model.Event{
					model.QuestionAdded{
						ID:    utils.NextID(),
						Key:   model.Key("new_key"),
						Value: model.Value("new_value"),
					},
				},
			},
		},
		{
			name: "duplicate question test",
			in: struct {
				Key   model.Key
				value model.Value
			}{
				"new_key",
				"any value",
			},
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewService(db)
			if err != nil {
				t.Fatalf("got = %v, want nil", err)
			}
			got, err := s.New(tt.in.Key, tt.in.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got = %v, want nil", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected question mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestServiceGet(t *testing.T) {
	db, clean := mustOpenDB(t)
	defer clean(t)
	utils.Generator = func() string {
		return "test_id_generator"
	}
	var testCases = []struct {
		name string
		in   struct {
			key model.Key
		}
		want    model.Question
		wantErr bool
	}{
		{
			name: "success retrieve",
			in: struct {
				key model.Key
			}{
				key: "retrieve_key",
			},
			want: model.Question{
				Id:    utils.NextID(),
				Key:   "retrieve_key",
				Value: "any value",
				History: []model.Event{
					model.QuestionAdded{
						ID: utils.NextID(),
						Key: "retrieve_key",
						Value: "any value",
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewService(db)
			if err != nil {
				t.Fatalf("got = %v, want nil", err)
			} else {
				if tt.want.Value != "" {
					s.New(tt.in.key, tt.want.Value)
				}
				got, err := s.Get(tt.in.key)
				if (err != nil) != tt.wantErr {
					t.Fatalf("got = %v, want nil", err)
				}
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("unexpected question mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	db, clean := mustOpenDB(t)
	defer clean(t)
	utils.Generator = func() string {
		return "test_id_generator"
	}
	var testCases = []struct {
		name string
		in   struct {
			newIn struct {
				key   model.Key
				value model.Value
			}
			updateIn struct {
				key   model.Key
				value model.Value
			}
		}
		want    model.Question
		wantErr bool
	}{
		{
			name: "success update question",
			in: struct {
				newIn struct {
					key   model.Key
					value model.Value
				}
				updateIn struct {
					key   model.Key
					value model.Value
				}
			}{
				newIn: struct {
					key   model.Key
					value model.Value
				}{
					key:   "new_key",
					value: "new value",
				},
				updateIn: struct {
					key   model.Key
					value model.Value
				}{
					key:   "new_key",
					value: model.Value("updated value"),
				},
			},
			want: model.Question{
				Id:    utils.NextID(),
				Key:   "new_key",
				Value: "updated value",
				History: []model.Event{
					model.QuestionAdded{
						ID:    utils.NextID(),
						Key:   model.Key("new_key"),
						Value: model.Value("new value"),
					},
					model.QuestionUpdate{
						Key:      model.Key("new_key"),
						NewValue: model.Value("updated value"),
					},
				},
				Version: 1,
			},
		},
		{
			name: "success update not existing question",
			in: struct {
				newIn struct {
					key   model.Key
					value model.Value
				}
				updateIn struct {
					key   model.Key
					value model.Value
				}
			}{
				newIn: struct {
					key   model.Key
					value model.Value
				}{
					key:   "another_key",
					value: "new value",
				},
				updateIn: struct {
					key   model.Key
					value model.Value
				}{
					key:   model.Key("not_found_key"),
					value: model.Value("updated value"),
				},
			},
			wantErr: true,
		},
		{
			name: "success update deleted question",
			in: struct {
				newIn struct {
					key   model.Key
					value model.Value
				}
				updateIn struct {
					key   model.Key
					value model.Value
				}
			}{
				newIn: struct {
					key   model.Key
					value model.Value
				}{
					key:   "name",
					value: "John",
				},
				updateIn: struct {
					key   model.Key
					value model.Value
				}{
					key:   model.Key("name"),
					value: model.Value("John Doe"),
				},
			},
			want: model.Question{
				Id: utils.NextID(),
				Key: model.Key("name"),
				Value: model.Value("John Doe"),
				History: []model.Event{
					model.QuestionAdded{
						ID: utils.NextID(),
						Key: model.Key("name"),
						Value: model.Value("John"),
					},
					model.QuestionUpdate{
						Key: model.Key("name"),
						NewValue: model.Value("John Doe"),
					},
				},
				Version: 1,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewService(db)
			if err != nil {
				t.Fatalf("got = %v, want nil", err)
			}
			if  _, err := s.New(tt.in.newIn.key, tt.in.newIn.value); err != nil {
				t.Fatalf("got = %v, want nil", err)
			}
			err = s.Update(tt.in.updateIn.key, tt.in.updateIn.value)
			if (err != nil) != tt.wantErr {
 				t.Fatalf("got = %v, want nil", err)
			}
			if err == nil {
				got, err := s.Get(tt.in.newIn.key)
				if (err != nil) != tt.wantErr {
					t.Fatalf("got = %v, want nil", err)
				}
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("unexpected question mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestServiceDelete(t *testing.T) {
	db, clean := mustOpenDB(t)
	defer clean(t)
	utils.Generator = func() string {
		return "test_id_generator"
	}
	var testCases = []struct {
		name    string
		in      struct {
			key model.Key
			value model.Value
		}
		wantErr bool
	}{
		{
			name: "success question delete",
			in:   struct{key model.Key; value model.Value}{
				model.Key("new_key"),
				model.Value("new value"),
			},
		},
		{
			name:    "success question delete",
			in:      struct{key model.Key; value model.Value}{
				key: model.Key("not found key"),
			},
			wantErr: true,
		},
		{
			name:    "success question delete",
			in:      struct{key model.Key; value model.Value}{
				key: model.Key("new_key"),
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewService(db)
			if err != nil {
				t.Fatalf("got = %v, want nil", err)
			}
			if tt.in.value != "" {
				if _, err := s.New(tt.in.key, tt.in.value); err != nil {
					t.Fatalf("got = %v, want nil", err)
				}
			}
			if err := s.Delete(tt.in.key); (err != nil) != tt.wantErr {
				t.Fatalf("got = %v, want nil", err)
			}
		})
	}
}

func TestServiceList(t *testing.T) {
	db, clean := mustOpenDB(t)
	defer clean(t)
	utils.Generator = func() string {
		return "test_id_generator"
	}
	var testCases = []struct {
		name    string
		in []struct{
			key model.Key
			value model.Value
		}
		want    []model.Question
		wantErr bool
	}{
		{
			name: "list questions",
			in: []struct{key model.Key; value model.Value}{
				{
					key: model.Key("name"),
					value: model.Value("John"),
				},
				{
					key: model.Key("last_name"),
					value: model.Value("Doe"),
				},
			},
			want: []model.Question{
				{
					Id:    utils.NextID(),
					Key:   model.Key("last_name"),
					Value: model.Value("Doe"),
					History: []model.Event{
						model.QuestionAdded{
							ID: utils.NextID(),
							Key: model.Key("last_name"),
							Value: model.Value("Doe"),
						},
					},
				},
				{
					Id:    utils.NextID(),
					Key:   model.Key("name"),
					Value: model.Value("John"),
					History: []model.Event{
						model.QuestionAdded{
							ID: utils.NextID(),
							Key: model.Key("name"),
							Value: model.Value("John"),
						},
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewService(db)
			if err != nil {
				t.Fatalf("got = %v, want nil", err)
			}
			for _, in := range tt.in {
				if _, err := s.New(in.key, in.value); err != nil {
					t.Fatalf("got = %v, want nil", err)
				}
			}
			got, err := s.List()
			if err != nil {
				t.Fatalf("got = %v, want nil", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected question mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
