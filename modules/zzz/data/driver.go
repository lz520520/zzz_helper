package data

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"slices"
	"zzz_helper/internal/db/db_zzz"
	"zzz_helper/modules/zzz/models"
)

func CollectDriverDisks(driverType []string, driverIds []string) (*models.DiskCollection, int, error) {
	collection := &models.DiskCollection{
		Disk1: make([]models.DriverDiskStat, 0),
		Disk2: make([]models.DriverDiskStat, 0),
		Disk3: make([]models.DriverDiskStat, 0),
		Disk4: make([]models.DriverDiskStat, 0),
		Disk5: make([]models.DriverDiskStat, 0),
		Disk6: make([]models.DriverDiskStat, 0),
	}
	count := 0

	driverCaches, _ := db_zzz.GetDriverCacheDB().Read(-1, -1)
	for _, cache := range driverCaches {
		if len(driverType) > 0 && !slices.Contains(driverType, cache.Name) {
			continue
		}
		if len(driverIds) > 0 && !slices.Contains(driverIds, cache.ID) {
			continue
		}

		var stat models.DriverDiskStat
		err := yaml.Unmarshal([]byte(cache.Data), &stat)
		if err != nil {
			return nil, 0, err
		}
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
			return nil, count, fmt.Errorf("position not supported: %d", stat.Position)
		}
		count++
	}
	return collection, count, nil

}
