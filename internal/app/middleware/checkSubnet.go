package middleware

import (
	"log"
	"net"
	"net/http"
)

// CheckTrustedSubnet проверяет IP-адрес клиента на вхождение в доверенную подсеть.
func CheckTrustedSubnet(next http.Handler, subnet string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if subnet == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		ipStr := r.Header.Get("X-Real-IP")
		ip := net.ParseIP(ipStr)

		if !ipNet.Contains(ip) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
