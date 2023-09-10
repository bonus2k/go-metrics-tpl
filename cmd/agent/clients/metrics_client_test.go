package clients

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	url2 "net/url"
	"testing"
)

func TestSendToCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
		url   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check client count aCount = 1",
			args: args{name: "aCount", value: 1, url: "/update/counter/aCount/1"},
		},
		{
			name: "check client count aCount = 2",
			args: args{name: "aCount", value: 2, url: "/update/counter/aCount/2"},
		},
		{
			name: "check client count zCount = 999",
			args: args{name: "zCount", value: 999, url: "/update/counter/zCount/999"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), tt.args.url)
				rw.Write([]byte(`OK`))
			}))

			parse, _ := url2.Parse(server.URL)
			port := parse.Port()
			client := Connect{Server: "127.0.0.1", Port: port, Protocol: "http"}
			defer server.Close()
			res, err := client.SendToCounter(tt.args.name, tt.args.value)
			require.NoError(t, err)
			assert.Equal(t, string(res), "OK")
		})
	}
}

func TestSendToGauge(t *testing.T) {
	type args struct {
		name  string
		value map[string]string
		url   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check client count aGauge = 0.999",
			args: args{name: "aGauge", value: map[string]string{"aGauge": "0.999"}, url: "/update/gauge/aGauge/0.999"},
		},
		{
			name: "check client count aCount = -0.999",
			args: args{name: "aGauge", value: map[string]string{"aGauge": "-0.999"}, url: "/update/gauge/aGauge/-0.999"},
		},
		{
			name: "check client count zCount = 999",
			args: args{name: "zGauge", value: map[string]string{"zGauge": "999"}, url: "/update/gauge/zGauge/999"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), tt.args.url)
				rw.Write([]byte(`OK`))
			}))

			parse, _ := url2.Parse(server.URL)
			port := parse.Port()
			client := Connect{Server: "127.0.0.1", Port: port, Protocol: "http"}
			defer server.Close()
			_, err := client.SendToGauge(tt.args.value)
			require.NoError(t, err)
		})
	}
}

func Test_getAddressUpdateCounter(t *testing.T) {
	type args struct {
		con   *Connect
		name  string
		value int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test#1 address for update counter",
			args: args{
				con:   &Connect{Server: "localhost", Port: "8888", Protocol: "http"},
				name:  "aCount",
				value: 99,
			},
			want: "http://localhost:8888/update/counter/aCount/99",
		},
		{
			name: "test#2 address for update counter",
			args: args{
				con:   &Connect{Server: "localhost", Port: "888", Protocol: "https"},
				name:  "aCount",
				value: 0,
			},
			want: "https://localhost:888/update/counter/aCount/0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAddressUpdateCounter(tt.args.con, tt.args.name, tt.args.value); got != tt.want {
				assert.Equal(t, tt.want, got, "getAddressUpdateCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAddressUpdateGauge(t *testing.T) {
	type args struct {
		con   *Connect
		name  string
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test#1 address for update gauge",
			args: args{
				con:   &Connect{Server: "localhost", Port: "8888", Protocol: "http"},
				name:  "aGauge",
				value: "0.99",
			},
			want: "http://localhost:8888/update/gauge/aGauge/0.99",
		},
		{
			name: "test#2 address for update counter",
			args: args{
				con:   &Connect{Server: "localhost", Port: "888", Protocol: "https"},
				name:  "aGauge",
				value: "-0.88",
			},
			want: "https://localhost:888/update/gauge/aGauge/-0.88",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAddressUpdateGauge(tt.args.con, tt.args.name, tt.args.value); got != tt.want {
				t.Errorf("getAddressUpdateGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}
