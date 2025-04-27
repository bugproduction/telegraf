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
var testPassword = "thisiscorrect"
var testGitlab_webhooks = "gitlab_webhook"
var testJobHookHeader = "Job Hook"
var testPipelineHookHeader = "Pipeline Hook"
var testMergeRequestHookHeader = "Merge Request Hook"

func readFile(t *testing.T, filePath string) string {
	data, err := os.ReadFile(filePath)
	require.NoErrorf(t, err, "could not read from file %s", filePath)
	return string(data)
}

func GitlabWebhookRequest(t *testing.T, input string, xGitlabEvent string) {
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
	acc.HasMeasurement(testGitlab_webhooks)
}

func GitlabWebhookRequestToken(t *testing.T, input string, xGitlabEvent string, token string) int {
	var acc testutil.Accumulator
	gl := &Webhook{Path: "/gitlab", acc: &acc, log: testutil.Logger{}, Secret: testPassword}
	jsonString := readFile(t, input)
	req, err := http.NewRequest("POST", "/gitlab", strings.NewReader(jsonString))
	require.NoError(t, err)
	req.Header.Add("X-Gitlab-Event", xGitlabEvent)
	req.Header.Add("X-Gitlab-Token", token)
	w := httptest.NewRecorder()
	gl.eventHandler(w, req)
	return w.Code
}

// ########################################
// Job hook event
// ########################################

func TestProjectJobHook(t *testing.T) {
	GitlabWebhookRequest(t, testDataPath+"job_hook.json", testJobHookHeader)
}

func TestProjectJobHookCorrectToken(t *testing.T) {
	code := GitlabWebhookRequestToken(t, testDataPath+"job_hook.json", testJobHookHeader, testPassword)
	if code != http.StatusOK {
		t.Errorf("POST with right password returned HTTP status code %v.\nExpected %v", code, http.StatusOK)
	}
}

func TestProjectJobHookWrongToken(t *testing.T) {
	code := GitlabWebhookRequestToken(t, testDataPath+"job_hook.json", testJobHookHeader, "thisiswrong")
	if code == http.StatusOK {
		t.Errorf("POST with wrong password returned HTTP status code %v.\nExpected failure", code)
	}
}

// ########################################
// Pipeline hook event
// ########################################

func TestProjectPipelineHook(t *testing.T) {
	GitlabWebhookRequest(t, testDataPath+"pipeline_hook.json", testPipelineHookHeader)
}

// ########################################
// Pipeline hook event
// ########################################

func TestMergeRequestHook(t *testing.T) {
	GitlabWebhookRequest(t, testDataPath+"merge_req_hook.json", testMergeRequestHookHeader)
}
