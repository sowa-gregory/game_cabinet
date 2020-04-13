#pragma once

#include <stdio.h>
#include <string>
#include <vector>

using namespace std;

struct input_event;

class JoyProxy {

  private:
    int input_fd_len_;
    int *input_fd_;
    void OpenInputDevs(vector<string> input_devs);

  public:
    JoyProxy(vector<string> input_devs);
    ~JoyProxy();
    void Start() const;
	void OnButtonEvent(const int joy_id, const input_event& ev) const;
};
