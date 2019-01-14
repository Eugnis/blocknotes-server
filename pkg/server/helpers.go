package server

import (
	"bytes"
	"net"
	"net/http"
	"strings"
)

//ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	ipRange{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	ipRange{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	ipRange{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	ipRange{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	ipRange{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	ipRange{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

func isPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func getIPAdress(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return ""
}

// func ValidateUserMiddleware(next http.HandlerFunc, us root.UserService) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		authorizationHeader := req.Header.Get("authorization")
// 		if authorizationHeader != "" {
// 			authToken := strings.Split(authorizationHeader, " ")
// 			// log.Println("walletpk", authToken[0])
// 			// log.Println("token", authToken[1])
// 			ip := getIPAdress(req)

// 			credentials := root.TokenCredentials{WalletPK: authToken[0], Token: authToken[1]}
// 			err, _ := us.TokenLogin(credentials, ip, false)
// 			if err == nil {
// 				next(w, req)
// 			} else {
// 				Error(w, http.StatusForbidden, err.Error())
// 			}

// 		} else {
// 			Error(w, http.StatusForbidden, "Authorization token required")
// 		}
// 	})
// }

// func ValidateAdminMiddleware(next http.HandlerFunc, as root.AdminService) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		authorizationHeader := req.Header.Get("authorization")
// 		if authorizationHeader != "" {
// 			authToken := strings.Split(authorizationHeader, " ")
// 			// log.Println("walletpk", authToken[0])
// 			// log.Println("token", authToken[1])
// 			ip := getIPAdress(req)

// 			credentials := root.TokenCredentials{WalletPK: authToken[0], Token: authToken[1]}
// 			_, err := as.TokenLogin(credentials, ip, false)
// 			if err == nil {
// 				next(w, req)
// 			} else {
// 				Error(w, http.StatusForbidden, err.Error())
// 			}

// 		} else {
// 			Error(w, http.StatusForbidden, "Authorization token required")
// 		}
// 	})
// }
