package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

type message interface{}

type closeSignal struct{}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

func NewLoop() *Loop {
	return &Loop{
		mq:   newMq(),
		stop: make(chan struct{}),
	}
}

var size = image.Pt(800, 800)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	go func() {
		defer close(l.stop)
		var err error
		l.next, err = s.NewTexture(size)
		if err != nil {
			panic("failed to create texture: " + err.Error())
		}
		defer l.next.Release()

		for {
			msg := l.mq.pull()

			switch m := msg.(type) {

			case updateOp:
				l.Receiver.Update(l.next)
				l.prev = l.next

			case Command:
				m.Do(l.next)

			case closeSignal:
				return
			}
		}
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	l.mq.push(closeSignal{})
	<-l.stop
}

type messageQueue struct {
	ch chan message
}

func newMq() messageQueue {
	return messageQueue{ch: make(chan message, 64)}
}

func (mq *messageQueue) push(m message) {
	mq.ch <- m
}

func (mq *messageQueue) pull() message {
	return <-mq.ch
}
