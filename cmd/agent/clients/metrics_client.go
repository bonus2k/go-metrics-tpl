package clients

import (
	"fmt"
	"io"
	"net/http"
)

type Connect struct {
	Server   string
	Protocol string
}

func (con *Connect) SendToGauge(m map[string]string) ([]byte, error) {
	client := http.DefaultClient
	defer client.CloseIdleConnections()
	var body []byte
	var err error
	var res *http.Response
	for k, v := range m {
		reqAddress := getAddressUpdateGauge(con, k, v)
		req, _ := http.NewRequest(http.MethodPost, reqAddress, nil)
		req.Header.Add("Content-Type", "text/plain")
		if res, err = client.Do(req); err != nil {
			defer res.Body.Close()
			body, _ = io.ReadAll(res.Body)
			return body, err
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
	} else {
		return nil, err
	}
}

func getAddressUpdateGauge(con *Connect, name string, value string) string {
	fmt.Sprintf("%v", con)
	sprintf := fmt.Sprintf("%s://%s/update/gauge/%s/%s", con.Protocol, con.Server, name, value)
	return sprintf
}

func getAddressUpdateCounter(con *Connect, name string, value int64) string {
	return fmt.Sprintf("%s://%s/update/counter/%s/%d", con.Protocol, con.Server, name, value)
}
