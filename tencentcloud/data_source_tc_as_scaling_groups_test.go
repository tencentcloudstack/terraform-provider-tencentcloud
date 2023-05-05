package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAsScalingGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingGroupsDataSource_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingGroupExists("tencentcloud_as_scaling_group.scaling_group"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.scaling_group_name", "tf-as-group-ds-basic"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.configuration_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.max_size", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.min_size", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.create_time"),

					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.0.scaling_group_name", "tf-as-group-ds-basic"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.0.configuration_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.0.max_size", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.0.min_size", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.0.vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_name", "scaling_group_list.0.create_time"),

					resource.TestMatchResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.scaling_group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.configuration_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.max_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.min_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.subnet_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups_tags", "scaling_group_list.0.tags.test", "test"),
				),
			},
		},
	})
}

func TestAccTencentCloudAsScalingGroupsDataSource_full(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingGroupsDataSource_full(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingGroupExists("tencentcloud_as_scaling_group.scaling_group"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.scaling_group_name", "tf-as-group-ds-full"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.configuration_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.max_size", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.min_size", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.project_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.default_cooldown", "400"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.desired_capacity", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.termination_policies.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.termination_policies.0", "NEWEST_INSTANCE"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.retry_policy", "INCREMENTAL_INTERVALS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_groups.scaling_groups", "scaling_group_list.0.create_time"),
				),
			},
		},
	})
}

// todo
func testAccAsScalingGroupsDataSource_basic() string {
	return defaultAsVariable + `
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
  configuration_name = "tf-as-configuration-ds-basic"
  image_id           = "img-2lr9q49h"
  instance_types     = [data.tencentcloud_instance_types.default.instance_types.0.instance_type]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-group-ds-basic"
  configuration_id   = tencentcloud_as_scaling_config.launch_configuration.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_as_scaling_groups" "scaling_groups" {
  scaling_group_id = tencentcloud_as_scaling_group.scaling_group.id
}

data "tencentcloud_as_scaling_groups" "scaling_groups_name" {
  scaling_group_name = tencentcloud_as_scaling_group.scaling_group.scaling_group_name
}

data "tencentcloud_as_scaling_groups" "scaling_groups_tags" {
  tags = tencentcloud_as_scaling_group.scaling_group.tags
}
`
}

func testAccAsScalingGroupsDataSource_full() string {
	return defaultAsVariable + `
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
  configuration_name = "tf-as-configuration-ds-full"
  image_id           = "img-2lr9q49h"
  instance_types     = [data.tencentcloud_instance_types.default.instance_types.0.instance_type]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-group-ds-full"
  configuration_id     = tencentcloud_as_scaling_config.launch_configuration.id
  max_size             = 1
  min_size             = 0
  vpc_id               = tencentcloud_vpc.vpc.id
  subnet_ids           = [tencentcloud_subnet.subnet.id]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = 1
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"
}

data "tencentcloud_as_scaling_groups" "scaling_groups" {
  scaling_group_id = tencentcloud_as_scaling_group.scaling_group.id
}
`
}
