package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbBackupTimeResource_basic -v
func TestAccTencentCloudMariadbBackupTimeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbBackupTime,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_backup_time.backup_time", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_backup_time.backup_time",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbBackupTime = `
resource "tencentcloud_mariadb_backup_time" "backup_time" {
  instance_id       = "tdsql-9vqvls95"
  start_backup_time = "01:00"
  end_backup_time   = "04:00"
}
`
