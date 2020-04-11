package main

import "time"

const (
	DEFAULT_TIME_DURATION = 1000 //s
)

type LeakBucket struct {
	Tokens   int64
	queue    chan int64
	internal time.Duration
}

func CreateRateLimit(tokens int64, duration time.Duration) *LeakBucket {
	l := &LeakBucket{
		Tokens:   tokens,
		queue:    make(chan int64, tokens),
		internal: time.Duration(duration),
	}
	go l.Generate()
	return l
}

func (s *LeakBucket) Generate() {

	inter := time.NewTicker(time.Millisecond * s.internal)
	for _ = range inter.C {
		s.queue <- 1
	}
}

func (s *LeakBucket) Take() bool {
	<-s.queue
	return true
}

func (s *LeakBucket) TakeT() time.Time {
	<-s.queue
	return time.Now()
}
