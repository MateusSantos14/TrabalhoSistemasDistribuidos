import SimulatedSensor
from CarLocLogic.CarLogic import CarLogic

car = CarLogic("CarLocLogic/coordinates.csv")

sensor = SimulatedSensor.SimulatedSensor(1,"224.0.0.1",9999,9998,car)

sensor.run()