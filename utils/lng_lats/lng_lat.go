package lng_lats

import (
	"math"
)

// 球面距离公式：https://baike.baidu.com/item/%E7%90%83%E9%9D%A2%E8%B7%9D%E7%A6%BB%E5%85%AC%E5%BC%8F/5374455?fr=aladdin ;
// GetDistance 计算地理距离，依次为两个坐标的纬度、经度、单位（默认：英里，K => 公里(Km)，N => 海里）
func GetDistance(lng1 float64, lat1 float64, lng2 float64, lat2 float64, unit ...string) float64 {
	radLat1 := math.Pi * lat1 / 180
	radLat2 := math.Pi * lat2 / 180
	theta := lng1 - lng2
	radTheta := math.Pi * theta / 180

	dist := math.Sin(radLat1)*math.Sin(radLat2) + math.Cos(radLat1)*math.Cos(radLat2)*math.Cos(radTheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}
