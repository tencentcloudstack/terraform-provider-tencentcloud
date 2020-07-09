package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudRedisInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_shard_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_replicas_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "mem_size", "8192"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "name", "terrform_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "status", "online"),
				),
			},
			{
				Config: testAccRedisInstanceTags(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "tags.test", "test"),
				),
			},
			{
				Config: testAccRedisInstanceTagsUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckNoResourceAttr("tencentcloud_redis_instance.redis_instance_test", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "tags.abc", "abc"),
				),
			},
			{
				Config: testAccRedisInstanceUpdateName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_shard_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_replicas_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "mem_size", "8192"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "name", "terrform_test_update"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "status", "online"),
				),
			},
			{
				Config: testAccRedisInstanceUpdateMemsizeAndPassword(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_shard_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "redis_replicas_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "mem_size", "12288"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "name", "terrform_test_update"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "status", "online"),
				),
			},
			{
				ResourceName:            "tencentcloud_redis_instance.redis_instance_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "type", "redis_shard_num", "redis_replicas_num", "force_delete"},
			},
			{
				Config: testAccRedisInstancePrepaidBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_prepaid_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_prepaid_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_prepaid_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "type_id", "2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "redis_shard_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "redis_replicas_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "mem_size", "8192"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "name", "terrform_prepaid_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "status", "online"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_prepaid_instance_test", "charge_type", "PREPAID"),
				),
			},
		},
	})
}

func testAccTencentCloudRedisInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		has, _, _, err := service.CheckRedisOnlineOk(ctx, rs.Primary.ID)
		if has {
			return nil
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("redis not exists.")
	}
}

func testAccTencentCloudRedisInstanceDestroy(s *terraform.State) error {

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_redis_instance" {
			continue
		}
		time.Sleep(5 * time.Second)
		has, _, info, err := service.CheckRedisOnlineOk(ctx, rs.Primary.ID)

		if !has {
			return nil
		}

		if info != nil {
			if *info.Status == REDIS_STATUS_ISOLATE || *info.Status == REDIS_STATUS_TODELETE {
				return nil
			}
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("redis not delete ok")
	}
	return nil
}

func testAccRedisInstanceBasic() string {
	return `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terrform_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
}`
}

func testAccRedisInstanceTags() string {
	return `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terrform_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  tags = {
    test = "test"
  }
}`
}

func testAccRedisInstanceTagsUpdate() string {
	return `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terrform_test"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  tags = {
    abc = "abc"
  }
}`
}

func testAccRedisInstanceUpdateName() string {
	return `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone  = "ap-guangzhou-3"
  type_id            = 2
  password           = "test12345789"
  mem_size           = 8192
  name               = "terrform_test_update"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1
  tags = {
    abc = "abc"
  }
}`
}

func testAccRedisInstanceUpdateMemsizeAndPassword() string {
	return `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone = "ap-guangzhou-3"
  type_id            = 2
  password           = "AAA123456BBB"
  mem_size           = 12288
  name               = "terrform_test_update"
  port               = 6379
  redis_shard_num    = 1
  redis_replicas_num = 1

  tags = {
    "abc" = "abc"
  }
}`
}

func testAccRedisInstancePrepaidBasic() string {
	return `
resource "tencentcloud_redis_instance" "redis_prepaid_instance_test" {
  availability_zone                   = "ap-guangzhou-3"
  type_id                             = 2
  password                            = "test12345789"
  mem_size                            = 8192
  name                                = "terrform_prepaid_test"
  port                                = 6379
  redis_shard_num                     = 1
  redis_replicas_num                  = 1
  charge_type                         = "PREPAID"
  prepaid_period                      = 2
  force_delete                        = true
}`
}
