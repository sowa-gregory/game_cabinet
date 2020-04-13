#pragma once

#include <vector>
#include <string>
#include <unistd.h>
#include <fstream>

#define PROC_PATH "/proc"

using namespace std;

class Process
{
	private:
	static string ReadCommFile(const string& proc_id_string);
	public:
	static bool FindById(const int proc_id);
	static int FindByCommName(const string& proc_name);
};

