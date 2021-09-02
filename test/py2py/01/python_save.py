from time import sleep

w = open("out.txt", 'w')
while True:
    print(input(), file=w)