package continuous

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	config "github.com/Continuous/config"
	"github.com/hashicorp/memberlist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, port int) (dir string, m *Messenger) {
	conf := memberlist.DefaultLocalConfig()
	conf.BindPort = port
	conf.AdvertisePort = port
	mConfig := config.CustomConfig(conf, true)
	m = NewMessenger(mConfig)
	_, err := m.Join(nil)
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	return
}

type CreateResponse struct {
	UUID string `json:"UUID"`
	code int    `json:"success"`
}

func TestAPI(t *testing.T) {
	port := "8080"
	intPort := 3125
	dir, _ := setup(t, intPort)
	defer os.RemoveAll(dir)
	URL := fmt.Sprintf("http://192.168.5.56:%s/api", port)
	jsonStr := []byte(`{
		"namespace": "Nathan Wong",
		"interval": "00:00:10" 
		}`)
	resp, err := http.Post(URL+"/v1/create", "enconding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var createresp CreateResponse
	if err := json.Unmarshal(body, &createresp); err != nil {
		t.Fatal(err)
	}
	//test get
	resp1, err := http.Get(URL + "/v1/" + "Nathan Wong/" + createresp.UUID)
	if err != nil {
		t.Fatal(err)
	}
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusAccepted, resp1.StatusCode)
	//test delete
	Req, err := http.NewRequest(http.MethodDelete, URL+"/v1/"+"Nathan Wong/"+createresp.UUID, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp2, err := client.Do(Req)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusAccepted, resp2.StatusCode)
}

func TestApiBad(t *testing.T) {
	port := "8080"
	intPort := 3125
	dir, _ := setup(t, intPort)
	defer os.RemoveAll(dir)
	URL := fmt.Sprintf("http://192.168.5.56:%s/api", port)
	jsonStr := []byte(`{
		"namespace": "Nathan Wong"
		}`)
	resp, err := http.Post(URL+"/v1/create", "enconding/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	var createresp CreateResponse
	if err := json.Unmarshal(body, &createresp); err != nil {
		t.Fatal(err)
	}
	//test get
	resp1, err := http.Get(URL + "/v1/" + "Nathan Wong/" + createresp.UUID)
	if err != nil {
		t.Fatal(err)
	}
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp1.StatusCode)
	//test delete
	Req, err := http.NewRequest(http.MethodDelete, URL+"/v1/"+"Nathan Wong/"+createresp.UUID, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp2, err := client.Do(Req)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)

}

func TestIP(t *testing.T) {
	a := GetOutboundIP()
	fmt.Println(a.String())
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
