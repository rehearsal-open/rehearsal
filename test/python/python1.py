from time import sleep

w = open("in.txt", 'w')

for i in range(10000):
    a = input()
    print(a, flush=True)
    print(a, file=w)
    
    # sleep(1)

