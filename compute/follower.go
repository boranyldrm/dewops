package main

import (
	"github.com/streadway/amqp"
	"net"
)

func listenQueues(server *Server) {
	conn := Connect("amqp://guest:guest@172.18.0.2:5672/")
	// Make sure we close the connection whenever the program is about to exit.
	defer conn.Close()

	ch := GetChannel(conn)
	// Make sure we close the channel whenever the program is about to exit.
	defer ch.Close()

	ip := GetLocalIP()

	// work queue
	work := ip + "_work"
	createAndSubscribeQueue(ch, QueueDef{
		Exchange: work,
		Queue:    work,
		Binding:  work,
	}, WORKERQ, "topic")

	// election queue
	elect := ip + "_election"
	createAndSubscribeQueue(ch, QueueDef{
		Exchange: elect,
		Queue:    elect,
		Binding:  elect,
	}, ELECTIONQ, "fanout")

	forever := make(chan bool)
	<-forever
}

func createAndSubscribeQueue(ch *amqp.Channel, def QueueDef, qn QueueName, exchangeType string) {
	DeclareExchange(ch, def.Exchange, exchangeType)
	DeclareQueue(ch, def)

	// Subscribe to the queue.
	msgs := Subscribe(ch, def.Queue)
	ListenQueue(msgs, qn)
}

// local ip is used for generating unique queue names for each service
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
