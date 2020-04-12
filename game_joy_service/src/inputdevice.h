#pragma once

#include <dirent.h>
#include <string>
#include <vector>

using namespace std;

struct InputDeviceList {
    string device_path;
    string device_name;
};

class InputDevice {
  private:
    vector<InputDeviceList> devices_;
    static void FreeRes(int num, dirent **name_list);
  public:
    vector<InputDeviceList> ScanDevices(void);
    void PrintDevices() const;
    string GetSingleDeviceByName(const string &name) const;
	vector<string> GetDevicesByName(const string &name) const;
};

