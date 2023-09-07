package tracer_test

import (
	"fmt"

	"github.com/itsubaki/tracer"
)

func ExampleParse() {
	// TRACE_ID/SPAN_ID;o=TRACE_TRUE
	xc, err := tracer.Parse("105445aa7843bc8bf206b12000100000/0000000000000001;o=1")
	if err != nil {
		panic(err)
	}

	fmt.Println(xc.TraceID)
	fmt.Println(xc.SpanID)
	fmt.Println(xc.TraceTrue)

	// Output:
	// 105445aa7843bc8bf206b12000100000
	// 0000000000000001
	// true
}
