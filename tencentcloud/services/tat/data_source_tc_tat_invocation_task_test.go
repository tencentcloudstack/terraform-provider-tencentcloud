package tat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatInvocationTaskDataSource_basic -v
func TestAccTencentCloudTatInvocationTaskDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvocationTaskDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tat_invocation_task.invocation_task"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_document.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_document.0.command_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_document.0.content"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_document.0.output_cos_bucket_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_document.0.output_cos_key_prefix"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_document.0.timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_document.0.username"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_document.0.working_directory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.command_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.invocation_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.invocation_source"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.invocation_task_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.task_result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.task_result.0.exec_end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.task_result.0.exec_start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.task_result.0.output_upload_cos_error_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.task_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_invocation_task.invocation_task", "invocation_task_set.0.updated_time"),
				),
			},
		},
	})
}

const testAccTatInvocationTaskDataSource = `

data "tencentcloud_tat_invocation_task" "invocation_task" {
	# invocation_task_ids = ["invt-a8bv0ip7"]
	filters {
	  name = "command-id"
	  values = ["cmd-rxbs7f5z"]
	}
	hide_output = true
}

`
