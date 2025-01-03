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
    │   │   └── messages_pb2.py
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
    ├── messages
    │   └── messages.proto
    └── README.md
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
    `Gateway` utilizam multicast para descoberta dentro de seus
    containers Docker.

-   **Dispositivos Simulados**: O sistema inclui atuadores simulados
    (`SimulatedActuator.py`) e sensores simulados (`SimulatedSensor.py`)
    para emular dispositivos do mundo real.

-   **Comunicação com Protobuf**: As mensagens são definidas usando
    Protocol Buffers, garantindo uma comunicação eficiente e estruturada
    entre os dispositivos.


# Containers Docker

-   Os dispositivos (`Device-CarLoc`, `Gateway`) estão encapsulados em
    containers Docker e utilizam multicast para comunicação.

-   Para rodar esses containers, o script `launch.sh` em cada diretório
    configura o ambiente necessário.

# Como Rodar
-   **Iniciar os Devices**: Rode todos os containers de device
-   **Iniciar o gateway**: Rode o container do gateway e confira o IP que ele ira printar
-   **Rode o cliente**: Altere no código o IP do destino e rode o cliente(PRVISÓRIO)

#   Mensagens:
ClientMessage.request
"GET_DEVICE_STATE|ID"
"CHANGE_DEVICE_STATE|ID|NEW_STATE"

message ClientResponse {
    string response = 1; // Resposta do Broker para o cliente
}
DiscoverMessage.requests:
"DISCOVER_DEVICES"

# RECADOS

Olá a todos, comecei a realização do trabalho, mas ainda falta muito a ser feito e conto com vocês.


# TODOS
-   **Linkar com cliente GUI**: Temos um cliente em linha de comando funcional, falta o cliente com interface web.
-   **Implementar disconnect do broker**: Quando fechar o broker, mandar multicast para os devices pararem de enviar mensagem e implementar
    essa função no device.
-   **Implementar lógica de HeadLight e AC**: A parte da comunicação está funcionando, porém, não temos a lógica dos dispositivos implementada, e o container de Device-AC também não, confiram como está feito em headlight e fazer a lógica dos dispositivos.
-   **Refatoração de código**: Como podem ver, o código foi feito num ritmo de prova de conceito, 
    por isso, seria interessante que fosse realizada a leitura do código e observada possíveis melhorias: Os devices são a parte mais simples e fácil por serem genéricos.
    Por sua vez, seria interessante o gateway ser mais robusto, observem o que pode ser melhorado nele(funções mais claras e as funcionaldiades que faltam)
-   **REVISAR ESQUEMA DE PORTAS ENTRE CONTAINERS**: A comunicação atual está funcionando, porém, seria interessante conferir se não vai haver algum problema na forma que as    portas foram alocadas.
