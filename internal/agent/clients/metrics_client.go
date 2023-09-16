package clients

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Connect struct {
	Server   string
	Protocol string
	Client   http.Client
}

func (con *Connect) SendToGauge(m map[string]string) ([]byte, error) {
	defer con.Client.CloseIdleConnections()
	var body []byte
	var res *http.Response
	for k, v := range m {
		reqAddress := getAddressUpdateGauge(con, k, v)
		req, err := http.NewRequest(http.MethodPost, reqAddress, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "[SendToCounter] %v", err)
		}
		req.Header.Add("Content-Type", "text/plain")
		if res, err = con.Client.Do(req); err != nil {
			if res != nil {
				defer res.Body.Close()
				body, _ = io.ReadAll(res.Body)
			}
			return body, err
		}
	}
	return nil, nil
}

func (con *Connect) SendToCounter(name string, value int64) ([]byte, error) {

	defer con.Client.CloseIdleConnections()
	reqAddress := getAddressUpdateCounter(con, name, value)
	req, err := http.NewRequest(http.MethodPost, reqAddress, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "[SendToCounter] %v", err)
	}
	req.Header.Add("Content-Type", "text/plain")
	if res, err := con.Client.Do(req); res != nil && err == nil {
		defer res.Body.Close()
		return io.ReadAll(res.Body)
	} else {
		return nil, err
	}
}

func getAddressUpdateGauge(con *Connect, name string, value string) string {
	return fmt.Sprintf("%s://%s/update/gauge/%s/%s", con.Protocol, con.Server, name, value)
}

func getAddressUpdateCounter(con *Connect, name string, value int64) string {
	return fmt.Sprintf("%s://%s/update/counter/%s/%d", con.Protocol, con.Server, name, value)
}
