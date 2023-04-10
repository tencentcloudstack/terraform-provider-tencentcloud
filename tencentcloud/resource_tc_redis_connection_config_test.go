package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRedisConnectionConfigResource_basic -v
func TestAccTencentCloudRedisConnectionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisConnectionConfigExists("tencentcloud_redis_connection_config.connection_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_connection_config.connection_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_connection_config.connection_config", "instance_id", defaultCrsInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_redis_connection_config.connection_config", "client_limit", "20000"),
					resource.TestCheckResourceAttr("tencentcloud_redis_connection_config.connection_config", "bandwidth", "20"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_redis_connection_config.connection_config",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccTencentCloudRedisConnectionConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		connectionConfig, err := service.DescribeRedisInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if connectionConfig == nil {
			return fmt.Errorf("redis connection config %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccRedisConnectionConfigVar = `
variable "instance_id" {
	default = "` + defaultCrsInstanceId + `"
}
`

const testAccRedisConnectionConfig = testAccRedisConnectionConfigVar + `

resource "tencentcloud_redis_connection_config" "connection_config" {
   instance_id = var.instance_id
   client_limit = "20000"
   bandwidth = "20"
}

`
