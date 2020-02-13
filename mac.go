package dist

import (
	//	"fmt"
	"time"

	Cpu "github.com/shirou/gopsutil/cpu"
	Host "github.com/shirou/gopsutil/host"
	Mem "github.com/shirou/gopsutil/mem"
)

// 服务器硬件信息
type MachMeta struct {
	MemTotal     uint64
	MemUsed      uint64
	MemAvailable uint64
	MemPercent   float64
	CpuNum       int
	CpuCoreNum   int32
	CpuModelName string
	CpuPercent   []float64
	Hostname     string
	UpTime       uint64
	BootTime     uint64
	OS           string
	Platform     string
	KernelVer    string
	Arch         string
}

func GetMacInfoImmutablePart(mac *MachMeta) {
	cpu, _ := Cpu.Info()
	mac.CpuNum = len(cpu)
	mac.CpuModelName = cpu[0].ModelName
	mac.CpuCoreNum = 0
	for _, core := range cpu {
		mac.CpuCoreNum += core.Cores
	}

	host, _ := Host.Info()
	//	fmt.Printf("Hostname:%v, uptime:%v, boottime:%v,OS:%v Platform:%v KernelVersion:%v KernelArch:%v\n",
	//		host.Hostname, host.Uptime, host.BootTime, host.OS, host.Platform, host.KernelVersion, host.KernelArch)
	mac.Hostname = host.Hostname
	mac.UpTime = host.Uptime
	mac.BootTime = host.BootTime
	mac.OS = host.OS
	mac.Platform = host.Platform
	mac.KernelVer = host.KernelVersion
	mac.Arch = host.KernelArch
}

func GetMachInfoMutablePart(mac *MachMeta) {
	cpup, _ := Cpu.Percent(time.Second, true)
	//	fmt.Println(len(cpup))
	//	fmt.Println(cpup[0])
	for _, p := range cpup {
		mac.CpuPercent = append(mac.CpuPercent, p)
	}

	mem, _ := Mem.VirtualMemory()
	//	fmt.Printf("Total:%vMB, Free:%vMB Used:%vMB Usage%f%%\n",
	//		mem.Total/1024/1024, mem.Available/1024/1024, mem.Used/1024/1024, mem.UsedPercent)
	mac.MemTotal = mem.Total
	mac.MemAvailable = mem.Available
	mac.MemUsed = mem.Used
	mac.MemPercent = mem.UsedPercent
}

func GetMachInfo(mac *MachMeta) {
	GetMacInfoImmutablePart(mac)
	GetMachInfoMutablePart(mac)
}
