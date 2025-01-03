import time
import socket
from messages import messages_pb2 as messages
from threading import Thread, Lock

class SimulatedSensor:
    def __init__(self, device_id, multicast_addr, multicast_port, port, simulator, periodicity=5):
        self.device_id = str(device_id)
        self.multicast_addr = multicast_addr
        self.multicast_port = multicast_port
        self.simulator = simulator
        self.port = port
        self.periodicity = periodicity
        self.type = "SENSOR"  # Type of device
        self.broker_ip = None
        self.broker_port = None
        self.udp_socket = None
        self.current_state = "ON"  # Add state tracking
        self.state_lock = Lock()    # Add lock for thread safety

    def start_tcp_server(self):
        """Start TCP server to handle GET/SET commands"""
        tcp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        tcp_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        tcp_socket.bind(('', self.port))
        tcp_socket.listen(5)
        print(f"TCP Server listening on port {self.port}", flush=True)
        
        while True:
            try:
                client_sock, addr = tcp_socket.accept()
                print(f"Accepted TCP connection from {addr}", flush=True)
                Thread(target=self.handle_tcp_client, args=(client_sock,), daemon=True).start()
            except Exception as e:
                print(f"Error accepting TCP connection: {e}", flush=True)

    def handle_tcp_client(self, client_sock):
        """Handle TCP client connection"""
        try:
            while True:
                data = client_sock.recv(1024)
                if not data:
                    break
                    
                device_msg = messages.DeviceMessage()
                device_msg.ParseFromString(data)
                
                print(f"Received TCP message: {device_msg.data}", flush=True)
                
                # Handle GET/SET commands
                if device_msg.data.startswith("GET_STATE"):
                    response = self.handle_get_state()
                elif device_msg.data.startswith("SET_STATE:"):
                    state_value = device_msg.data.split(":")[1]
                    response = self.handle_set_state(state_value)
                else:
                    response = self.create_device_message(f"Unknown command: {device_msg.data}")
                
                # Send response back
                client_sock.sendall(response.SerializeToString())
        except Exception as e:
            print(f"Error handling TCP client: {e}", flush=True)
        finally:
            client_sock.close()

    def handle_get_state(self):
        """Handle GET_STATE command"""
        with self.state_lock:
            current_value = self.simulator.get_data()
            state_info = f"STATE:{self.current_state}|VALUE:{current_value}"
            return self.create_device_message(state_info)

    def handle_set_state(self, new_state):
        """Handle SET_STATE command"""
        with self.state_lock:
            self.current_state = new_state
            return self.create_device_message(f"State set to: {new_state}")
        
    def create_device_message(self, data):
        """Create a DeviceMessage with the given data"""
        message = messages.DeviceMessage()
        message.device_id = self.device_id
        message.data = data
        return message

    def listen_multicast(self):
        # Listen for multicast messages
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        sock.bind(("", self.multicast_port))

        mreq = socket.inet_aton(self.multicast_addr) + socket.inet_aton("0.0.0.0")
        sock.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)

        while True:
            data, addr = sock.recvfrom(1024)
            Thread(target=self.process_message, args=(data, addr)).start()

    def process_message(self, data, addr):
        try:
            # Parse the message using Protobuf
            discover_msg = messages.DiscoverMessage()
            discover_msg.ParseFromString(data)

            if discover_msg.request == "DISCOVERY_REQUEST":
                print(f"Received DISCOVERY_REQUEST from {addr}, Data: {discover_msg}", flush=True)
                self.send_discovery_response(addr)
                self.setup_udp_connection(discover_msg.ip, discover_msg.port)
            else:
                print(f"Received unknown message from {addr}, Request: {discover_msg.request}", flush=True)
        except Exception as e:
            print(f"Error processing multicast message from {addr}: {e}", flush=True)

    def send_discovery_response(self, addr):
        # Send a discovery response
        response = messages.DiscoverResponse()
        response.device_id = self.device_id
        response.ip = socket.gethostbyname(socket.gethostname())
        response.port = self.port
        response.type = 0
        data = response.SerializeToString()

        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
        sock.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, 255)
        sock.sendto(data, (self.multicast_addr, self.multicast_port))
        print(f"Sent discovery response to {self.multicast_addr}:{self.multicast_port}", flush=True)
        sock.close()

    def setup_udp_connection(self, ip, port):
        # Initialize UDP connection to the broker
        self.broker_ip = ip
        self.broker_port = port
        self.udp_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        print(f"UDP connection setup with broker at {ip}:{port}", flush=True)

        # Start periodic message sending
        Thread(target=self.send_periodic_messages, daemon=True).start()

    def send_periodic_messages(self):
        # Periodically send messages to the broker
        while True:
            if self.broker_ip and self.broker_port:
                try:
                    with self.state_lock:
                        if self.current_state == "ON":
                            message = messages.DeviceMessage()
                            message.device_id = self.device_id
                            message.data = self.simulator.get_data()  # Simulate sensor data
                            print(message)
                            data = message.SerializeToString()

                            self.udp_socket.sendto(data, (self.broker_ip, self.broker_port))
                            print(f"Sent sensor data to broker at {self.broker_ip}:{self.broker_port}", flush=True)
                except Exception as e:
                    print(f"Error sending sensor data: {e}", flush=True)

            time.sleep(self.periodicity)

    def run(self):
        # Start the TCP server for handling GET/SET commands
        Thread(target=self.start_tcp_server, daemon=True).start()

        # Start the multicast listener in a separate thread
        Thread(target=self.listen_multicast, daemon=True).start()
        
        print("SimulatedSensor is running...", flush=True)
        while True:
            time.sleep(1)
