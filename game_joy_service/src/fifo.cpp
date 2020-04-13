
int fifo_fd_;

#include <iostream>
#include <stdio.h>
#include <string>
#include <vector>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <fcntl.h>
#include <optional>
#include "fifo.h"

using namespace std;

void Fifo::Open()
{
	fifo_fd_ = open( "/tmp/arcade_mon1", O_RDONLY);
}

optional<string> Fifo::ReadLine()
{
	fd_set set;
	FD_ZERO(&set);
	FD_SET(fifo_fd_, &set);
	
}
