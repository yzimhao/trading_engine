package response

type Code int64

const (
	CodeSuccess     Code = 0
	CodeFailUnknown Code = 1 + iota
)

type Option func(opts *options)

type options struct {
	Code    int64  `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func WithError(code int64, message string) Option {
	return func(opts *options) {
		opts.Code = code
		opts.Message = message
	}
}

func WithData(data any) Option {
	return func(opts *options) {
		opts.Data = data
	}
}

func defaultSuccessOptions() *options {
	return &options{
		Code: int64(CodeSuccess),
		Data: nil,
	}
}

func defaultFailOptions() *options {
	return &options{
		Code:    int64(CodeFailUnknown),
		Message: "unknown error",
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func Success() *options {
	return defaultSuccessOptions()
}

func (o *options) WithData(data any) any {
	o.Data = data
	return o
}

func Fail() *options {
	return defaultFailOptions()
}

func (o *options) WithError(code int64, message string) any {
	o.Code = code
	o.Message = message
	return o
}
