from random import randrange

# キーの作成
n   = randrange(1,10,1)
print(n, flush=True)
key = ("%0"+str(n)+"d") % (randrange(0, 10**n - 1, 1))
# print(key) # テスト用

loop = True
while loop:
    ans = input() # キーを受け取る
    match = 0
    for i in range(n): # 結果をチェック
        if key[i] == ans[i]:
            match += 1

    if n == match:
        loop = False
    
    print(match, flush=True)
