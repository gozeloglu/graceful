package graceful

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraceful_RegisterCleanupFunctions(t *testing.T) {
	type args struct {
		functions []CleanUpFunc
	}
	testCases := []struct {
		name   string
		args   args
		expLen int
	}{
		{
			name: "Don't pass any function",
			args: args{
				functions: make([]CleanUpFunc, 0),
			},
			expLen: 0,
		},
		{
			name:   "Pass multiple cleanup functions",
			args:   args{functions: []CleanUpFunc{func() {}, func() {}, func() {}}},
			expLen: 3,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := &Graceful{}
			g.RegisterCleanupFunctions(tc.args.functions...)
			assert.Len(t, g.functions, tc.expLen)
		})
	}
}
