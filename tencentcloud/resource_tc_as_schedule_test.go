package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudAsSchedule(t *testing.T) {
	startTime := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
	endTime := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsSchedule(startTime, endTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScheduleExists("tencentcloud_as_schedule.schedule"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_schedule.schedule", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "schedule_action_name", "tf-as-schedule"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "max_size", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "min_size", "0"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "desired_capacity", "0"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "start_time", startTime),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "end_time", endTime),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "recurrence", "0 0 */1 * *"),
				),
			},
			// test update case
			{
				Config: testAccAsSchedule_update(startTime, endTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScheduleExists("tencentcloud_as_schedule.schedule"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_schedule.schedule", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "schedule_action_name", "tf-as-schedule-update"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "max_size", "2"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "desired_capacity", "0"),
					resource.TestCheckResourceAttr("tencentcloud_as_schedule.schedule", "recurrence", "1 1 */1 * *"),
				),
			},
		},
	})
}

func testAccCheckAsScheduleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling schedule %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling schedule id is not set")
		}
		asService := AsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, err := asService.DescribeScheduledActionById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAsScheduleDestroy(s *terraform.State) error {
	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	asService := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_schedule" {
			continue
		}

		_, err := asService.DescribeScheduledActionById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("auto scaling schedule still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsSchedule(startTime, endTime string) string {
	return fmt.Sprintf(`
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

resource "tencentcloud_as_schedule" "schedule" {
  scaling_group_id     = "${tencentcloud_as_scaling_group.scaling_group.id}"
  schedule_action_name = "tf-as-schedule"
  max_size             = 1
  min_size             = 0
  desired_capacity     = 0
  start_time           = "%s"
  end_time             = "%s"
  recurrence           = "0 0 */1 * *"
}
`, startTime, endTime)
}

func testAccAsSchedule_update(startTime, endTime string) string {
	return fmt.Sprintf(`
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

resource "tencentcloud_as_schedule" "schedule" {
  scaling_group_id     = "${tencentcloud_as_scaling_group.scaling_group.id}"
  schedule_action_name = "tf-as-schedule-update"
  max_size             = 2
  min_size             = 0
  desired_capacity     = 0
  start_time           = "%s"
  end_time             = "%s"
  recurrence           = "1 1 */1 * *"
}
`, startTime, endTime)
}
