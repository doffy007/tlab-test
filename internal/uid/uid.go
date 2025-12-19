package uid

import (
	"errors"
	"net"
	"time"
)

var ErrAssignNewID = errors.New("failed to assign new id")

func New() (uint64, error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	now := time.Now().UnixMilli()
	if sf.timestamp == now {
		sf.sequence = (sf.sequence + 1) & sequenceMax
		if sf.sequence == 0 {
			for now <= sf.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		sf.sequence = 0
	}

	sf.timestamp = now

	id := ((now - epoch) << timeShift) |
		(sf.nodeID << nodeShift) |
		sf.sequence

	return uint64(id), nil
}

func getNodeIDFromIP(ipStr string) (int64, error) {
	if ipStr == "" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return 0, err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				ipStr = ipnet.IP.String()
				break
			}
		}
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, errors.New("invalid IP")
	}

	return int64(ip[len(ip)-1]) & nodeMax, nil
}
