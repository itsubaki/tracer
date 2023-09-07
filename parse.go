package tracer

import (
	"fmt"
	"strconv"
	"strings"
)

// https://cloud.google.com/trace/docs/setup#force-trace
type XCloudTraceContext struct {
	TraceID   string
	SpanID    string
	TraceTrue bool
}

func Parse(xCloudTraceContext string) (*XCloudTraceContext, error) {
	// https://cloud.google.com/trace/docs/setup
	// The header specification is:
	// "X-Cloud-Trace-Context: TRACE_ID/SPAN_ID;o=TRACE_TRUE"
	ids := strings.Split(strings.Split(xCloudTraceContext, ";")[0], "/")

	// https://cloud.google.com/trace/docs/setup
	// SPAN_ID is the decimal representation of the (unsigned) span ID.
	i, err := strconv.ParseUint(ids[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse uint=%v: %v", ids[1], err)
	}

	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md#retrieving-the-traceid-and-spanid
	// MUST be a 16-hex-character lowercase string
	spanID := fmt.Sprintf("%016x", i)

	// https://cloud.google.com/trace/docs/setup
	// TRACE_TRUE must be 1 to trace this request. Specify 0 to not trace the request.
	var traceTrue bool
	if len(strings.Split(xCloudTraceContext, ";")) > 1 && strings.Split(xCloudTraceContext, ";")[1] == "o=1" {
		traceTrue = true
	}

	return &XCloudTraceContext{
		TraceID:   ids[0],
		SpanID:    spanID,
		TraceTrue: traceTrue,
	}, nil
}
