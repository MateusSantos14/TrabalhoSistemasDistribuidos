import xml.etree.ElementTree as ET
import csv
import math



class CarLogic:
    def __init__(self, csv_file,step = 1):
        self.coordinates = []
        self.step = step
        self.index = 0
        with open(csv_file, mode="r") as file:
            reader = csv.reader(file)
            next(reader)  # Skip the header
            for row in reader:
                x, y = map(float, row)
                self.coordinates.append((x, y))

        for i in range(len(self.coordinates)-1,-1,-1):
            self.coordinates.append(self.coordinates[i])

    def get_data(self):
        self.index+=self.step
        data = self.coordinates[self.index%len(self.coordinates)]
        return f"{data[0]}|{data[1]}"