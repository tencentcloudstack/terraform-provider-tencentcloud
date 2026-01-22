package apm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudApmApplicationConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApmApplicationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_application_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_application_config.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_application_config.example", "service_name"),
				),
			},
			{
				Config: testAccApmApplicationConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_application_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_application_config.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_application_config.example", "service_name"),
				),
			},
			{
				ResourceName:      "tencentcloud_apm_application_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApmApplicationConfig = `
resource "tencentcloud_apm_application_config" "example" {
  instance_id                           = tencentcloud_apm_instance.example.id
  service_name                          = "java-order-serive"
  url_convergence_switch                = 1
  agent_enable                          = true
  disable_cpu_used                      = 90
  disable_memory_used                   = 90
  enable_dashboard_config               = false
  enable_log_config                     = false
  enable_security_config                = false
  enable_snapshot                       = false
  event_enable                          = false
  is_delete_any_file_analysis           = 0
  is_deserialization_analysis           = 0
  is_directory_traversal_analysis       = 0
  is_expression_injection_analysis      = 0
  is_include_any_file_analysis          = 0
  is_instrumentation_vulnerability_scan = 1
  is_jndi_injection_analysis            = 0
  is_jni_injection_analysis             = 0
  is_memory_hijacking_analysis          = 0
  is_read_any_file_analysis             = 0
  is_related_dashboard                  = 0
  is_related_log                        = 0
  is_remote_command_execution_analysis  = 0
  is_script_engine_injection_analysis   = 0
  is_sql_injection_analysis             = 0
  is_template_engine_injection_analysis = 0
  is_upload_any_file_analysis           = 0
  is_webshell_backdoor_analysis         = 0
  log_index_type                        = 0
  log_source                            = "CLS"
  snapshot_timeout                      = 2000
  trace_squash                          = true
  url_auto_convergence_enable           = false
  url_convergence_threshold             = 1000
  url_long_segment_threshold            = 40
  url_number_segment_threshold          = 5

  agent_operation_config_view {
    retention_valid = false
  }

  instrument_list {
    enable = true
    name   = "apm-spring-annotations"
  }

  instrument_list {
    enable = true
    name   = "dubbo"
  }

  instrument_list {
    enable = true
    name   = "googlehttpclient"
  }

  instrument_list {
    enable = true
    name   = "grpc"
  }

  instrument_list {
    enable = true
    name   = "httpclient3"
  }

  instrument_list {
    enable = true
    name   = "httpclient4"
  }

  instrument_list {
    enable = true
    name   = "hystrix"
  }

  instrument_list {
    enable = true
    name   = "lettuce"
  }

  instrument_list {
    enable = true
    name   = "mongodb"
  }

  instrument_list {
    enable = true
    name   = "mybatis"
  }

  instrument_list {
    enable = true
    name   = "mysql"
  }

  instrument_list {
    enable = true
    name   = "okhttp"
  }

  instrument_list {
    enable = true
    name   = "redis"
  }

  instrument_list {
    enable = true
    name   = "rxjava"
  }

  instrument_list {
    enable = true
    name   = "spring-webmvc"
  }

  instrument_list {
    enable = true
    name   = "tomcat"
  }
}
`

const testAccApmApplicationConfigUpdate = `
resource "tencentcloud_apm_application_config" "example" {
  instance_id                           = tencentcloud_apm_instance.example.id
  service_name                          = "java-order-serive"
  url_convergence_switch                = 0
  agent_enable                          = true
  disable_cpu_used                      = 90
  disable_memory_used                   = 90
  enable_dashboard_config               = false
  enable_log_config                     = false
  enable_security_config                = false
  enable_snapshot                       = false
  event_enable                          = false
  is_delete_any_file_analysis           = 0
  is_deserialization_analysis           = 0
  is_directory_traversal_analysis       = 0
  is_expression_injection_analysis      = 0
  is_include_any_file_analysis          = 0
  is_instrumentation_vulnerability_scan = 1
  is_jndi_injection_analysis            = 0
  is_jni_injection_analysis             = 0
  is_memory_hijacking_analysis          = 0
  is_read_any_file_analysis             = 0
  is_related_dashboard                  = 0
  is_related_log                        = 0
  is_remote_command_execution_analysis  = 0
  is_script_engine_injection_analysis   = 0
  is_sql_injection_analysis             = 0
  is_template_engine_injection_analysis = 0
  is_upload_any_file_analysis           = 0
  is_webshell_backdoor_analysis         = 0
  log_index_type                        = 0
  log_source                            = "CLS"
  snapshot_timeout                      = 2000
  trace_squash                          = true
  url_auto_convergence_enable           = false
  url_convergence_threshold             = 1000
  url_long_segment_threshold            = 40
  url_number_segment_threshold          = 5

  agent_operation_config_view {
    retention_valid = false
  }

  instrument_list {
    enable = true
    name   = "apm-spring-annotations"
  }

  instrument_list {
    enable = true
    name   = "dubbo"
  }

  instrument_list {
    enable = true
    name   = "googlehttpclient"
  }

  instrument_list {
    enable = true
    name   = "grpc"
  }

  instrument_list {
    enable = true
    name   = "httpclient3"
  }

  instrument_list {
    enable = true
    name   = "httpclient4"
  }

  instrument_list {
    enable = true
    name   = "hystrix"
  }

  instrument_list {
    enable = true
    name   = "lettuce"
  }

  instrument_list {
    enable = true
    name   = "mongodb"
  }

  instrument_list {
    enable = true
    name   = "mybatis"
  }

  instrument_list {
    enable = true
    name   = "mysql"
  }

  instrument_list {
    enable = true
    name   = "okhttp"
  }

  instrument_list {
    enable = true
    name   = "redis"
  }

  instrument_list {
    enable = true
    name   = "rxjava"
  }

  instrument_list {
    enable = true
    name   = "spring-webmvc"
  }

  instrument_list {
    enable = true
    name   = "tomcat"
  }
}
`
