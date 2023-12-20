package as_test

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var scalingConfigNameRE = regexp.MustCompile("tf-as-config-basic")
var scalingConfigNameFullRE = regexp.MustCompile("tf-as-config-full")

func TestAccTencentCloudAsScalingConfigsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfigsDataSource_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.configuration_id"),
					resource.TestMatchResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.configuration_name", scalingConfigNameRE),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.image_id", tcacctest.DefaultTkeOSImageId),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.instance_types.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.instance_types.0", "SA1.SMALL1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.configuration_id"),
					resource.TestMatchResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.configuration_name", scalingConfigNameRE),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.image_id", tcacctest.DefaultTkeOSImageId),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.instance_types.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs_name", "configuration_list.0.instance_types.0", "SA1.SMALL1"),
				),
			},
		},
	})
}

func TestAccTencentCloudAsScalingConfigsDataSource_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfigsDataSource_full(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.#", "1"),

					resource.TestMatchResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.configuration_name", scalingConfigNameFullRE),
					resource.TestCheckResourceAttr("data.tencentcloud_as_scaling_configs.scaling_configs", "configuration_list.0.image_id", tcacctest.DefaultTkeOSImageId),
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
	return fmt.Sprintf(`
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-config-basic-%d"
  image_id           = "%s"
  instance_types     = ["SA1.SMALL1"]
}

data "tencentcloud_as_scaling_configs" "scaling_configs" {
  configuration_id = tencentcloud_as_scaling_config.launch_configuration.id
}

data "tencentcloud_as_scaling_configs" "scaling_configs_name" {
  configuration_name = tencentcloud_as_scaling_config.launch_configuration.configuration_name
}
`, rand.Intn(1000), tcacctest.DefaultTkeOSImageId)
}

func testAccAsScalingConfigsDataSource_full() string {
	return fmt.Sprintf(`
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-config-full-%d"
  image_id           = "%s"
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
  configuration_id = tencentcloud_as_scaling_config.launch_configuration.id
}
`, rand.Intn(1000), tcacctest.DefaultTkeOSImageId)
}
