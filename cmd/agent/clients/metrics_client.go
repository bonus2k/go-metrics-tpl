package clients

import (
	"fmt"
	"io"
	"net/http"
)

type Connect struct {
	Server   string
	Port     string
	Protocol string
}

func (con *Connect) SendToGauge(m map[string]string) ([]byte, error) {
	client := http.DefaultClient
	defer client.CloseIdleConnections()
	for k, v := range m {
		reqAddress := getAddressUpdateGauge(con, k, v)
		req, _ := http.NewRequest(http.MethodPost, reqAddress, nil)
		req.Header.Add("Content-Type", "text/plain")
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			return io.ReadAll(res.Body)
			fmt.Println(res.Status, " ", k, "=", v)
		} else {
			fmt.Println(err)
			return nil, err
		}
	}
	return nil, nil
}

func (con *Connect) SendToCounter(name string, value int64) ([]byte, error) {

	client := http.DefaultClient
	defer client.CloseIdleConnections()
	reqAddress := getAddressUpdateCounter(con, name, value)
	req, _ := http.NewRequest(http.MethodPost, reqAddress, nil)
	req.Header.Add("Content-Type", "text/plain")
	if res, err := client.Do(req); err == nil {
		defer res.Body.Close()
		return io.ReadAll(res.Body)
		fmt.Println(res.Status, value)
	} else {
		return nil, err
		fmt.Println(err)
	}
	return nil, nil
}

func getAddressUpdateGauge(con *Connect, name string, value string) string {
	return fmt.Sprintf("%s://%s:%s/update/gauge/%s/%s", con.Protocol, con.Server, con.Port, name, value)
}

func getAddressUpdateCounter(con *Connect, name string, value int64) string {
	return fmt.Sprintf("%s://%s:%s/update/counter/%s/%d", con.Protocol, con.Server, con.Port, name, value)
}
