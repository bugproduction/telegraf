# GitLab webhooks

To receive infromation from GitLab projects or groups using webhooks, configure the GitLab webhooks for the wanted projects or groups in GitLab in the group or project settings - Webhooks -page (for example `https://<gitlab-url>/<group>/<project>/-/hooks`). Add a new webhook by clicking the "Add new webhook" button.

Set the webhook target URL to your telegraf address, `https://<telegraf>:1619/gitlab` and if wanted, add a secret token that is used to verify the origin of the webhook requests. Add a name and description for the webhook configuration.

The "Job events", "Pipeline events", and "Merge request events" triggers are supported. You can enable one or more of them in a single webhook configuration.


## Metrics

The following events from GitLab are creating metrics with the following tags and fields. The format for the tags and fields is following:

```md
Tags:
- 'tagKey' = `tagValue` type

Fields
- 'fieldKey' = `fieldValue` type
```

The events will be stored in `gitlab_webhooks` measurement.

### [Pipeline events](https://docs.gitlab.com/user/project/integrations/webhook_events/#pipeline-events)

Tags:
- 'event' = `pipeline_event` string
- 'pipeline_status' = `object_attributes.status` string
- 'pipeline_source' = `object_attributes.Source` string
- 'project_id' =  `project.id` string
- 'user_id' = `user.id` string

Fields:
- 'pipeline_id' = `object_attributes.id` int
- 'pipeline_name' = `object_attributes.name` string
- 'ref' = `object_attributes.ref` string
- 'before_sha' = `object_attributes.before_sha` string
- 'sha' = `object_attributes.sha` string
- 'pipeline_created_at' = `object_attributes.created_at` string
- 'pipeline_finished_at' = `object_attributes.finished_at` string
- 'pipeline_duration' = `object_attributes.duration` int


### [Job events](https://docs.gitlab.com/user/project/integrations/webhook_events/#job-events)

Tags:
- 'event' = `job_event` string
- 'job_name' = `build_name` string
- 'job_stage' = `build_stage` string
- 'job_status' = `build_status` string
- 'job_failure_reason' = `build_failure_reason` string
- 'allow_failure' = `build_allow_failure` string
- 'is_tag' = `tag` string
- 'project_id' = `project.id` string
- 'user_id' = `user.id` string
- 'runner_id' = `runner.id` string

Fields:
- 'job_id' = `build_id` int
- 'pipeline_id' = `pipeline_id` int
- 'ref' = `ref` string
- 'before_sha' = `before_sha` string
- 'sha' = `sha` string
- 'job_created_at' = `build_created_at` string
- 'job_started_at' = `build_started_at` string
- 'job_finished_at' = `build_finished_at` string
- 'job_duration' = `build_duration` float64
- 'job_queued_duration' = `build_queued_duration` float64

### [Merge request events](https://docs.gitlab.com/user/project/integrations/webhook_events/#merge-request-events)

Tags:
- 'event' = `merge_request_event` string
- 'project_id' = `project.id` string
- 'user_id' = `user.id` string
- 'author_id' = `object_attributes.author_id` string
- 'blocking_discussions_resolved' = `object_attributes.blocking_discussions_resolved` string
- 'work_in_progress' = `object_attributes.work_in_progress` string
- 'draft' = `object_attributes.draft` string
- 'merge_status' = `object_attributes.detailed_merge_status` string

Fields:
- 'merge_request_id' = `object_attributes.id` int
- 'title' = `object_attributes.title` string
- 'description' = `object_attributes.description` string
- 'source_branch' = `object_attributes.source_branch` string
- 'target_branch' = `object_attributes.target_branch` string
- 'created_at' = `object_attributes.created_at` string
- 'updated_at' = `object_attributes.updated_at` string
