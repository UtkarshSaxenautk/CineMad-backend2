package middleware

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

type HeaderKey string

func (h HeaderKey) String() string {
	return string(h)
}

const (
	TraceIDContextKey = ContextKey("Trace-Id")
	TraceIDHeaderKey  = HeaderKey("x-b3-traceid")
)
