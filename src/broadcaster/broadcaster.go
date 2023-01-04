package broadcaster

import (
	"go/src/payment"
)

type broadcaster struct {
	input chan payment.Payment
	reg   chan chan<- interface{}
	unreg chan chan<- interface{}

	outputs map[chan<- interface{}]bool
}

type Broadcaster interface {
	Register(chan<- interface{})
	Unregister(chan<- interface{})
	Close() error
	Submit(payment.Payment) bool
}

func (bc *broadcaster) broadcast(payment payment.Payment) {
	for listener := range bc.outputs {
		listener <- payment
	}
}

func (bc *broadcaster) run() {
	for {
		select {
		case pm := <-bc.input:
			bc.broadcast(pm)
		case ch, ok := <-bc.reg:
			if ok {
				bc.outputs[ch] = true
			} else {
				return
			}
		case ch := <-bc.unreg:
			delete(bc.outputs, ch)
		}
	}
}

func NewBroadcaster(buflen int) Broadcaster {
	bc := &broadcaster{
		input:   make(chan payment.Payment, buflen),
		reg:     make(chan chan<- interface{}),
		unreg:   make(chan chan<- interface{}),
		outputs: make(map[chan<- interface{}]bool),
	}

	go bc.run()

	return bc
}

func (bc *broadcaster) Register(newch chan<- interface{}) {
	bc.reg <- newch
}

func (bc *broadcaster) Unregister(newch chan<- interface{}) {
	bc.unreg <- newch
}

func (bc *broadcaster) Close() error {
	close(bc.reg)
	close(bc.unreg)
	return nil
}

func (bc *broadcaster) Submit(payment payment.Payment) bool {
	if bc == nil {
		return false
	}
	select {
	case bc.input <- payment:
		return true
	default:
		return false
	}
}
