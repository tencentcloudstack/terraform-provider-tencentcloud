package as_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_as_scaling_config
	resource.AddTestSweepers("tencentcloud_as_scaling_config", &resource.Sweeper{
		Name: "tencentcloud_as_scaling_config",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta)
			asService := svcas.NewAsService(client.GetAPIV3Conn())
			configs, err := asService.DescribeLaunchConfigurationByFilter(ctx, "", "")
			if err != nil {
				return err
			}
			for _, config := range configs {
				instanceName := *config.LaunchConfigurationName
				now := time.Now()
				createTime := tccommon.StringToTime(*config.CreatedTime)
				interval := now.Sub(createTime).Minutes()

				if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
					continue
				}

				if tccommon.NeedProtect == 1 && int64(interval) < 30 {
					continue
				}

				ee := asService.DeleteLaunchConfiguration(ctx, *config.LaunchConfigurationId)
				if ee != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudAsScalingConfig_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "configuration_name", "tf-as-basic"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "image_id", "img-2lr9q49h"),
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

func TestAccTencentCloudAsScalingConfig_loginSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfigLoginSetting("skey-58gbxolb"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.login_setting"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.login_setting", "key_ids.0", "skey-58gbxolb"),
				),
			},
			{
				Config: testAccAsScalingConfigLoginSetting("skey-i55cwgvl"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.login_setting"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.login_setting", "key_ids.0", "skey-i55cwgvl"),
				),
			},
		},
	})
}

func TestAccTencentCloudAsScalingConfig_full(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfig_full(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "configuration_name", "tf-as-full"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "image_id", "img-2lr9q49h"),
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
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "host_name_settings.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "host_name_settings.0.host_name", "host-name"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "host_name_settings.0.host_name_style", "ORIGINAL"),
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
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "public_ip_assigned", "true"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "enhanced_security_service", "true"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "enhanced_monitor_service", "true"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "user_data", "dGVzdA=="),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "host_name_settings.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "host_name_settings.0.host_name", "host-name-test"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "host_name_settings.0.host_name_style", "UNIQUE"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_tags.tag", "as"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_tags.test", "update"),
				),
			},
		},
	})
}

func TestAccTencentCloudAsScalingConfig_charge(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAsScalingConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScalingConfig_charge(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_charge_type", "POSTPAID_BY_HOUR"),
				),
			},
			{
				Config: testAccAsScalingConfig_charge_spot(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_charge_type", "SPOTPAID"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "spot_instance_type", "one-time"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "spot_max_price", "1000"),
				),
			},
			{
				Config: testAccAsScalingConfig_charge_prepaid(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsScalingConfigExists("tencentcloud_as_scaling_config.launch_configuration"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_charge_type_prepaid_period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_as_scaling_config.launch_configuration", "instance_charge_type_prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
		},
	})
}

func testAccCheckAsScalingConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling configuration %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling configuration id is not set")
		}
		asService := svcas.NewAsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := asService.DescribeLaunchConfigurationById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has < 1 {
			return fmt.Errorf("auto scaling configuration not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckAsScalingConfigDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	asService := svcas.NewAsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_scaling_config" {
			continue
		}

		_, has, err := asService.DescribeLaunchConfigurationById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has > 0 {
			return fmt.Errorf("auto scaling configuration still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsScalingConfig_basic() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
	configuration_name = "tf-as-basic"
	image_id = "img-2lr9q49h"
	instance_types = ["SA1.SMALL1"]
}
	`
}

func testAccAsScalingConfig_full() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-full"
  image_id           = "img-2lr9q49h"
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

  host_name_settings {
	host_name       = "host-name"
	host_name_style = "ORIGINAL"
  }
  
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
  image_id           = "img-2lr9q49h"
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
  public_ip_assigned         = true
  password                   = "test123#"
  enhanced_security_service  = true
  enhanced_monitor_service   = true
  user_data                  = "dGVzdA=="

  host_name_settings {
	host_name       = "host-name-test"
	host_name_style = "UNIQUE"
  }
  
  instance_tags = {
    tag  = "as"
    test = "update"
  }
  
}
	`
}

func testAccAsScalingConfig_charge() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
	configuration_name = "tf-as-basic-charge"
	image_id = "img-2lr9q49h"
	instance_types = ["SA1.SMALL1"]
	instance_charge_type = "POSTPAID_BY_HOUR"
}
	`
}

func testAccAsScalingConfig_charge_spot() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
	configuration_name = "tf-as-basic-charge-spot"
	image_id = "img-2lr9q49h"
	instance_types = ["SA1.SMALL1"]
	instance_charge_type = "SPOTPAID"
	spot_instance_type = "one-time"
	spot_max_price = "1000"
}
	`
}

func testAccAsScalingConfig_charge_prepaid() string {
	return `
resource "tencentcloud_as_scaling_config" "launch_configuration" {
	configuration_name = "tf-as-basic-charge-prepaid"
	image_id = "img-2lr9q49h"
	instance_types = ["SA1.SMALL1"]
	instance_charge_type = "PREPAID"
	instance_charge_type_prepaid_period = 1
	instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}
	`
}

func testAccAsScalingConfigLoginSetting(keyId string) string {
	return fmt.Sprintf(`
resource "tencentcloud_as_scaling_config" "login_setting" {
  configuration_name   = "test-as-login-setting"
  image_id = "img-2lr9q49h"
  instance_types = ["S5.SMALL2"]
  instance_charge_type = "POSTPAID_BY_HOUR"
  key_ids              = ["%s"]
}
`, keyId)
}
