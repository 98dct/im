package job

import "time"

type RetryOptions func(opts *retryOptions)
type retryOptions struct {
	timeout     time.Duration   // 重试总的超时时间
	retryNums   int             // 重试次数
	isRetry     IsRetryFunc     // 是否重试
	retryJetLag RetryJetLagFunc // 重试间隔
}

func newOptions(opts ...RetryOptions) *retryOptions {
	opt := &retryOptions{
		timeout:     DefaultRetryTimeout,
		retryNums:   DefaultRetryNums,
		isRetry:     RetryAlways,
		retryJetLag: RetryJetLagAlways,
	}

	for _, option := range opts {
		option(opt)
	}

	return opt
}

func WithRetryTimeout(timeout time.Duration) RetryOptions {
	return func(opts *retryOptions) {
		opts.timeout = timeout
	}
}

func WithRetryNums(nums int) RetryOptions {
	return func(opts *retryOptions) {
		opts.retryNums = 1

		if nums > 1 {
			opts.retryNums = nums
		}
	}
}

func WithIsRetryFunc(retryFunc IsRetryFunc) RetryOptions {
	return func(opts *retryOptions) {
		if retryFunc != nil {
			opts.isRetry = retryFunc
		}
	}
}

func WithRetryJetLagFunc(retryJetLagFunc RetryJetLagFunc) RetryOptions {
	return func(opts *retryOptions) {
		if retryJetLagFunc != nil {
			opts.retryJetLag = retryJetLagFunc
		}
	}
}
