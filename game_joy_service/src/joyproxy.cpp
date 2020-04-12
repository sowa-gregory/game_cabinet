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

JoyProxy::JoyProxy(vector<string> input_devs) {
    // open input devices
    input_fd_len_ = input_devs.size();
    input_fd_ = new int[input_fd_len_];
    for( int i=0; i<=input_fd_len_; i++) input_fd_[i]=-1;

    OpenInputDevs(input_devs);
}

JoyProxy::~JoyProxy() {
    cout << "Closing devices" << endl;
    for(int i=0; i<input_fd_len_; i++) close(input_fd_[i]);
    delete []input_fd_;
    input_fd_=NULL;
}

void JoyProxy::OpenInputDevs(vector<string> input_devs) {
    for( int i = 0 ; i< input_fd_len_; i++) {
        auto current_dev = input_devs[i];
        int fd = open(current_dev.c_str(), O_RDONLY);
        if( fd < 0 ) throw "Cannot open input device:" + current_dev;
        input_fd_[i]=fd;
        cout << "Input device opened:" << current_dev << endl;
    }
}



void JoyProxy::start() {
    fd_set set;
    input_event ev;
    timeval timeout;
    unsigned int size;
    long counter = 0;
    int current_fd;

    // maximum value of observed file descriptor
    int max_fd = 0;
    for(int i=0; i< input_fd_len_ ; i++)
        if( input_fd_[i] > max_fd) max_fd=input_fd_[i];


    int rv;
    while(true) {
        FD_ZERO(&set);

        // set file descriptors to observe
        for( int i=0; i<input_fd_len_; i++) FD_SET(input_fd_[i], &set);

        timeout.tv_sec = 0;
        timeout.tv_usec = 200000;

        rv = select( max_fd+1, &set, NULL, NULL, &timeout);
        if( rv < 0 ) {
            cerr << "error ################" << endl;
            continue;
        }
        if( rv == 0 ) { // timeout
            cout << "timeout" << endl;
            continue;
        }

        // read selected fd's
        for( int i=0; i<input_fd_len_; i++) {
            current_fd = input_fd_[i];
            if(FD_ISSET(current_fd, &set)) {
                size = read(current_fd, &ev, sizeof(input_event));
                //cout << current_fd << " :  " << counter++ << " " << ev.time.tv_sec << " " << ev.type << " " << ev.code << " " << ev.value << endl;
            }
        }
    }

}
