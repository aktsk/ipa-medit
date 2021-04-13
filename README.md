# ipa-medit
Ipa-medit is a memory search and patch tool for resigned ipa without jailbreak. It was created for mobile game security testing.

## Motivation
Memory modification is the easiest way to cheat in games, it is one of the items to be checked in the security test.
There are also cheat tools that can be used casually like GameGem and iGameGuardian.
However, there were no tools available for un-jailbroken device and CUI.
So I made it as a security testing tool.
Android version is [aktsk/apk-medit](https://github.com/aktsk/apk-medit).

## Demo


## Requirements
- macOS
  - You need to have a valid iOS Development certificate installed
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


## Usage


## Trouble shooting

### failed to get reply to handshake packet
If you get the error `lldb: failed to connect to remote target /private/var/containers/Bundle/Application/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/hoge.app: error: failed to get reply to handshake packet` and can't communicate properly with iOS device and lldb, launch Xcode and build some app, and it will work.

### Could not start com.apple.debugserver
This can be fixed by installing the latest unversioned code of libimobiledevice and ideviceinstaller by adding the `--HEAD` option when doing `brew install`.

- Reference: [Could not start com.apple.debugserver! ios 14.1 xcode 12.2 MacOS 10.15.4 iphone12 · Issue #1104 · libimobiledevice/libimobiledevice](https://github.com/libimobiledevice/libimobiledevice/issues/1104)
