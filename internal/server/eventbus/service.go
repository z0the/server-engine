package eventbus

type Bus interface {
	Publish(topic string, payload BusMsg)
	PublishWithToPrefix(topic string, msg BusMsg)
	PublishWithFromPrefix(topic string, msg BusMsg)
	Subscribe(topic string, newSub MsgHandler)
	SubscribeWithToPrefix(topic string, newSub MsgHandler)
	SubscribeWithFromPrefix(topic string, newSub MsgHandler)
}
