#include <iostream>
#include <vector>
#include <sys/types.h>
#include <dirent.h>
#include <ctype.h>
#include <functional>
#include <sstream>
#include <sys/stat.h>
#include <stdlib.h>
#include <linux/sched.h>
#include "process.h"

#define PROC_PATH "/proc"
#define MAX_COMM_LEN 15

using namespace std;

string Process::ReadCommFile(const string& proc_id_string)
{
	ifstream input_stream(string(PROC_PATH)+"/"+proc_id_string+"/comm");
	ostringstream os;
	os<<input_stream.rdbuf();
	string out = os.str();
	out.pop_back();
	return out;		
}

int Process::FindByCommName(const string& proc_name) {
	DIR *dir;
    dirent *ent;
	if(proc_name.length()>MAX_COMM_LEN) throw "FindByCommName proc_name > " + to_string(MAX_COMM_LEN);
	dir = opendir(PROC_PATH);
	if(dir==NULL) throw "cannot open /proc dir";

	auto res = [=]()mutable->int
	{
		while((ent=readdir(dir))!=NULL)
		{
			// skip non numeric entries
			if(!isdigit(ent->d_name[0]))continue;
			auto temp_proc_name = ReadCommFile(ent->d_name);
			if(temp_proc_name.find(proc_name)!=string::npos) return stoi(ent->d_name);
		}
		return 0;
	}();

	cout << res << endl;
	closedir(dir);
	return res;
}


bool Process::FindById(int proc_id) {
	struct stat sb;
	auto path = string(PROC_PATH)+"/"+to_string(proc_id);
	auto res = stat(path.c_str() , &sb);
	if( res==-1) return false;
	return true;
}

