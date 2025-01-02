import socket
import threading
from messages import messages_pb2 as messages
import inquirer

class GatewayClient:
    def __init__(self, gateway_ip, tcp_port):
        self.gateway_ip = gateway_ip
        self.tcp_port = tcp_port
        self.running = True

    def send_tcp_message(self, message):
        """Serializa e envia mensagem Protobuf via TCP."""
        tcp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        try:
            print(f"Attempting to connect to {self.gateway_ip}:{self.tcp_port}")
            tcp_socket.connect((self.gateway_ip, self.tcp_port))
            print("Connection established successfully")
            
            data = message.SerializeToString()
            tcp_socket.sendall(data)
            print("Data sent successfully")
            
            response_data = tcp_socket.recv(1024)
            response = messages.ClientResponse()
            response.ParseFromString(response_data)
            print(f"Resposta TCP do gateway: {response.response}")
            
        except ConnectionResetError:
            print("Connection reset by peer - the server closed the connection")
            print("Possible causes:")
            print("1. Protocol mismatch")
            print("2. Invalid message format")
            print("3. Server rejected the connection")
        except Exception as e:
            print(f"Error: {type(e).__name__}: {str(e)}")
        finally:
            tcp_socket.close()

    def run(self):
        """Loop principal para enviar mensagens e escutar respostas."""

        questions = [
            inquirer.List(
                'option',
                message="Choose",
                choices=[
                    ('Get Device State', 'GET_DEVICE_STATE'),
                    ('Set Device State', 'SET_DEVICE_STATE'),
                    ('Exit', 'EXIT'),
                ],
            ),
        ]

        try:
            while self.running:
                answers = inquirer.prompt(questions)
                action = answers['option']

                if action == 'GET_DEVICE_STATE':
                    device_num = input('Type the device number: ')
                    message = messages.ClientMessage()
                    message.request = f"{action}|{device_num}"
                    self.send_tcp_message(message)
                elif action == 'SET_DEVICE_STATE':
                    try:
                        res = input('Type the device number and state value like "num, state": ').split(',')
                        device_num, state = [item.strip() for item in res]
                        message = messages.ClientMessage()
                        message.request = f"{action}|{device_num}|{state}"
                        self.send_tcp_message(message)
                    except ValueError:
                        print("Invalid input. Please try again.")
                elif action == 'EXIT':
                    self.running = False
                    print("Client finished.")
        except KeyboardInterrupt:
            self.running = False
            print("\nEncerrando cliente...")

# Configurações do Gateway
GATEWAY_IP = "172.17.0.3"  # Substitua pelo IP do gateway
TCP_PORT = 9990

if __name__ == "__main__":
    client = GatewayClient(GATEWAY_IP, TCP_PORT)
    client.run()
