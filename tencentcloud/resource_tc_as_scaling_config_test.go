package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudAsScalingConfig_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "configuration_name", "tf-as-basic"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "image_id", "img-9qabwvbn"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_types.0", "SA1.SMALL1"),
				),
			},
			{
				ResourceName:            "tencentcloud_as_scaling_config.launch_configuration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccTencentCloudAsScalingConfig_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfig_full(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "configuration_name", "tf-as-full"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "image_id", "img-9qabwvbn"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_types.0", "SA1.SMALL1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "system_disk_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "data_disk.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "data_disk.0.disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "data_disk.0.disk_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "internet_charge_type", "TRAFFIC_POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "internet_max_bandwidth_out", "10"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "public_ip_assigned", "true"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "enhanced_security_service", "false"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "enhanced_monitor_service", "false"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "user_data", "test"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_tags.tag", "as"),
				),
			},
			// todo: update test
			{
				Config: testAccAsScalingConfig_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "configuration_name", "tf-as-full-update"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_types.0", "S4.SMALL2"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "system_disk_size", "60"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "data_disk.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "data_disk.0.disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "data_disk.0.disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "internet_max_bandwidth_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "public_ip_assigned", "false"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "enhanced_security_service", "true"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "enhanced_monitor_service", "true"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "user_data", "dGVzdA=="),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_tags.tag", "as"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_tags.test", "update"),
				),
			},
		},
	})
}

func testAccCheckAsScalingConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling configuration %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling configuration id is not set")
		}
		asService := AsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, _, err := asService.DescribeLaunchConfigurationById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAsScalingConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	asService := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_scaling_config" {
			continue
		}

		_, _, err := asService.DescribeLaunchConfigurationById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("auto scaling configuration still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsScalingConfig_basic() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
	configuration_name = "tf-as-basic"
	image_id = "img-9qabwvbn"
	instance_types = ["SA1.SMALL1"]
}
	`
}

func testAccAsScalingConfig_full() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-full"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
  project_id         = 0
  system_disk_type   = "CLOUD_PREMIUM"
  system_disk_size   = "50"
  
  data_disk   {
    disk_type = "CLOUD_PREMIUM"
    disk_size = 50
  }
  
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 10
  public_ip_assigned         = true
  password                   = "test123#"
  enhanced_security_service  = false
  enhanced_monitor_service   = false
  user_data                  = "test"
  
  instance_tags = {
    tag = "as"
  }
  
}
	`
}

func testAccAsScalingConfig_update() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-full-update"
  image_id           = "img-9qabwvbn"
  instance_types     = ["S4.SMALL2"]
  project_id         = 0
  system_disk_type   = "CLOUD_PREMIUM"
  system_disk_size   = "60"
  
  data_disk   {
    disk_type = "CLOUD_SSD"
    disk_size = 100
  }
  
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 20
  public_ip_assigned         = false
  password                   = "test123#"
  enhanced_security_service  = true
  enhanced_monitor_service   = true
  user_data                  = "dGVzdA=="
  
  instance_tags = {
    tag  = "as"
    test = "update"
  }
  
}
	`
}
