import SimulatedActuator
from ACLogic.CarACLogic import CarACLogic

# Configuração do atuador
device_id = 1
multicast_addr = "224.0.0.1"
multicast_port = 9999
port = 9996
ac = CarACLogic() ## AC

# Instanciação do SimulatedActuator
sensor = SimulatedActuator.SimulatedActuator(
    device_id=device_id,
    multicast_addr=multicast_addr,
    multicast_port=multicast_port,
    port=port,
    simulator=ac
)

# Executa o SimulatedActuator
sensor.run()