package job

import (
	"context"
	"errors"
	"time"
)

var ErrJobTimeout = errors.New("任务执行超时！")

type RetryJetLagFunc func(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration

func RetryJetLagAlways(ctx context.Context, retryCount int, lastTime time.Duration) time.Duration {
	return DefaultRetryJetLag
}

type IsRetryFunc func(ctx context.Context, retryCount int, err error) bool

func RetryAlways(ctx context.Context, retryCount int, err error) bool {
	return true
}

func WithRetry(ctx context.Context, handler func(ctx context.Context) error, opts ...RetryOptions) error {

	opt := newOptions(opts...)

	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.timeout)
		defer cancel()
	}

	var (
		herr        error
		retryJetLag time.Duration
		ch          = make(chan error, 1)
	)

	for i := 0; i < opt.retryNums; i++ {
		go func() {
			ch <- handler(ctx)
		}()
		select {
		case herr = <-ch:
			if herr == nil {
				return nil
			}

			if !opt.isRetry(ctx, i, herr) {
				return herr
			}

			retryJetLag = opt.retryJetLag(ctx, i, retryJetLag)
			time.Sleep(retryJetLag)
		case <-ctx.Done():
			return ErrJobTimeout
		}

	}

	return herr
}
