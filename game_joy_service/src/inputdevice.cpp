#include <iostream>
#include <vector>
#include <fcntl.h>
#include <pwd.h>
#include <grp.h>
#include <errno.h>
#include <unistd.h>
#include <linux/input.h>
#include "inputdevice.h"

#define INPUT_DEV_PATH "/dev/input"

void InputDevice::FreeRes( int num, dirent **name_list ) {
    for(int i=0; i<num; i++) free(name_list[i]);
    free(name_list);
}

vector<InputDeviceList> InputDevice::ScanDevices(void) {
    dirent **name_list;

    int num = scandir(INPUT_DEV_PATH, &name_list, [](auto ent)-> int{return ent->d_type==DT_CHR;}, alphasort);

    devices_.clear();
    for( int i=0; i<num; i++) {
        char name[256];
        name[0]=0;
        auto path = string(INPUT_DEV_PATH)+"/"+name_list[i]->d_name;
        auto fd = open(path.c_str(), O_RDONLY);

        if(fd<0) {
            InputDevice::FreeRes(num, name_list);
            throw "cannot open device:" +path;
        }
        int ret = ioctl(fd, EVIOCGNAME(sizeof(name)), name);
        close(fd);
        devices_.push_back(InputDeviceList{path, name});
    }
    InputDevice::FreeRes(num, name_list);
    return devices_;
}

void InputDevice::PrintDevices() const {
    for( auto &dev : devices_)
        cout << dev.device_path << " " << dev.device_name << endl;
}

vector<string> InputDevice::GetDevicesByName( const string &name ) const {
    vector<string> match_devs;
    for( auto &dev : devices_)
        if( dev.device_name.find(name)!=string::npos) match_devs.push_back(dev.device_path);
    return match_devs;
}

string InputDevice::GetSingleDeviceByName(const string &name) const {
    for( auto &dev : devices_)
        if( dev.device_name.find(name)!=string::npos) return dev.device_path;
    throw "GetSingleDeviceByName:\"" + name + "\" not found";
}

