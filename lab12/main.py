import random

class Router:
    def __init__(self, ip):
        self.ip = ip
        self.neighbors = set()
        self.routing_table = {ip: (ip, 0)}

    def update_table(self, neighbor_ip, neighbor_table):
        updated = False
        for dest, (_, metric) in neighbor_table.items():
            if dest == self.ip:
                continue
            new_metric = metric + 1
            if (dest not in self.routing_table or 
                new_metric < self.routing_table[dest][1]):
                self.routing_table[dest] = (neighbor_ip, new_metric)
                updated = True
        return updated

    def print_table(self, step=None):
        if step is not None:
            print(f"Simulation step {step} of router {self.ip}")
        else:
            print(f"Final state of router {self.ip} table:")
        print(f"{'Source IP':<18}{'Destination IP':<20}{'Next Hop':<18}{'Metric'}")
        for dest, (next_hop, metric) in sorted(self.routing_table.items()):
            print(f"{self.ip:<18}{dest:<20}{next_hop:<18}{metric}")
        print()
    

n = 5
routers = {}
ips = [".".join(str(random.randint(1, 254)) for _ in range(4)) for _ in range(n)]
for ip in ips:
    routers[ip] = Router(ip)
    
for ip in ips:
    for neighbor in random.sample([i for i in ips if i != ip], random.randint(1, n - 1)):
        routers[ip].neighbors.add(neighbor)
        routers[neighbor].neighbors.add(ip)
        routers[ip].routing_table[neighbor] = (neighbor, 1)
        routers[neighbor].routing_table[ip] = (ip, 1)



max_steps = 10
for step in range(1, max_steps + 1):
    updated = False
    for router in routers.values():
        for ip in router.neighbors:
            if router.update_table(ip, routers[ip].routing_table.copy()):
                updated = True
    for router in routers.values():
        router.print_table(step)
    if not updated:
        break

for router in routers.values():
    router.print_table()