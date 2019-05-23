package tencentcloud

import (
	"context"
	"fmt"
	"time"

	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type", "master_slave_redis"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "mem_size", "8192"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "name", "terrform_test"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "status", "online"),
				),
			},
			{
				Config: testAccRedisInstanceUpdateName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisInstanceExists("tencentcloud_redis_instance.redis_instance_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_instance.redis_instance_test", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "port", "6379"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type", "master_slave_redis"),
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
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "type", "master_slave_redis"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "mem_size", "12288"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "name", "terrform_test_update"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_redis_instance.redis_instance_test", "status", "online"),
				),
			},
		},
	})
}

func testAccTencentCloudRedisInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		has, _, _, err := service.CheckRedisCreateOk(ctx, rs.Primary.ID)
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

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_redis_instance" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, _, info, err := service.CheckRedisCreateOk(ctx, rs.Primary.ID)

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
resource "tencentcloud_redis_instance" "redis_instance_test"{
	availability_zone="ap-guangzhou-3"
	type="master_slave_redis"
	password="test12345789"
	mem_size=8192
	name="terrform_test"
	port=6379
}`
}

func testAccRedisInstanceUpdateName() string {
	return `
resource "tencentcloud_redis_instance" "redis_instance_test"{
	availability_zone="ap-guangzhou-3"
	type="master_slave_redis"
	password="test12345789"
	mem_size=8192
	name="terrform_test_update"
	port=6379
}`
}
func testAccRedisInstanceUpdateMemsizeAndPassword() string {
	return `
resource "tencentcloud_redis_instance" "redis_instance_test"{
	availability_zone="ap-guangzhou-3"
	type="master_slave_redis"
	password="AAA123456BBB"
	mem_size=12288
	name="terrform_test_update"
	port=6379
}`
}
