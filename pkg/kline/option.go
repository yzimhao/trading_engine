package kline

import "go.uber.org/zap"

type options struct {
	logger            *zap.Logger
	pricePrecision    int32
	quantityPrecision int32
	amountPrecision   int32
}

type Option func(opts *options)

func defaultOptions() *options {
	return &options{
		logger:            zap.NewNop(),
		pricePrecision:    2,
		quantityPrecision: 2,
		amountPrecision:   2,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithPricePrecision(precision int32) Option {
	return func(o *options) {
		o.pricePrecision = precision
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithQuantityPrecision(precision int32) Option {
	return func(o *options) {
		o.quantityPrecision = precision
	}
}

func WithAmountPrecision(precision int32) Option {
	return func(o *options) {
		o.amountPrecision = precision
	}
}
