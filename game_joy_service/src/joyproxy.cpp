#include <iostream>
#include <stdio.h>
#include <string>
#include <vector>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <fcntl.h>
#include <linux/input.h>
#include "joyproxy.h"

using namespace std;

JoyProxy::JoyProxy(vector<string> input_devs)
{
	// open input devices
	input_fd_len_ = input_devs.size();
	input_fd_ = new int[input_fd_len_];
	for( int i=0;i<=input_fd_len_;i++) input_fd_[i]=-1;

	OpenInputDevs(input_devs);
}

JoyProxy::~JoyProxy()
{
	cout << "Closing devices" << endl;
	for(int i=0;i<input_fd_len_;i++) close(input_fd_[i]);
	delete []input_fd_;
	input_fd_=NULL;
}

void JoyProxy::OpenInputDevs(vector<string> input_devs)
{
	for( int i = 0 ; i< input_fd_len_; i++)
{
	auto current_dev = input_devs[i];
	int fd = open(current_dev.c_str(), O_RDONLY);
	if( fd < 0 ) throw "Cannot open input device:" + current_dev;
	input_fd_[i]=fd;
	cout << "Input device opened:" << current_dev << endl;
}
}
