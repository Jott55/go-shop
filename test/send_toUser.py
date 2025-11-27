import requests
import threading

def getUser(token):
    t = f"jwtToken: {token}"
    print(requests.get("http://localhost:8069/user", headers={
        "Authorization": t
    }))


with open("token.txt") as f:
    astring = f.read()
    threads: list[threading.Thread] = [] 
    for i in range(100): 
        threads.append(threading.Thread(target=getUser, args=(astring, )))
    for t in threads:
        t.start()
