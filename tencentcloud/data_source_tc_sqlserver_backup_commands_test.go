package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverBackupCommands_basic -v
func TestAccTencentCloudSqlserverBackupCommands_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBackupCommands,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_backup_commands.example", "list.#"),
				),
			},
		},
	})
}

const testAccSqlserverBackupCommands = `
data "tencentcloud_sqlserver_backup_commands" "example" {
  backup_file_type = "FULL"
  data_base_name   = "keep-publish-instance"
  is_recovery      = "NO"
}
`
