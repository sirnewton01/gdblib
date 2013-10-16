# Introduction

This library provides an interface between Go and the gdb/MI debugger. The goal is to allow any Go program to provide a GUI for gdb.

# Installation Notes

Gdblib uses the gdb MI (Machine Interface) to debug yor application. The MI changes from time to time. This version of gdblib should work with gdb versions 7.5 and 7.6. Newer versions of Linux will often come with these versions of gdb but Windows and Mac need a little extra setup.

## Windows
Gdb is available on Windows in either MinGW or Cygwin. To install the MinGW version visit http://www.mingw.org/ to download and install the tool suite (mingw-get-setup.exe). Once MingW is installed run the "MinGW Installer" to add the mingw32-gdb package (under "All Packages"). Make sure to add the "C:\MinGW\bin" directory to your PATH so that gdblib can pick it up.

Note that neither MinGW nor Cygwin versions of gdb support the Go language extensions. If you are adventurous you can install the complete MinGW msys component (gcc compiler, autoconf, automake, etc.) and compile the latest gdb yourself. After installing msys install python 2.7.x from the intaller on the python website. Run the "msys.bat" in the c:\mingw\msys\x.y directory. From the command prompt unzip the gdb source code and run "./configure --with-expat --with-python=/c/Python27 --prefix=/. && make && make install". After the install is finished you can run gdb from the c:\mingw\msys\x.y\bin directory. You may find that gdb dies with a segmentation fault. A patch for the source code is described in this (bug)[https://sourceware.org/bugzilla/show_bug.cgi?id=15924]. You will probably need to fix the safe-load path with an "-iex" parameter.

## Mac OS X
The version of gdb on Mac OS X as part of Xcode is very old and will not work with gdblib. Instead, you can download and compile the latest version of gdb from https://www.gnu.org/software/gdb/download/. Note that there is a bug in all versions of gdb including the latest (7.6.1 at the time of writing this) that prevents the Go language support from being loaded. A patch is available [here](http://sourceware-org.1504.n7.nabble.com/Path-Add-support-for-mach-o-reader-to-be-aware-of-debug-gdb-scripts-td238372.html). After applying the patch you compile it with the Xcode compiler using "./configure --with-expat --with-python && make"

### Mac Codesigning Problem
Mac OS X requires that the debugger binary is signed with a trusted certificate before it can take control of another process. If you see a message in the gdb console similar to "Unable to find Mach task port for process-id 12345: (os/kern) failure (0x5). (please check gdb is codesigned - see taskgated(8))" then you will need to follow these steps.

* Start the Keychain Access application (you can use Spotlight to find it)
* Select Keychain Access -> Certificate Assistant -> Create a Certificate...
    + Choose a name for the certificate
    + Set Identity Type to Self Signed Root
    + Set Certificate Type to Code Signing
    + Activate the "Let me override defaults" option
* Continue on to the "Specify a Location For The Certificate" page
    + Set Keychain to System
* Continue and create the certificate
* Double click on the newly created certificate
    + Set When using this certificate to Always Trust
* Restart the computer (yes, this is a required step)
* Sign the gdb binary by executing the following command
    + codesign -f -s "gdb-cert-name" "location-of-gdb-binary"

## Go Runtime Support

GDB versions 7.6+ come with increased safety precautions for auto-loading scripts, including the Go language helper script. This script not only helps gdb to pretty print variables but also helps gdb to avoid analyzing unitialized variables, which can cause gdb to become unresponsive.

If the Go language helper script initializes properly you will see the message "Loading Go Runtime support." Otherwise, you might encounter the following message.

warning: File "/usr/local/go/src/pkg/runtime/runtime-gdb.py" auto-loading has been declined by your `auto-load safe-path' set to "$debugdir:$datadir/auto-load".
To enable execution of this file add
add-auto-load-safe-path /usr/local/go/src/pkg/runtime/runtime-gdb.py
line to your configuration file "/home/cmcgee/.gdbinit".
...

If you follow this instruction then your debugging experience should be much better. 

# Debug session sometimes freezes and gdb process consumes alot of CPU

There is a problem with the string pretty-printer in the standard Go runtime library for gdb, which causes it to attempt to parse uninitialized strings. If there are alot of uninitialized strings then gdb attempts to transfer alot of target memory to satisfy the pretty printer. There is a small tweak you can make to your runtime-gdb.py script (located in your $GOROOT/src/pkg/runtime directory). Find the StringTypePrinter class in the file and change the to_string() method body to look like this (be mindful of the tabs):

        def to_string(self):
                l = int(self.val['len'])
                if l < 1024 and l > -1:
                        return self.val['str'].string("utf-8", "ignore", l)
                return self.val['len']

Note that this tweak will print out the length of really big strings instead of their value.
