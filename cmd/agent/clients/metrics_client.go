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
	for k, v := range m {
		reqAddress := getAddressUpdateGauge(con, k, v)
		req, _ := http.NewRequest(http.MethodPost, reqAddress, nil)
		req.Header.Add("Content-Type", "text/plain")
		res, err := client.Do(req)
		fmt.Println(res.Status, err)
	}

}

func SendToCounter(name string, value int64) {
	con := connect{"localhost", "8080", "http"}
	client := http.DefaultClient
	reqAddress := getAddressUpdateCounter(con, name, value)
	req, _ := http.NewRequest(http.MethodPost, reqAddress, nil)
	req.Header.Add("Content-Type", "text/plain")
	res, err := client.Do(req)
	fmt.Println(res.Status, err)
}

func getAddressUpdateGauge(con connect, name string, value string) string {
	return fmt.Sprintf("%s://%s:%s/update/gauge/%s/%s", con.protocol, con.server, con.port, name, value)
}

func getAddressUpdateCounter(con connect, name string, value int64) string {
	return fmt.Sprintf("%s://%s:%s/update/counter/%s/%d", con.protocol, con.server, con.port, name, value)
}
