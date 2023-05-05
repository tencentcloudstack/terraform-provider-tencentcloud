package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_as_scaling_group", &resource.Sweeper{
		Name: "tencentcloud_as_scaling_group",
		F:    testSweepAsScalingGroups,
	})
}

func testSweepAsScalingGroups(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(*TencentCloudClient)

	asService := AsService{
		client: client.apiV3Conn,
	}
	scalingGroups, err := asService.DescribeAutoScalingGroupByFilter(ctx, "", "", "", nil)
	if err != nil {
		return fmt.Errorf("list scaling group error: %s", err.Error())
	}

	for _, v := range scalingGroups {
		scalingGroupId := *v.AutoScalingGroupId
		scalingGroupName := *v.AutoScalingGroupName
		if !strings.HasPrefix(scalingGroupName, "tf-as-") {
			continue
		}

		if err = asService.DeleteScalingGroup(ctx, scalingGroupId); err != nil {
			log.Printf("[ERROR] delete scaling group %s error: %s", scalingGroupId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudAsScalingGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingGroup_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingGroupExists("tencentcloud_as_scaling_group.scaling_group"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "scaling_group_name", "tf-as-group-basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "configuration_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "max_size", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "min_size", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "subnet_ids.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "create_time"),
				),
			},
			{
				ResourceName:      "tencentcloud_as_scaling_group.scaling_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudAsScalingGroup_full(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingGroup_full(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingGroupExists("tencentcloud_as_scaling_group.scaling_group"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "scaling_group_name", "tf-as-group-full"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "configuration_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "max_size", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "min_size", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "subnet_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "default_cooldown", "400"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "desired_capacity", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "termination_policies.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "termination_policies.0", "NEWEST_INSTANCE"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "retry_policy", "INCREMENTAL_INTERVALS"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "instance_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_group.scaling_group", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "tags.test", "test"),
				),
			},
			{
				Config: testAccAsScalingGroup_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingGroupExists("tencentcloud_as_scaling_group.scaling_group"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "scaling_group_name", "tf-as-group-update"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "max_size", "2"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "default_cooldown", "300"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "desired_capacity", "0"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "termination_policies.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "termination_policies.0", "OLDEST_INSTANCE"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "retry_policy", "IMMEDIATE_RETRY"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "scaling_mode", "WAKE_UP_STOPPED_SCALING"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "replace_monitor_unhealthy", "true"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "replace_load_balancer_unhealthy", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_as_scaling_group.scaling_group", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_group.scaling_group", "tags.abc", "abc"),
				),
			},
		},
	})
}

func testAccCheckAsScalingGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling group id is not set")
		}
		asService := AsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := asService.DescribeAutoScalingGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has < 1 {
			return fmt.Errorf("auto scaling group not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckAsScalingGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	asService := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_scaling_group" {
			continue
		}

		_, has, err := asService.DescribeAutoScalingGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has > 0 {
			return fmt.Errorf("auto scaling group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsScalingGroup_basic() string {
	return fmt.Sprintf(`
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration-basic"
  image_id           = "img-2lr9q49h"
  instance_types     = ["SA1.SMALL1","SA2.SMALL1","SA2.SMALL2","SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
    # instance_name_style
  }
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-group-basic"
  configuration_id   = tencentcloud_as_scaling_config.launch_configuration.id
  max_size           = 1
  min_size           = 0
  vpc_id             = "%s"
  subnet_ids         = ["%s"]
}
`, defaultVpcId, defaultSubnetId)
}

func testAccAsScalingGroup_full() string {
	return fmt.Sprintf(`

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration-full"
  image_id           = "img-2lr9q49h"
  instance_types     = ["SA1.SMALL1","SA2.SMALL1","SA2.SMALL2","SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name-full"
  }
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-group-full"
  configuration_id     = tencentcloud_as_scaling_config.launch_configuration.id
  max_size             = 1
  min_size             = 0
  vpc_id               = "%s"
  subnet_ids           = ["%s"]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = 1
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"

  tags = {
    "test" = "test"
  }
}
`, defaultVpcId, defaultSubnetId)
}

func testAccAsScalingGroup_update() string {
	return fmt.Sprintf(`

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration-full"
  image_id           = "img-2lr9q49h"
  instance_types     = ["SA1.SMALL1","SA2.SMALL1","SA2.SMALL2","SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name-full"
  }
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-group-update"
  configuration_id     = tencentcloud_as_scaling_config.launch_configuration.id
  max_size             = 2
  min_size             = 0
  vpc_id               = "%s"
  subnet_ids           = ["%s"]
  project_id           = 0
  default_cooldown     = 300
  desired_capacity     = 0
  termination_policies = ["OLDEST_INSTANCE"]
  retry_policy         = "IMMEDIATE_RETRY"
  scaling_mode		   = "WAKE_UP_STOPPED_SCALING"
  replace_monitor_unhealthy       = true
  replace_load_balancer_unhealthy = true

  tags = {
    "abc" = "abc"
  }
}
`, defaultVpcId, defaultSubnetId)
}
