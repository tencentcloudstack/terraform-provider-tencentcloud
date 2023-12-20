package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisMaintenanceWindowResource_basic -v
func TestAccTencentCloudRedisMaintenanceWindowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisMaintenanceWindow,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_maintenance_window.maintenance_window", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_maintenance_window.maintenance_window", "instance_id", tcacctest.DefaultCrsInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_redis_maintenance_window.maintenance_window", "start_time", "17:00"),
					resource.TestCheckResourceAttr("tencentcloud_redis_maintenance_window.maintenance_window", "end_time", "19:00"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_maintenance_window.maintenance_window",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisMaintenanceWindowVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultCrsInstanceId + `"
}
`

const testAccRedisMaintenanceWindow = testAccRedisMaintenanceWindowVar + `

resource "tencentcloud_redis_maintenance_window" "maintenance_window" {
  instance_id = var.instance_id
  start_time = "17:00"
  end_time = "19:00"
}

`
