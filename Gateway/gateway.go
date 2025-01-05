package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/username/gateway/messages"
	"google.golang.org/protobuf/proto"
)

// Device represents a discovered device.
type Device struct {
	ID        string
	IP        string
	Port      int
	Type      int    // Type of device (0 or 1)
	LastState string // Último estado recebido do dispositivo
}

// Gateway holds the state of devices and clients.
type Gateway struct {
	devices map[string]Device   // Keyed by device ID
	clients map[string]net.Conn // Keyed by client address
	mutex   sync.RWMutex
}

// NewGateway initializes and returns a new Gateway instance.
func NewGateway() *Gateway {
	return &Gateway{
		devices: make(map[string]Device),
		clients: make(map[string]net.Conn),
	}
}

// handleClient escuta na porta TCP e processa mensagens dos clientes
func (g *Gateway) handleClient(port int) {
	// Inicia o listener TCP
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Erro ao iniciar listener TCP: %v", err)
	}
	defer listener.Close()

	log.Printf("Gateway ouvindo em TCP na porta %d", port)

	for {
		// Aceita conexões de clientes
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}

		// Processa a conexão em uma nova goroutine
		go g.processClient(conn)
	}
}

// processClient lê e processa mensagens de um cliente
func (g *Gateway) processClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024) // Buffer para armazenar mensagens

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Erro ao ler do cliente: %v", err)
			return
		}

		// Deserializa a mensagem recebida usando Protobuf
		var clientMsg messages.ClientMessage
		err = proto.Unmarshal(buf[:n], &clientMsg)
		if err != nil {
			log.Printf("Erro ao deserializar mensagem do cliente: %v", err)
			continue
		}

		response, err := g.processClientMessage(&clientMsg)

		// Serializa a resposta
		serializedResp, err := proto.Marshal(response)
		if err != nil {
			log.Printf("Erro ao serializar resposta para o cliente: %v", err)
			continue
		}

		// Envia a resposta ao cliente
		_, err = conn.Write(serializedResp)
		if err != nil {
			log.Printf("Erro ao enviar resposta ao cliente: %v", err)
			return
		}
	}
}

func (g *Gateway) processClientMessage(clientMsg *messages.ClientMessage) (*messages.ClientResponse, error) {
	g.mutex.RLock() // Use read lock to allow concurrent reads
	defer g.mutex.RUnlock()

	log.Printf("Processing ClientMessage: Request=%s", clientMsg.Request)

	// Split the request to parse command and parameters
	parts := strings.Split(clientMsg.Request, "|")
	if len(parts) != 2 {
		log.Printf("invalid request format, expected 'COMMAND|PARAM'")
	}

	command := parts[0]
	deviceID := parts[1]

	switch command {
	case "GET_DEVICE_STATE":
		// Check if the device exists
		device, exists := g.devices[deviceID]
		if !exists {
			return &messages.ClientResponse{
				Response: fmt.Sprintf("Device ID=%s not found", deviceID),
			}, nil
		}
		// Return the device's last state
		return &messages.ClientResponse{
			Response: fmt.Sprintf("Device ID=%s, LastState=%s", deviceID, device.LastState),
		}, nil
	case "SET_DEVICE_STATE":
		log.Printf("Setting new device state")
		device, exists := g.devices[deviceID]
		if !exists {
			return &messages.ClientResponse{
				Response: fmt.Sprintf("Device ID=%s not found", deviceID),
			}, nil
		}
		// Set new device state
		if len(parts) != 3 {
			log.Printf("invalid request format, expected 'COMMAND|PARAM'")
			return &messages.ClientResponse{
				Response: fmt.Sprintf("Invalid request format COMMAND|PARAM|DATA"),
			}, nil
		}
		if device.Type == 0 {
			log.Printf("invalid request format, Sensor cannot change state")
			return &messages.ClientResponse{
				Response: fmt.Sprintf("Sensors cannot change state"),
			}, nil
		}
		new_data := parts[2]
		device.LastState = new_data
		g.sendMessageToDevice(deviceID, new_data)
		return &messages.ClientResponse{
			Response: fmt.Sprintf("Device ID=%s, LastStateChanged=%s ", deviceID, device.LastState),
		}, nil

	default:
		return &messages.ClientResponse{
			Response: fmt.Sprintf("Unknown command: %s", command),
		}, nil
	}
}

