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
		cityInfo.CountryName = "Internal network"
		return cityInfo
	}
	cityInfo.CountryName = infos[1]
	cityInfo.RegionName = infos[2]
	cityInfo.CityName = infos[3]
	return cityInfo
}
