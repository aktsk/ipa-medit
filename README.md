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


### Commands
Here are the commands available in an interactive prompt.

#### find
Search the specified integer on memory.

#### filter
Filter previous search results that match the current search results.

#### patch
Write the specified value on the address found by search.

#### ps
Get information about the target process.

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
thread #12: tid = 0x54622, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Job.Worker 3'
thread #13: tid = 0x54623, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Job.Worker 4'
thread #14: tid = 0x54624, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 0'
thread #15: tid = 0x54625, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 1'
thread #16: tid = 0x54626, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 2'
thread #17: tid = 0x54627, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 3'
thread #18: tid = 0x54628, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 4'
thread #19: tid = 0x54629, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 5'
thread #20: tid = 0x5462a, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 6'
thread #21: tid = 0x5462b, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 7'
thread #22: tid = 0x5462c, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 8'
thread #23: tid = 0x5462d, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 9'
thread #24: tid = 0x5462e, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 10'
thread #25: tid = 0x5462f, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 11'
thread #26: tid = 0x54630, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 12'
thread #27: tid = 0x54631, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 13'
thread #28: tid = 0x54632, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 14'
thread #29: tid = 0x54633, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Background Job.Worker 15'
thread #30: tid = 0x54634, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'BatchDeleteObjects'
thread #31: tid = 0x54635, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Loading.AsyncRead'
thread #32: tid = 0x5463f, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'UnityGfxDeviceWorker'
thread #33: tid = 0x54641, 0x00000001bd6552d0 libsystem_kernel.dylib`mach_msg_trap + 8, name = 'AVAudioSession Notify Thread'
thread #34: tid = 0x54658, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8
thread #35: tid = 0x54659, 0x00000001bd6552d0 libsystem_kernel.dylib`mach_msg_trap + 8, name = 'AURemoteIO::IOThread'
thread #36: tid = 0x54662, 0x00000001bd679814 libsystem_kernel.dylib`__semwait_signal + 8
thread #37: tid = 0x54663, 0x00000001bd6552d0 libsystem_kernel.dylib`mach_msg_trap + 8, name = 'com.apple.CoreMotion.MotionThread'
thread #38: tid = 0x54664, 0x00000001bd65530c libsystem_kernel.dylib`semaphore_wait_trap + 8, name = 'Loading.PreloadManager'
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

#### exit
To exit medit, use the exit command.

```
> exit
Bye!
```

## Trouble shooting
### failed to get reply to handshake packet
If you get the error `/private/var/containers/Bundle/Application/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX/hoge.app: error: failed to get reply to handshake packet` and can't communicate properly with iOS device and lldb, launch Xcode and build some app, and it will work.

### Could not start com.apple.debugserver
This can be fixed by installing the latest unversioned code of libimobiledevice and ideviceinstaller by adding the `--HEAD` option when doing `brew install`.

- Reference: [Could not start com.apple.debugserver! ios 14.1 xcode 12.2 MacOS 10.15.4 iphone12 · Issue #1104 · libimobiledevice/libimobiledevice](https://github.com/libimobiledevice/libimobiledevice/issues/1104)

## License
MIT License
