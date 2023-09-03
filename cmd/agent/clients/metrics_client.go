package clients

import (
	"fmt"
	"net/http"
)

type connect struct {
	server   string
	port     string
	protocol string
}

func SendToGauge(m map[string]string) {
	con := connect{"localhost", "8080", "http"}
	client := http.DefaultClient
	defer client.CloseIdleConnections()
	for k, v := range m {
		reqAddress := getAddressUpdateGauge(con, k, v)
		req, _ := http.NewRequest(http.MethodPost, reqAddress, nil)
		req.Header.Add("Content-Type", "text/plain")
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			fmt.Println(res.Status, " ", k, "=", v)
		} else {
			fmt.Println(err)
		}

	}

}

func SendToCounter(name string, value int64) {
	con := connect{"localhost", "8080", "http"}
	client := http.DefaultClient
	defer client.CloseIdleConnections()
	reqAddress := getAddressUpdateCounter(con, name, value)
	req, _ := http.NewRequest(http.MethodPost, reqAddress, nil)
	req.Header.Add("Content-Type", "text/plain")
	if res, err := client.Do(req); err == nil {
		defer res.Body.Close()
		fmt.Println(res.Status, value)
	} else {
		fmt.Println(err)
	}
}

func getAddressUpdateGauge(con connect, name string, value string) string {
	return fmt.Sprintf("%s://%s:%s/update/gauge/%s/%s", con.protocol, con.server, con.port, name, value)
}

func getAddressUpdateCounter(con connect, name string, value int64) string {
	return fmt.Sprintf("%s://%s:%s/update/counter/%s/%d", con.protocol, con.server, con.port, name, value)
}
