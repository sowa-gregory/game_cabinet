#pragma once

#include <stdio.h>
#include <string>
#include <vector>

using namespace std;

class JoyProxy {

  private:
    int input_fd_len_;
    int *input_fd_;
    void OpenInputDevs(vector<string> input_devs);

  public:
    JoyProxy(vector<string> input_devs);
    ~JoyProxy();
    void start();
};
