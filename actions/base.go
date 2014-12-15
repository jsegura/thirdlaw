package actions

import (
	"fmt"
	"github.com/karlseguin/thirdlaw/core"
	"gopkg.in/karlseguin/typed.v1"
	"strings"
	"time"
	"log"
)

type Base struct {
	retries int
	name string
	delay time.Duration
	action core.Action
}

func (b *Base) Run() error {
	if b.run() {
		return nil
	}
	for i := 0; i < b.retries; i++ {
		time.Sleep(b.delay)
		if b.run() {
			return nil
		}
	}
	return fmt.Errorf("%s giving up, attempts %d", b.name, b.retries+1)
}

func (b *Base) run() bool {
	err := b.action.Run()
	if err == nil {
		return true
	}
	log.Println(fmt.Sprintf("action %s failure: %v", b.name, err))
	return false
}

func New(name string, t typed.Typed) core.Action {
	switch strings.ToLower(t.String("type")) {
	case "shell":
		return build(name, t, NewShell(t))
	default:
		panic(fmt.Errorf("invalid action type %v", string(t.MustBytes(""))))
	}
}

func build(name string, t typed.Typed, action core.Action) core.Action {
	c := &Base{
		name: name,
		action: action,
		retries: t.IntOr("retries", 0),
		delay: time.Millisecond * time.Duration(t.IntOr("delay", 5000)),
	}
	return c
}
