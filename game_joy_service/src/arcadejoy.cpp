#include <iostream>
#include <fcntl.h>
#include <pwd.h>
#include <dirent.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>
#include <linux/input.h>

using namespace std;

void get_user_groups()
{
	__uid_t uid = getuid();
	passwd *pw = getpwuid(uid);
	if(pw==NULL)
	{
		perror("getpwuid error");
	}
}

int is_event_device( const dirent *pDir )
{
	return strncmp( "event", pDir->d_name, 5) == 0;
}

void scan_devices(void)
{
	dirent **pNameList;

	int iNDev = scandir("/dev/input", &pNameList, is_event_device, alphasort);
	cout << iNDev << endl;

	for( int iIndex=0;iIndex<iNDev;iIndex++)
	{
		char name[256];
		auto strDevice = string("/dev/input/")+pNameList[iIndex]->d_name;
		auto fd = open(strDevice.c_str(), O_RDONLY);
		cout << strDevice;
		if( fd>=0)
		{
			ioctl(fd, EVIOCGNAME(sizeof(name)), name);
			close(fd);
			cout << fd << endl;
		}
		printf("%s\n", name);
		cout << strDevice << "  " << name << endl;
	}
}

int main(void)
{
	get_user_groups();
	scan_devices();
	return 0;
}
