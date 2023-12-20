package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisRenewInstanceOperationResource_basic -v
func TestAccTencentCloudRedisRenewInstanceOperationResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
				Config:    testAccRedisRenewInstanceOperation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_renew_instance_operation.renew_instance_operation", "id"),
				),
			},
		},
	})
}

func testAccRedisRenewInstanceOperation() string {
	return testAccRedisInstancePrepaidBasic() + `
resource "tencentcloud_redis_renew_instance_operation" "renew_instance_operation" {
	instance_id = tencentcloud_redis_instance.redis_prepaid_instance_test.id
	period = 1
	modify_pay_mode = "prepaid"
}`
}
