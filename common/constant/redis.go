package constant

// redis 不同类型 key 的定义
const (
	RedisKeyConnid       = "UMS:CONNID:%d:%d:%d"             // "UMS:CONNID:appid:uid:plat"
	RedisKeyToken        = "PASSPORT:TOKEN:%d:%d:%d"         // "PASSPORT:TOKEN:appid:uid:plat"
	RedisKeyRefreshToken = "PASSPORT:REFRESH_TOKEN:%d:%d:%d" // "PASSPORT:REFRESH_TOKEN:appid:uid:plat"
)
