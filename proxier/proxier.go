package proxier

import "time"

var IPs []string
var ipIndex = 0
var lastIteration time.Time = time.Now()
var numberOfCalls float64 = 0

// var mx sync.Mutex

func SetIPs(list []string) {
	IPs = list
}
func ChangeIP() string {
	if len(IPs) == 0 {
		return ""
	}
	ipIndex++
	if ipIndex > len(IPs)-1 {
		ipIndex = 0
	}
	return IPs[ipIndex]
}
func GetCurrentIP() string {
	if len(IPs) == 0 {
		return ""
	}
	numberOfCalls++
	rate := numberOfCalls / float64(time.Since(lastIteration).Minutes())
	if rate > 30 {
		numberOfCalls = 0
		lastIteration = time.Now()
		return ChangeIP()
	}
	return IPs[ipIndex]
}
