# remove_mature.py - script to remove mame roms with games for adults.
#
# It reads gamelist.xml index, identifies adult games by checking if genre contains "*mature*".
# Rom files are renamed by adding .adult extension. Corresponding game screenshot is remove from images directory (if exists).
#
# This script must be invoked inside /home/pi/RetroPie/roms/mame-libretro directory.
# 

import xmltodict
import os
with open ("gamelist.xml") as fd:
    doc = xmltodict.parse(fd.read())

mature_games= [game['path'] for game in doc['gameList']['game'] if game['genre'] and "*Mature*" in game['genre']]

for game in mature_games:
    print(game)
    os.rename( game, game+".adult")
