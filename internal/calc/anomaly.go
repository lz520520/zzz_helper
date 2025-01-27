package calc

import "zzz_helper/internal/data"

func AnomalyImprove(anomaly float64, count int) float64 {
	out := anomaly + data.BaseDriverDiskSubStat.AnomalyProficiency*float64(count)
	return out/anomaly - 1
}

func AnomalyExpected(improve float64, count int) (anomaly float64) {
	improve = (data.BaseDriverDiskSubStat.AnomalyProficiency * float64(count)) / anomaly
	anomaly = (data.BaseDriverDiskSubStat.AnomalyProficiency * float64(count)) / improve
	return
}
