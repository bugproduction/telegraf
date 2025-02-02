package gitlab

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/influxdata/telegraf"
)

type Webhook struct {
	Path   string
	Secret string
	acc    telegraf.Accumulator
	log    telegraf.Logger
}

func (gl *Webhook) Register(router *mux.Router, acc telegraf.Accumulator, log telegraf.Logger) {
	router.HandleFunc(gl.Path, gl.eventHandler).Methods("POST")
	gl.log = log
	gl.log.Infof("Started gitlab_webhooks on %s", gl.Path)
	gl.acc = acc
}

func (gl *Webhook) eventHandler(w http.ResponseWriter, r *http.Request) {
	const gitlabHeaderEvent = "X-Gitlab-Event"
	const gitlabHeaderToken = "X-Gitlab-Token"

	var headers = map[string]event{
		"Job Hook":      &jobEventType{},
		"Pipeline Hook": &pipelineEventType{},
	}

	defer r.Body.Close()

	eventType := r.Header.Get(gitlabHeaderEvent)
	ev := headers[eventType]
	if ev == nil {
		gl.log.Infof("No %s found in headers.", gitlabHeaderEvent)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if gl.Secret != "" && gl.Secret != r.Header.Get(gitlabHeaderToken) {
		gl.log.Error("Fail to check the gitlab webhook secret token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)

	if err != nil {
		gl.log.Info("Failed to read the request body")
		gl.log.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	gl.log.Debugf("New %s event received", eventType)

	e, err := generateEvent(data, ev)

	if err != nil || e == nil {
		gl.log.Errorf("Error parsing the data \"%s\"", err)
		gl.log.Error(string(data[:]))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p := e.NewMetric()
	gl.acc.AddFields(p.Name(), p.Fields(), p.Tags(), p.Time())
	w.WriteHeader(http.StatusOK)
}

func generateEvent(data []byte, event event) (event, error) {
	err := json.Unmarshal(data, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}
