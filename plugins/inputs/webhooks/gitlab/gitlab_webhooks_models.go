package gitlab

import (
	"strconv"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
)

type event interface {
	NewMetric() telegraf.Metric
}

// ########################################
// Job Hook event
// ########################################

const jobEvents = "job_event"

type jobEventType struct {
	ObjectKind          string      `json:"object_kind"`
	Ref                 string      `json:"ref"`
	Tag                 bool        `json:"tag"`
	BeforeSha           string      `json:"before_sha"`
	Sha                 string      `json:"sha"`
	RetriesCount        int         `json:"retries_count"`
	BuildID             int         `json:"build_id"`
	BuildName           string      `json:"build_name"`
	BuildStage          string      `json:"build_stage"`
	BuildStatus         string      `json:"build_status"`
	BuildCreatedAt      string      `json:"build_created_at"`
	BuildStartedAt      string      `json:"build_started_at"`  // Can be null
	BuildFinishedAt     string      `json:"build_finished_at"` // Can be null
	BuildDuration       *float64    `json:"build_duration"`    // Can be null
	BuildQueuedDuration *float64    `json:"build_queued_duration"`
	BuildAllowFailure   bool        `json:"build_allow_failure"`
	BuildFailureReason  string      `json:"build_failure_reason"` // Defaults always to unknown_failure
	PipelineID          int         `json:"pipeline_id"`
	Runner              runner      `json:"runner"` // Can be null
	ProjectID           int         `json:"project_id"`
	ProjectName         string      `json:"project_name"` // Currently namespace + name, separated with " / "
	User                user        `json:"user"`
	Commit              commit      `json:"commit"`
	Repository          repository  `json:"repository"`
	Project             project     `json:"project"`
	Environment         environment `json:"environment"` // Can be null
}

func (t *jobEventType) NewMetric() telegraf.Metric {
	tags := map[string]string{
		"job_name":           t.BuildName,
		"job_stage":          t.BuildStage,
		"job_status":         t.BuildStatus,
		"job_failure_reason": t.BuildFailureReason,
		"allow_failure":      strconv.FormatBool(t.BuildAllowFailure),
		"is_tag":             strconv.FormatBool(t.Tag),
		"retries":            strconv.Itoa(t.RetriesCount),
		"project_id":         strconv.Itoa(t.Project.ID),
		"project_name":       t.Project.Name,
		"user_id":            strconv.Itoa(t.User.ID),
		"user":               t.User.Name,
		"runner_id":          strconv.Itoa(t.Runner.ID),
		"runner_description": t.Runner.Description,
		"environment":        t.Environment.Name,
	}
	fields := map[string]interface{}{
		"pipeline_id":         t.PipelineID,
		"ref":                 t.Ref,
		"job_id":              t.BuildID,
		"before_sha":          t.BeforeSha,
		"sha":                 t.Sha,
		"job_created_at":      t.BuildCreatedAt,
		"job_started_at":      t.BuildStartedAt,      // not there
		"job_finished_at":     t.BuildFinishedAt,     // not there
		"job_duration":        t.BuildDuration,       // not there
		"job_queued_duration": t.BuildQueuedDuration, // not there
	}
	n := metric.New(jobEvents, tags, fields, time.Now())
	return n
}

type runner struct {
	ID          int      `json:"id"`
	Description string   `json:"description"`
	RunnerType  string   `json:"runner_type"`
	Active      bool     `json:"active"`
	IsShared    bool     `json:"is_shared"`
	Tags        []string `json:"tags"`
}

type environment struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	//TODO Deployment tier? in pipeline yes but not in job?
}

type project struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"` // Can be null
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"` // Can be null
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int    `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	CiConfigPath      string `json:"ci_config_path"` // Can be null
}

type user struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type repository struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Description     string `json:"description"` // Can be null
	HomePage        string `json:"homepage"`
	GitHTTPURL      string `json:"git_http_url"`
	GitSSHURL       string `json:"git_ssh_url"`
	VisibilityLevel int    `json:"visibility_level"`
}

type commit struct {
	ID          int    `json:"id"`
	Name        string `json:"name"` // Can be null
	Sha         string `json:"sha"`
	Message     string `json:"message"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	AuthorURL   string `json:"author_url"`
	Status      string `json:"status"`
	Duration    int    `json:"duration"`
	StartedAt   string `json:"started_at"`
	FinishedAt  string `json:"finished_at"`
}

