package rest

import (
	"github.com/bonus2k/go-metrics-tpl/internal/middleware/logger"
	m "github.com/bonus2k/go-metrics-tpl/internal/models"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
)

type TrustSubnet struct {
	subnet string
}

func NewTrustSubnet(subnet string) *TrustSubnet {
	return &TrustSubnet{subnet: subnet}
}

// AddRealIp добавляет X-Real-IP заголовок к request
func AddRealIp(c *resty.Client, r *http.Request, ip string) error {
	r.Header.Add(m.KeyXRealIP, ip)
	return nil
}

// CheckRealIp проверяет заголовок X-Real-IP
func (n *TrustSubnet) CheckRealIp(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(n.subnet) == 0 {
			h.ServeHTTP(w, r)
			return
		}
		ip := r.Header.Get(m.KeyXRealIP)
		_, ipNetA, err := net.ParseCIDR(n.subnet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ipB := net.ParseIP(ip)
		logger.Log.Debugf("[CheckRealIp] subnet: %s, ip: %s, contains: %v", ipNetA, ipB, ipNetA.Contains(ipB))
		if !(ipNetA.Contains(ipB)) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
