#include <iostream>
#include <stdio.h>
#include <string>
#include <vector>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <fcntl.h>
#include <optional>
#include <string.h>
#include "fifo.h"

using namespace std;

Fifo::Fifo(const string &path)
{
	memset(buffer_, 0, FIFO_BUFFER_SIZE);
	buffer_pos_ = 0;
	fifo_fd_ = open(path.c_str(), O_RDWR);
}

Fifo::~Fifo()
{
	close(fifo_fd_);
}

int Fifo::FindNewLine()
{
	int ind = 0;
	while (ind < buffer_pos_ && buffer_[ind] != '\n')
		ind++;

	if (ind == buffer_pos_)
		return -1;
	return ind;
}

void Fifo::ReadToBuffer()
{
	fd_set set;
	FD_ZERO(&set);
	FD_SET(fifo_fd_, &set);
	
	int res = select(fifo_fd_ + 1, &set, NULL, NULL, NULL);
	if (res <= 0)
	{
		cerr << "select error" << endl;
		return;
	}

	// reached end of the buffer
	int buf_free_space = FIFO_BUFFER_SIZE - buffer_pos_;
	if (!buf_free_space)
	{
		cerr << "buffer overflow" << endl;
		return;
	}
	int read_len = read(fifo_fd_, buffer_ + buffer_pos_, buf_free_space);
	buffer_pos_ += read_len;
}

optional<string> Fifo::GetLine()
{
	// find new line in buf
	int nl_pos;
	if ((nl_pos = FindNewLine()) < 0)
		return {};

	string str = string(buffer_, nl_pos);
	memcpy(buffer_, buffer_ + nl_pos + 1, buffer_pos_ - nl_pos - 1);
	buffer_pos_ -= nl_pos + 1;
	return str;
}

optional<string> Fifo::ReadLine()
{
	if (auto res = GetLine())
		return *res;
	ReadToBuffer();
	return GetLine();
}
