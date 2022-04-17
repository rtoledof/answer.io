package utils

import (
	"encoding/base32"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIDString(t *testing.T) {
	var testCases = []struct {
		name string
		in   ID
		want string
	}{
		{
			name: "id string test",
			in:   ID("test_id"),
			want: base32.HexEncoding.EncodeToString([]byte("test_id")),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, tt.in.String()); diff != "" {
				t.Errorf("unexpected question mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
