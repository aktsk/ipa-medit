#!/usr/bin/env python3
# coding: UTF-8

import multiprocessing
import os
import signal
import struct
import sys
import time

platform, local_bin, device_bin_or_pid = sys.argv[1], sys.argv[2], sys.argv[3]

env = []
for k, v in os.environ.items():
    env.append(k + '=' + v)

sys.path.append('/Applications/Xcode.app/Contents/SharedFrameworks/LLDB.framework/Resources/Python3')

try:
    import lldb
except ModuleNotFoundError:
    print('LLDB library not found... Please install Xcode.')
    sys.exit()


def signal_handler(signal, frame):
    process.Signal(signal)


def run_program(target):
    # Forward SIGQUIT to the program.
    signal.signal(signal.SIGQUIT, signal_handler)
    # Tell the Go driver that the program is running and should not be retried.
    process = target.GetProcess()
    process.Continue()
    print('lldb: running program....')


def run_prompt(target, listener, debugger):
    addr_cache = []
    process = target.GetProcess()
    while True:
        input_text = input('\033[1m\033[34m> \033[36m')
        print('\033[0m', end='')  # reset color
        input_text_list = input_text.split(' ')
        cmd = input_text_list[0]
        state = process.GetState()
        if cmd == 'attach':
            if state == lldb.eStateStopped:
                print('Already attached...')
            else:
                attach(target, listener)
        elif cmd == 'detach':
            if state == lldb.eStateRunning:
                print('Already detached.')
            else:
                detach(target, listener)
        elif cmd == 'ps':
            info(target)
        elif cmd == 'find':
            if len(input_text_list) < 2:
                print('Target value cannot be specified.')
                continue
            if state != lldb.eStateStopped:
                attach(target, listener)
            value = int(input_text_list[1])
            addr_cache = start_search_process(target, value)
        elif cmd == 'filter':
            if len(input_text_list) < 2:
                print('Target value cannot be specified.')
                continue
            if state != lldb.eStateStopped:
                attach(target, listener)
            value = int(input_text_list[1])
            addr_cache = filter_addr(process, value, addr_cache)
        elif cmd == 'patch':
            if len(input_text_list) < 2:
                print('Target value cannot be specified.')
                continue
            if state != lldb.eStateStopped:
                attach(target, listener)
            value = int(input_text_list[1])
            patch(process, value, addr_cache)
        elif cmd == 'exit':
            print('Bye!')
            lldb_exit(target, debugger)
        else:
            print('Command not found...')


def info(target):
    process = target.GetProcess()
    print(process)

    state = process.GetState()
    if state == lldb.eStateStopped:
        print('State: Stopped')
    elif process.GetState() == lldb.eStateRunning:
        print('State: Running')

    for i in range(process.GetNumThreads()):
        print(process.GetThreadAtIndex(i))


def attach(target, listener):
    process = target.GetProcess()
    process.SendAsyncInterrupt()
    while listener.WaitForEvent(2, event):
        pass
    if process.GetState() == lldb.eStateStopped:
        print('Success to halt process')
    else:
        print('Failed to halt process')


def detach(target, listener):
    process = target.GetProcess()
    process.Continue()
    while listener.WaitForEvent(2, event):
        pass
    if process.GetState() == lldb.eStateRunning:
        print('Success to continue process')
    else:
        print('Failed to continue process')


def lldb_exit(target, debugger):
    process = target.GetProcess()
    process.Kill()
    debugger.Terminate()
    sys.exit()


def start_search_process(target, pattern):
    start = time.time()
    process = target.GetProcess()
    memory_regions = process.GetMemoryRegions()
    memory_regions_size = memory_regions.GetSize()

    manager = multiprocessing.Manager()
    manager_list = manager.list()

    int_pattern, int_type = int_to_byte(pattern, None)
    int_pattern_lengh = len(int_pattern)
    int_hash = pow(16777619, int_pattern_lengh - 1) % 999999937

    string_pattern, string_type = int_to_byte(pattern, 'string')
    string_pattern_lengh = len(int_pattern)
    string_hash = pow(16777619, int_pattern_lengh - 1) % 999999937

    search_jobs = []
    for i in range(memory_regions_size):
        memory_region_info = lldb.SBMemoryRegionInfo()
        success = memory_regions.GetMemoryRegionAtIndex(i, memory_region_info)
        if success:
            begin_addr, end_addr = parse_memory_region(memory_region_info)
            if (begin_addr is not None) and (end_addr is not None):
                if begin_addr < 0x1d0000000:
                    print('Scanning: 0x{:016x}-0x{:016x}'.format(begin_addr, end_addr))
                    err = lldb.SBError()
                    memory_length = end_addr - begin_addr
                    memory_bytes = process.ReadMemory(begin_addr, memory_length, err)
                    if err.Success():
                        search_int_process = multiprocessing.Process(target=find_bytes_memory_region, args=(
                            memory_bytes, begin_addr, memory_length, int_pattern, int_pattern_lengh, int_hash, int_type,
                            manager_list))
                        search_int_process.start()
                        search_jobs.append(search_int_process)
                        search_string_process = multiprocessing.Process(target=find_bytes_memory_region, args=(
                            memory_bytes, begin_addr, memory_length, string_pattern, string_pattern_lengh, string_hash, string_type,
                            manager_list))
                        search_string_process.start()
                        search_jobs.append(search_string_process)
                    if len(manager_list) > 500000:
                        print('Too many addresses with target data found...')
                        [j.join() for j in search_jobs]
                        return manager_list

    [j.join() for j in search_jobs]
    elapsed_time = time.time() - start
    print('elapsed_time:{0}'.format(elapsed_time) + '[sec]')
    print('Found: {0}!!'.format(len(manager_list)))
    if len(manager_list) < 10:
        for addr, _ in manager_list:
            print('Address: {0}'.format(hex(addr)))
    return manager_list


