import I2C_driver
import datetime
import cpuload
import signal
import time
import RPi.GPIO as GPIO

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
        
        lcd.lcd_display_string(chr(3)+" "+str(cpu_load)+"%",2)
  
def do_exit():
    lcd.lcd_clear()
    lcd.backlight(0)
    cpu.stop()
    exit()
      
def on_sigkill(signumber, frame):  
  print("kill")                
  do_exit()

lcd = I2C_driver.lcd()
signal.signal(signal.SIGTERM, on_sigkill)
cpu = cpuload.CpuLoad()

lcd.lcd_display_string("Czas gry minal",1)

lcd.lcd_display_string("Wrzuc monete...",2)
time.sleep(30)
lcd.lcd_clear()
GPIO.setmode(GPIO.BCM)
GPIO.setup(4, GPIO.OUT)
GPIO.output(4, False)
show_system_status()




while True:
  try:
    show_system_status()
    time.sleep(1)
  except OSError:
    print("power button press")
  except KeyboardInterrupt:
    do_exit()

