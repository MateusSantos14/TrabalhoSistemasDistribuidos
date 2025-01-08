import sys
import SimulatedActuator
from HeadlightLogic.CarHeadlightLogic import CarHeadlightLogic

# Configurações padrão
DEFAULT_DEVICE_ID = "HL-1"
DEFAULT_PORT = 9998
multicast_addr = "224.0.0.1"
multicast_port = 9999

print("Argumentos recebidos:", sys.argv)

# Obtém os argumentos da linha de comando
if len(sys.argv) > 1:
    device_id = sys.argv[1]  # Primeiro argumento: ID do dispositivo
else:
    print(f"Nenhum ID fornecido. Usando o ID padrão {DEFAULT_DEVICE_ID}.")
    device_id = DEFAULT_DEVICE_ID

if len(sys.argv) > 2:
    try:
        port = int(sys.argv[2])  # Segundo argumento: Porta
    except ValueError:
        print(f"Porta inválida: {sys.argv[2]}. Usando a porta padrão {DEFAULT_PORT}.")
        port = DEFAULT_PORT
else:
    print(f"Nenhuma porta fornecida. Usando a porta padrão {DEFAULT_PORT}.")
    port = DEFAULT_PORT

print(f"ID do dispositivo: {device_id}")
print(f"Porta: {port}")

# Instanciação da lógica do farol
headlights=CarHeadlightLogic() 

# Instanciação do SimulatedActuator
sensor = SimulatedActuator.SimulatedActuator(
    device_id=device_id,
    multicast_addr=multicast_addr,
    multicast_port=multicast_port,
    port=port,
    simulator=headlights
)

# Executa o SimulatedActuator
sensor.run()

