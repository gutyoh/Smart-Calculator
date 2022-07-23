# write your code here
class Calculator:
    def __init__(self):
        self.run = True
        self.store = {}

    def is_command(self, user_input):
        return user_input.startswith('/')

    def is_assignment(self, user_input):
        return '=' in user_input

    def assign(self, user_input):
        try:
            var, val = [x.strip() for x in user_input.split('=')]
        except Exception:
            return 'Invalid assignment'
        if not var.isalpha():
            return 'Invalid identifier'
        if not val.isnumeric():
            if val not in self.store:
                return 'Invalid assignment'
            else:
                val = self.store[val]

        self.store[var] = val
        return None

    def get_command(self, user_input):
        if user_input == '/exit':
            self.run = False
            return "Bye!"
        elif user_input == '/help':
            return 'The program calculates the sum of numbers'
        return 'Unknown command'

    def get_sign(self, symbol):
        if '-' in symbol:
            return 1 if len(symbol) % 2 == 0 else -1
        return 1

    def get_total(self, user_input):
        if isinstance(user_input, str):
            return user_input
        try:
            sign = 1
            output = []
            for idx, symbol in enumerate(user_input):
                if idx % 2 == 0:
                    output.append(sign * int(symbol))
                else:
                    sign = self.get_sign(symbol)
            return sum(output)
        except (SyntaxError, ValueError):
            return 'Invalid expression'

    def get_expression(self, expression):
        parsed_exp = []
        for val in expression.split():
            if val.isalpha():
                if val in self.store:
                    val = self.store[val]
                else:
                    return 'Unknown variable'
            parsed_exp.append(val)
        return parsed_exp

    def run_calculator(self):
        while self.run:
            user_input = input()
            if user_input:
                if self.is_command(user_input):
                    output = self.get_command(user_input)
                elif self.is_assignment(user_input):
                    output = self.assign(user_input)
                else:
                    expression = self.get_expression(user_input)
                    output = self.get_total(expression)

                if output is not None:
                    print(output)


calculator = Calculator()
calculator.run_calculator()