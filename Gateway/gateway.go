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
	Type    int      // Type of device (0 or 1)
	TCPConn net.Conn // Optional TCP connection
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
	ticker := time.NewTicker(30 * time.Second) // Set the interval for discovery (30 seconds)
	defer ticker.Stop()

	for {
		// Create a new discovery message
		discoverMsg := &messages.DiscoverMessage{
			Request: "DISCOVERY_REQUEST",
		}

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
	// Lock for safe concurrent access
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Print the received DiscoverResponse for debugging purposes
	log.Printf("Received DiscoverResponse: DeviceId=%s, Ip=%s, Port=%d, Type=%d",
		discoverResp.DeviceId, discoverResp.Ip, discoverResp.Port, discoverResp.Type)

	// Create a Device object from the response data
	device := Device{
		ID:   discoverResp.DeviceId,
		IP:   discoverResp.Ip,
		Port: int(discoverResp.Port),
		Type: int(discoverResp.Type), // Assuming discoverResp has Type field
	}

	// If device type is 1, establish a TCP connection
	if device.Type == 1 {
		tcpConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", device.IP, device.Port))
		if err == nil {
			device.TCPConn = tcpConn
		}
	}

	// Save the discovered device by ID
	g.devices[device.ID] = device
	log.Printf("Discovered device: ID=%s, IP=%s, Port=%d, Type=%d", device.ID, device.IP, device.Port, device.Type)
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