def parse_memory_region(memory_region_info):
    begin_addr = memory_region_info.GetRegionBase()
    end_addr = memory_region_info.GetRegionEnd()
    if memory_region_info.IsReadable() and memory_region_info.IsWritable() and memory_region_info.IsMapped():
        return begin_addr, end_addr
    else:
        return None, None


def find_bytes_memory_region(memory_bytes, base_addr, memory_length, search_pattern, search_pattern_length, search_hash,
                             search_type, manager_list):
    memory_region_index = find_all_bytes_by_rabin_karp(memory_bytes, memory_length, search_pattern,
                                                       search_pattern_length, search_hash)
    search_result = list(map(lambda x: (base_addr + x, search_type), memory_region_index))
    manager_list.extend(search_result)


def find_all_bytes_by_rabin_karp(text, n, pattern, m, h):
    p = 0
    t = 0
    d = 16777619
    q = 999999937
    result = []
    for i in range(m):  # preprocessing
        p = (d * p + pattern[i]) % q
        t = (d * t + text[i]) % q
    for i in range(n - m + 1):  # note the +1
        if p == t:  # check character by character
            match = True
            for j in range(m):
                if pattern[j] != text[i + j]:
                    match = False
                    break
            if match:
                result.append(i)
        if i < n - m:
            t = (t - h * text[i]) % q  # remove letter s
            t = (t * d + text[i + m]) % q  # add letter s+m
            t = (t + q) % q  # make sure that t >= 0
    return result


def int_to_byte(integer, int_type):
    if int_type is None:
        if integer <= 255:
            return struct.pack('1B', integer), 'uint8'
        elif 255 < integer <= 65535:
            return struct.pack('1H', integer), 'uint16'
        elif 65535 < integer <= 4294967295:
            return struct.pack('1I', integer), 'uint32'
        else:
            return struct.pack('1Q', integer), 'uint64'
    elif int_type == 'uint8':
        return struct.pack('1B', integer), 'uint8'
    elif int_type == 'uint16':
        return struct.pack('1H', integer), 'uint16'
    elif int_type == 'uint32':
        return struct.pack('1I', integer), 'uint32'
    elif int_type == 'uint64':
        return struct.pack('1Q', integer), 'uint64'
    elif int_type == 'string':
        return bytes(str(integer), encoding="UTF-8"), 'string'


def patch(process, search_pattern, addr_cache):
    err = lldb.SBError()
    for addr, search_type in addr_cache:
        target_bytes, _ = int_to_byte(search_pattern, search_type)
        result = process.WriteMemory(addr, target_bytes, err)
        if not err.Success():
            print('Failed to write memory')
        else:
            print('Successfully patched!')


def filter_addr(process, search_pattern, addr_cache):
    result = []
    for begin_addr, search_type in addr_cache:
        target_bytes, _ = int_to_byte(search_pattern, search_type)
        err = lldb.SBError()
        memory_length = len(target_bytes)
        memory_bytes = process.ReadMemory(begin_addr, memory_length, err)
        if target_bytes == memory_bytes:
            result.append((begin_addr, search_type))
    print('Found: {0}!!'.format(len(result)))
    if len(result) < 10:
        for addr, _ in result:
            print('Address: {0}'.format(hex(addr)))
    return result


if __name__ == '__main__':
    debugger = lldb.SBDebugger.Create()
    debugger.SetAsync(True)
    debugger.SkipLLDBInitFiles(True)

    err = lldb.SBError()
    event = lldb.SBEvent()

    target = debugger.CreateTarget(local_bin, None, platform, True, err)
    if not target.IsValid() or not err.Success():
        print('lldb: failed to setup up target: %s' % (err))
        sys.exit(1)

    listener = debugger.GetListener()

    if platform == 'remote-ios':
        target.modules[0].SetPlatformFileSpec(lldb.SBFileSpec(device_bin_or_pid))
        process = target.ConnectRemote(listener, 'connect://localhost:3222', None, err)
    else:
        process = target.AttachToProcessWithID(listener, int(device_bin_or_pid), err)

    if not err.Success():
        print('lldb: failed to connect to remote target %s: %s' % (device_bin_or_pid, err))
        sys.exit(1)

    # Don't stop on signals.
    sigs = process.GetUnixSignals()
    for i in range(0, sigs.GetNumSignals()):
        sig = sigs.GetSignalAtIndex(i)
        sigs.SetShouldStop(sig, False)
        sigs.SetShouldNotify(sig, False)

    if platform != 'remote-ios':
        run_program(target)

    while True:
        if (not listener.WaitForEvent(1, event)) or (not lldb.SBProcess.EventIsProcessEvent(event)):
            run_prompt(target, listener, debugger)
            continue

        state = process.GetStateFromEvent(event)
        if state == lldb.eStateConnected:
            if platform == 'remote-ios':
                process.RemoteLaunch([], env, None, None, None, None, 0, False, err)
                if not err.Success():
                    print('lldb: failed to launch remote process: %s' % (err))
                    process.Kill()
                    debugger.Terminate()
                    sys.exit(1)
            run_program(target)

    exitStatus = process.GetExitStatus()
    exitDesc = process.GetExitDescription()
    process.Kill()
    debugger.Terminate()
    sys.exit(exitStatus)
