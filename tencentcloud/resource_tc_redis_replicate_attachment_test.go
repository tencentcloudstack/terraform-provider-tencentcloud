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
					resource.TestCheckResourceAttr("tencentcloud_redis_replicate_attachment.replicate_attachment", "group_id", "crs-rpl-orfiwmn5"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replicate_attachment.replicate_attachment", "master_instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_replicate_attachment.replicate_attachment", "instance_ids.#", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_replicate_attachment.replicate_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisReplicateAttachmentUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replicate_attachment.replicate_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_replicate_attachment.replicate_attachment", "group_id", "crs-rpl-orfiwmn5"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replicate_attachment.replicate_attachment", "master_instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_replicate_attachment.replicate_attachment", "instance_ids.#", "2"),
				),
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
const testAccRedisReplicateInstance = testAccRedisReplicateAttachmentVar + `

resource "tencentcloud_redis_instance" "redis_cluster_rw" {
	count = 3

	availability_zone = "ap-guangzhou-6"
	type_id            = 7
	password           = "AAA123456BBB"
	mem_size           = 4096
	name               = "terraform_test_${count.index}"
	port               = 6379
	redis_shard_num    = 1
	redis_replicas_num = 1
	vpc_id 			   = var.vpc_id
	subnet_id		   = var.subnet_id
	security_groups    = [var.security_groups]
}

`

const testAccRedisReplicateAttachment = testAccRedisReplicateInstance + `

resource "tencentcloud_redis_replicate_attachment" "replicate_attachment" {
  group_id           = var.group_id
  master_instance_id = tencentcloud_redis_instance.redis_cluster_rw.0.id
  instance_ids       = "${tencentcloud_redis_instance.redis_cluster_rw.*.id}"
}

`

const testAccRedisReplicateAttachmentUp = testAccRedisReplicateInstance + `

resource "tencentcloud_redis_replicate_attachment" "replicate_attachment" {
  group_id           = var.group_id
  master_instance_id = tencentcloud_redis_instance.redis_cluster_rw.1.id
  instance_ids       = [tencentcloud_redis_instance.redis_cluster_rw.0.id, tencentcloud_redis_instance.redis_cluster_rw.1.id]
}

`
