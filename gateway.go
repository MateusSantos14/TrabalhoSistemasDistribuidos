package main

import (
	"fmt"
	"log"
	"net"

	"messages"

	"github.com/golang/protobuf/proto"
)

func handleClientMessage(conn net.Conn) {
	defer conn.Close()
	var message messages.ClientMessage
	err := proto.UnmarshalFromReader(conn, &message)
	if err != nil {
		log.Println("Error decoding ClientMessage:", err)
		return
	}

	// Processar mensagem do cliente
	fmt.Println("Received ClientMessage:", message.Request)
	// Enviar resposta de volta
	response := &messages.ClienteResponse{Response: "Message received"}
	data, err := proto.Marshal(response)
	if err != nil {
		log.Println("Error encoding ClientResponse:", err)
		return
	}
	conn.Write(data)
}

func handleDeviceMessage(conn net.Conn) {
	defer conn.Close()
	var message messages.DeviceMessage
	err := proto.UnmarshalFromReader(conn, &message)
	if err != nil {
		log.Println("Error decoding DeviceMessage:", err)
		return
	}

	// Processar mensagem do dispositivo
	fmt.Println("Received DeviceMessage from", message.DeviceId, ":", message.Data)
	// Enviar resposta de volta
	response := &messages.DeviceResponse{DeviceId: message.DeviceId, Response: "Data received"}
	data, err := proto.Marshal(response)
	if err != nil {
		log.Println("Error encoding DeviceResponse:", err)
		return
	}
	conn.Write(data)
}

func listenTCP(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error listening on port:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleClientMessage(conn)
		go handleDeviceMessage(conn)
	}
}

func main() {
	// Iniciar servidor
	go listenTCP(8888)

	// Aguarda indefinidamente
	select {}
}
