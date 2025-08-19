package ratelimiter

import (
	"sync/atomic"
	"time"
)

const SIZE int64 = 3000

type TokenOwner struct {
	fillerInterval time.Duration
	tokens         int64
	tick           *time.Ticker
}

func InitTokenOwner() *TokenOwner {
	t := &TokenOwner{
		fillerInterval: 1 * time.Second,
		tokens:         0,
	}

	go t.Filler()

	return t
}

func (t *TokenOwner) Filler() {
	t.tick = time.NewTicker(t.fillerInterval)
	defer t.tick.Stop()

	for {
		select {
		case <-t.tick.C:
			t.push()
		default:
			break
		}
	}
}

func (t *TokenOwner) push() {
	currentValue := atomic.LoadInt64(&t.tokens)
	if currentValue >= SIZE {
		return
	}
	atomic.AddInt64(&t.tokens, 1)
}

func (t *TokenOwner) TryTakeToken() bool {
	tokens := atomic.LoadInt64(&t.tokens)
	if tokens >= SIZE {
		return false
	}
	t.takeToken()
	return true
}

func (t *TokenOwner) takeToken() {
	atomic.AddInt64(&t.tokens, -1)
}
