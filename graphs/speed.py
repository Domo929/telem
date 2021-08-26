import numpy as np
import requests
from matplotlib import pyplot as plt

req = requests.get("http://localhost:8080/telemetry/2021/1/")
if req.status_code != 200:
    exit(1)

laps = req.json()

# Speeds
ham_speeds = []
ver_speeds = []
for lap in laps[1:]:
    ham_speed = lap["Lines"]["44"]["Speeds"]["ST"]["Value"]
    ver_speed = lap["Lines"]["33"]["Speeds"]["ST"]["Value"]

    ham_speeds.append(int(ham_speed) if ham_speed != '' else None)
    ver_speeds.append(int(ver_speed) if ver_speed != '' else None)

ham_speeds = np.array(ham_speeds).astype(np.double)
h_mask = np.isfinite(ham_speeds)

ver_speeds = np.array(ver_speeds).astype(np.double)
v_mask = np.isfinite(ver_speeds)

speed_lap_numbers = np.arange(1, len(laps))

fig_speeds, ax_speeds = plt.subplots()
ax_speeds.plot(speed_lap_numbers[v_mask], ver_speeds[v_mask], color='blue')
ax_speeds.plot(speed_lap_numbers[h_mask], ham_speeds[h_mask], color='cyan')
ax_speeds.set_title("VER vs HAM")
ax_speeds.set_xlabel("Lap Number")
ax_speeds.set_ylabel("Speed in Speed trap")

plt.show()
