package pubsub

const (
	msgTypeSubmit = iota
	msgTypeSubscribe
	msgTypeUnsubscribe
)

type message struct {
	Type  byte
	Value interface{}
}

// PubSub is a publish–subscribe channel.
type PubSub struct {
	in   chan message
	outs map[chan<- interface{}]struct{}
}

// New returns a new publish–subscribe channel.
func New(bufLen int) *PubSub {
	p := &PubSub{
		in:   make(chan message, bufLen),
		outs: make(map[chan<- interface{}]struct{}),
	}
	go p.run()
	return p
}

// Publish sends v to the subscribers.
func (p *PubSub) Publish(v interface{}) {
	p.in <- message{msgTypeSubmit, v}
}

// Subscribe subscribes c on the channel.
func (p *PubSub) Subscribe(c chan<- interface{}) {
	p.in <- message{msgTypeSubscribe, c}
}

// Unsubscribe unsubscribes c from the channel.
func (p *PubSub) Unsubscribe(c chan<- interface{}) {
	p.in <- message{msgTypeUnsubscribe, c}
}

// Close closes the publish–subscribe channel.
func (p *PubSub) Close() {
	close(p.in)
}

func (p *PubSub) run() {
	for {
		msg, ok := <-p.in
		if !ok {
			return
		}
		switch msg.Type {
		case msgTypeSubmit:
			for c := range p.outs {
				c <- msg.Value
			}

		case msgTypeSubscribe:
			p.outs[msg.Value.(chan<- interface{})] = struct{}{}

		case msgTypeUnsubscribe:
			delete(p.outs, msg.Value.(chan<- interface{}))
		}
	}
}