func (g *Gateway) sendMessageToDevice(deviceID, message string) error {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	device, exists := g.devices[deviceID]
	if !exists {
		return fmt.Errorf("device ID=%s not found", deviceID)
	}

	log.Printf("Sending change message to device ID=%s.", deviceID)

	switch device.Type {
	case 0: // Sensor
		log.Printf("Sensor can't receive messages.")
	case 1: // Actuator
		if device.IP == "" || device.Port == 0 {
			return fmt.Errorf("actuator ID=%s has invalid address or port", deviceID)
		}

		// Create a new TCP connection for the message
		localAddress := fmt.Sprintf(":%d", device.Port)
		remoteAddress := fmt.Sprintf("%s:%d", device.IP, device.Port)
		/*
			log.Printf("Sending to ID=%s in %s.", deviceID, address)

			conn, err := net.Dial("tcp", address)
			if err != nil {
				return fmt.Errorf("failed to connect to actuator ID=%s: %v", deviceID, err)
			}
			defer conn.Close()
		*/
		localAddr, err := net.ResolveTCPAddr("tcp", localAddress)
		if err != nil {
			log.Fatalf("Failed to resolve local address: %v", err)
		}

		// Dial the remote server using the specified source port
		dialer := &net.Dialer{LocalAddr: localAddr}
		conn, err := dialer.Dial("tcp", remoteAddress)
		if err != nil {
			log.Fatalf("Failed to connect to %s from %s: %v", remoteAddress, localAddress, err)
		}
		defer conn.Close()
		// Send the message
		deviceResponse := &messages.DeviceResponse{
			DeviceId: deviceID,
			Response: message,
		}

		serializedResp, err := proto.Marshal(deviceResponse)
		_, err = conn.Write(serializedResp)
		if err != nil {
			return fmt.Errorf("failed to send message to actuator ID=%s: %v", deviceID, err)
		}

		log.Printf("Message sent to actuator: ID=%s, Address=%s, Message=%s", deviceID, remoteAddress, message)
	default:
		return fmt.Errorf("unknown device type: ID=%s, Type=%d", deviceID, device.Type)
	}

	return nil
}

// discoverDevices sends a discovery request periodically and listens for responses.
func (g *Gateway) discoverDevices(multicastAddr string) {
    // Listen to every network interface on 0.0.0.0 at port 9990
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", "0.0.0.0", 9990))
	if err != nil {
		log.Printf("Failed to listen on port %s", fmt.Sprintf("%s:%d", "0.0.0.0", 9990))
		return
	}

	// Aqui criamos um listener UDP na porta que foi fornecida pela resposta do discovery
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("Failed to establish UDP listener.")
		return
	}

	go g.handleUDPConnection(udpConn)

	go g.handleClient(9991)

	ticker := time.NewTicker(5 * time.Second) // Set the interval for discovery (30 seconds)
	defer ticker.Stop()

	for {
		// Create a new discovery message
		discoverMsg := &messages.DiscoverMessage{
			Request: "DISCOVERY_REQUEST",
			Ip:      getLocalIP(),
			Port:    9990,
		}

		log.Printf("Discovered message: Request=%s, IP=%s, Port=%d", discoverMsg.Request, discoverMsg.Ip, discoverMsg.Port)

		// Serialize the message
		serializedMsg, err := proto.Marshal(discoverMsg)
		if err != nil {
			log.Fatalf("Failed to marshal DiscoverMessage: %v", err)
		}

		// Setup the multicast address
		conn, err := net.Dial("udp", multicastAddr)
		if err != nil {
			log.Fatalf("Failed to dial multicast address: %v", err)
		}
		defer conn.Close()

		// Send the serialized message over multicast
		_, err = conn.Write(serializedMsg)
		if err != nil {
			log.Fatalf("Failed to send multicast message: %v", err)
		}
		log.Println("Multicast discover sent")

		// Listen for responses
		g.listenForResponses(multicastAddr)

		// Wait for the next tick before sending another discovery message
		<-ticker.C
	}
}

