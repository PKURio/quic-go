package node

/*
功能一：在开机时刻发现本机的 IP 地址
*/

import (
	"code.byted.org/videoarch/pcdn_lab_node/pkg/tc"
	"errors"
	"github.com/PKURio/quic-go/log"
	"net"
	"time"
)

const (
	UpdateLocalIPTimeout = time.Second * 10
)

var (
	ErrorTimeoutGettingLocalIP = errors.New("timeout getting local ip")
)

var (
	LocalIP    string
	LocalPort  int
	RemoteIP   string
	RemotePort int
	Delay      tc.Delayer
	Loss       tc.Losser
	Reorder    tc.Reorder
)

// discoverLocalIP 通过查询全部DNS服务器，获取和更新本 Node 的IP地址
// 多个DNS的返回结果先到先得
// 阻塞式调用，无法联通或超时，则不更新IP，并返回错误
func discoverLocalIP() (net.IP, error) {
	dnsAddrs := []string{"114.114.114.114", "8.8.4.4", "223.5.5.5"}
	ip := make(chan *net.UDPAddr, 1)
	for _, dnsAddr := range dnsAddrs {
		go func(dnsAddr string) {
			conn, err := net.Dial("udp", dnsAddr+":53")
			if err != nil {
				log.GetLogger().Warningf("cannot connect to DNS server: %s", dnsAddrs)
				return
			}
			defer conn.Close()
			ip <- conn.LocalAddr().(*net.UDPAddr)
		}(dnsAddr)
	}
	select {
	case r := <-ip:
		return r.IP, nil
	case <-time.After(UpdateLocalIPTimeout):
		log.GetLogger().Warnf(ErrorTimeoutGettingLocalIP.Error())
		return nil, ErrorTimeoutGettingLocalIP
	}
}

func UpdateIPAndPort() {
	ip, err := discoverLocalIP()
	if err != nil {
		log.GetLogger().WithField("ip", ip).Fatalln("failed to get local ip")
		return
	}
	LocalIP = ip.String()
}
