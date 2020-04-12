#pragma once

#include <string>
#include <vector>

using namespace std;

class UserGroups
{
  public:
	static string get_user_name();
	static vector<string> get_groups();
	static bool is_member_of(const string& group_name);
};

