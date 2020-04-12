#include "usergroups.h"
#include <pwd.h>
#include <grp.h>
#include <errno.h>
#include <unistd.h>

string UserGroups::get_user_name() {
    const auto pw = getpwuid(getuid());
    if(pw==NULL) throw "getpwuid error";
    return pw->pw_name;
}

vector<string> UserGroups::get_groups() {
    const auto pw = getpwuid(getuid());
    if(pw==NULL) throw "getpwuid error";

    int ngroups=0; // must be initialized to 0 to get number of groups
    getgrouplist(pw->pw_name, pw->pw_gid, NULL, &ngroups);

    const auto groups = new __gid_t[ngroups];
    getgrouplist(pw->pw_name, pw->pw_gid, groups, &ngroups);

    vector<string> user_groups;
    for( int i=0; i<ngroups; i++) {
        auto gr = getgrgid(groups[i]);
        if(gr==NULL) {
            delete []groups;
            throw "getgrgid error";
        }
        user_groups.push_back(gr->gr_name);
    }
    delete []groups;
    return user_groups;
}

bool UserGroups::is_member_of(const string &group_name) {
    for( auto &grp : UserGroups::get_groups())
        if(grp.compare(group_name)==0) return true;

    return false;
}

