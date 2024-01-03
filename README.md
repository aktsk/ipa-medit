# ipa-medit

[![GitHub release](https://img.shields.io/github/v/release/aktsk/ipa-medit.svg)](https://github.com/aktsk/ipa-medit/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/aktsk/ipa-medit/blob/main/LICENSE)
[![](https://img.shields.io/badge/Black%20Hat%20Arsenal-USA%202021-blue.svg)](https://www.blackhat.com/us-21/arsenal/schedule/index.html#ipa-medit-memory-search-and-patch-tool-for-ipa-without-jailbreaking-24072)
![](https://github.com/aktsk/ipa-medit/workflows/test/badge.svg?branch=main)

Ipa-medit is a memory search and patch tool for resigned ipa without jailbreaking. 
It supports iOS apps running on iPhone and Apple Silicon Mac.
It was created for mobile game security testing.
Many mobile games have jailbreak detection, but ipa-medit does not require jailbreaking, so memory modification can be done without bypassing the jailbreak detection.

## Motivation
Memory modification is the easiest way to cheat in games, it is one of the items to be checked in the security test.
There are also cheat tools that can be used casually like GameGem and iGameGuardian.
However, there were no tools available for un-jailbroken device and CUI, Apple Silicon Mac.
So I made it as a security testing tool.
Android version is [aktsk/apk-medit](https://github.com/aktsk/apk-medit).

## Demo
<img src="screenshots/desktop.gif" width=850px>

## Requirements
- macOS
  - You need to have a valid iOS Development certificate installed.
- Only when targeting iOS apps running on an iPhone.
  - Xcode
  - [libimobiledevice/libimobiledevice](https://github.com/libimobiledevice/libimobiledevice)
  - [libimobiledevice/ideviceinstaller](https://github.com/libimobiledevice/ideviceinstaller)

```
$ brew install --HEAD libplist
$ brew install --HEAD usbmuxd
$ brew install --HEAD libimobiledevice
$ brew install --HEAD ideviceinstaller
```

## Installation
### Binary(Intell Mac Only)
Download the binary from [GitHub Releases](https://github.com/aktsk/ipa-medit/releases/) and drop it in your $PATH.

### Manually Build
You can build it by using the make command.
Go compiler is required to build.
If you are targeting an iOS app that runs on an Apple Silicon Mac, you will need to sign it, but `script/codesign.sh` will be executed and signed automatically.

```
$ git clone git@github.com:aktsk/ipa-medit.git
$ cd ipa-medit
$ make build
```

## Usage

The target .ipa file must be signed with a certificate installed on your computer. 
If you want to modify memory on third-party applications, please use a tool such as [ipautil](https://github.com/aktsk/ipautil) for re-signing.

```
$ ipautil decode tap1000000.ipa # unzip
$ ipautil build Payload         # re-sign and generate .ipa file
```

### Targeting the iOS app on iPhone

To launch it, you need to specify the executable file path contained in the .ipa file with `-bin` and the bundle id with `-id`.

```
$ unzip tap1000000.ipa
$ ipa-medit -bin="./Payload/tap1000000.app/tap1000000" -id="jp.hoge.tap1000000"
```

### Targeting the iOS app on Apple Silicon Mac

To launch it, you need to specify the process name with `-name` or the pid with `-pid`. 
The process name and pid of the iOS app can be checked in the Activity Monitor.

```
$ ipa-medit -name <process name>
```

### Commands
Here are the commands available in an interactive prompt.

#### find
Search the specified integer on memory.

```
> find 999986
Success to halt process
Scanning: 0x00000001025e4000-0x00000001025e8000
Scanning: 0x00000001025f4000-0x00000001025fc000
Scanning: 0x0000000102604000-0x0000000102608000
....
Scanning: 0x000000016eb34000-0x000000016ebbc000
Scanning: 0x000000016ebc0000-0x000000016ebe8000
Scanning: 0x000000016ebec000-0x000000016ec74000
Scanning: 0x000000016ec78000-0x000000016ed00000
Found: 1!!
Address: 0x10a2feea0
```

By default, only integers are searched when targeting iOS apps running on iPhone, because the LLDB API is slow.
When targeting an iOS app running on Apple Silicon Mac, strings will also be searched.

You can also specify datatype such as string, word, dword, qword.

```
> find dword 999994
Search Double Word...
Target Value: 999994([58 66 15 0])
Found: 1!!
Address: 0x11378aea0
```

#### filter
Filter previous search results that match the current search results.

```
> filter 999842
Success to halt process
Found: 1!!
Address: 0x1087beea0
```

#### patch
Write the specified value on the address found by search.

```
> patch 10
Successfully patched!
```

#### attach
Attach to the target process.

```
> attach
Success to halt process
```

#### detach
Detach from the attached process.

```
> detach
Success to continue process
```

#### ps
Get information about the target process.
It will only work if you are targeting an iOS app running on an iPhone.

```
> ps
SBProcess: pid = 926, state = running, threads = 37, executable = tap1000000
State: Running
thread #1: tid = 0x545ee, 0x00000001bd6552d0 libsystem_kernel.dylib`mach_msg_trap + 8, queue = 'com.apple.main-thread'
thread #3: tid = 0x54619, 0x00000001bd67a184 libsystem_kernel.dylib`__workq_kernreturn + 8
thread #4: tid = 0x5461a, 0x00000001bd67a184 libsystem_kernel.dylib`__workq_kernreturn + 8
thread #5: tid = 0x5461b, 0x00000001bd67a184 libsystem_kernel.dylib`__workq_kernreturn + 8
thread #6: tid = 0x5461c, 0x00000001bd67a184 libsystem_kernel.dylib`__workq_kernreturn + 8
thread #7: tid = 0x5461d, 0x00000001bd6552d0 libsystem_kernel.dylib`mach_msg_trap + 8, name = 'com.apple.uikit.eventfetch-thread'
thread #8: tid = 0x5461e, 0x00000001bd6791ac libsystem_kernel.dylib`__psynch_cvwait + 8, name = 'GC Finalizer'
thread #9: tid = 0x5461f, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Job.Worker 0'
thread #10: tid = 0x54620, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Job.Worker 1'
thread #11: tid = 0x54621, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Job.Worker 2'
...
thread #35: tid = 0x54659, 0x00000001bd6552d0 libsystem_kernel.dylib`mach_msg_trap + 8, name = 'AURemoteIO::IOThread'
thread #36: tid = 0x54662, 0x00000001bd679814 libsystem_kernel.dylib`__semwait_signal + 8
thread #37: tid = 0x54663, 0x00000001bd6552d0 libsystem_kernel.dylib`mach_msg_trap + 8, name = 'com.apple.CoreMotion.MotionThread'
thread #38: tid = 0x54664, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Loading.PreloadManager'
```

#### exit
To exit medit, use the `exit` command or `Ctrl-D`.

```
> exit
Bye!
```

## Trouble shooting
### Failed to get reply to handshake packet
If you get the error `/private/var/containers/Bundle/Application/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/hoge.app: error: failed to get reply to handshake packet` and can't communicate properly with iOS device and lldb, launch Xcode and build some app, and it will work.

### Could not connect to lockdownd
If you get the error `Could not connect to lockdownd.` and can't communicate properly with iOS device and ideviceinstaller, launch Xcode and build some app, and it will work.
If this does not solve the problem, please update ideviceinstaller and libimobiledevice to the latest versions using the example of commands in the Requirements section.

### Could not start com.apple.debugserver
This can be fixed by installing the latest unversioned code of libimobiledevice and ideviceinstaller by adding the `--HEAD` option when doing `brew install`.

- Reference: [Could not start com.apple.debugserver! ios 14.1 xcode 12.2 MacOS 10.15.4 iphone12 · Issue #1104 · libimobiledevice/libimobiledevice](https://github.com/libimobiledevice/libimobiledevice/issues/1104)

## License
MIT License
