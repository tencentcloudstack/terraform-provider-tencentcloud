package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresqlRecoveryTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlRecoveryTimeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_recovery_time.recovery_time")),
			},
		},
	})
}

const testAccPostgresqlRecoveryTimeDataSource = `

data "tencentcloud_postgresql_recovery_time" "recovery_time" {
  d_b_instance_id = ""
      tags = {
    "createdBy" = "terraform"
  }
}

`
