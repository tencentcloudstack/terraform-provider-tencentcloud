package tat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatCommandDataSource -v
func TestAccTencentCloudTatCommandDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTatCommand,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tat_command.command"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.command_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.command_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.command_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.content"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.created_by"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.enable_parameter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.formatted_description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.updated_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.username"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_command.command", "command_set.0.working_directory"),
				),
			},
		},
	})
}

const testAccDataSourceTatCommand = `

data "tencentcloud_tat_command" "command" {
	# command_id = ""
	# command_name = ""
	command_type = "SHELL"
	created_by = "TAT"
}

`
