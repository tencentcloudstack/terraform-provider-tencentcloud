package crs_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudRedisAuditLogResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisAuditLog,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_audit_log.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_audit_log.example", "log_sub_type", "all"),
					resource.TestCheckResourceAttr("tencentcloud_redis_audit_log.example", "log_expire_day", "7"),
					resource.TestCheckResourceAttr("tencentcloud_redis_audit_log.example", "high_log_expire_day", "7"),
					resource.TestCheckResourceAttr("tencentcloud_redis_audit_log.example", "degrade_strategy", "500"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_audit_log.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"log_sub_type",
					"log_expire_day",
					"high_log_expire_day",
					"degrade_strategy",
				},
			},
		},
	})
}

const testAccRedisAuditLog = `
resource "tencentcloud_redis_audit_log" "example" {
  instance_id          = "crs-6eqwe3lt"
  log_sub_type         = "all"
  log_expire_day       = 7
  high_log_expire_day  = 7
  degrade_strategy     = 500
}
`
