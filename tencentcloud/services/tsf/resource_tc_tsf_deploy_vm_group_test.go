package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfDeployVmGroupResource_basic -v
func TestAccTencentCloudTsfDeployVmGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDeployVmGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "deploy_batch.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "deploy_beta_enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "deploy_desc"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "deploy_exe_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "deploy_wait_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "enable_health_check"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "force_start"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "jdk_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "jdk_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "pkg_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "startup_parameters"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "update_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.action_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.failure_threshold"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.initial_delay_seconds"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.path"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.period_seconds"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.port"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.scheme"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.success_threshold"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_vm_group.deploy_vm_group", "health_check_settings.0.readiness_probe.0.timeout_seconds"),
				),
			},
		},
	})
}

const testAccTsfDeployVmGroup = `

resource "tencentcloud_tsf_deploy_vm_group" "deploy_vm_group" {
	group_id            = "group-vzd97zpy"
	pkg_id              = "pkg-131bc1d3"
	startup_parameters  = "-Xms128m -Xmx512m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=512m"
	deploy_desc         = "deploy test"
	force_start         = false
	enable_health_check = true
	health_check_settings {
	  readiness_probe {
		action_type           = "HTTP"
		initial_delay_seconds = 10
		timeout_seconds       = 2
		period_seconds        = 10
		success_threshold     = 1
		failure_threshold     = 3
		scheme                = "HTTP"
		port                  = "80"
		path                  = "/"
	  }
	}
	update_type = 0
	jdk_name    = "konaJDK"
	jdk_version = "8"
}

`
