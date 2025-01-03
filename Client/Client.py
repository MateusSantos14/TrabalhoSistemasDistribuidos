import socket
from messages import messages_pb2 as messages
import threading


class GatewayClient:
    def __init__(self, gateway_ip, tcp_port):
        self.gateway_ip = gateway_ip
        self.tcp_port = tcp_port
        self.running = True

    def send_tcp_message(self, message):
        """Serializa e envia mensagem Protobuf via TCP."""
        tcp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        tcp_socket.connect((self.gateway_ip, self.tcp_port))
        data = message.SerializeToString()
        tcp_socket.sendall(data)
        print(f"Mensagem TCP enviada para {self.gateway_ip}:{self.tcp_port}")
        response_data = tcp_socket.recv(1024)
        response = messages.ClientResponse()
        response.ParseFromString(response_data)
        print(f"Resposta TCP do gateway: {response.response}")
        tcp_socket.close()

    def run(self):
        """Loop principal para enviar mensagens e escutar respostas."""

        # Loop para enviar mensagens
        try:
            while self.running:
                print("\nEscolha:")
                print("1 - Enviar mensagem TCP")
                print("2 - Sair")
                choice = input("Opção: ").strip()
                if choice == '1':
                    message = messages.ClientMessage()
                    message.request = input("Digite a mensagem TCP: ")
                    self.send_tcp_message(message)
                elif choice == '2':
                    self.running = False
                    print("Cliente encerrado.")
                else:
                    print("Opção inválida!")
        except KeyboardInterrupt:
            self.running = False
            print("\nEncerrando cliente...")

# Configurações do Gateway
GATEWAY_IP = "172.17.0.3"  # Substitua pelo IP do gateway
TCP_PORT = 9991

if __name__ == "__main__":
    client = GatewayClient(GATEWAY_IP, TCP_PORT)
    client.run()
