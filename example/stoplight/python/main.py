import time
import random
from threading import Thread

import car as c
import stoplight as s

class Car:
    def __init__(self, events):
        self.events = events
        self.action = None

    def Brake(self):
        self.action = "applying brake"

    def Gas(self):
        self.action = "applying gas"

    def Decide(self):
        if random.random() < 0.5:
            self.events.append(c.Go)
        else:
            self.events.append(c.Stop)


def light(stoplight, states, events):
    ssm = stoplight.Machine(None, states[0])
    while True:
        light_state = ssm.tick(s.Timer)
        events.append(light_state)

        if light_state == states[0]:
            time.sleep(7)
        elif light_state == states[1]:
            time.sleep(2)
        elif light_state == states[2]:
            time.sleep(5)


def main():
    light_states = [s.Red, s.Yellow, s.Green]
    events = []

    Thread(target=light, args=(s, light_states, events)).start()

    handlers = Car(events)
    car = c.Machine(handlers, c.Stopped)

    while True:
        time.sleep(1)
        if events:
            event = events.pop(0) 
            state = car.tick(event)
            print(f"{event} | {handlers.action} | car is {state}")


if __name__ == "__main__":
    main()
