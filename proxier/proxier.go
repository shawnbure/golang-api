package proxier

var IPs []string
var ipIndex = 0

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
	return IPs[ipIndex]
}
