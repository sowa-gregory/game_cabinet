import os
import zipfile
import tempfile
import shutil
import subprocess

"""
AmiBerry amiga emulator supports whdload enabled images compressed in lha format only.
The script converts zip images into lha. Before conversion, internal structure of zip archive is verified if it contains 
whdload rom. 
Lha command must be installed in operating system - the script uses it to create lha archive.
Sources for lha are located in lha.tgz file.

Output lha archive must have specific structure to be loaded by amiberry.
Following dump represents correct structure of ShadowOfTheBeast archive.
Please notice, that, there is no ShadowOfTheBeast root directory entry itself - ShadowOfTheBeast.
Archive just contains subdirectories and file entries.

-rwxrwxrwx  1000/1000     5871   13317  44.1% -lh5- c998 Apr 17  2003 ShadowOfTheBeast/Manual
-rwxrwxrwx  1000/1000      881    1396  63.1% -lh5- 613e Sep 12  2006 ShadowOfTheBeast/Manual.info
-rwxrwxrwx  1000/1000     1392    2455  56.7% -lh5- 2ebd May 11  2003 ShadowOfTheBeast/ReadMe
-rwxrwxrwx  1000/1000      883    1396  63.3% -lh5- 76aa Sep 12  2006 ShadowOfTheBeast/ReadMe.info
-rwxrwxrwx  1000/1000     1583    2880  55.0% -lh5- 2498 Apr 29  2003 ShadowOfTheBeast/ShadowOfTheBeast.Slave
-rwxrwxrwx  1000/1000     8116   18243  44.5% -lh5- e73d Sep 12  2006 ShadowOfTheBeast/ShadowOfTheBeast.info
-rwxrwxrwx  1000/1000     2232    6641  33.6% -lh5- 48f4 Apr 17  2003 ShadowOfTheBeast/Solution
-rwxrwxrwx  1000/1000      883    1396  63.3% -lh5- 7c6a Sep 12  2006 ShadowOfTheBeast/Solution.info
-rwxrwxrwx  1000/1000   563431  731600  77.0% -lh5- bc7d Sep 12  2006 ShadowOfTheBeast/disk.1
-rwxrwxrwx  1000/1000   852538  985800  86.5% -lh5- 37e6 Sep 12  2006 ShadowOfTheBeast/disk.2
-rwxrwxrwx  1000/1000      962    1024  93.9% -lh5- 43dc Sep 12  2006 ShadowOfTheBeast/init

Notice:
AmiBerry allows for automatic loading of emulator configuration individually for specific rom.
This functionality works only if internal rom name (root directory inside lha archive is identical to archive filename).
eg. ShadowOfTheBeast.lha
Changes to names of directories inside archive are impossible (there are references inside .info files to .slave files).
The only option is to change archive name. This approach is prone to conflict of output file names.
If conversion is about to generate duplicate archive file, operation is aborted and original zip file name is extended *.duplicate'.


"""

def is_zip_file(file) -> bool:
    return file.is_file() and file.name.endswith(".zip")


def is_root_level_name(name: str) -> bool:
    if((name.endswith("/") and name.count("/") == 1) or name.count("/") == 0):
        return True
    return False


class InvalidContentException(Exception):
    pass

def check_content(game: dict, zip_file: zipfile.ZipFile) -> str:
    content = [name for name in zip_file.namelist()
               if is_root_level_name(name)]

    if len(content) != 2:
        raise InvalidContentException()
    
    if content[0].endswith("/"):
        basename=content[0][:-1] 
    else:
        basename=content[1][:-1]

    if not basename+".info" in content:
        print("!!!")
        raise InvalidContentException("invalid content")
        
    return basename             
    


def convert_game(game: dict):
    
    with tempfile.TemporaryDirectory() as temp_dir:
        with zipfile.ZipFile(game["path"], "r") as zip_file:
            try:            
                base_name=check_content(game, zip_file)
                zip_file.extractall(temp_dir)
                
                res = subprocess.run("lha a /tmp/out.lha "+base_name+"/* "+base_name+".info", shell=True, cwd=temp_dir)            
                if res.returncode!=0:
                    shutil.move(game["path"], game["path"]+".problem")
                    return
                out_name = OUTPUT_DIR+"/"+base_name+".lha"
                if os.path.isfile(out_name):
                    print("confilict of output archive names:"+out_name)
                    shutil.move( game["path"], game["path"]+".duplicate")
                else:
                    shutil.move( "/tmp/out.lha", OUTPUT_DIR+"/"+base_name+".lha")
                    shutil.move( game["path"], game["path"]+".done")                
            except InvalidContentException as err:
                print("!!!!!!!!!!!!!",err)
                shutil.move( game["path"], "roms_problems")

def convert_all():
    files = [{"name": file.name[:file.name.rfind(".")], "path":file.path}
             for file in os.scandir(INPUT_DIR+"/") if is_zip_file(file)]
    total = len(files)
    count = 1
    for game in files:
        print("Processing:"+str(count)+"/"+str(total)+"  "+game["name"])
        count=count+1
        convert_game(game)

INPUT_DIR = "roms"
OUTPUT_DIR = "roms"
convert_all()
