# tracer


```go
// Here's an example of using gin-gonic/gin.
func SetTraceID(c *gin.Context) {
	value := c.GetHeader("X-Cloud-Trace-Context")

	xc, err := tracer.Parse(value)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set("trace_id", xc.TraceID)
	c.Set("span_id", xc.SpanID)
	c.Set("trace_true", xc.TraceTrue)

	c.Next()
}
```

```go
func Func(c *gin.Context) {
	traceID := c.GetString("trace_id")
	spanID := c.GetString("span_id")
	traceTrue := c.GetBool("trace_true")

	parent, err := tracer.Context(c.Request.Context(), traceID, spanID, traceTrue)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	func() {
		_, s := tr.Start(parent, "something to do")
		defer s.End()

		...
	}()
...
```
