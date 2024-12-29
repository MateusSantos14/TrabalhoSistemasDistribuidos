import time
import socket
from messages import messages_pb2 as messages
from threading import Thread

class SimulatedSensor:
    def __init__(self, device_id, multicast_addr, multicast_port, simulator, periodicity=5):
        self.device_id = device_id
        self.multicast_addr = multicast_addr
        self.multicast_port = multicast_port
        self.simulator = simulator
        self.periodicity = periodicity
        self.type = "SENSOR"  # Atualize conforme o tipo de dispositivo

    def listen_multicast(self):
        # Configura o sensor para escutar a rede multicast
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
            # Tenta fazer o unmarshal da mensagem usando o Protobuf
            discover_msg = messages.DiscoverMessage()  # Substitua pelo nome correto da sua mensagem
            discover_msg.ParseFromString(data)

            # Verifica o conteúdo do campo 'request'
            if discover_msg.request == "DISCOVERY_REQUEST":
                print(f"Received DISCOVERY_REQUEST from {addr}, Data: {discover_msg.request}")
                self.send_discovery_response(addr)
            else:
                print(f"Received unknown message from {addr}, Request: {discover_msg.request}")
        except Exception as e:
            print(f"Error processing multicast message from {addr}: {e}")
    
    def send_discovery_response(self, addr):
        # Responde a solicitação de descoberta dos dispositivos via multicast
        response = messages.DiscoverResponse()
        response.device_id = str(self.device_id)
        response.ip = socket.gethostbyname(socket.gethostname())  # Atualize com o IP correto do dispositivo
        response.port = 9998
        response.type = 0
        data = response.SerializeToString()

        # Create the socket for multicast
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)

        # Set socket options for multicast
        sock.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, 255)

        # Send the response as a multicast message
        sock.sendto(data, (self.multicast_addr,self.multicast_port))
        print(f"Sent discovery response to { (self.multicast_addr,self.multicast_port)}")

        sock.close()
    def run(self):
        # Inicia threads para escutar e enviar dados
        Thread(target=self.listen_multicast, daemon=True).start()
        while True:
            time.sleep(1)
