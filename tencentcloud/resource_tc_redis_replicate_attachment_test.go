package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisReplicateAttachmentResource_basic -v
func TestAccTencentCloudRedisReplicateAttachmentResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisReplicateAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replicate_attachment.replicate_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replicate_attachment.replicate_attachment", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_replicate_attachment.replicate_attachment", "group_id", "crs-rpl-orfiwmn5"),
					resource.TestCheckResourceAttr("tencentcloud_redis_replicate_attachment.replicate_attachment", "instance_role", "rw"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_replicate_attachment.replicate_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisReplicateAttachmentVar = `
variable "security_groups" {
  default = "` + defaultCrsSecurityGroups + `"
}
variable "vpc_id" {
	default = "` + defaultCrsVpcId + `"
}
variable "subnet_id" {
	default = "` + defaultCrsSubnetId + `"
}
variable "group_id" {
	default = "` + defaultCrsGroupId + `"
}
`

const testAccRedisReplicateAttachment = testAccRedisReplicateAttachmentVar + `

resource "tencentcloud_redis_instance" "redis_cluster" {
	availability_zone = "ap-guangzhou-6"
	type_id            = 7
	password           = "AAA123456BBB"
	mem_size           = 4096
	name               = "terraform_cluster"
	port               = 6379
	redis_shard_num    = 1
	redis_replicas_num = 1
	vpc_id 			 = var.vpc_id
	subnet_id			 = var.subnet_id
	security_groups    = [var.security_groups]
}

resource "tencentcloud_redis_replicate_attachment" "replicate_attachment" {
	instance_id = tencentcloud_redis_instance.redis_cluster.id
	group_id = var.group_id
	instance_role = "rw"
}

`
