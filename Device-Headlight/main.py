import SimulatedActuator

class Actuator():
    def __init__(self,name):
        self.name = name
        self.data = 0
    def get_data(self):
        return self.data
    def set_data(self,data):
        self.data = data



sensor = SimulatedActuator.SimulatedActuator(1,"224.0.0.1",9999,9998,Actuator("2"))

sensor.run()