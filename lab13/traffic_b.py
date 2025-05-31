from scapy.all import sniff, IP, TCP, UDP
from collections import defaultdict
import time
import threading
from scapy.arch import get_if_addr


traffic_stats = defaultdict(lambda: {"sent": 0, "recv": 0})

def process(packet):
    if IP not in packet:
        return
    if TCP in packet:
            l4 = packet[TCP]
    elif UDP in packet:
        l4 = packet[UDP]
    else:
        return
    if packet[IP].src == get_if_addr('en0'):
        port = l4.dport
    else: 
        port = l4.sport
    if packet[IP].src == get_if_addr('en0'):
        traffic_stats[port]['sent'] += len(packet)
    else: 
        traffic_stats[port]['recv'] += len(packet)

def print_report():
    while True:
        time.sleep(5)
        print("\nОтчет по трафику")
        print("Port - Sent - Recv")
        for port, stats in traffic_stats.items():
            print(port,"-", stats["sent"], "B","-", stats["recv"], "B")

threading.Thread(target=print_report, daemon=True).start()
sniff(filter="ip", prn=process, store=0)


