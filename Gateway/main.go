package main

import (
	"fmt"
	"net"
	"sync"

	"your_project_path/broker"

	"github.com/golang/protobuf/proto"
)

type Broker struct {
	devices    map[string]*DeviceInfo
	devicesMux sync.Mutex
	clientAddr string
	clientPort int
}

type DeviceInfo struct {
	deviceID string
	ip       string
	port     int
	conn     net.Conn
}

// Cria uma nova instância do Broker
func NewBroker(clientAddr string, clientPort int) *Broker {
	return &Broker{
		devices:    make(map[string]*DeviceInfo),
		clientAddr: clientAddr,
		clientPort: clientPort,
	}
}

// Função que lida com a descoberta de dispositivos
func (b *Broker) handleDiscoveryRequest(conn net.Conn, message *broker.DiscoverMessage) {
	fmt.Println("Recebido Discovery Request")

	b.devicesMux.Lock()
	defer b.devicesMux.Unlock()

	// Respondendo com todos os dispositivos registrados
	response := &broker.DiscoverResponse{}
	for _, device := range b.devices {
		response.Devices = append(response.Devices, &broker.DiscoverResponse_Device{
			DeviceId: device.deviceID,
			Ip:       device.ip,
			Port:     int32(device.port),
		})
	}

	// Serializando a resposta e enviando para o cliente/dispositivo
	data, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("Erro ao serializar DiscoverResponse:", err)
		return
	}
	conn.Write(data)
}

// Função que lida com a mensagem do cliente
func (b *Broker) handleClientMessage(conn net.Conn, message *broker.ClientMessage) {
	fmt.Printf("Recebido ClientMessage: %s\n", message.Request)

	// Respondendo ao cliente
	response := &broker.ClienteResponse{
		Response: "Mensagem recebida: " + message.Request,
	}

	data, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("Erro ao serializar ClienteResponse:", err)
		return
	}
	conn.Write(data)
}

// Função que lida com as mensagens do dispositivo
func (b *Broker) handleDeviceMessage(conn net.Conn, message *broker.DeviceMessage) {
	fmt.Printf("Recebido DeviceMessage: %s -> %s\n", message.DeviceId, message.Data)

	// Respondendo ao dispositivo
	response := &broker.DeviceResponse{
		DeviceId: message.DeviceId,
		Response: "Dados recebidos com sucesso",
	}

	data, err := proto.Marshal(response)
	if err != nil {
		fmt.Println("Erro ao serializar DeviceResponse:", err)
		return
	}
	conn.Write(data)
}

// Função que trata a conexão com o cliente
func (b *Broker) handleClientConnection(conn net.Conn) {
	defer conn.Close()
	var message broker.ClientMessage

	// Recebe a mensagem do cliente
	err := proto.UnmarshalFromReader(conn, &message)
	if err != nil {
		fmt.Println("Erro ao ler ClientMessage:", err)
		return
	}

	// Chama o handler do cliente
	b.handleClientMessage(conn, &message)
}

// Função que lida com a conexão do dispositivo
func (b *Broker) handleDeviceConnection(conn net.Conn) {
	defer conn.Close()
	var message broker.DeviceMessage

	// Recebe a mensagem do dispositivo
	err := proto.UnmarshalFromReader(conn, &message)
	if err != nil {
		fmt.Println("Erro ao ler DeviceMessage:", err)
		return
	}

	// Chama o handler do dispositivo
	b.handleDeviceMessage(conn, &message)
}

// Função principal que roda o Broker
func (b *Broker) Run() {
	// Escuta conexões do cliente
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", b.clientAddr, b.clientPort))
	if err != nil {
		fmt.Println("Erro ao iniciar o Broker:", err)
		return
	}
	defer listener.Close()

	// Escuta conexões de dispositivos
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Erro ao aceitar conexão:", err)
				continue
			}

			// Identifica se a mensagem é do cliente ou do dispositivo
			var messageType string
			// Aqui, adicione lógica para verificar o tipo da mensagem e direcionar para o handler
			if messageType == "CLIENT" {
				go b.handleClientConnection(conn)
			} else {
				go b.handleDeviceConnection(conn)
			}
		}
	}()

	// Escuta multicast para dispositivos
	go b.listenMulticast()

	// Aguarda indefinidamente
	select {}
}

// Função que lida com a escuta de dispositivos via multicast
func (b *Broker) listenMulticast() {
	conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP("224.0.0.1"),
		Port: 9999,
	})
	if err != nil {
		fmt.Println("Erro ao escutar multicast:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Erro ao ler dados multicast:", err)
			continue
		}
		fmt.Printf("Recebido pacote de %s: %s\n", addr, string(buffer[:n]))

		// Descobre dispositivos
		if string(buffer[:n]) == "DISCOVERY_REQUEST" {
			b.handleDiscoveryRequest(conn, &broker.DiscoverMessage{Request: "DISCOVER_DEVICES"})
		}
	}
}

func main() {
	// Criação do Broker
	broker := NewBroker("localhost", 8888)

	// Inicializa o Broker
	go broker.Run()

	// O Broker continua executando indefinidamente
	select {}
}
