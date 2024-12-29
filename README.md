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
**Iniciar os Containers**: Você pode iniciar os containers
executando o script `launch.sh`.

            ./launch.sh

# TODOS TEÓRICOS ESSENCIAIS O QUANTO ANTES
-   **Documentar mensagens**: Defina os tipos de mensagens que serão trocadas entre clientes e broker e entre broker e devices, além dos dados, existem mensagens de controle essenciais como conexão, desconexão etc. Defini-las antes da implementação seria melhor

# TODOS
-   **Conexão do cliente**: Implemente a lógica de conexão do cliente, para enviar as mensagens periodicas
    e estabeleça a conexão TCP para uso de actuator
-   **Envio UDP de Broker para cliente**: Implemente a lógica do envio periodico de mensagens
    recebidas dos devices para os clientes
-   **Envio UDP de Devices**: Implemente a lógica do envio periodico de mensagens
    via UDP no Sensors e Actuator e implemente a lógica no broker para recebimento
-   **Atualização do Protocolo de Descoberta**: Modifique a definição do
    `DiscoverMessage` para incluir o endereço do broker. Será necessário pois
    os atuadores estabelecem conexão TCP com o broker além da UDP.
    Atualize o arquivo `messages.proto` para incluir:

            message DiscoverMessage {
                repeated string broker_ips = 1;
            }
-   **Simular o Atuador**: Implemente o `SimulatedActuator.py` com base
    na estrutura do `SimulatedSensor.py`. Além das funcionalidades do Sensor,
    ele deve estabelecer uma conexão TCP com o broker para receber comandos.
-   **Fazer os próximos atuadores**: Implemente o `SimulatedActuator.py` com base
    na estrutura do `SimulatedSensor.py`. Além das funcionalidades do Sensor,
    ele deve estabelecer uma conexão TCP com o broker para receber comandos.

# TODOS EXTRAS
-   **Multicast periodico**: Seria legal se os multicasts fossem periodicos e se um dispositivo não estivesse ele fechasse o socket