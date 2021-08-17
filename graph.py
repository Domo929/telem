import requests
from matplotlib import pyplot as plt

req = requests.get("http://localhost:8080/telemetry/2021/1/")
if req.status_code != 200:
    exit(1)

laps = req.json()

ham_speeds = []
ver_speeds = []
for lap in laps[1:]:
    ham_speed = lap["Lines"]["44"]["Speeds"]["ST"]["Value"]
    ver_speed = lap["Lines"]["33"]["Speeds"]["ST"]["Value"]
    if ham_speed != '' and ver_speed != '':
        ham_speeds.append(int(ham_speed))
        ver_speeds.append(int(ver_speed))

lap_numbers = range(1, len(ham_speeds))

fig, ax = plt.subplots()
ax.plot(lap_numbers, ver_speeds, color='blue')
ax.plot(lap_numbers, ham_speeds, color='cyan')
ax.set_title("VER vs HAM")
ax.set_xlabel("Lap Number")
ax.set_ylabel("Speed in Speed trap")
plt.show()
