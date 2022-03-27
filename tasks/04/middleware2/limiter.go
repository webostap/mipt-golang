package middleware

import (
	"net/http"
	"sync"
)

func Limit(l Limiter) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if l.TryAcquire() {
				handler.ServeHTTP(w, req)
				l.Release()
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
			}
		})
	}
}

type Limiter interface {
	TryAcquire() bool
	Release()
}

type MutexLimiter struct {
	cons  int
	maxC  int
	mutex sync.Mutex
}

func NewMutexLimiter(count int) *MutexLimiter {
	return &MutexLimiter{maxC: count}
}

func (l *MutexLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.cons < l.maxC {
		l.cons += 1
		return true
	} else {
		return false
	}
}

func (l *MutexLimiter) Release() {
	l.mutex.Lock()
	l.cons -= 1
	l.mutex.Unlock()
}

type ChanLimiter struct {
	limit chan int
}

func NewChanLimiter(count int) *ChanLimiter {
	return &ChanLimiter{limit: make(chan int, count)}
}

func (l *ChanLimiter) TryAcquire() bool {
	select {
	case l.limit <- 1:
		return true
	default:
		return false
	}
}

func (l *ChanLimiter) Release() {
	<-l.limit
}
