package server

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
)

func GetDiskSpace() (map[string]string, error) { // return full and free disk memory on usb and hdd || ssd
	disks := make(map[string]string, 1)
	partitions, err := disk.Partitions(false)

	if err != nil {
		return nil, err
	}

	for _, partition := range partitions {
		disks[partition.Device] = fmt.Sprintf("Device: %v\nMountpoint: %v\nFstype: %v\nOpts: %v\n",
			partition.Device, partition.Mountpoint, partition.Fstype, partition.Opts)

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return nil, err
		}

		disks[partition.Device] += fmt.Sprintf("Total: %v GB\nUsed: %v GB\nFree: %v GB\nPercent: %0.3f %%\n",
			usage.Total/1e9, usage.Used/1e9, usage.Free/1e9, usage.UsedPercent)
	}

	return disks, nil
}
