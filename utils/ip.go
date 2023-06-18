package utils

import (
	"strings"
)

// GetIpToAddress 根据ip地址获取到地址
func GetIpToAddress(ip string) *CityInfo {
	var cityInfo = &CityInfo{}
	var ipStr string
	var infos []string
	p, _ := NewIpdb("./config/qqzeng-ip-utf8.dat")
	ipStr = p.Get(ip)
	infos = strings.Split(ipStr, "|")
	if infos[1] == "保留" {
		cityInfo.Country = "Internal network"
		return cityInfo
	}
	cityInfo.Country = infos[1]
	cityInfo.Province = infos[2]
	cityInfo.City = infos[3]
	return cityInfo
}
