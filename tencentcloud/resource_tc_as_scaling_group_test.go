package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_as_scaling_group", &resource.Sweeper{
		Name: "tencentcloud_as_scaling_group",
		F:    testSweepAsScalingGroups,
	})
}

func testSweepAsScalingGroups(region string) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("geting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(TencentCloudClient)

	asService := AsService{
		client: client.apiV3Conn,
	}
	scalingGroups, err := asService.DescribeAutoScalingGroupByFilter(ctx, "", "", "")
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
				),
			},
		},
	})
}

func testAccCheckAsScalingGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling group id is not set")
		}
		asService := AsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, err := asService.DescribeAutoScalingGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAsScalingGroupDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	asService := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_scaling_group" {
			continue
		}

		_, err := asService.DescribeAutoScalingGroupById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("auto scaling group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsScalingGroup_basic() string {
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
  configuration_name = "tf-as-configuration-basic"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-group-basic"
  configuration_id   = "${tencentcloud_as_scaling_config.launch_configuration.id}"
  max_size           = 1
  min_size           = 0
  vpc_id             = "${tencentcloud_vpc.vpc.id}"
  subnet_ids         = ["${tencentcloud_subnet.subnet.id}"]
}
`
}

func testAccAsScalingGroup_full() string {
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
  configuration_name = "tf-as-configuration-full"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-group-full"
  configuration_id     = "${tencentcloud_as_scaling_config.launch_configuration.id}"
  max_size             = 1
  min_size             = 0
  vpc_id               = "${tencentcloud_vpc.vpc.id}"
  subnet_ids           = ["${tencentcloud_subnet.subnet.id}"]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = 1
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"
}
`
}

func testAccAsScalingGroup_update() string {
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
  configuration_name = "tf-as-configuration-full"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-group-update"
  configuration_id     = "${tencentcloud_as_scaling_config.launch_configuration.id}"
  max_size             = 2
  min_size             = 0
  vpc_id               = "${tencentcloud_vpc.vpc.id}"
  subnet_ids           = ["${tencentcloud_subnet.subnet.id}"]
  project_id           = 0
  default_cooldown     = 300
  desired_capacity     = 0
  termination_policies = ["OLDEST_INSTANCE"]
  retry_policy         = "IMMEDIATE_RETRY"
}
`
}
