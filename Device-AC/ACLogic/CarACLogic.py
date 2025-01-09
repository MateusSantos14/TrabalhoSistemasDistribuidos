import random

class CarACLogic:
    def __init__(self, step=1):
        self.step = step
        self.index = 0
        self.current_state = random.choice([1, 2, 3])  # Estado inicial aleatório

    def get_data(self):
        self.index += self.step
        temp = self.calculate_temperature()  # Calcula a temperatura baseada no estado atual
        print(f"AC|{self.current_state}|{temp:.1f}" ,flush=True)
        return f"AC|{self.current_state}|{temp:.1f}"  # Formata com uma casa decimal

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

    def calculate_temperature(self):
        # Define uma faixa de temperatura para cada estado
        base_temp = 30  # Temperatura base em graus Celsius
        reduction_per_state = 5  # Redução por nível do estado
        temp = base_temp - (self.current_state - 1) * reduction_per_state
        noise = random.uniform(-1.0, 1.0)  # Ruído aleatório entre -1 e 1
        return temp + noise
