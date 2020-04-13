#include <iostream>
#include <string>
#include <thread>
#include <chrono>
#include "inputdevice.h"
#include "usergroups.h"
#include "joyproxy.h"
#include "colormod.h"
#include "process.h"
#include "fifo.h"

using namespace std;

#define INPUT_GROUP "input"

Color::Modifier blue(Color::FG_BLUE);
Color::Modifier def(Color::FG_DEFAULT);

void joy() {

    if(!UserGroups::is_member_of(INPUT_GROUP)) {
        auto user_name = UserGroups::get_user_name();
        cerr << "User:" << user_name<< " is not in input group!!!" << endl;
        cerr << "use command: sudo usergroup -a -G " << INPUT_GROUP << " " << user_name << endl;
        exit(-1);
    }
    InputDevice input_device;

    input_device.ScanDevices();
    cout <<  "Detecting input devices..." << endl << blue;
    input_device.PrintDevices();

    cout << def << "Looking for DragonRise joysticks..." << endl;
    vector<string> devices = input_device.GetDevicesByName( "DragonRise");

    auto joy_proxy = JoyProxy(devices);
    joy_proxy.start();
}


void fifotest()
{
  Fifo fifo("/tmp/fif");
  cout << "adfs" << endl;
  while(true)
  {
    if( auto out = fifo.ReadLine() )
    cout << *out << endl;
  }
}

int main(void) {
	fifotest();
	
    try {

	exit(-1);
        joy();
    } catch( string  exc ) {
        cerr << "Exception:"+exc << endl;
    }
    return 0;
}
