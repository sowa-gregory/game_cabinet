#include <iostream>
#include <string>
#include "inputdevice.h"
#include "usergroups.h"

using namespace std;

#define INPUT_GROUP "input"

void joy() {

    if(!UserGroups::is_member_of(INPUT_GROUP)) {
        auto user_name = UserGroups::get_user_name();
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
