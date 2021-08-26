import numpy as np
import requests
from matplotlib import pyplot as plt

req = requests.get("http://localhost:8080/telemetry/2021/3/")
if req.status_code != 200:
    exit(1)

laps = req.json()

ham_pos = []
lec_pos = []
ver_pos = []
nor_pos = []

for lap in laps:
    ham_pos.append(lap["Lines"]["44"]["Position"])
    lec_pos.append(lap["Lines"]["16"]["Position"])
    ver_pos.append(lap["Lines"]["33"]["Position"])
    nor_pos.append(lap["Lines"]["4"]["Position"])

ham_pos = np.array(ham_pos).astype(np.double)
h_mask = np.isfinite(ham_pos)

ver_pos = np.array(ver_pos).astype(np.double)
v_mask = np.isfinite(ver_pos)

lec_pos = np.array(lec_pos).astype(np.double)
l_mask = np.isfinite(lec_pos)

nor_pos = np.array(nor_pos).astype(np.double)
n_mask = np.isfinite(nor_pos)

pos_lap_numbers = np.arange(len(laps))

fig_pos, ax_pos = plt.subplots()
ax_pos.plot(pos_lap_numbers[h_mask], ham_pos, color='cyan', linestyle='solid', label="HAM")
ax_pos.plot(pos_lap_numbers[v_mask], ver_pos, color='blue', linestyle='solid', label="VER")
ax_pos.plot(pos_lap_numbers[l_mask], lec_pos, color='red', linestyle='solid', label="LEC")
ax_pos.plot(pos_lap_numbers[n_mask], nor_pos, color='orange', linestyle='solid', label="NOR")
ax_pos.set_title("Positions")
ax_pos.set_xlabel("Lap Number")
ax_pos.set_ylim([1, 20])
ax_pos.set_ylabel("Position")
ax_pos.legend()

plt.show()
