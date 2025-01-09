import time
import socket
from messages import messages_pb2 as messages
from threading import Thread

class SimulatedActuator:
    def __init__(self, device_id, multicast_addr, multicast_port, port, simulator, periodicity=5):
        self.device_id = str(device_id)
        self.multicast_addr = multicast_addr
        self.multicast_port = multicast_port
        self.simulator = simulator
        self.port = port
        self.periodicity = periodicity
        self.type = "ACTUATOR"  # Type of device
        self.brokers_address = []

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
                address = f"{discover_msg.ip}:{discover_msg.port}"
                if address not in self.brokers_address:
                    self.brokers_address.append(address)
                    self.send_discovery_response()
                    Thread(target=self.handle_gateway_tcp_communication, 
                           args=(discover_msg.ip, discover_msg.port), 
                           daemon=True).start()
                    self.setup_udp_connection(discover_msg.ip, discover_msg.port)
                    
            else:
                print(f"Received unknown message from {addr}, Request: {discover_msg.request}", flush=True)
        except Exception as e:
            print(f"Error processing multicast message from {addr}: {e}", flush=True)

    def send_discovery_response(self):
        # Send a discovery response
        response = messages.DiscoverResponse()
        response.device_id = self.device_id
        response.ip = socket.gethostbyname(socket.gethostname())
        response.port = self.port
        response.type = 1
        data = response.SerializeToString()

        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
        sock.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, 255)
        sock.sendto(data, (self.multicast_addr, self.multicast_port))
        print(f"Sent discovery response to {self.multicast_addr}:{self.multicast_port}", flush=True)
        sock.close()

    def setup_udp_connection(self, ip, port):
        # Initialize UDP connection to the broker
        udp_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        print(f"UDP connection setup with broker at {ip}:{port}", flush=True)
        # Start periodic message sending
        while True:
            try:
                message = messages.DeviceMessage()
                message.device_id = self.device_id
                message.data = self.simulator.get_data()  # Simulate sensor data
                print(message)
                data = message.SerializeToString()

                udp_socket.sendto(data, (ip, port))
                print(f"Sent sensor data to broker at {ip}:{port}", flush=True)
            except Exception as e:
                print(f"Error sending sensor data: {e}", flush=True)

            time.sleep(self.periodicity)
        
    def handle_gateway_tcp_communication(self,ip,port):
        """
        Listen for incoming TCP connections from the gateway.
        This will run on a separate thread for handling incoming connections.
        """
        while True:
            try:
                # Cria um socket TCP para escutar
                server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
                #server_socket.bind((self.broker_ip, self.port))
                server_socket.bind(("0.0.0.0", self.port))
                server_socket.listen(5)

                print(f"Actuator listening on {socket.gethostbyname(socket.gethostname())}:{self.port}", flush=True)

                while True:
                    # Aceita conexões de entrada
                    client_socket, addr = server_socket.accept()
                    print(f"Received connection from {addr}", flush=True)

                    # Processa a conexão em uma thread separada
                    Thread(target=self.handle_gateway_connection, args=(client_socket, addr), daemon=True).start()

            except Exception as e:
                continue
                print(f"Error setting up TCP listener: {e}", flush=True)

    def handle_gateway_connection(self, client_socket, addr):
        """
        Handle communication with a connected gateway.
        Process incoming messages from the gateway.
        """
        try:
            while True:
                # Receive data from the gateway
                data = client_socket.recv(1024)
                if not data:
                    print(f"Connection closed by {addr}", flush=True)
                    break

                # Parse the received message
                device_response = messages.DeviceResponse()
                device_response.ParseFromString(data)
                print(f"Received DeviceResponse from {addr}: Device ID: {device_response.device_id}, Response: {device_response.response}", flush=True)
                # You can process the received data or update it in the simulator here
                self.simulator.set_data(device_response.response)

        except Exception as e:
            print(f"Error handling connection from {addr}: {e}", flush=True)
        finally:
            client_socket.close()


    def run(self):
        # Start the multicast listener in a separate thread
        Thread(target=self.listen_multicast, daemon=True).start()
        
        print("SimulatedActuator is running...", flush=True)
        while True:
            time.sleep(1)
