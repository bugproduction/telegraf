package gitlab

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

var testDataPath = "testdata/"

var testJobHookHeader = "Job Hook"
var testPipelineHookHeader = "Pipeline Hook"

func readFile(t *testing.T, filePath string) string {
	data, err := os.ReadFile(filePath)
	require.NoErrorf(t, err, "could not read from file %s", filePath)
	return string(data)
}

func GitlabWebhookRequest(t *testing.T, input string, xGitlabEvent string, measurement string) {
	var acc testutil.Accumulator
	gl := &Webhook{Path: "/gitlab", acc: &acc, log: testutil.Logger{}}
	jsonString := readFile(t, input)
	req, err := http.NewRequest("POST", "/gitlab", strings.NewReader(jsonString))
	require.NoError(t, err)
	req.Header.Add("X-Gitlab-Event", xGitlabEvent)
	w := httptest.NewRecorder()
	gl.eventHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("POST returned HTTP status code %v.\nExpected %v", w.Code, http.StatusOK)
	}
	acc.HasMeasurement(measurement)
}

// ########################################
// Job hook event
// ########################################

func TestProjectJobHook(t *testing.T) {
	GitlabWebhookRequest(t, testDataPath+"job_hook.json", testJobHookHeader, "job_event")
}

// ########################################
// Pipeline hook event
// ########################################

func TestProjectPipelineHook(t *testing.T) {
	GitlabWebhookRequest(t, testDataPath+"pipeline_hook.json", testPipelineHookHeader, "pipeline_event")
}
