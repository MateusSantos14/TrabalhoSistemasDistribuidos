# Visão Geral do Projeto

Este projeto consiste em uma série de componentes que trabalham juntos
para simular e controlar diversos dispositivos, como atuadores, sensores
e um gateway, todos se comunicando por meio de mensagens multicast em
containers Docker. Os dispositivos interagem através de Protocol Buffers
(protobufs), e o sistema foi projetado para permitir a descoberta
dinâmica e comunicação entre os dispositivos.

# Estrutura do Projeto
C:.
├───.venv
│   └───bin
├───Client
│   └───messages
│       └───__pycache__
├───Device-AC
│   ├───ACLogic
│   └───messages
├───Device-CarLoc
│   ├───CarLocLogic
│   ├───messages
│   └───__pycache__
├───Device-Headlight
│   ├───HeadlightLogic
│   └───messages
├───DeviceClasses
├───Gateway
│   └───messages
└───messages
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
"GET_DEVICE_STATE|ID"
"set_DEVICE_STATE|ID|NEW_STATE"

