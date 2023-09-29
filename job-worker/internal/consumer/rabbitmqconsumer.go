package consumer

type rabbitMQConsumer[T any] struct {
}

func newRabbitMQConsumer[T any]() *rabbitMQConsumer[T] {
	return &rabbitMQConsumer[T]{}
}

func (r *rabbitMQConsumer[T]) consume(callback consumeCallback[T]) error {
	return nil
}
