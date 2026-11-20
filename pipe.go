package cancer

import "container/list"

type Pipeline struct {
	channel  *Channel
	handlers *list.List // handler linked
}

func newPipe(channel *Channel) *Pipeline {
	pipe := new(Pipeline)
	pipe.channel = channel
	pipe.handlers = list.New()
	return pipe
}

func (p *Pipeline) AddLast(handler Handler) {
	p.handlers.PushBack(handler)
}
func (p *Pipeline) HandlerMsg(msg Msg) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Handler message error :%s", err)
		}
	}()
	handlers := p.handlers
	frontHandler := handlers.Front()
	context := HandlerContext{
		pipeline: p,
		channel:  p.channel,
		current:  frontHandler,
	}
	frontHandler.Value.(Handler)(&context, msg)
}
