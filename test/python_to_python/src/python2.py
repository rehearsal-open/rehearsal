from time import sleep

w = open("out.txt", 'w')

for i in range(10000):
    print(i, flush=True)
    print(i, file=w)
    sleep(0.001)

