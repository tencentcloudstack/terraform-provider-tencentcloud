package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisReplicateAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisReplicateAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_replicate_attachment.replicate_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_replicate_attachment.replicate_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisReplicateAttachment = `

resource "tencentcloud_redis_replicate_attachment" "replicate_attachment" {
  instance_id = "crs-c1nl9rpv"
  group_id = "crs-rpl-c1nl9rpv"
  instance_role = "rw"
}

`
