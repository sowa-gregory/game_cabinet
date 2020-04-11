#include <iostream>
#include <pwd.h>
#include <grp.h>
#include <errno.h>
#include <unistd.h>
#include <linux/input.h>
#include "inputdevice.h"

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

    auto groups = new __gid_t[ngroups];
    getgrouplist(pw->pw_name, pw->pw_gid, groups, &ngroups);

    vector<string> user_groups;
    for( int i=0; i<ngroups; i++) {
        auto gr = getgrgid(groups[i]);
        if(gr==NULL) {
            delete []groups;
            throw "getgrgid error";
        }
        user_groups.push_back(gr->gr_name);
    }
    delete []groups;
    return user_groups;
}

bool has_user_group(const string &group_name) {
    for( auto &grp : get_user_groups()) {
        if(grp.compare(group_name)==0) return true;
    }
    return false;
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
    cout << input_device.GetDeviceByName( "Sleep Button") << endl;

}
int main(void) {
    try {
        joy();
    } catch( string  exc ) {
        cerr << "Exception:"+exc << endl;
    }
    return 0;
}
