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

const gitlabWebhooks = "gitlab_webhooks"
const jobEvents = "job_event"

func gitLabTimeToTime(timestamp string) time.Time {
	const timeformat = "2006-01-02 15:04:05 MST"
	val, _ := time.Parse(timeformat, timestamp)
	return val
}

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
		"event":              jobEvents,
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
		"user_name":          t.User.Name,
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
		"job_created_at":      gitLabTimeToTime(t.BuildCreatedAt),
		"job_started_at":      gitLabTimeToTime(t.BuildStartedAt),
		"job_finished_at":     gitLabTimeToTime(t.BuildFinishedAt),
		"job_duration":        t.BuildDuration,
		"job_queued_duration": t.BuildQueuedDuration,
	}
	n := metric.New(gitlabWebhooks, tags, fields, time.Now())
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
	BeforeSha  string     `json:"before_sha"`
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
		"event":           pipelineEvents,
		"pipeline_status": t.ObjectAttributes.Status,
		"project_id":      strconv.Itoa(t.Project.ID),
		"project_name":    t.Project.Name,
		"user_id":         strconv.Itoa(t.User.ID),
		"user_name":       t.User.Username,
	}
	fields := map[string]interface{}{
		"pipeline_id":          t.ObjectAttributes.ID,
		"pipeline_name":        t.ObjectAttributes.Name,
		"ref":                  t.ObjectAttributes.Ref,
		"before_sha":           t.ObjectAttributes.BeforeSha,
		"sha":                  t.ObjectAttributes.Sha,
		"pipeline_created_at":  gitLabTimeToTime(t.ObjectAttributes.CreatedAt),
		"pipeline_finished_at": gitLabTimeToTime(t.ObjectAttributes.FinishedAt),
		"pipeline_duration":    t.ObjectAttributes.Duration,
		"url":                  t.ObjectAttributes.URL,
	}
	n := metric.New(gitlabWebhooks, tags, fields, time.Now())
	return n
}

// ########################################
// Merge request event
// ########################################

const mergeRequestEvents = "merge_request_event"

type label struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	ProjectID   int    `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Template    bool   `json:"template"`
	Description string `json:"description"`
	Type        string `json:"type"`
	GroupID     int    `json:"group_id"`
}

type mrObjectAttributes struct {
	ID                          int         `json:"id"`
	Iid                         int         `json:"iid"`
	TargetBranch                string      `json:"target_branch"`
	SourceBranch                string      `json:"source_branch"`
	SourceProjectID             int         `json:"source_project_id"`
	AuthorID                    int         `json:"author_id"`
	AssigneeIDs                 []int       `json:"assignee_ids"`
	AssigneeID                  int         `json:"assignee_id"` // Can be null
	ReviewerIDs                 []int       `json:"reviewer_ids"`
	Title                       string      `json:"title"`
	CreatedAt                   string      `json:"created_at"`
	UpdatedAt                   string      `json:"updated_at"`
	LastEditedAt                string      `json:"last_edited_at"`    // Can be null
	LastEditedByID              int         `json:"last_edited_by_id"` // Can be null
	MilestoneID                 int         `json:"milestone_id"`      // Can be null
	StateID                     int         `json:"state_id"`
	State                       string      `json:"state"`
	BlockingDiscussionsResolved bool        `json:"blocking_discussions_resolved"`
	WorkInProgress              bool        `json:"work_in_progress"`
	Draft                       bool        `json:"draft"`
	FirstContribution           bool        `json:"first_contribution"`
	MergeStatus                 string      `json:"merge_status"`
	TargetProjectID             int         `json:"target_project_id"`
	Description                 string      `json:"description"`
	PreparedAt                  string      `json:"prepared_at"`
	TotalTimeSpent              int         `json:"total_time_spent"`
	TimeChange                  int         `json:"time_change"`
	HumanTotalTimeSpent         string      `json:"human_total_time_spent"` // Can be null
	HumanTimeChange             string      `json:"human_time_change"`      // Can be null
	HumanTimeEstimate           string      `json:"human_time_estimate"`    // Can be null
	URL                         string      `json:"url"`
	Source                      project     `json:"source"`
	Target                      project     `json:"target"`
	LastCommit                  shortCommit `json:"last_commit"`
	Labels                      []label     `json:"labels"`
	Action                      string      `json:"action"`
	DetailedMergeStatus         string      `json:"detailed_merge_status"`
}

type updatedByID struct {
	Previous int `json:"previous"`
	Current  int `json:"current"`
}

type draft struct {
	Previous bool `json:"previous"`
	Current  bool `json:"current"`
}

type updatedAt struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type labels struct {
	Previous []label `json:"previous"`
	Current  []label `json:"current"`
}

type changes struct {
	UpdatedByID    updatedByID `json:"updated_by_id"`
	Draft          draft       `json:"draft"`
	UpdatedAt      updatedAt   `json:"updated_at"`
	Labels         labels      `json:"labels"`
	LastEditedAt   updatedAt   `json:"last_edited_at"`
	LastEditedByID updatedByID `json:"last_edited_by_id"`
}

type mergeRequestEventType struct {
	ObjectKind       string             `json:"object_kind"`
	EventType        string             `json:"event_type"`
	User             user               `json:"user"`
	Project          project            `json:"project"`
	Repository       repository         `json:"repository"`
	ObjectAttributes mrObjectAttributes `json:"object_attributes"`
	Labels           []label            `json:"labels"`
	Changes          changes            `json:"changes"`
	Assignees        []user             `json:"assignees"`
	Reviewers        []user             `json:"reviewers"`
}

func (t *mergeRequestEventType) NewMetric() telegraf.Metric {
	tags := map[string]string{
		"event":                         mergeRequestEvents,
		"project_id":                    strconv.Itoa(t.Project.ID),
		"project_name":                  t.Project.Name,
		"user_id":                       strconv.Itoa(t.User.ID),
		"user_name":                     t.User.Name,
		"target_branch":                 t.ObjectAttributes.TargetBranch,
		"author_id":                     strconv.Itoa(t.ObjectAttributes.AuthorID),
		"blocking_discussions_resolved": strconv.FormatBool(t.ObjectAttributes.BlockingDiscussionsResolved),
		"work_in_progress":              strconv.FormatBool(t.ObjectAttributes.WorkInProgress),
		"draft":                         strconv.FormatBool(t.ObjectAttributes.Draft),
		"detailed_merge_status":         t.ObjectAttributes.DetailedMergeStatus,
	}
	fields := map[string]interface{}{
		"mr_id":         t.ObjectAttributes.ID,
		"title":         t.ObjectAttributes.Title,
		"source_branch": t.ObjectAttributes.SourceBranch,
		"description":   t.ObjectAttributes.Description,
		"created_at":    gitLabTimeToTime(t.ObjectAttributes.CreatedAt),
		"updated_at":    gitLabTimeToTime(t.ObjectAttributes.UpdatedAt),
		"url":           t.ObjectAttributes.URL,
	}
	n := metric.New(gitlabWebhooks, tags, fields, time.Now())
	return n
}
