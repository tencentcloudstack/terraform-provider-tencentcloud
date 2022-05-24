package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAsLifecycleHook(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsLifecycleHookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsLifecycleHook(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsLifecycleHookExists("tencentcloud_as_lifecycle_hook.lifecycle_hook"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_lifecycle_hook.lifecycle_hook", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "lifecycle_hook_name", "tf-as-lifecycle-hook"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "lifecycle_transition", "INSTANCE_LAUNCHING"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "default_result", "CONTINUE"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "heartbeat_timeout", "500"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "notification_metadata", "tf test"),
				),
			},
			// test update case
			{
				Config: testAccAsLifecycleHook_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsLifecycleHookExists("tencentcloud_as_lifecycle_hook.lifecycle_hook"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_lifecycle_hook.lifecycle_hook", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "lifecycle_hook_name", "tf-as-lifecycle-test"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "lifecycle_transition", "INSTANCE_TERMINATING"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "default_result", "ABANDON"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "heartbeat_timeout", "300"),
					resource.TestCheckResourceAttr("tencentcloud_as_lifecycle_hook.lifecycle_hook", "notification_metadata", "tf lifecycle test"),
				),
			},
		},
	})
}

func testAccCheckAsLifecycleHookExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling lifecycle hook %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling lifecycle hook id is not set")
		}
		asService := AsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := asService.DescribeLifecycleHookById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has < 1 {
			return fmt.Errorf("auto scaling lifecycle hook not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckAsLifecycleHookDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	asService := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_lifecycle_hook" {
			continue
		}

		_, has, err := asService.DescribeLifecycleHookById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has > 0 {
			return fmt.Errorf("auto scaling lifecycle hook still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsLifecycleHook() string {
	return `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration-lifecycle-hook"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-scaling-group-lifecycle-hook"
  configuration_id   = tencentcloud_as_scaling_config.launch_configuration.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_as_lifecycle_hook" "lifecycle_hook" {
  scaling_group_id      = tencentcloud_as_scaling_group.scaling_group.id
  lifecycle_hook_name   = "tf-as-lifecycle-hook"
  lifecycle_transition  = "INSTANCE_LAUNCHING"
  default_result        = "CONTINUE"
  heartbeat_timeout     = 500
  notification_metadata = "tf test"
}
`
}

func testAccAsLifecycleHook_update() string {
	return `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration-lifecycle-hook"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-scaling-group-lifecycle-hook"
  configuration_id   = tencentcloud_as_scaling_config.launch_configuration.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_as_lifecycle_hook" "lifecycle_hook" {
  scaling_group_id      = tencentcloud_as_scaling_group.scaling_group.id
  lifecycle_hook_name   = "tf-as-lifecycle-test"
  lifecycle_transition  = "INSTANCE_TERMINATING"
  default_result        = "ABANDON"
  heartbeat_timeout     = 300
  notification_metadata = "tf lifecycle test"
}
`
}
