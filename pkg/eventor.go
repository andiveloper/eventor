package pkg

import "context"

type Eventor struct {
	config *EventorConfig
}

func NewEventor(config EventorConfig) *Eventor {
	return &Eventor{&config}
}

func (eventor *Eventor) Run(ctx context.Context) {
	return
}
