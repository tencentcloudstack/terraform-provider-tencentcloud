package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudAsScalingConfigsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfigsDataSource_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.configuration_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.configuration_name", "tf-as-config-basic"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.image_id", "img-9qabwvbn"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.instance_types.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.instance_types.0", "SA1.SMALL1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.configuration_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.configuration_name", "tf-as-config-basic"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.image_id", "img-9qabwvbn"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.instance_types.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.instance_types.0", "SA1.SMALL1"),
				),
			},
		},
	})
}

func TestAccTencentCloudAsScalingConfigsDataSource_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfigsDataSource_full(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.#", "1"),

					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.configuration_name", "tf-as-config-full"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.image_id", "img-9qabwvbn"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.instance_types.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.instance_types.0", "SA1.SMALL1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.project_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.system_disk_size", "50"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.data_disk.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.data_disk.0.disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.data_disk.0.disk_size", "50"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.internet_charge_type", "TRAFFIC_POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.internet_max_bandwidth_out", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.public_ip_assigned", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.enhanced_security_service", "false"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.enhanced_monitor_service", "false"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.user_data", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.instance_tags.tag", "as"),
				),
			},
		},
	})
}

func testAccAsScalingConfigsDataSource_basic() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-config-basic"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

data "tencentcloud_as_scaling_configs" "scaling_configs" {
  configuration_id = "${tencentcloud_as_scaling_config.launch_configuration.id}"
}

data "tencentcloud_as_scaling_configs" "scaling_configs_name" {
  configuration_name = "${tencentcloud_as_scaling_config.launch_configuration.configuration_name}"
}
`
}

func testAccAsScalingConfigsDataSource_full() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-config-full"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
  project_id         = 0
  system_disk_type   = "CLOUD_PREMIUM"
  system_disk_size   = "50"
  
  data_disk {
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

data "tencentcloud_as_scaling_configs" "scaling_configs" {
  configuration_id = "${tencentcloud_as_scaling_config.launch_configuration.id}"
}
`
}
