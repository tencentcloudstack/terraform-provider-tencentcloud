package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbLogFileRetentionPeriodResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbLogFileRetentionPeriod,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_log_file_retention_period.log_file_retention_period", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_log_file_retention_period.log_file_retention_period",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbLogFileRetentionPeriod = `

resource "tencentcloud_mariadb_log_file_retention_period" "log_file_retention_period" {
  instance_id = &lt;nil&gt;
  days = &lt;nil&gt;
}

`
