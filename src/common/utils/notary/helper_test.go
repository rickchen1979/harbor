package notary

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	notarytest "github.com/vmware/harbor/src/common/utils/notary/test"

	"net/http/httptest"
	"os"
	"path"
	"testing"
)

var endpoint = "10.117.4.142"
var notaryServer *httptest.Server

func TestMain(m *testing.M) {
	notaryServer = notarytest.NewNotaryServer(endpoint)
	defer notaryServer.Close()
	notaryCachePath = "/tmp/notary"
	result := m.Run()
	if result != 0 {
		os.Exit(result)
	}
}

func TestGetTargets(t *testing.T) {
	targets, err := GetTargets(notaryServer.URL, "admin", path.Join(endpoint, "notary-demo/busybox"))
	assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
	assert.Equal(t, 1, len(targets), "")
	assert.Equal(t, "1.0", targets[0].Tag, "")

	targets, err = GetTargets(notaryServer.URL, "admin", path.Join(endpoint, "notary-demo/notexist"))
	assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
	assert.Equal(t, 0, len(targets), "Targets list should be empty for non exist repo.")
}

func TestGetDigestFromTarget(t *testing.T) {
	str := ` {
		      "tag": "1.0",
			  "hashes": {
			        "sha256": "E1lggRW5RZnlZBY4usWu8d36p5u5YFfr9B68jTOs+Kc="
				}
		}`

	var t1 Target
	err := json.Unmarshal([]byte(str), &t1)
	if err != nil {
		panic(err)
	}
	hash2 := make(map[string][]byte)
	t2 := Target{"2.0", hash2}
	d1, err1 := DigestFromTarget(t1)
	assert.Nil(t, err1, "Unexpected error: %v", err1)
	assert.Equal(t, "sha256:1359608115b94599e5641638bac5aef1ddfaa79bb96057ebf41ebc8d33acf8a7", d1, "digest mismatch")
	_, err2 := DigestFromTarget(t2)
	assert.NotNil(t, err2, "")

}
