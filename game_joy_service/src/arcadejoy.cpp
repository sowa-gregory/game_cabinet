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

struct input_device
{
   string device_path;
   string device_name;
};

//vector<struct input_device> scan_devices(void);


vector<input_device> scan_devices(void) {
    dirent **pNameList;

    int iNDev = scandir("/dev/input", &pNameList, [](auto ent)-> int{return ent->d_type==DT_CHR;}, alphasort);

	vector<input_device> devices;
    for( int iIndex=0; iIndex<iNDev; iIndex++) {
        char name[256];
        auto path = string("/dev/input/")+pNameList[iIndex]->d_name;
        auto fd = open(path.c_str(), O_RDONLY);

        if(fd<0) {
            string msg = "cannot open device:" +path;
            perror(msg.c_str());
            exit(-1);
        }
        ioctl(fd, EVIOCGNAME(sizeof(name)), name);
        close(fd);
        devices.push_back(input_device{path, name});
    }
    return devices;
}

void print_devices(vector<input_device> devices)
{
	for( auto &dev : devices)
		cout << dev.device_path << " " << dev.device_name << endl;
}

#define INPUT_GROUP "input"

int main(void) {
    auto res = has_user_group(INPUT_GROUP);
    if(!res) {
        auto user_name = get_user_name();
        cerr << "User:" << user_name<< " is not in input group!!!" << endl;
        cerr << "use command: sudo usergroup -a -G " << INPUT_GROUP << " " << user_name << endl;
        exit(-1);
    }
    cout << res << endl;

    auto devices = scan_devices();
	print_devices(devices);	
    return 0;
}
