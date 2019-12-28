import I2C_driver
import datetime
import cpuload
import time

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
        
                

lcd = I2C_driver.lcd()
cpu = cpuload.CpuLoad()
show_system_status()

try:
    while True:
        show_system_status()
        time.sleep(1) # 2 sec delay
except KeyboardInterrupt:
        lcd.lcd_clear()
        lcd.backlight(0)
        cpu.stop()
        exit()




# let's define a custom icon, consisting of 6 individual characters
# 3 chars in the first row and 3 chars in the second row
fontdata1 = [
        # Char 0 - Upper-left
        [ 0x00, 0x00, 0x03, 0x04, 0x08, 0x19, 0x11, 0x10 ],
        # Char 1 - Upper-middle
        [ 0x00, 0x1F, 0x00, 0x00, 0x00, 0x11, 0x11, 0x00 ],
        # Char 2 - Upper-right
        [ 0x00, 0x00, 0x18, 0x04, 0x02, 0x13, 0x11, 0x01 ],
        # Char 3 - Lower-left
        [ 0x12, 0x13, 0x1b, 0x09, 0x04, 0x03, 0x00, 0x00 ],
        # Char 4 - Lower-middle
        [ 0x00, 0x11, 0x1f, 0x1f, 0x0e, 0x00, 0x1F, 0x00 ],
        # Char 5 - Lower-right
        [ 0x09, 0x19, 0x1b, 0x12, 0x04, 0x18, 0x00, 0x00 ],
        # Char 6 - my test
	[ 0x1f,0x0,0x4,0xe,0x0,0x1f,0x1f,0x1f],
]

# Load logo chars (fontdata1)
mylcd.lcd_load_custom_chars(fontdata1)


# Write first three chars to row 1 directly
mylcd.lcd_write(0x80)
mylcd.lcd_write_char(0)
mylcd.lcd_write_char(1)
mylcd.lcd_write_char(2)
# Write next three chars to row 2 directly
mylcd.lcd_write(0xC0)
mylcd.lcd_write_char(3)
mylcd.lcd_write_char(4)
mylcd.lcd_write_char(5)
sleep(2)

mylcd.lcd_clear()

mylcd.lcd_display_string_pos("Testing",1,1) # row 1, column 1
sleep(1)
mylcd.lcd_display_string_pos("Testing",2,3) # row 2, column 3
sleep(1)
mylcd.lcd_clear()

# Now let's define some more custom characters
fontdata2 = [
        # Char 0 - left arrow
        [ 0x1,0x3,0x7,0xf,0xf,0x7,0x3,0x1 ],
        # Char 1 - left one bar 
        [ 0x10,0x10,0x10,0x10,0x10,0x10,0x10,0x10 ],
        # Char 2 - left two bars
        [ 0x18,0x18,0x18,0x18,0x18,0x18,0x18,0x18 ],
        # Char 3 - left 3 bars
        [ 0x1c,0x1c,0x1c,0x1c,0x1c,0x1c,0x1c,0x1c ],
        # Char 4 - left 4 bars
        [ 0x1e,0x1e,0x1e,0x1e,0x1e,0x1e,0x1e,0x1e ],
        # Char 5 - left start
        [ 0x0,0x1,0x3,0x7,0xf,0x1f,0x1f,0x1f ],
        # Char 6 - 
        # [ ],
]

# Load logo chars from the second set
mylcd.lcd_load_custom_chars(fontdata2)

block = chr(255) # block character, built-in

# display two blocks in columns 5 and 6 (i.e. AFTER pos. 4) in row 1
# first draw two blocks on 5th column (cols 5 and 6), starts from 0
mylcd.lcd_display_string_pos(block * 2,1,4)

# 
pauza = 0.2 # define duration of sleep(x)
#
# now draw cust. chars starting from col. 7 (pos. 6)

pos = 6
mylcd.lcd_display_string_pos(chr(1),1,6)
sleep(pauza)

mylcd.lcd_display_string_pos(chr(2),1,pos)
sleep(pauza)

mylcd.lcd_display_string_pos(chr(3),1,pos)
sleep(pauza)

mylcd.lcd_display_string_pos(chr(4),1,pos)
sleep(pauza)

mylcd.lcd_display_string_pos(block,1,pos)
sleep(pauza)

# and another one, same as above, 1 char-space to the right
pos = pos +1 # increase column by one

mylcd.lcd_display_string_pos(chr(1),1,pos)
sleep(pauza)
mylcd.lcd_display_string_pos(chr(2),1,pos)
sleep(pauza)
mylcd.lcd_display_string_pos(chr(3),1,pos)
sleep(pauza)
mylcd.lcd_display_string_pos(chr(4),1,pos)
sleep(pauza)
mylcd.lcd_display_string_pos(block,1,pos)
sleep(pauza)


#
# now again load first set of custom chars - smiley
mylcd.lcd_load_custom_chars(fontdata1)

mylcd.lcd_display_string_pos(chr(0),1,9)
mylcd.lcd_display_string_pos(chr(1),1,10)
mylcd.lcd_display_string_pos(chr(2),1,11)
mylcd.lcd_display_string_pos(chr(3),2,9)
mylcd.lcd_display_string_pos(chr(4),2,10)
mylcd.lcd_display_string_pos(chr(5),2,11)

sleep(2)
mylcd.lcd_clear()
sleep(1)
mylcd.backlight(0)
