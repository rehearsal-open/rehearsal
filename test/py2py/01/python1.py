from time import sleep

w = open("in.txt", 'w')

for i in range(10000):
    a = int(input())
    print(a * 2, flush=True)
    print(a * 2, file=w)
    
    # sleep(1)

