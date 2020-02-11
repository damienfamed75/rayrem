package msg

type Generic struct {
	data    interface{}
	msgType string
}

func NewGenericMsg(msgType string, data interface{}) *Generic {
	return &Generic{
		data:    data,
		msgType: msgType,
	}
}

func (g *Generic) Data() interface{} {
	return g.data
}

func (g *Generic) Type() string {
	return g.msgType
}
