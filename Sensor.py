class Sensor:
    def __init__(self, name, device):
        self.name = name
        self.device = device
    
    def getData(self):
        return self.device.getData()