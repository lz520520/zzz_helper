package data

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"zzz_helper/internal/models"
	"zzz_helper/internal/utils/file2"
)

func CollectDriverDisks(name string) (*models.DiskCollection, error) {
	collection := &models.DiskCollection{
		Disk1: make([]models.DriverDiskStat, 0),
		Disk2: make([]models.DriverDiskStat, 0),
		Disk3: make([]models.DriverDiskStat, 0),
		Disk4: make([]models.DriverDiskStat, 0),
		Disk5: make([]models.DriverDiskStat, 0),
		Disk6: make([]models.DriverDiskStat, 0),
	}
	content, err := file2.ReadFileBytes(name)
	if err != nil {
		return nil, err
	}
	m := make(map[string]models.DriverDiskStat)
	err = yaml.Unmarshal(content, &m)
	if err != nil {
		return nil, err
	}
	for _, stat := range m {
		switch stat.Position {
		case 1:
			collection.Disk1 = append(collection.Disk1, stat)
		case 2:
			collection.Disk2 = append(collection.Disk2, stat)
		case 3:
			collection.Disk3 = append(collection.Disk3, stat)
		case 4:
			collection.Disk4 = append(collection.Disk4, stat)
		case 5:
			collection.Disk5 = append(collection.Disk5, stat)
		case 6:
			collection.Disk6 = append(collection.Disk6, stat)
		default:
			return nil, fmt.Errorf("position not supported: %d", stat.Position)
		}
	}
	return collection, nil

}
