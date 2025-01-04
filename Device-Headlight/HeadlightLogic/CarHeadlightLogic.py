import random
    
class CarHeadlightLogic:
    def __init__(self, step=1):
        self.step = step
        self.index = 0
        self.current_state = "on" if random.randint(0, 1) == 1 else "off"  # Default state is random

    def get_data(self):
        self.index += self.step
        # Mantém o estado atual ao retornar
        return f"Headlight|{self.current_state}"

    def set_data(self, data):
        if data in ["on", "off"]:
            self.current_state = data
            print(f"CarHeadlightLogic: Estado atualizado para {self.current_state}")
        else:
            print("CarHeadlightLogic: Valor inválido. O estado deve ser 'on' ou 'off'.")