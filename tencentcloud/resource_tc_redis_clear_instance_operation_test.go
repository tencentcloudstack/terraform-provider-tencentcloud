package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisClearInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisClearInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_clear_instance_operation.clear_instance_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_clear_instance_operation.clear_instance_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisClearInstanceOperation = `

resource "tencentcloud_redis_clear_instance_operation" "clear_instance_operation" {
  instance_id = "crs-c1nl9rpv"
  password = &lt;nil&gt;
}

`
