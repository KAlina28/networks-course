import socket
import pickle
import threading
from tkinter import Tk, Canvas

def main():    
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.bind(('localhost', 12345))
    server.listen(1)
    
    root = Tk()
    root.title("сервер")
    canvas = Canvas(root, width=600, height=400, bg='white')
    canvas.pack()
    
    last_x, last_y = None, None
    color = 'black'        
    
    def handle(client):
        try:
            while True:
                data = client.recv(4096)
                if not data:
                    break
                points = pickle.loads(data)
                nonlocal last_x, last_y
                for x, y, action in points:
                    if action == 'start':
                        last_x, last_y = x, y
                    elif action == 'move' and last_x is not None and last_y is not None:
                        canvas.create_line(last_x, last_y, x, y, fill=color, width=2)
                        last_x, last_y = x, y
        except Exception as e:
            print(f"ошибка - {e}")
        finally:
            client.close()
    
    def accepts():
        while True:
            client, addr = server.accept()
            threading.Thread(target=handle, args=(client,)).start()
    
    client_thread = threading.Thread(target=accepts)
    client_thread.daemon = True
    client_thread.start()
    
    root.mainloop()

if __name__ == "__main__":
    main()