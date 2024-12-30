package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/username/gateway/messages"
	"google.golang.org/protobuf/proto"
)

// Device represents a discovered device.
type Device struct {
	ID      string
	IP      string
	Port    int
	Type    int          // Type of device (0 or 1)
	UDPSock *net.UDPConn // UDP socket for communication
	TCPSock net.Conn     // Optional TCP connection
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

// discoverDevices sends a discovery request periodically and listens for responses.
func (g *Gateway) discoverDevices(multicastAddr string) {
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

	ticker := time.NewTicker(30 * time.Second) // Set the interval for discovery (30 seconds)
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

// listenForResponses listens for responses to the multicast discovery request.
func (g *Gateway) listenForResponses(multicastAddr string) {
	// Set up the UDP listener to listen on the multicast address
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

	for {
		// Receive the message
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP connection: %v", err)
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

	// Se o tipo do dispositivo for 1, tentamos estabelecer uma conexão TCP
	if device.Type == 1 {
		tcpConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", device.IP, device.Port))
		if err == nil {
			device.TCPSock = tcpConn
		} else {
			log.Printf("Failed to establish TCP connection for device %s: %v", discoverResp.DeviceId, err)
		}
	}

	// Salva o dispositivo descoberto
	g.devices[device.ID] = device
	log.Printf("Discovered device: ID=%s, IP=%s, Port=%d, Type=%d", device.ID, device.IP, device.Port, device.Type)

	if device.TCPSock != nil {
		go g.handleTCPConnection(device.TCPSock, device.ID)
	}
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

		// Unmarshal the Protobuf message
		var deviceMsg messages.DeviceMessage
		err = proto.Unmarshal(buf[:n], &deviceMsg)
		if err != nil {
			log.Printf("Failed to unmarshal UDP message")
			continue
		}

		log.Printf("Received UDP message from %s: ID=%s, Data=%s",
			addr.String(), deviceMsg.DeviceId, deviceMsg.Data)

		// Process the DeviceMessage
		g.processDeviceMessage(&deviceMsg)
	}
}

func (g *Gateway) handleTCPConnection(conn net.Conn, deviceID string) {
	defer conn.Close()

	buf := make([]byte, 1024)
	localAddr := conn.LocalAddr().String()
	log.Printf("Gateway listening on TCP: %s", localAddr)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("TCP read error for device %s: %v", deviceID, err)
			return
		}

		// Unmarshal the Protobuf message
		var deviceMsg messages.DeviceMessage
		err = proto.Unmarshal(buf[:n], &deviceMsg)
		if err != nil {
			log.Printf("Failed to unmarshal TCP message for device %s: %v", deviceID, err)
			continue
		}

		log.Printf("Received TCP message from device %s: ID=%s, Data=%s",
			deviceID, deviceMsg.DeviceId, deviceMsg.Data)

		// Process the DeviceMessage
		g.processDeviceMessage(&deviceMsg)
	}
}

func (g *Gateway) processDeviceMessage(deviceMsg *messages.DeviceMessage) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	log.Printf("Processing DeviceMessage: ID=%s, Data=%s", deviceMsg.DeviceId, deviceMsg.Data)

	// Lógica de processamento segura
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

	// This is a simple way to keep the main function alive while the goroutine runs.
	select {}
}
