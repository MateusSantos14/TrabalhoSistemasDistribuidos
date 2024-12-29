---
title: Visão Geral do Projeto
---

# Visão Geral do Projeto

Este projeto consiste em uma série de componentes que trabalham juntos
para simular e controlar diversos dispositivos, como atuadores, sensores
e um gateway, todos se comunicando por meio de mensagens multicast em
containers Docker. Os dispositivos interagem através de Protocol Buffers
(protobufs), e o sistema foi projetado para permitir a descoberta
dinâmica e comunicação entre os dispositivos.

# Estrutura do Projeto

    ├── comandsProtoc
    ├── Device-AC
    │   ├── main.py
    │   └── SimulatedActuator.py
    ├── Device-CarLoc
    │   ├── Dockerfile
    │   ├── launch.sh
    │   ├── main.py
    │   ├── messages
    │   │   ├── messages_pb2.py
    │   │   └── __pycache__
    │   │       └── messages_pb2.cpython-312.pyc
    │   ├── __pycache__
    │   │   └── SimulatedSensor.cpython-312.pyc
    │   ├── requirements.txt
    │   └── SimulatedSensor.py
    ├── Device-Headlight
    │   └── main.py
    ├── Gateway
    │   ├── Dockerfile
    │   ├── gateway
    │   ├── gateway.go
    │   ├── go.mod
    │   ├── go.sum
    │   ├── launch.sh
    │   └── messages
    │       ├── messages.pb.go
    │       └── messages.proto
    ├── launch.sh
    └── messages
        ├── messages_pb2_grpc.py
        └── messages.proto

# Componentes

-   **Device-AC**: Simula um dispositivo atuador,
    `SimulatedActuator.py`.

-   **Device-CarLoc**: Simula o sensor de localização de um veículo
    (`SimulatedSensor.py`), e contém as definições de protobuf para
    comunicação com outros componentes.

-   **Device-Headlight**: Simula um dispositivo de farol.

-   **Gateway**: Atua como o gateway de comunicação para o sistema,
    trocando dados e servindo mensagens multicast.

-   **Messages**: Contém as definições de protobuf para as mensagens
    trocadas entre os dispositivos e o gateway.

# Funcionalidades

-   **Mensagens Multicast**: Tanto o componente `Device-CarLoc` quanto o
    `Gateway` utilizam multicast para comunicação dentro de seus
    containers Docker.

-   **Dispositivos Simulados**: O sistema inclui atuadores simulados
    (`SimulatedActuator.py`) e sensores simulados (`SimulatedSensor.py`)
    para emular dispositivos do mundo real.

-   **Comunicação com Protobuf**: As mensagens são definidas usando
    Protocol Buffers, garantindo uma comunicação eficiente e estruturada
    entre os dispositivos.

# Tarefas Pendentes

-   **Ajustar o `DiscoverMessage`**: O `DiscoverMessage` precisa ser
    atualizado para incluir os IPs do broker, permitindo que os
    atuadores descubram e salvem esses IPs dinamicamente.

-   **Simular o Atuador**: A base para o sensor simulado
    (`SimulatedSensor.py`) foi criada, agora é necessário desenvolver
    uma estrutura similar para o atuador simulado
    (`SimulatedActuator.py`).

# Containers Docker

-   Os dispositivos (`Device-CarLoc`, `Gateway`) estão encapsulados em
    containers Docker e utilizam multicast para comunicação.

-   Para rodar esses containers, o script `launch.sh` em cada diretório
    configura o ambiente necessário.

# Como Rodar

1.  **Instalar Dependências**: As dependências de cada componente estão
    especificadas no `requirements.txt` (para Python) ou `go.mod` (para
    Go). Para componentes Python:

            pip install -r requirements.txt

    Para componentes Go:

            go mod tidy

2.  **Construir os Containers Docker**: Navegue até os diretórios
    respectivos e construa os containers Docker.

            docker build -t device-car-loc .
            docker build -t gateway .

3.  **Iniciar os Containers**: Você pode iniciar os containers
    executando o script `launch.sh`.

            ./launch.sh

# Ajustes Necessários

-   **Atualização do Protocolo de Descoberta**: Modifique a definição do
    `DiscoverMessage` para incluir os IPs do broker. Isso permitirá que
    os atuadores descubram e salvem o IP do broker para comunicação.
    Atualize o arquivo `messages.proto` para incluir:

            message DiscoverMessage {
                repeated string broker_ips = 1;
            }

-   **Simular o Atuador**: Implemente o `SimulatedActuator.py` com base
    na estrutura do `SimulatedSensor.py`. O atuador precisará simular o
    comportamento do mundo real e enviar as mensagens multicast
    apropriadas para comunicar com os outros componentes. Exemplo de
    `SimulatedActuator.py`:

            import socket
            import struct
            from messages import messages_pb2

            class SimulatedActuator:
                def __init__(self, ip, port):
                    self.ip = ip
                    self.port = port
                    self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

                def send_message(self, message):
                    self.sock.sendto(message.encode(), (self.ip, self.port))

                def receive_message(self):
                    data, addr = self.sock.recvfrom(1024)
                    return data

                def discover_broker(self):
                    # Envia uma mensagem multicast de descoberta
                    discover_message = messages_pb2.DiscoverMessage()
                    discover_message.broker_ips.append(self.ip)
                    self.send_message(discover_message.SerializeToString())

            if __name__ == "__main__":
                actuator = SimulatedActuator('192.168.1.100', 5005)
                actuator.discover_broker()

# Contribuindo

Sinta-se à vontade para submeter um pull request caso queira contribuir
com o projeto. Certifique-se de seguir as convenções de codificação e
documentar quaisquer mudanças feitas.

# Licença

Este projeto é licenciado sob a Licença MIT - veja o arquivo `LICENSE`
para mais detalhes.
