package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbLogFileRetentionPeriod_basic -v
func TestAccTencentCloudMariadbLogFileRetentionPeriod_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbLogFileRetentionPeriod,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_log_file_retention_period.logFileRetentionPeriod", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_log_file_retention_period.logFileRetentionPeriod", "days", "8"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_log_file_retention_period.logFileRetentionPeriod",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbLogFileRetentionPeriod = testAccMariadbHourDbInstance + `

resource "tencentcloud_mariadb_log_file_retention_period" "logFileRetentionPeriod" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
  days = "8"
}

`
