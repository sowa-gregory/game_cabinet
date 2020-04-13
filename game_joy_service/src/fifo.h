#pragma once

#include <optional>
#include <string>

using namespace std;

#define FIFO_BUFFER_SIZE 1000

class Fifo
{
private:
    char buffer_[FIFO_BUFFER_SIZE];
    int buffer_pos_;
    int fifo_fd_;
    int FindNewLine();
    void ReadToBuffer();
    optional<string>GetLine();


public:
    Fifo(const string &path);
    ~Fifo();
    optional<string> ReadLine();
};