// listenForResponses listens for responses to the multicast discovery request for a limited duration.
func (g *Gateway) listenForResponses(multicastAddr string) {
	// Resolve the UDP address
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	// Create a UDP connection to listen for responses
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to listen on multicast address: %v", err)
	}
	defer conn.Close()

	// Buffer to read incoming data
	buf := make([]byte, 1024)

	// Channel to signal timeout
	timeout := time.After(2 * time.Second)

	for {
		select {
		case <-timeout:
			log.Println("Timeout reached: Stopping response listening")
			return
		default:
			// Set a short read timeout to allow checking the timeout channel
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			// Receive the message
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				if !isTimeoutError(err) {
					log.Printf("Error reading from UDP connection: %v", err)
				}
				continue
			}

			// Unmarshal the received message into a DiscoverResponse
			var discoverResp messages.DiscoverResponse
			err = proto.Unmarshal(buf[:n], &discoverResp)
			if err != nil {
				log.Printf("Failed to unmarshal DiscoverResponse: %v", err)
				continue
			}

			// Process the discovered device
			g.processDevice(&discoverResp)
		}
	}
}

// isTimeoutError checks if the error is a timeout error.
func isTimeoutError(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func (g *Gateway) processDevice(discoverResp *messages.DiscoverResponse) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Inicializa o dispositivo
	device := Device{
		ID:   discoverResp.DeviceId,
		IP:   discoverResp.Ip,
		Port: int(discoverResp.Port),
		Type: int(discoverResp.Type),
	}

	// Salva o dispositivo descoberto
	g.devices[device.ID] = device
	log.Printf("Discovered device: ID=%s, IP=%s, Port=%d, Type=%d", device.ID, device.IP, device.Port, device.Type)
}

func (g *Gateway) handleDeviceConnection(buf []byte, n int, addr *net.UDPAddr) {
    // Unmarshal the Protobuf message
    var deviceMsg messages.DeviceMessage
    err := proto.Unmarshal(buf[:n], &deviceMsg)
    if err != nil {
        log.Printf("Failed to unmarshal UDP message")
        return
    }

    log.Printf("Received UDP message from %s: ID=%s, Data=%s",
        addr.String(), deviceMsg.DeviceId, deviceMsg.Data)

    // Process the DeviceMessage
    g.processDeviceMessage(&deviceMsg)
}

func (g *Gateway) handleUDPConnection(conn *net.UDPConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	localAddr := conn.LocalAddr().String()
	log.Printf("Gateway listening on UDP: %s", localAddr)
	for {
		log.Printf("Waiting message")
		n, addr, err := conn.ReadFromUDP(buf) //TALVEZ PRECISE DE UM MUTEX AQUI(RECEBER VÁRIOS UDP)
		if err != nil {
			log.Printf("UDP read error")
			return
		}

        g.handleDeviceConnection(buf, n, addr)
	}
}

func (g *Gateway) processDeviceMessage(deviceMsg *messages.DeviceMessage) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	log.Printf("Processing DeviceMessage: ID=%s, Data=%s", deviceMsg.DeviceId, deviceMsg.Data)

	device, exists := g.devices[deviceMsg.DeviceId]
	if exists {
		device.LastState = deviceMsg.Data
		g.devices[deviceMsg.DeviceId] = device
		log.Printf("Device ID=%s already discovered. Updated LastState to: %s", deviceMsg.DeviceId, deviceMsg.Data)
	} else {
		log.Printf("Device ID=%s not discovered. Ignoring message or add it if required.", deviceMsg.DeviceId)
	}
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("Failed to determine local IP: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// main function
func main() {
	// Initialize the Gateway
	gateway := NewGateway()

	// Start multicast discovery in a separate goroutine.
	go gateway.discoverDevices("224.0.0.1:9999")

	println("Gateway IP = %s", getLocalIP())

	// This is a simple way to keep the main function alive while the goroutine runs.
	select {}
}
