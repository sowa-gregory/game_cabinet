#pragma once

#include <optional>
#include <string>

using namespace std;

class Fifo
{
    private:
        int fifo_fd_;

	public:
        void Open();
        optional<string> ReadLine();
};
