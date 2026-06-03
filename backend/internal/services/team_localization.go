package services

type localizedTeam struct {
	Name      string
	Continent string
}

var localizedTeams = map[string]localizedTeam{
	"ALG": {Name: "阿尔及利亚", Continent: "非洲"},
	"ARG": {Name: "阿根廷", Continent: "南美洲"},
	"AUS": {Name: "澳大利亚", Continent: "大洋洲"},
	"AUT": {Name: "奥地利", Continent: "欧洲"},
	"BEL": {Name: "比利时", Continent: "欧洲"},
	"BIH": {Name: "波黑", Continent: "欧洲"},
	"BRA": {Name: "巴西", Continent: "南美洲"},
	"CAN": {Name: "加拿大", Continent: "北美洲"},
	"CHI": {Name: "智利", Continent: "南美洲"},
	"COL": {Name: "哥伦比亚", Continent: "南美洲"},
	"CRC": {Name: "哥斯达黎加", Continent: "北美洲"},
	"CRO": {Name: "克罗地亚", Continent: "欧洲"},
	"CZE": {Name: "捷克", Continent: "欧洲"},
	"DEN": {Name: "丹麦", Continent: "欧洲"},
	"ECU": {Name: "厄瓜多尔", Continent: "南美洲"},
	"EGY": {Name: "埃及", Continent: "非洲"},
	"ENG": {Name: "英格兰", Continent: "欧洲"},
	"ESP": {Name: "西班牙", Continent: "欧洲"},
	"FRA": {Name: "法国", Continent: "欧洲"},
	"GER": {Name: "德国", Continent: "欧洲"},
	"GHA": {Name: "加纳", Continent: "非洲"},
	"GRE": {Name: "希腊", Continent: "欧洲"},
	"HUN": {Name: "匈牙利", Continent: "欧洲"},
	"IRN": {Name: "伊朗", Continent: "亚洲"},
	"ITA": {Name: "意大利", Continent: "欧洲"},
	"JPN": {Name: "日本", Continent: "亚洲"},
	"KOR": {Name: "韩国", Continent: "亚洲"},
	"KSA": {Name: "沙特阿拉伯", Continent: "亚洲"},
	"MAR": {Name: "摩洛哥", Continent: "非洲"},
	"MEX": {Name: "墨西哥", Continent: "北美洲"},
	"NED": {Name: "荷兰", Continent: "欧洲"},
	"NGA": {Name: "尼日利亚", Continent: "非洲"},
	"NOR": {Name: "挪威", Continent: "欧洲"},
	"NZL": {Name: "新西兰", Continent: "大洋洲"},
	"PAR": {Name: "巴拉圭", Continent: "南美洲"},
	"PER": {Name: "秘鲁", Continent: "南美洲"},
	"POL": {Name: "波兰", Continent: "欧洲"},
	"POR": {Name: "葡萄牙", Continent: "欧洲"},
	"QAT": {Name: "卡塔尔", Continent: "亚洲"},
	"RSA": {Name: "南非", Continent: "非洲"},
	"SCO": {Name: "苏格兰", Continent: "欧洲"},
	"SEN": {Name: "塞内加尔", Continent: "非洲"},
	"SRB": {Name: "塞尔维亚", Continent: "欧洲"},
	"SUI": {Name: "瑞士", Continent: "欧洲"},
	"SWE": {Name: "瑞典", Continent: "欧洲"},
	"TUN": {Name: "突尼斯", Continent: "非洲"},
	"TUR": {Name: "土耳其", Continent: "欧洲"},
	"UKR": {Name: "乌克兰", Continent: "欧洲"},
	"URU": {Name: "乌拉圭", Continent: "南美洲"},
	"USA": {Name: "美国", Continent: "北美洲"},
	"WAL": {Name: "威尔士", Continent: "欧洲"},
}

func localizeTeam(code, fallbackName string) (name, continent string) {
	if item, ok := localizedTeams[code]; ok {
		return item.Name, item.Continent
	}
	return fallbackName, ""
}
