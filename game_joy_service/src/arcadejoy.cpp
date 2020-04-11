#include <iostream>
#include <vector>
#include <fcntl.h>
#include <pwd.h>
#include <grp.h>
#include <dirent.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>
#include <linux/input.h>

using namespace std;

string get_user_name() {
    passwd *pw = getpwuid(getuid());
    if(pw==NULL) {
        perror("getpwuid error");
        exit(-1);
    }
    return string(pw->pw_name);
}

vector<string> get_user_groups() {
    passwd *pw = getpwuid(getuid());
    if(pw==NULL) {
        perror("getpwuid error");
        exit(-1);
    }

    int ngroups=0; // must be initialized to 0 to get number of groups
    getgrouplist(pw->pw_name, pw->pw_gid, NULL, &ngroups);

    __gid_t groups[ngroups];
    getgrouplist(pw->pw_name, pw->pw_gid, groups, &ngroups);

    vector<string> user_groups;
    for( int i=0; i<ngroups; i++) {
        auto gr = getgrgid(groups[i]);
        if(gr==NULL) perror("getgrgid error");
        user_groups.push_back(gr->gr_name);
    }
    return user_groups;
}

bool has_user_group(const string &group_name) {
    for( auto &grp : get_user_groups()) {
        if(grp.compare(group_name)==0) return true;
    }
    return false;
}

struct InputDeviceList {
    string device_path;
    string device_name;
};

class InputDevice {
  private:
    vector<InputDeviceList> devices_;
  public:
    vector<InputDeviceList> ScanDevices(void);
    void PrintDevices();
    string GetDeviceByName(const string &name);
};


#define INPUT_DEV_PATH "/dev/input"

vector<InputDeviceList> InputDevice::ScanDevices(void) {
    dirent **name_list;

    int num = scandir(INPUT_DEV_PATH, &name_list, [](auto ent)-> int{return ent->d_type==DT_CHR;}, alphasort);

    devices_.clear();
    for( int i=0; i<num; i++) {
        char name[256];
        auto path = string(INPUT_DEV_PATH)+"/"+name_list[i]->d_name;
        auto fd = open(path.c_str(), O_RDONLY);

        if(fd<0) {
            string msg = "cannot open device:" +path;
            perror(msg.c_str());
            exit(-1);
        }
        ioctl(fd, EVIOCGNAME(sizeof(name)), name);
        close(fd);
        devices_.push_back(InputDeviceList{path, name});
    }
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

#define INPUT_GROUP "input"


void joy() {
    if(!has_user_group(INPUT_GROUP)) {
        auto user_name = get_user_name();
        cerr << "User:" << user_name<< " is not in input group!!!" << endl;
        cerr << "use command: sudo usergroup -a -G " << INPUT_GROUP << " " << user_name << endl;
        exit(-1);
    }

    InputDevice input_device;

    auto devices = input_device.ScanDevices();
    input_device.PrintDevices();
    cout << input_device.GetDeviceByName( "Sleep Buttona");

}
int main(void) {
    try {
        joy();
    } catch( string  exc ) {
        cerr << "Exception:"+exc << endl;
    }
    return 0;
}
