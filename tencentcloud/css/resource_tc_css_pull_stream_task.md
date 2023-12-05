Provides a resource to create a css pull_stream_task

Example Usage

```hcl
resource "tencentcloud_css_pull_stream_task" "pull_stream_task" {
  source_type = "source_type"
  source_urls = ["source_urls"]
  domain_name = "domain_name"
  app_name = "app_name"
  stream_name = "stream_name"
  start_time = "2022-11-16T22:09:28Z"
  end_time = "2022-11-16T22:09:28Z"
  operator = "admin"
  comment = "comment."
  }

```
Import

css pull_stream_task can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_pull_stream_task.pull_stream_task pullStreamTask_id
```