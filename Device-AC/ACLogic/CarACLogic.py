import random

class CarACLogic:
    def __init__(self, step=1):
        self.step = step
        self.index = 0
        self.current_state = random.coiche([1, 2, 3])  # Default state is random

    def get_data(self):
        self.index += self.step
        # Mantém o estado atual ao retornar
        return f"AC|{self.current_state}"

    def set_data(self, data):
        try:
            # Atualiza o estado do ar condicionado se for 1, 2 ou 3
            data = int(data)
            if data in [1, 2, 3]:
                self.current_state = data
                print(f"CarACLogic: Estado atualizado para {self.current_state}")
            else:
                print("CarACLogic: Valor inválido. O estado deve ser 1, 2 ou 3.")
        except ValueError:
            print("CarACLogic: Valor inválido recebido. Deve ser um número inteiro.")