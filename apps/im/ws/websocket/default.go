package websocket

import (
	"math"
	"time"
)

const (
	defaultMaxIdleTime = time.Duration(math.MaxInt64)
	defaultAckTimeout  = 90 * time.Second
)
