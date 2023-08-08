package domain

type ChannelTuple struct {
	Events chan []string
	Status chan []string
}

func NewChannelTuple(depth int) *ChannelTuple {
	return &ChannelTuple{
		Events: make(chan []string, depth),
		Status: make(chan []string, depth),
	}
}

func (ct *ChannelTuple) Close() {
	close(ct.Events)
	close(ct.Status)
}
