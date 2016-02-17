package logrus_amqp

import (
	"github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

type AMQPHook struct {
	AMQPUrl   string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Mandatory    bool
	Immediate    bool
	Durable      bool
	Internal     bool
	NoWait       bool
	AutoDeleted  bool
}

func NewAMQPHook(url, exchange, routingKey string) *AMQPHook {
	hook := AMQPHook{}

	hook.AMQPUrl = url
	hook.Exchange = exchange
	hook.ExchangeType = "direct"
	hook.Durable = true
	hook.RoutingKey = routingKey

	return &hook
}

// Fire is called when an event should be sent to the message broker
func (hook *AMQPHook) Fire(entry *logrus.Entry) error {	
	conn, err := amqp.Dial(hook.AMQPUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		hook.Exchange,
		hook.ExchangeType,
		hook.Durable,
		hook.AutoDeleted,
		hook.Internal,
		hook.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := entry.String()
	if err != nil {
		return err
	}

	err = ch.Publish(
		hook.Exchange,
		hook.RoutingKey,
		hook.Mandatory,
		hook.Immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}

	return nil
}

// Levels is available logging levels.
func (hook *AMQPHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
