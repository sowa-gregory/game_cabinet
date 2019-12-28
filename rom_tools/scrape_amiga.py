
"""
list of the games is available at the following url:

http://www.lemonamiga.com/games/list.php?&lineoffset=756

lineoffset parameter starts at value:0
each time lineoffset should be increased by 27
to obtain more consistent html structure header of request must contain cookie: "games_display=list".

result html contains sections, which refers to thumbnails of games - which looks this way:
messages['4'] = new Array('/games/screenshots/thumbs/10_out_of_10_-_dinosaurs.png');

another section contains name of games:

<td bgcolor="#D3DAEE"><a href="details.php?id=19" onmouseover="doTooltip(event,14)" onmouseout="hideTip()"><b>1943: The Battle of Midway</b></a><br></td>


http://www.lemonamiga.com/games/screenshots/full/1_across_2_down_01.png
http://www.lemonamiga.com/games/screenshots/thumbs/1_across_2_down.png


http://www.lemonamiga.com/games/screenshots/full/century_03.png


<td bgcolor="#D3DAEE"><a href="details.php?id=19" onmouseover="doTooltip(event,14)" onmouseout="hideTip()"><b>1943: The Battle of Midway</b></a><br></td>

"""

import requests
import re
import json
import stringdist
import os

def download_game_page(offset: int)->list:
    base_url = "http://www.lemonamiga.com/games/list.php?&lineoffset="
    url = base_url+str(offset)
    print("Download url:"+url)
    res = requests.get(url, cookies={"games_display": "list"})
    body = res.content.decode("utf-8", errors='ignore')
    screenshots = re.findall(r'Array\(\'.*/(.*)\'\);', body)
    names = re.findall(r'<td bgcolor=.*><b>(.*)</b></a>.*<br></td>', body)
    if len(screenshots) != len(names):
        raise Exception("error")
    return list(zip(names, screenshots))


def download_game_list()->list:
    offset = 0
    std_list_len = 27
    games = []
    while True:
        chunk = download_game_page(offset)
        games = games + chunk
        print(len(chunk), len(games))
        if len(chunk) < std_list_len:
            break
        offset = offset + len(chunk)
    return games


def load_games_json(file_name:str)->list:
    with open(file_name, "r") as file:
        content = file.read()
        return json.loads(content)

def store_games():
    games = download_game_list()
    with open("games.json", "wt") as file:
        file.write(json.dumps(games))

def permut(name):
    out=[]
    out.append(name)
    if "III" in name:
        out.append(name.replace("III", "3"))
    elif "II" in name:
        out.append(name.replace("II", "2"))
        
    return out

def find_best_match( orig_name, names_list:list )->dict:
    min_res=1.0
    best_name =""
    postfixes_list = ["AGA", "Levels", "Activision", "Sega", "Demo", "Demo1", "Demo2", "CD32", "CDTV", "Psygnosis", "NTSC", "Disk", "DemoPlay"]
    name1=orig_name
    for postfix in postfixes_list:
        if name1.endswith(postfix): name1=name1[:-len(postfix)]
     
    
    if name1.endswith("Fr") or name1.endswith("De") or name1.endswith("Pl"): name1=name1[:-2]
      
    for name2 in names_list:
        if min_res==0.0: break
        for temp_name in permut(name2):
            res = stringdist.levenshtein_norm(name1, temp_name)
            if temp_name.replace(" ", "").lower().startswith(name1.lower()):
                min_res = 0
                best_name = name2
                break
            if res<min_res:
                min_res = res
                best_name = name2
    return {"res":min_res,"retro":orig_name,"lemon":best_name}

def compare( names_list1:list, names_list2:list )->list:
    res_sum = 0
    comp_list=[]
    total = len(names_list1)
    count = 1
    for name1 in names_list1:
        match = find_best_match(name1, names_list2)
        comp_list.append(match)
        res_sum = res_sum + match["res"]
        print(count,total, res_sum*1000/count,  match)
        count=count+1
    return comp_list
    
def download_screenshots():
    games_list=load_games_json("comp_list.json")
    total=len(games_list)
    count=1
    for game in games_list:
        name = game["retro"]
        print("------------------------")
        print(count,total)
        count=count+1
        screenshot = game["screenshot"]
        if screenshot=="adult.gif": 
            print("skip")
            continue
        file_name = "images/"+name+"-image.png"
        if os.path.isfile(file_name):
            print("Already exists:"+file_name)
            continue
        screenshot=screenshot[:-4]
        url = "http://www.lemonamiga.com/games/screenshots/full/"+screenshot+"_03.png"
        print("Download:"+url)
        response = requests.get(url)
        if(response.status_code==200):
            print("Save:"+file_name)
            with open(file_name, "wb") as file:
                for chunk in response:
                    file.write(chunk)
        else:
            print("Error - file doesn't exist")
            


def do_compare():
    games=load_games_json("amiga_lemon.json")
    
    lemon=[name[0] for name in games]
    retropie = load_games_json("amiga_retropie.json")

    #print(retropie[0])
    comp_list=compare(retropie, lemon)
    for item in comp_list:
        screen = None
        lemon_name = item["lemon"]
        screen = [temp[1] for temp in games if temp[0]==lemon_name][0]
        item["screenshot"]=screen
        #extend_comp_list.append(item)
        
    print(comp_list)
        

    with open("comp_list.json", "wt") as file:
        file.write(json.dumps(comp_list))
#do_compare()
download_screenshots()