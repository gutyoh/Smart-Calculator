# write your code here
import operator
from collections import deque


class Calculator:
    OP_RANKING = {
        '+': 1,
        '-': 1,
        '*': 2,
        '/': 2,
        '^': 3
    }
    OPERATOR_DICT = {
        '+': operator.add,
        '-': operator.sub,
        '*': operator.mul,
        '/': operator.truediv,
        '^': operator.pow
    }

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

    def get_total(self, postfix):
        if isinstance(postfix, str):
            return postfix
        try:
            stack = deque()
            for val in postfix:
                if val.isnumeric():
                    stack.append(val)
                else:
                    b, a = stack.pop(), stack.pop()
                    stack.append(self.evaluate_binary(a, b, val))
            return stack.pop()
        except (SyntaxError, ValueError, TypeError):
            return 'Invalid expression'

    def evaluate_binary(self, a, b, op, operator_fnc=OPERATOR_DICT.get):
        return operator_fnc(op)(int(a), int(b))

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

    def get_postfix(self, user_input: str):
        stack = deque()
        prev_sym = None
        postfix = []
        for sym in user_input:
            if sym == ' ':
                continue
            elif sym.isalpha():
                if sym in self.store:
                    prev_sym = sym
                    sym = self.store[sym]
                    postfix.append(sym)
                else:
                    return 'Unknown variable'
            elif sym.isnumeric():
                if prev_sym and prev_sym.isnumeric():
                    postfix.append(postfix.pop() + sym)
                else:
                    postfix.append(sym)
                    prev_sym = sym
            elif sym in ['+', '-', '*', '/', '(', ')']:
                if prev_sym is None:
                    prev_sym = sym
                    stack.append(sym)
                    continue
                if prev_sym == sym:
                    if '+' == sym:
                        continue
                    elif '-' == sym:
                        prev_sym = '+'
                        stack.pop()
                        stack.append('+')
                    elif '*' == sym or '/' == sym:
                        return 'Invalid expression'
                else:
                    prev_sym = sym
                stack, postfix = self.stack_operator(stack, postfix, sym)

        while stack:
            postfix.append(stack.pop())
        return postfix

    def stack_operator(self, stack, postfix, sym):
        if not stack or stack[-1] == '(' or sym == '(':
            stack.append(sym)
            return stack, postfix
        if sym == ')':
            while stack[-1] != '(':
                postfix.append(stack.pop())
            stack.pop()
            return stack, postfix
        if self.higher_precedence(stack[-1], sym):
            stack.append(sym)
        else:
            while stack and not self.higher_precedence(stack[-1], sym):
                postfix.append(stack.pop())
            stack.append(sym)

        return stack, postfix

    @classmethod
    def higher_precedence(cls, stackpop, op):
        if stackpop == '(' or cls.OP_RANKING[op] > cls.OP_RANKING[stackpop]:
            return True
        return False

    def check_parenthesis(self, user_input):
        return user_input.count('(') != user_input.count(')')

    def run_calculator(self):
        while self.run:
            user_input = input()
            if user_input:
                if self.is_command(user_input):
                    output = self.get_command(user_input)
                elif self.is_assignment(user_input):
                    output = self.assign(user_input)
                else:
                    if self.check_parenthesis(user_input):
                        output = 'Invalid expression'
                    else:
                        postfix = self.get_postfix(user_input)
                        output = self.get_total(postfix)

                if output is not None:
                    print(output)


calculator = Calculator()
calculator.run_calculator()
