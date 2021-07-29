//+build darwin

package memory

/*
#include <stdlib.h>
#include <sys/types.h>
#include <mach/mach.h>
#include <mach/mach_vm.h>
#include <mach/thread_info.h>

task_t get_task_for_pid(int pid) {
	task_t task = 0;
	mach_port_t self = mach_task_self();
	task_for_pid(self, pid, &task);
	return task;
}

int write_memory(task_t task, mach_vm_address_t addr, void *d, mach_msg_type_number_t len) {
	kern_return_t kret = mach_vm_write((vm_map_t)task, addr, (vm_offset_t)d, len);
	if (kret != KERN_SUCCESS) return -1;
	return 0;
}

int read_memory(task_t task, mach_vm_address_t addr, void *d, mach_msg_type_number_t len) {
	pointer_t data;
	mach_msg_type_number_t count;
	kern_return_t kret = mach_vm_read((vm_map_t)task, addr, len, &data, &count);
	if (kret != KERN_SUCCESS) return -1;
	memcpy(d, (void *)data, len);
	return count;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func GetTaskForPid(pid int) C.task_t {
	return C.get_task_for_pid(C.int(pid))
}

func WriteMemory(task C.task_t, addr int, data []byte) error {
	var (
		vmData = unsafe.Pointer(&data[0])
		vmAddr = C.mach_vm_address_t(addr)
		length = C.mach_msg_type_number_t(len(data))
	)

	if ret := C.write_memory(task, vmAddr, vmData, length); ret < 0 {
		return fmt.Errorf("could not write memory")
	}
	return nil
}

func ReadMemory(task C.task_t, buf []byte, beginAddr int, endAddr int) error {
	var (
		vmData = unsafe.Pointer(&buf[0])
		vmAddr = C.mach_vm_address_t(beginAddr)
		length = C.mach_msg_type_number_t(endAddr - beginAddr)
	)

	ret := C.read_memory(task, vmAddr, vmData, length)
	if ret < 0 {
		return fmt.Errorf("could not read memory")
	}
	return nil
}
