import psutil
import time

net_io = psutil.net_io_counters()
sent_old, recv_old = net_io.bytes_sent, net_io.bytes_recv

try:
    while True:
        time.sleep(1)
        net_io_new = psutil.net_io_counters()
        sent_new, recv_new = net_io_new.bytes_sent, net_io_new.bytes_recv
        print(f"out - {(sent_new - sent_old) / 1024:.2f} KB; in - {(recv_new - recv_old) / 1024:.2f} KB")

        sent_old, recv_old = sent_new, recv_new
except KeyboardInterrupt:
    print()

