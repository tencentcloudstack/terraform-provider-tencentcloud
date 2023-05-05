package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudRedisBackupConfig(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudRedisBackupConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisBackupConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisBackupConfigExists("tencentcloud_redis_backup_config.redis_backup_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_backup_config.redis_backup_config", "redis_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_backup_config.redis_backup_config", "backup_time", "06:00-07:00"),
					// backup_period doesn't work right now by design, skipped1
					//resource.TestCheckResourceAttr("tencentcloud_redis_backup_config.redis_backup_config", "backup_period.#", "1"),
					//resource.TestCheckResourceAttr("tencentcloud_redis_backup_config.redis_backup_config", "backup_period.1970423419", "Wednesday"),
				),
			},
			{
				Config: testAccRedisBackupConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccTencentCloudRedisBackupConfigExists("tencentcloud_redis_backup_config.redis_backup_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_backup_config.redis_backup_config", "redis_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_backup_config.redis_backup_config", "backup_time", "01:00-02:00"),
					//resource.TestCheckResourceAttr("tencentcloud_redis_backup_config.redis_backup_config", "backup_period.#", "2"),
					//resource.TestCheckResourceAttr("tencentcloud_redis_backup_config.redis_backup_config", "backup_period.1075549138", "Sunday"),
					//resource.TestCheckResourceAttr("tencentcloud_redis_backup_config.redis_backup_config", "backup_period.3286956037", "Saturday"),
				),
			},
			{
				ResourceName:      "tencentcloud_redis_backup_config.redis_backup_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTencentCloudRedisBackupConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, _, err := service.DescribeAutoBackupConfig(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccTencentCloudRedisBackupConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_redis_backup_config" {
			continue
		}
		time.Sleep(5 * time.Second)
		has, isolated, err := service.CheckRedisDestroyOk(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has || isolated {
			return nil
		}
		return fmt.Errorf("instance %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccRedisBackupConfigUpdate() string {
	return defaultVpcSubnets + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone = "ap-guangzhou-3"
  type_id           = 2 
  password          = "test12345789"
  mem_size          = 8192
  name              = "terraform_test"
  vpc_id 			 = local.vpc_id
  subnet_id			 = local.subnet_id
}
resource "tencentcloud_redis_backup_config" "redis_backup_config" {
  redis_id      = tencentcloud_redis_instance.redis_instance_test.id
  backup_time   = "01:00-02:00"
}`
}

func testAccRedisBackupConfig() string {
	return defaultVpcSubnets + `
resource "tencentcloud_redis_instance" "redis_instance_test" {
  availability_zone = "ap-guangzhou-3"
  type_id           = 2 
  password          = "test12345789"
  mem_size          = 8192
  name              = "terraform_test"
  vpc_id 			 = local.vpc_id
  subnet_id			 = local.subnet_id
}
resource "tencentcloud_redis_backup_config" "redis_backup_config" {
  redis_id      = tencentcloud_redis_instance.redis_instance_test.id
  backup_time   = "06:00-07:00"
}`
}
