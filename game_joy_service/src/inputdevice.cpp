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
        auto path = string(INPUT_DEV_PATH)+"/"+name_list[i]->d_name;
        auto fd = open(path.c_str(), O_RDONLY);

        if(fd<0) {
            FreeRes(num, name_list);
            throw "cannot open device:" +path;
        }
        ioctl(fd, EVIOCGNAME(sizeof(name)), name);
        close(fd);
        devices_.push_back(InputDeviceList{path, name});
    }
    InputDevice::FreeRes(num, name_list);
    return devices_;
}

void InputDevice::PrintDevices() {
    for( auto &dev : devices_)
        cout << dev.device_path << " " << dev.device_name << endl;
}

string InputDevice::GetDeviceByName(const string &name) {
    for( auto &dev : devices_)
        if( dev.device_name.compare(name)==0) return dev.device_path;
    throw "get_device_by_name:\"" + name + "\" not found";
}

