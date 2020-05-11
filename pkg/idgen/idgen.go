package idgen

import (
	"dmicro/pkg/idgen/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	st := sonyflake.Settings{
		MachineID: getMachineId,
	}
	sf = sonyflake.NewSonyflake(st)
}

/*
1. AWS VPC和 Docker
提供了返回Amazon实例的底层 16位 private IP地址的函数 AmazonEC2MachineID。 它还可以通过检索实例元数据在 Docker 上正确工作。
如果each的netmask instance在/28 和/16. 之间分配一个 netmask，那么地址中的每个EC2实例都有唯一的private IP地址，地址也是惟一的。
在这种常见情况下，你可以使用AmazonEC2MachineID作为 Settings.MachineID。
2. 直接配置 datacenter id + worker id， worker id 就是 docker  的唯一 id，同一个 datacenter 中 worker id 不重复。
3. IP地址有网络号与主机号, 拿出16位主机号即可实现。这可以通过掩码来实现。

func machineID() (uint16, error) {
  ipStr := os.Getenv("MY_IP")
  if len(ipStr) == 0 {
    return 0, errors.New("'MY_IP' environment variable not set")
  }
  ip := net.ParseIP(ipStr)
  if len(ip) < 4 {
    return 0, errors.New("invalid IP")
  }
  return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

*/
func getMachineId() (uint16, error) {
	// TODO
	return 1, nil
}

func Next() (id int64) {
	var i uint64
	if sf != nil {
		i, _ = sf.NextID()
		id = int64(i)
	}
	return
}

func GetOne() int64 {
	return Next()
}

func GetMulti(n int) (ids []int64) {
	for i := 0; i < n; i++ {
		ids = append(ids, Next())
	}
	return
}
