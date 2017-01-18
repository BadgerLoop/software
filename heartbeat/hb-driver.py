import riffle
import time
import os
import sys

fabric = "ws://192.168.1.99:9000"
interval_ms = int(sys.argv[1])
current_milli_time = lambda: int(round(time.time() * 1000))
class HB(riffle.Domain):

    def onJoin(self):
        print("Connected to Exis Node")
        while True:
            self.publish("cmd","5C0#1301")
            self.publish("can",[[str(current_milli_time()),"100","00","00 00"],[str(current_milli_time()),"100","00","00 00"]]) # Using for tests
            self.publish("hb","5C0#1301")
            print("Sent WCM heartbeat")
            time.sleep(float(interval_ms*.001))
            print("waiting...")

if __name__ == '__main__':
    print("Starting HB Driver")
    riffle.SetFabric(fabric)
    HB('xs.node').join()