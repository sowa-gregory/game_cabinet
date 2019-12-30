import os
import zipfile
import tempfile
import shutil
import subprocess
import binascii

"""
Detects broken whdload lha archives. Those, which contains more than one directory inside.
"""


def is_root_level_name(name: str) -> bool:
    if((name.endswith(b"/") and name.count(b"/") == 1) or name.count(b"/") == 0):
        return True
    return False


class InvalidContentException(Exception):
    pass

def check_content(game, files_list:list):
    base_name = [name for name in files_list if name.endswith(b".info") and not b"/" in name]
    
    
    if len(base_name) != 1: return False
    print(game.name, base_name)
    if not game.name[:-4] in base_name[0].decode("ascii"): raise Exception()
    
    base_name = base_name[0]
    check_name = [name for name in files_list if name!=base_name and not name.startswith(base_name[:-5]+b"/")]
    return len(check_name)==0
    


def get_files_list(game)->list:
    files_list=[]    
    with tempfile.TemporaryDirectory() as temp_dir:
        #base_name=check_content(game, zip_file)
        res = subprocess.run("lha v "+game.name, shell=True, cwd=INPUT_DIR, capture_output=True)
        files_start = False
        for line in res.stdout.split(b"\n"):
          columns=line.split()
          if not columns: continue
          if b"------" in columns[0]: 
            files_start = not files_start
            continue
          if files_start:
            files_list.append( columns[10])
    return files_list
      
def check_game(game):
    files = get_files_list(game)
    return check_content(game,files)
    
def is_lha_file(file)->bool:
    return file.is_file() and file.name.endswith(".lha") 

def check_all():
    files = [file for file in os.scandir(INPUT_DIR+"/") if is_lha_file(file)]
    total = len(files)
    count = 1
    print("start")
    with open("invalid.out", "w") as file:
    
      for game in files:
          print("Processing:"+str(count)+"/"+str(total)+"  "+game.name)
          count=count+1
          res = check_game(game)
          print(res)
          if not res: shutil.move( INPUT_DIR+"/"+game.name, OUTPUT_DIR)          

#INPUT_DIR = "/home/pi/RetroPie/roms/amiga"
INPUT_DIR = "roms_done"
OUTPUT_DIR = "romsi"
check_all()
