package main

import (
	"context"
	"flag"
	"log"
	neturl "net/url"
	"strings"
	"time"

	"pack.ag/amqp"
)

func main() {
	var uri string
	flag.StringVar(&uri, "uri", "", "amqp broke uri (including auth)")

	var operation string
	flag.StringVar(&operation, "operation", "", "supported: send, receive")

	var address string
	flag.StringVar(&address, "address", "", "queue or topic name")

	var data string
	flag.StringVar(&data, "data", "", "message contents (payload)")

	var limit uint
	flag.UintVar(&limit, "limit", 1, "max number of messages to fetch")

	var ack bool
	flag.BoolVar(&ack, "ack", false, "acknowledge and remove messages from the queue")

	flag.Parse()

	var client *amqp.Client
	{
		var addr, username, password string
		{
			u, err := neturl.Parse(uri)
			if err != nil {
				log.Fatalf("error unable to parse uri: %s", err)
			}

			username = u.User.Username()
			password, _ = u.User.Password()
			addr = strings.ReplaceAll(uri, strings.Join([]string{u.User.String(), "@"}, ""), "")
		}

		c, err := amqp.Dial(addr, amqp.ConnSASLPlain(username, password))
		if err != nil {
			log.Fatalf("error initialising client: %s", err)
		}
		client = c
	}
	defer client.Close()

	var session *amqp.Session
	{
		s, err := client.NewSession()
		if err != nil {
			log.Fatalf("error establishing session: %s", err)
		}
		session = s
	}
	defer session.Close(context.Background())

	switch operation {
	case "send":
		linkAddr := amqp.LinkTargetAddress(address)
		var sender *amqp.Sender
		{
			s, err := session.NewSender(linkAddr)
			if err != nil {
				log.Fatalf("error connecting to %s: %s", address, err)
			}
			sender = s
		}
		defer sender.Close(context.Background())

		message := amqp.NewMessage([]byte(data))

		ctx := context.Background()
		err := sender.Send(ctx, message)
		if err != nil {
			log.Fatalf("error sending message: %s", err)
		}

		break
	case "receive":
		linkAddr := amqp.LinkSourceAddress(address)
		var receiver *amqp.Receiver
		{
			r, err := session.NewReceiver(
				linkAddr,
				amqp.LinkBatching(true),
				amqp.LinkCredit(uint32(limit)),
			)
			if err != nil {
				log.Fatalf("error connecting to %s: %s", address, err)
			}
			receiver = r
		}
		defer receiver.Close(context.Background())

		err := receive(receiver, limit, ack)
		if err != nil {
			log.Fatalf("error receiving: %s", err)
		}
		break
	default:
		log.Printf("unknown operation: %s", operation)
	}
}

func receive(r *amqp.Receiver, max uint, ack bool) error {
	var received uint
	for {
		if received >= max {
			return nil
		}
		var message *amqp.Message
		{
			msg, err := r.Receive(context.Background())
			if err != nil {
				log.Printf("error receive: %s", err)
				time.Sleep(1 * time.Second)
				continue
			}
			message = msg
		}

		received = received + 1

		log.Printf("%b %s %s\n", message.DeliveryTag, message.Properties.MessageID, message.GetData())

		if ack {
			err := message.Accept()
			if err != nil {
				log.Printf("error acknowleding message: %s", err)
			}
		} else {
			err := message.Release()
			// err := message.Modify(true, true, nil) <- does not work with RabbitMQ
			if err != nil {
				log.Printf("error releasing message: %s", err)
			}
		}
	}
}
