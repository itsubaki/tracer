package tracer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/itsubaki/tracer"
)

func TestMustPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if !ok {
				t.Fail()
			}

			if err.Error() != "something went wrong" {
				t.Fail()
			}
		}
	}()

	tracer.Must(func() error { return nil }, fmt.Errorf("something went wrong"))
	t.Fail()
}

func TestContext(t *testing.T) {
	cases := []struct {
		traceID   string
		spanID    string
		traceTrue bool
		hasErr    bool
	}{
		{"105445aa7843bc8bf206b12000100000", "0000000000000001", true, false},
		{"105445aa7843bc8bf206b12000100000", "0000000000000001", false, false},
		{"hoge", "0000000000000001", false, true},
		{"105445aa7843bc8bf206b12000100000", "hoge", false, true},
	}

	for _, c := range cases {
		if _, err := tracer.Context(
			context.Background(),
			c.traceID,
			c.spanID,
			c.traceTrue,
		); (err != nil) != c.hasErr {
			t.Errorf("err: %v", err)
			continue
		}
	}
}
