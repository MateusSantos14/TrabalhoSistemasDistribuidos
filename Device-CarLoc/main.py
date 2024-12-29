import SimulatedSensor

class Simulator:
    def get_data(self):
        return 1

a = Simulator()

sensor = SimulatedSensor.SimulatedSensor(1,"224.0.0.1",9999,a)

sensor.run()