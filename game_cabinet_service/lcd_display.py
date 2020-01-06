import I2C_driver
import subprocess
import datetime
import cpuload
import signal
import time
import RPi.GPIO as GPIO


class Relay:
  _initialized = False
  
  def _initialize():
    if not Relay._initialized:
      GPIO.setmode(GPIO.BCM)
      GPIO.setup(4, GPIO.OUT)
      Relay._initialized = True

  def __init__(self):
    Relay._initialize()
    
  def switch(self, state:bool):
    Relay._initialize()
    GPIO.output(4, not state)
         
    

def check_running_process():
    pass

def get_cpu_load()->int:
    load = cpu.getlastload()
    if load: return int(load['cpu'])
    return 0



def get_system_temperature()->str:
    with open("/sys/class/thermal/thermal_zone0/temp") as file:
        return str(int(file.read())//1000)
    
def show_system_status():
        cpu_load = get_cpu_load()
        time_str = datetime.datetime.now().strftime("%H:%M:%S")
        temp_str = get_system_temperature()

        
        font_data=[
        ## thermometer
        [0b00100, 0b01010,  0b01010, 0b01110, 0b01110, 0b11111,  0b11111, 0b01110],
        
  ## clock      
  [0b01110,
        0b10001,
        0b10001,
        0b10101,
        0b10101,
        0b11001,
        0b10001,
        0b01110],
## celsius
 [0b11000,  0b11000, 0b00011,   0b00100,  0b00100,  0b00100,  0b00100,  0b00011]
  ,
## cpu
  [0x0A, 0x1F,  0x0E,  0x1F,  0x1F,  0x0E,  0x1F,  0x0A]
  ]
        lcd.lcd_load_custom_chars(font_data)
        lcd.lcd_display_string(chr(1)+" "+time_str+"  "+chr(0)+temp_str+chr(2),1)
        
        lcd.lcd_display_string(chr(3)+" "+str(cpu_load)+"% ",2)
  
def do_exit():
    lcd.lcd_clear()
    lcd.backlight(0)
    cpu.stop()
    print("switch off relay")
    Relay().switch(False)
    exit()
      
def on_sigkill(signumber, frame):  
  print("kill")                
  do_exit()


def wait_time_sync():
    pos=0
    lcd.lcd_display_string("wait time sync..",1)
    
    while True:
      lcd.lcd_display_string_pos(".", 2, pos)
      if pos<16: pos=pos+1
      res=subprocess.run("/usr/bin/chronyc tracking", shell=True, capture_output=True)
      for line in res.stdout.split(b"\n"):
        line = line.decode("ascii")
        line = " ".join(line.split())
        if "Leap status : Normal" in line:
           return
      time.sleep(1)
    
    

lcd = I2C_driver.lcd()

signal.signal(signal.SIGTERM, on_sigkill)
print("signal set")

# wait for timed synchronization
lcd.lcd_clear()
wait_time_sync()
cpu = cpuload.CpuLoad()
Relay().switch(True)

lcd.lcd_clear()
show_system_status()




while True:
  try:
    show_system_status()
    time.sleep(1)
  except OSError:
    print("power button press")
  except KeyboardInterrupt:
    do_exit()

