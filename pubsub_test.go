package pubsub

import (
	"sync"
	"testing"
	"time"
)

const (
	val     = "This is a test value."
	timeOut = time.Second * 5
)

func TestSingle(t *testing.T) {
	var (
		channel  = New(2)
		recvChan = make(chan interface{})
	)
	defer channel.Close()

	channel.Subscribe(recvChan)
	defer channel.Unsubscribe(recvChan)
	channel.Publish(val)
	if testVal := <-recvChan; testVal != val {
		t.Error("incorrect value")
	}
}

func TestMulti(t *testing.T) {
	var (
		channel = New(16)
		wg      sync.WaitGroup
	)
	defer channel.Close()

	for i := 0; i < 1000; i++ {
		recvChan := make(chan interface{}, 1)
		channel.Subscribe(recvChan)
		wg.Add(1)
		s := i
		go func() {
			select {
			case testVal := <-recvChan:
				if testVal != val {
					t.Errorf("sub %4d: incorrect value", s)
				} else {
					t.Logf("sub %4d: ok", s)
				}

			case <-time.After(timeOut):
				t.Errorf("sub %4d: timed out", s)
			}
			channel.Unsubscribe(recvChan)
			wg.Done()
		}()
	}

	channel.Publish(val)

	wg.Wait()
}
