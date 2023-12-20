package as_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAsScalingPoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAsScalingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingPoliciesDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingPolicyExists("tencentcloud_as_scaling_policy.scaling_policy"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.scaling_group_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.policy_name", "tf-as-scaling-policy"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.adjustment_type", "EXACT_CAPACITY"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.adjustment_value", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.comparison_operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.metric_name", "CPU_UTILIZATION"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.threshold", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.period", "300"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.continuous_time", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.statistic", "AVERAGE"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies", "scaling_policy_list.0.cooldown", "360"),

					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.scaling_group_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.policy_name", "tf-as-scaling-policy"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.adjustment_type", "EXACT_CAPACITY"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.adjustment_value", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.comparison_operator", "GREATER_THAN"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.metric_name", "CPU_UTILIZATION"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.threshold", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.period", "300"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.continuous_time", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.statistic", "AVERAGE"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_policies.scaling_policies_name", "scaling_policy_list.0.cooldown", "360"),
				),
			},
		},
	})
}

// todo
func testAccAsScalingPoliciesDataSource() string {
	return tcacctest.DefaultAsVariable + `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = var.availability_zone
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration-ds"
  image_id           = "img-9qabwvbn"
  instance_types     = [data.tencentcloud_instance_types.default.instance_types.0.instance_type]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-scaling-group-datasource"
  configuration_id   = tencentcloud_as_scaling_config.launch_configuration.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_as_scaling_policy" "scaling_policy" {
  scaling_group_id    = tencentcloud_as_scaling_group.scaling_group.id
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

data "tencentcloud_as_scaling_policies" "scaling_policies" {
  scaling_policy_id = tencentcloud_as_scaling_policy.scaling_policy.id
}

data "tencentcloud_as_scaling_policies" "scaling_policies_name" {
  policy_name = tencentcloud_as_scaling_policy.scaling_policy.policy_name
}
`
}
