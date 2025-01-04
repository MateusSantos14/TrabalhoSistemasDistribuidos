import SimulatedActuator
from HeadlightLogic.CarHeadlightLogic import CarHeadlightLogic

# Configuração do atuador
device_id = 1
multicast_addr = "224.0.0.1"
multicast_port = 9999
port = 9996
headlights=CarHeadlightLogic() ## Headlight

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