package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisSslResource_basic -v
func TestAccTencentCloudRedisSslResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisSsl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_ssl.ssl", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_ssl.ssl", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_ssl.ssl", "ssl_config", "enabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_ssl.ssl",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisSslUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_ssl.ssl", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_ssl.ssl", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_ssl.ssl", "ssl_config", "disabled"),
				),
			},
		},
	})
}

const testAccRedisSslVar = `
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[1].zone
  name              = "tf_redis_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_redis_instance" "foo" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[1].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[1].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[1].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[1].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}
`

const testAccRedisSsl = testAccRedisSslVar + `

resource "tencentcloud_redis_ssl" "ssl" {
	instance_id = tencentcloud_redis_instance.foo.id
	ssl_config = "enabled"
  }

`

const testAccRedisSslUpdate = testAccRedisSslVar + `

resource "tencentcloud_redis_ssl" "ssl" {
	instance_id = tencentcloud_redis_instance.foo.id
	ssl_config = "disabled"
  }

`
