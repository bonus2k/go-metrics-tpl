package clients

//
//import (
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//	"net/url"
//	"testing"
//)
//
//func TestSendToCounter(t *testing.T) {
//	type args struct {
//		name  string
//		value int64
//		url   string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "check client count aCount = 1",
//			args: args{name: "aCount", value: 1, url: "/update/counter/aCount/1"},
//		},
//		{
//			name: "check client count aCount = 2",
//			args: args{name: "aCount", value: 2, url: "/update/counter/aCount/2"},
//		},
//		{
//			name: "check client count zCount = 999",
//			args: args{name: "zCount", value: 999, url: "/update/counter/zCount/999"},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
//				assert.Equal(t, req.URL.String(), tt.args.url)
//				rw.Write([]byte(`OK`))
//			}))
//			url, _ := url.Parse(server.URL)
//			client := Connect{Server: url.Host, Protocol: "http"}
//			defer server.Close()
//			res, err := client.SendToCounter(tt.args.name, tt.args.value)
//			require.NoError(t, err)
//			assert.Equal(t, string(res), "OK")
//		})
//	}
//}
//
//func TestSendToGauge(t *testing.T) {
//	type args struct {
//		name  string
//		value map[string]string
//		url   string
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "check client count aGauge = 0.999",
//			args: args{name: "aGauge", value: map[string]string{"aGauge": "0.999"}, url: "/update/gauge/aGauge/0.999"},
//		},
//		{
//			name: "check client count aCount = -0.999",
//			args: args{name: "aGauge", value: map[string]string{"aGauge": "-0.999"}, url: "/update/gauge/aGauge/-0.999"},
//		},
//		{
//			name: "check client count zCount = 999",
//			args: args{name: "zGauge", value: map[string]string{"zGauge": "999"}, url: "/update/gauge/zGauge/999"},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
//				assert.Equal(t, req.URL.String(), tt.args.url)
//				rw.Write([]byte(`OK`))
//			}))
//			url, _ := url.Parse(server.URL)
//			client := Connect{Server: url.Host, Protocol: "http"}
//			defer server.Close()
//			_, err := client.SendToGauge(tt.args.value)
//			require.NoError(t, err)
//		})
//	}
//}
