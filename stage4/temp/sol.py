while True:
    s = input()
    if s == '/exit':
        break
    elif s == '/help':
        print('The program calculates the sum of numbers')
        continue
    elif s == '':
        continue
    s = s.split()
    res, sign = 0, 1
    for el in s:
        if el.isnumeric() or (el[0] == '-' and el[1:].isnumeric()):
            res, sign = res + sign * int(el), 1
        else:
            for ch in el:
                if ch == '-':
                    sign *= -1
    print(res)
print('Bye!')