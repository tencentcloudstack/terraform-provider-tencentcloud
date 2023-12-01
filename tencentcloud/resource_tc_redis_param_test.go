package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudRedisParamResource_basic -v
func TestAccTencentCloudRedisParamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisParam,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisParamExists("tencentcloud_redis_param.param"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_param.param", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_id", defaultCrsInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.cluster-node-timeout", "15000"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.disable-command-list", "\"\""),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.hash-max-ziplist-entries", "512"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.hash-max-ziplist-value", "64"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.hz", "20"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.lazyfree-lazy-eviction", "yes"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.lazyfree-lazy-expire", "yes"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.lazyfree-lazy-server-del", "yes"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.maxmemory-policy", "noeviction"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.notify-keyspace-events", "\"\""),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.proxy-slowlog-log-slower-than", "500"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.replica-lazy-flush", "yes"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.sentineauth", "no"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.set-max-intset-entries", "512"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.slowlog-log-slower-than", "10"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.timeout", "31536000"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.zset-max-ziplist-entries", "128"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.zset-max-ziplist-value", "64"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.lazyfree-lazy-user-del", "yes"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_param.param",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRedisParamUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisParamExists("tencentcloud_redis_param.param"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_param.param", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_id", defaultCrsInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.cluster-node-timeout", "15000"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.disable-command-list", "\"\""),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.hash-max-ziplist-entries", "512"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.hash-max-ziplist-value", "64"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.hz", "10"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.lazyfree-lazy-eviction", "yes"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.lazyfree-lazy-expire", "yes"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.lazyfree-lazy-server-del", "yes"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.maxmemory-policy", "noeviction"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.notify-keyspace-events", "\"\""),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.proxy-slowlog-log-slower-than", "500"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.replica-lazy-flush", "yes"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.sentineauth", "no"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.set-max-intset-entries", "512"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.slowlog-log-slower-than", "10"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.timeout", "31536000"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.zset-max-ziplist-entries", "128"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.zset-max-ziplist-value", "64"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param.param", "instance_params.lazyfree-lazy-user-del", "yes"),
				),
			},
		},
	})
}

func testAccTencentCloudRedisParamExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		params, err := service.DescribeRedisParamById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if len(params) == 0 {
			return fmt.Errorf("redis param %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccRedisParamVar = `
variable "instance_id" {
	default = "` + defaultCrsInstanceId + `"
}
`

const testAccRedisParam = testAccRedisParamVar + `

resource "tencentcloud_redis_param" "param" {
    instance_id     = var.instance_id
    instance_params = {
        "cluster-node-timeout"          = "15000"
        "disable-command-list"          = "\"\""
        "hash-max-ziplist-entries"      = "512"
        "hash-max-ziplist-value"        = "64"
        "hz"                            = "20"
        "lazyfree-lazy-eviction"        = "yes"
        "lazyfree-lazy-expire"          = "yes"
        "lazyfree-lazy-server-del"      = "yes"
        "maxmemory-policy"              = "noeviction"
        "notify-keyspace-events"        = "\"\""
        "proxy-slowlog-log-slower-than" = "500"
        "replica-lazy-flush"            = "yes"
        "sentineauth"                   = "no"
        "set-max-intset-entries"        = "512"
        "slowlog-log-slower-than"       = "10"
        "timeout"                       = "31536000"
        "zset-max-ziplist-entries"      = "128"
        "zset-max-ziplist-value"        = "64"
		"lazyfree-lazy-user-del"		= "yes"
    }
}

`

const testAccRedisParamUpdate = testAccRedisParamVar + `

resource "tencentcloud_redis_param" "param" {
    instance_id     = var.instance_id
    instance_params = {
        "cluster-node-timeout"          = "15000"
        "disable-command-list"          = "\"\""
        "hash-max-ziplist-entries"      = "512"
        "hash-max-ziplist-value"        = "64"
        "hz"                            = "10"
        "lazyfree-lazy-eviction"        = "yes"
        "lazyfree-lazy-expire"          = "yes"
        "lazyfree-lazy-server-del"      = "yes"
        "maxmemory-policy"              = "noeviction"
        "notify-keyspace-events"        = "\"\""
        "proxy-slowlog-log-slower-than" = "500"
        "replica-lazy-flush"            = "yes"
        "sentineauth"                   = "no"
        "set-max-intset-entries"        = "512"
        "slowlog-log-slower-than"       = "10"
        "timeout"                       = "31536000"
        "zset-max-ziplist-entries"      = "128"
        "zset-max-ziplist-value"        = "64"
		"lazyfree-lazy-user-del"		= "yes"
    }
}

`
