import time
import subprocess
import RPi.GPIO as GPIO

class GPIOButton:

    def __init__(self, button_id, long_press_millis=1000):
        self.__start_time = 0
        self.__long_press_millis = long_press_millis
        GPIO.setmode(GPIO.BCM)
        GPIO.setup(button_id, GPIO.IN)
        GPIO.add_event_detect(button_id, GPIO.BOTH, callback=self.__on_button)
            
    def __on_button(self, button_id):
        state = not GPIO.input(button_id)
        # button is pulled up, state must be reversed    
        if state:
            self.__start_time=time.time()*1000
        else:
            if self.__start_time==0: return
            delay = time.time()*1000 - self.__start_time
            self.__start_time = 0
            if delay > self.__long_press_millis:
                self.on_long_press(button_id)

    def on_long_press(self, button_id):
        print("long press", button_id)

class PowerButton(GPIOButton):
    def on_long_press(self, button_id):
        print("switching off")
        sleep(2)
        subprocess.run("killall emulationstation", shell=True)
        subprocess.run("sudo shutdown -h now", shell=True)
        
        
button = PowerButton(3, 3000)

while True:
    time.sleep(1)

GPIO.cleanup()