// ########################################
// Pipeline hook event
// ########################################

const pipelineEvents = "pipeline_event"

type pipelineEventType struct {
	ObjectKind       string           `json:"object_kind"`
	ObjectAttributes objectAttributes `json:"object_attributes"`
	MergeRequest     mergeRequest     `json:"merge_request"` // Can be null
	User             user             `json:"user"`
	Project          project          `json:"project"`
	Commit           shortCommit      `json:"commit"`
	SourcePipeline   sourcePipeline   `json:"source_pipeline"` // Can be missing
	Builds           []build          `json:"builds"`
}
type shortCommit struct {
	ID        string `json:"id"` // This is the hash
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Author    author `json:"author"`
}

type author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type objectAttributes struct {
	ID         int        `json:"id"`
	Iid        int        `json:"iid"`
	Name       string     `json:"name"` // Can be null
	Ref        string     `json:"ref"`
	Tag        bool       `json:"tag"`
	Sha        string     `json:"sha"`
	BeforeSha  string     `json:"before_Sha"`
	Source     string     `json:"source"`
	Status     string     `json:"status"`
	Stages     []string   `json:"stages"`
	CreatedAt  string     `json:"created_at"`
	FinishedAt string     `json:"finished_at"`
	Duration   int        `json:"duration"`
	Variables  []variable `json:"variables"` // Can be empty
	URL        string     `json:"url"`
}

type variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type mergeRequest struct {
	ID                  int    `json:"id"`
	Iid                 int    `json:"iid"`
	Title               string `json:"title"`
	SourceBranch        string `json:"source_branch"`
	SourceBranchID      int    `json:"source_branch_id"`
	TargetBranch        string `json:"target_branch"`
	TargetBranchID      int    `json:"target_branch_id"`
	State               string `json:"state"`
	MergeStatus         string `json:"merge_status"`
	DetailedMergeStatus string `json:"detailed_merge_status"`
	URL                 string `json:"url"`
}

type sourcePipeline struct {
	Project    shortProject `json:"project"`
	PipelineID int          `json:"pipeline_id"`
	JobID      int          `json:"job_id"`
}

type shortProject struct {
	ID                int    `json:"id"`
	WebURL            string `json:"web_url"`
	PathWithNamespace string `json:"path_with_namespace"`
}

type build struct {
	ID             int           `json:"id"`
	Stage          string        `json:"stage"`
	Name           string        `json:"name"`
	Status         string        `json:"status"`
	CreatedAt      string        `json:"created_at"`
	StartedAt      string        `json:"started_at"`      // Can be null
	FinishedAt     string        `json:"finished_at"`     // Can be null
	Duration       *float64      `json:"duration"`        // Can be null
	QueuedDuration *float64      `json:"queued_duration"` // Can be null
	FailureReason  string        `json:"failure_reason"`  // Can be null
	When           string        `json:"when"`
	Manual         bool          `json:"manual"`
	AllowFailure   bool          `json:"allow_failure"`
	User           user          `json:"user"`
	Runner         runner        `json:"runner"` // Can be null
	ArtifactsFile  artifactsFile `json:"artifacts_file"`
	Environment    environment   `json:"environment"` // Can be null
}

type artifactsFile struct {
	Filename string `json:"filename"`
	Size     *int   `json:"size"`
}

func (t *pipelineEventType) NewMetric() telegraf.Metric {
	tags := map[string]string{
		"pipeline_status": t.ObjectAttributes.Status,
		"project_id":      strconv.Itoa(t.Project.ID),
		"project_name":    t.Project.Name,
		"user_id":         strconv.Itoa(t.User.ID),
		"user":            t.User.Username,
	}
	fields := map[string]interface{}{
		"pipeline_id":     strconv.Itoa(t.ObjectAttributes.ID),
		"pipeline_name":   t.ObjectAttributes.Name,
		"ref":             t.ObjectAttributes.Ref,
		"before_sha":      t.ObjectAttributes.BeforeSha,
		"sha":             t.ObjectAttributes.Sha,
		"job_created_at":  t.ObjectAttributes.CreatedAt,
		"job_finished_at": t.ObjectAttributes.FinishedAt, // not there
		"job_duration":    t.ObjectAttributes.Duration,   // not there
		"url":             t.ObjectAttributes.URL,
	}
	n := metric.New(pipelineEvents, tags, fields, time.Now())
	return n
}
