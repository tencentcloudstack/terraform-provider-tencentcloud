package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfDeployContainerGroupResource_basic -v
func TestAccTencentCloudTsfDeployContainerGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDeployContainerGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "tag_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "instance_num"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "server"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "reponame"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "cpu_limit"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "mem_limit"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "jvm_opts"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "cpu_request"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "mem_request"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "do_not_start"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "update_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "update_ivl"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "agent_cpu_request"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "agent_cpu_limit"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "agent_mem_request"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "agent_mem_limit"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "max_surge"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "max_unavailable"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.action_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.initial_delay_seconds"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.timeout_seconds"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.period_seconds"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.success_threshold"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.failure_threshold"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.scheme"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.port"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.path"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "health_check_settings.0.readiness_probe.0.type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.access_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.protocol_ports.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.protocol_ports.0.protocol"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.protocol_ports.0.port"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.protocol_ports.0.target_port"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.protocol_ports.0.node_port"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.disable_service"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.headless_service"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.allow_delete_service"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.open_session_affinity"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "service_setting.0.session_affinity_timeout_seconds"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "deploy_agent"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "scheduling_strategy.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "scheduling_strategy.0.type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "repo_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "volume_clean"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "warmup_setting.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_container_group.deploy_container_group", "warmup_setting.0.enabled"),
				),
			},
		},
	})
}

const testAccTsfDeployContainerGroupVar = `
variable "group_id" {
	default = "` + tcacctest.DefaultTsfContainerGroupId + `"
}
`

const testAccTsfDeployContainerGroup = testAccTsfDeployContainerGroupVar + `

resource "tencentcloud_tsf_deploy_container_group" "deploy_container_group" {
	group_id          = var.group_id
	cpu_request       = "0.25"
	mem_request       = "640"
	server            = "ccr.ccs.tencentyun.com"
	reponame          = "tsf_100011913960/terraform"
	tag_name          = "terraform-only-1"
	do_not_start      = false
	instance_num      = 1
	update_type       = 1
	update_ivl        = 10
	mem_limit         = "1280"
	cpu_limit         = "0.5"
	agent_cpu_request = "0.1"
	agent_cpu_limit   = "0.2"
	agent_mem_request = "125"
	agent_mem_limit   = "400"
	max_surge         = "25%"
	max_unavailable   = "0"
	service_setting {
		access_type = 1
		protocol_ports {
			protocol    = "TCP"
			port        = 18081
			target_port = 18081
			node_port   = 30001
		}
		subnet_id						 = ""
		disable_service                  = false
		headless_service                 = false
		allow_delete_service             = true
		open_session_affinity            = false
		session_affinity_timeout_seconds = 10800
	
	}
	health_check_settings {
		readiness_probe {
			action_type           = "TCP"
			initial_delay_seconds = 0
			timeout_seconds       = 3
			period_seconds        = 30
			success_threshold     = 1
			failure_threshold     = 3
			scheme                = "HTTP"
			port                  = 80
			path                  = "/"
			type                  = "TSF_DEFAULT"
		}
	}
	scheduling_strategy {
		type = "NONE"
	}
	deploy_agent = true
	repo_type = "personal"
	volume_clean = false
	jvm_opts          = "-Xms128m -Xmx512m -XX:MetaspaceSize=128m -XX:MaxMetaspaceSize=512m"
	warmup_setting {
		enabled = false
	}
}

`
