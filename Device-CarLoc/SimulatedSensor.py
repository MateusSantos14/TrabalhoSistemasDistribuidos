import time
import socket
from messages import messages_pb2 as messages
from threading import Thread

class SimulatedSensor:
    def __init__(self, device_id, multicast_addr, multicast_port, port, simulator, periodicity=5):
        self.device_id = str(device_id)
        self.multicast_addr = multicast_addr
        self.multicast_port = multicast_port
        self.simulator = simulator
        self.port = port
        self.periodicity = periodicity
        self.type = "SENSOR"  # Type of device
        self.brokers_address = []
        self.last_received_time = {}

    def listen_multicast(self):
        # Listen for multicast messages
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        sock.bind(("", self.multicast_port))

        mreq = socket.inet_aton(self.multicast_addr) + socket.inet_aton("0.0.0.0")
        sock.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)
    
        while True:
            data, addr = sock.recvfrom(1024)
            print("")
            Thread(target=self.process_message, args=(data, addr)).start()

    def process_message(self, data, addr):
        try:
            # Parse the message using Protobuf
            discover_msg = messages.DiscoverMessage()
            discover_msg.ParseFromString(data)

            if discover_msg.request == "DISCOVERY_REQUEST":
                print(f"Received DISCOVERY_REQUEST from {addr}, Data: {discover_msg}", flush=True)
                address = f"{discover_msg.ip}:{discover_msg.port}"
                self.last_received_time[address] = time.time()
                if address not in self.brokers_address:
                    self.brokers_address.append(address)
                    self.send_discovery_response()
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
        response.type = 0
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
                print("Vou ver o If")
                print(f"Verificando timeout para {ip}:{port}", flush=True)
                print(f"Último tempo recebido: {self.last_received_time.get(f'{ip}:{port}', 'Não encontrado')}", flush=True)
                print(f"Tempo atual: {time.time()}", flush=True)
                print(f"Intervalo: {time.time() - self.last_received_time.get(f'{ip}:{port}', 0)}", flush=True)
                if time.time() - self.last_received_time[f"{ip}:{port}"] > 7:
                    print(f"Gateway {ip}:{port} timeout",flush=True)
                    self.brokers_address.remove(f"{ip}:{port}")
                    udp_socket.close()
                    break
                print("passei do IF",flush=True)
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
        

    def run(self):
        # Start the multicast listener in a separate thread
        Thread(target=self.listen_multicast, daemon=True).start()
        print("SimulatedSensor is running...", flush=True)
        while True:
            time.sleep(1)
