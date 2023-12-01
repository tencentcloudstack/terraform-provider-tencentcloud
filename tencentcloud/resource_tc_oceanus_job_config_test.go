package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusJobConfigResource_basic -v
func TestAccTencentCloudNeedFixOceanusJobConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusJobConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "job_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "entrypoint_class"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "program_args"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "default_parallelism"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "log_collect"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "job_manager_spec"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "task_manager_spec"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "cls_logset_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "cls_topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "log_collect_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "work_space_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "log_level"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "auto_recover"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_config.example", "expert_mode_on"),
				),
			},
		},
	})
}

const testAccOceanusJobConfig = `
resource "tencentcloud_oceanus_job_config" "example" {
  job_id           = "cql-4xwincyn"
  entrypoint_class = "tf_example"
  program_args     = "--conf Key=Value"
  remark           = "remark."
  resource_refs {
    resource_id = "resource-q22ntswy"
    version     = 1
    type        = 1
  }
  default_parallelism = 1
  properties {
    key   = "pipeline.max-parallelism"
    value = "2048"
  }
  log_collect       = true
  job_manager_spec  = "1"
  task_manager_spec = "1"
  cls_logset_id     = "cd9adbb5-6b7d-48d2-9870-77658959c7a4"
  cls_topic_id      = "cec4c2f1-0bf3-470e-b1a5-b1c451e88838"
  log_collect_type  = 2
  work_space_id     = "space-2idq8wbr"
  log_level         = "INFO"
  auto_recover      = 1
  expert_mode_on    = false
}
`
