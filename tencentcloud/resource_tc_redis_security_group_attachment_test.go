package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisSecurityGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisSecurityGroupAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_security_group_attachment.security_group_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_security_group_attachment.security_group_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisSecurityGroupAttachment = `

resource "tencentcloud_redis_security_group_attachment" "security_group_attachment" {
  product = "redis"
  instance_ids = 
  security_group_id = "crs-c1nl9rpv"
}

`
