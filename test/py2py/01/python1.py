from time import sleep

w = open("in.txt", 'w')

for i in range(10000):
    a = int(input())
    print(a * a, flush=True)
    print(a * a, file=w)
    
    # sleep(1)
