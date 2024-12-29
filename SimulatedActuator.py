import time
import socket
from threading import Thread
import broker_pb2  # Protocolo gerado pelo protobuf

class SimulatedActuator:
    def __init__(self, device_id, port, multicast_addr, multicast_port, simulator):
        self.device_id = device_id
        self.port = port
        self.multicast_addr = multicast_addr
        self.multicast_port = multicast_port
        self.simulator = simulator
        self.type = "ACTUATOR"  # Atualizar conforme seu tipo de dispositivo

    def send_discovery_response(self, addr):
        response = broker_pb2.DiscoverResponse()
        device = response.device.add()
        device.device_id = self.device_id
        device.ip = "127.0.0.1"  # Endereço de exemplo
        device.port = self.port
        data = response.SerializeToString()

        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        sock.sendto(data, addr)
        sock.close()

    def listen_multicast(self):
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        sock.bind(("", self.multicast_port))

        mreq = socket.inet_aton(self.multicast_addr) + socket.inet_aton("0.0.0.0")
        sock.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)

        while True:
            data, addr = sock.recvfrom(1024)
            if data == b"DISCOVERY_REQUEST":
                self.send_discovery_response(addr)

    def handle_tcp_connection(self, conn):
        data = conn.recv(1024)
        message = broker_pb2.DeviceMessage()
        message.ParseFromString(data)

        print(f"Device {message.device_id} sent data: {message.data}")
        conn.close()

    def listen_tcp(self):
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.bind(("0.0.0.0", self.port))
        sock.listen(5)

        while True:
            conn, _ = sock.accept()
            Thread(target=self.handle_tcp_connection, args=(conn,)).start()

    def send_periodic_data(self):
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        while True:
            try:
                message = broker_pb2.DeviceMessage()
                message.device_id = self.device_id
                message.data = self.simulator.get_data()
                sock.sendto(message.SerializeToString(), ("127.0.0.1", 8888))  # Broker address
                print(f"[{self.device_id}] Periodic data sent.")
                time.sleep(5)
            except Exception as e:
                print(f"[{self.device_id}] Error sending periodic data: {e}")

    def run(self):
        Thread(target=self.listen_multicast, daemon=True).start()
        Thread(target=self.send_periodic_data, daemon=True).start()
        self.listen_tcp()
