package dlog

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

const (
	debugMsg = "debugMsg"
	infoMsg  = "infoMsg"
	errorMsg = "errorMsg"
)

func TestSetLogLevel(t *testing.T) {
	oldLevel := getLevel()
	defer func() { state.currentLevel = oldLevel }()

	tests := []struct {
		name      string
		newLevel  string
		wantLevel int
	}{
		{name: "invalid level", newLevel: "xyz", wantLevel: DebugLevel},
		{name: "debug level", newLevel: "debug", wantLevel: DebugLevel},
		{name: "info level", newLevel: "info", wantLevel: InfoLevel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetLevel(test.newLevel)
			gotLevel := getLevel()
			if test.wantLevel != gotLevel {
				t.Errorf("Error: want=%s, got=%s", toString(test.wantLevel), toString(gotLevel))
			}
		})
	}

}

// Do not run in parallel. It overrides logger with mockLogger.
func TestLogLevel(t *testing.T) {
	oldLogger := state.logger
	defer func() { state.logger = oldLogger }()
	l := &fakeLogger{}
	state.logger = l
	// logs below info(like debug) won't print
	SetLevel("info")

	tests := []struct {
		name     string
		logFunc  func(context.Context, any)
		logMsg   string
		expected bool
	}{
		{name: "debug", logFunc: Debug, logMsg: debugMsg, expected: false},
		{name: "info", logFunc: Info, logMsg: infoMsg, expected: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.logFunc(context.Background(), test.logMsg)
			logs := l.logs
			got := strings.Contains(logs, test.logMsg)

			if got != test.expected {
				t.Errorf("Error: want=%v, got %v", test.expected, got)
			}
		})
	}
}

type fakeLogger struct {
	logs string
}

func (l *fakeLogger) log(ctx context.Context, level int, payload any) {
	l.logs += fmt.Sprintf("%s: %+v", toString(level), payload)
}
