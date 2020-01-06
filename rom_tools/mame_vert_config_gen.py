"""

The script generates configuration files for vertical games roms.
Individual configuration for each rom is located in roms directory:
/home/pi/RetroPie/roms/mame-libretro/
.cfg extension is added to original rom file name.
Configuration file is identical for each vertical game and just contains
reference to general configuration for vertical games located at:
"/opt/retropie/configs/mame-libretro/retroarch-vert.cfg"

Detection of game orientation relies on proportions of game image. 
Images of games are located in subdirectory images of roms directory.
Height of image must be larger than width to categorize game as vertical.
Naming convention for image is romname-image.jpg.
The script scans all game's images without checking, if corresponding game exists.
 
"""
import os
import re
CONFIGS_DIR="/home/pi/RetroPie/roms/mame-libretro/"
IMAGES_DIR="/home/pi/RetroPie/roms/mame-libretro/images"
IMAGE_NAME_POSTFIX="-image.jpg"
CONFIG_TEMPLATE='#include "/opt/retropie/configs/mame-libretro/retroarch-vert.cfg"'
CONFIG_POSTFIX=".zip.cfg"

def save_config(file_name:str):
    with open(CONFIGS_DIR+convert_name(file_name), "w") as file:
        file.write(CONFIG_TEMPLATE)
        
def convert_name(file_name:str)->str:
    return file_name[:-len(IMAGE_NAME_POSTFIX)]+CONFIG_POSTFIX

counter=1
for file_entry in os.scandir(IMAGES_DIR):
    if not file_entry.name.endswith(IMAGE_NAME_POSTFIX): continue
    with os.popen("file "+file_entry.path) as cmd:
        info = cmd.read()
        res = re.findall( r".* JPEG image data, [a-z|0-9|,| ]+ ([0-9]{3,4})x([0-9]{3,4}),.*", info)
        if not res: continue
        width=res[0][0]
        height=res[0][1]
        if width>height: continue
        print(counter, file_entry.name)
        save_config(file_entry.name)
        counter=counter+1
        

