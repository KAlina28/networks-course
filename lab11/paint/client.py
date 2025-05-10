import socket
import pickle
from tkinter import Tk, Canvas

def main():
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client.connect(('localhost', 12345))
    
    root = Tk()
    root.title("клиент")
    canvas = Canvas(root, width=600, height=400, bg='white')
    canvas.pack()
    
    last_x, last_y = None, None
    color = 'black'
    points = []
    
    def start_drawing(event):
        nonlocal last_x, last_y
        last_x, last_y = event.x, event.y
        points.append((event.x, event.y, 'start'))
    
    def drawing(event):
        nonlocal last_x, last_y
        if last_x and last_y:
            canvas.create_line(last_x, last_y, event.x, event.y, fill=color, width=2)
            points.append((event.x, event.y, 'move'))
            last_x, last_y = event.x, event.y
            try:
                data = pickle.dumps(points)
                client.sendall(data)
            except Exception as e:
                print(f"ошибка - {e}")
    
    def stop_drawing(event):
        nonlocal last_x, last_y, points
        last_x, last_y = None, None
        points.append((event.x, event.y, 'stop'))
        try:
            data = pickle.dumps(points)
            client.sendall(data)
        except Exception as e:
            print(f"ошибка - {e}")
        points.clear()        
    
    canvas.bind("<Button-1>", start_drawing)
    canvas.bind("<B1-Motion>", drawing)
    canvas.bind("<ButtonRelease-1>", stop_drawing)
    
    root.mainloop()

if __name__ == "__main__":
    main()