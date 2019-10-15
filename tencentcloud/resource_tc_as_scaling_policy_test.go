package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudAsScalingPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingPolicy(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingPolicyExists("tencentcloud_as_scaling_policy.scaling_policy"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_policy.scaling_policy", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "policy_name", "tf-as-scaling-policy"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "adjustment_type", "EXACT_CAPACITY"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "adjustment_value", "0"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "comparison_operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "metric_name", "CPU_UTILIZATION"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "threshold", "80"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "period", "300"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "continuous_time", "10"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "statistic", "AVERAGE"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "cooldown", "360"),
				),
			},
			// test update case
			{
				Config: testAccAsScalingPolicy_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingPolicyExists("tencentcloud_as_scaling_policy.scaling_policy"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_scaling_policy.scaling_policy", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "adjustment_type", "CHANGE_IN_CAPACITY"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "adjustment_value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "comparison_operator", "GREATER_THAN_OR_EQUAL_TO"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "metric_name", "MEM_UTILIZATION"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "threshold", "85"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "period", "60"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "continuous_time", "9"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "statistic", "MAXIMUM"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_policy.scaling_policy", "cooldown", "300"),
				),
			},
		},
	})
}

func testAccCheckAsScalingPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling policy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling policy id is not set")
		}
		asService := AsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, _, err := asService.DescribeScalingPolicyById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAsScalingPolicyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	asService := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_scaling_policy" {
			continue
		}

		_, _, err := asService.DescribeScalingPolicyById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("auto scaling policy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsScalingPolicy() string {
	return `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-scaling-group"
  configuration_id   = "${tencentcloud_as_scaling_config.launch_configuration.id}"
  max_size           = 1
  min_size           = 0
  vpc_id             = "${tencentcloud_vpc.vpc.id}"
  subnet_ids         = ["${tencentcloud_subnet.subnet.id}"]
}

resource "tencentcloud_as_scaling_policy" "scaling_policy" {
  scaling_group_id    = "${tencentcloud_as_scaling_group.scaling_group.id}"
  policy_name         = "tf-as-scaling-policy"
  adjustment_type     = "EXACT_CAPACITY"
  adjustment_value    = 0
  comparison_operator = "GREATER_THAN"
  metric_name         = "CPU_UTILIZATION"
  threshold           = 80
  period              = 300
  continuous_time     = 10
  statistic           = "AVERAGE"
  cooldown            = 360
}
`
}

func testAccAsScalingPolicy_update() string {
	return `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-scaling-group"
  configuration_id   = "${tencentcloud_as_scaling_config.launch_configuration.id}"
  max_size           = 1
  min_size           = 0
  vpc_id             = "${tencentcloud_vpc.vpc.id}"
  subnet_ids         = ["${tencentcloud_subnet.subnet.id}"]
}

resource "tencentcloud_as_scaling_policy" "scaling_policy" {
  scaling_group_id    = "${tencentcloud_as_scaling_group.scaling_group.id}"
  policy_name         = "tf-as-scaling-policy"
  adjustment_type     = "CHANGE_IN_CAPACITY"
  adjustment_value    = 1
  comparison_operator = "GREATER_THAN_OR_EQUAL_TO"
  metric_name         = "MEM_UTILIZATION"
  threshold           = 85
  period              = 60
  continuous_time     = 9
  statistic           = "MAXIMUM"
  cooldown            = 300
}
`
}
