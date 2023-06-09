package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisChangeMasterOperationResource_basic -v
func TestAccTencentCloudRedisChangeMasterOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisChangeMasterOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_change_master_operation.change_master_operation", "id")),
			},
		},
	})
}

const testAccRedisChangeMasterOperationVar = `
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

const testAccRedisChangeMasterOperation = testAccRedisChangeMasterOperationVar + `

resource "tencentcloud_redis_instance" "redis_cluster_rw" {
	availability_zone  = "ap-guangzhou-6"
	type_id            = 7
	password           = "AAA123456BBB"
	mem_size           = 4096
	name               = "terraform_cluster"
	port               = 6379
	redis_shard_num    = 1
	redis_replicas_num = 1
	vpc_id 			   = var.vpc_id
	subnet_id		   = var.subnet_id
	security_groups    = [var.security_groups]
}

resource "tencentcloud_redis_replicate_attachment" "replicate_attachment_rw" {
	instance_id = tencentcloud_redis_instance.redis_cluster_rw.id
	group_id = var.group_id
	instance_role = "rw"
}

resource "tencentcloud_redis_instance" "redis_cluster_r" {
	availability_zone  = "ap-guangzhou-6"
	type_id            = 7
	password           = "AAA123456BBB"
	mem_size           = 4096
	name               = "terraform_cluster"
	port               = 6379
	redis_shard_num    = 1
	redis_replicas_num = 1
	vpc_id 			   = var.vpc_id
	subnet_id		   = var.subnet_id
	security_groups    = [var.security_groups]
}

resource "tencentcloud_redis_replicate_attachment" "replicate_attachment_r" {
	instance_id = tencentcloud_redis_instance.redis_cluster_r.id
	group_id = var.group_id
	instance_role = "r"

	depends_on = [
		tencentcloud_redis_replicate_attachment.replicate_attachment_rw
  	]
}

resource "tencentcloud_redis_change_master_operation" "change_master_operation" {
	instance_id = tencentcloud_redis_instance.redis_cluster_r.id
	group_id = var.group_id

	depends_on = [
		tencentcloud_redis_replicate_attachment.replicate_attachment_rw,
    	tencentcloud_redis_replicate_attachment.replicate_attachment_r
  	]
}

resource "tencentcloud_redis_change_master_operation" "change_master_operation_rw" {
	instance_id = tencentcloud_redis_instance.redis_cluster_rw.id
	group_id = var.group_id

	depends_on = [
		tencentcloud_redis_change_master_operation.change_master_operation
  	]
}

`
