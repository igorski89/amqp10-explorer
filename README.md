# amqp10-explorer

`amqp10-explorer` is a simple command-line tool helping you work with [AMQP1.0](https://www.amqp.org/resources/specifications)-compatible message brokers such as [Azure Service Bus](https://azure.microsoft.com/fr-fr/services/service-bus/), [Amazon SQS](https://aws.amazon.com/sqs/), [RabbitMQ  AMQP1.0 plugin](https://github.com/rabbitmq/rabbitmq-amqp1.0), etc.

## Usage

```bash
amqp10-explorer -uri $PROTO://$USER:$PASS@$HOST:$PORT -address $QUEUE -operation $OPERATION [params]
```

For example:

```bash
amqp10-explorer -uri amqp://guest:guest@localhost:5672 -address test -operation send -data '{"foo":"bar"}'
```

## Features

### `send`
Submits a message to the specified queue or topic. Message body supplied via the `data` parameter.

### `receive`
Receive messages in the specified queue or topic. Additional options:
* `limit` (`int`, default: 1) – max number of messages to fetch
* `ack` (`bool`, default: `false`) – acknowledge received messages. _Important_: acknowledging messages removes them from the queue!.

Missing feature you would like to see? Please file-in a feture request in the issues!
