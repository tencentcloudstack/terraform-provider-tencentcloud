package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTcmAccessLogConfigResource_basic -v
func TestAccTencentCloudTcmAccessLogConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrometheusAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmAccessLogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcmAccessLogConfigExists("tencentcloud_tcm_access_log_config.access_log_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_access_log_config.access_log_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_access_log_config.access_log_config", "mesh_name"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "address", "10.0.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "enable_server", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "enable_stdout", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "encoding", "JSON"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "format", "{\n\t\"authority\": \"%REQ(:AUTHORITY)%\",\n\t\"bytes_received\": \"%BYTES_RECEIVED%\",\n\t\"bytes_sent\": \"%BYTES_SENT%\",\n\t\"downstream_local_address\": \"%DOWNSTREAM_LOCAL_ADDRESS%\",\n\t\"downstream_remote_address\": \"%DOWNSTREAM_REMOTE_ADDRESS%\",\n\t\"duration\": \"%DURATION%\",\n\t\"istio_policy_status\": \"%DYNAMIC_METADATA(istio.mixer:status)%\",\n\t\"method\": \"%REQ(:METHOD)%\",\n\t\"path\": \"%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%\",\n\t\"protocol\": \"%PROTOCOL%\",\n\t\"request_id\": \"%REQ(X-REQUEST-ID)%\",\n\t\"requested_server_name\": \"%REQUESTED_SERVER_NAME%\",\n\t\"response_code\": \"%RESPONSE_CODE%\",\n\t\"response_flags\": \"%RESPONSE_FLAGS%\",\n\t\"route_name\": \"%ROUTE_NAME%\",\n\t\"start_time\": \"%START_TIME%\",\n\t\"upstream_cluster\": \"%UPSTREAM_CLUSTER%\",\n\t\"upstream_host\": \"%UPSTREAM_HOST%\",\n\t\"upstream_local_address\": \"%UPSTREAM_LOCAL_ADDRESS%\",\n\t\"upstream_service_time\": \"%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%\",\n\t\"upstream_transport_failure_reason\": \"%UPSTREAM_TRANSPORT_FAILURE_REASON%\",\n\t\"user_agent\": \"%REQ(USER-AGENT)%\",\n\t\"x_forwarded_for\": \"%REQ(X-FORWARDED-FOR)%\"\n}\n"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "template", "istio"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "cls.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "cls.0.enable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "selected_range.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_access_log_config.access_log_config", "selected_range.0.all", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_access_log_config.access_log_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTcmAccessLogConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if mesh.Mesh.Config.AccessLog == nil {
			return fmt.Errorf("tcm accessLog %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTcmAccessLogConfig = `

resource "tencentcloud_tcm_mesh" "basic" {
	display_name = "test_mesh"
	mesh_version = "1.12.5"
	type = "HOSTED"
	config {
	  istio {
		outbound_traffic_policy = "ALLOW_ANY"
		disable_policy_checks = true
		enable_pilot_http = true
		disable_http_retry = true
		smart_dns {
		  istio_meta_dns_capture = true
		  istio_meta_dns_auto_allocate = true
		}
	  }
	  tracing {
		  enable = true
		  sampling = 1
		  apm {
			  enable = false
		  }
		  zipkin {
			  address = "10.0.0.1:1000"
		  }
	  }
	  prometheus {
		  custom_prom {
			  url = "https://10.0.0.1:1000"
			  auth_type = "none"
			  vpc_id = "vpc-j9yhbzpn"
		  }
	  }
	}
	tag_list {
	  key = "key"
	  value = "value"
	  passthrough = false
	}
  }

resource "tencentcloud_tcm_access_log_config" "access_log_config" {
    address       = "10.0.0.1"
    enable        = true
    enable_server = true
    enable_stdout = true
    encoding      = "JSON"
    format        = "{\n\t\"authority\": \"%REQ(:AUTHORITY)%\",\n\t\"bytes_received\": \"%BYTES_RECEIVED%\",\n\t\"bytes_sent\": \"%BYTES_SENT%\",\n\t\"downstream_local_address\": \"%DOWNSTREAM_LOCAL_ADDRESS%\",\n\t\"downstream_remote_address\": \"%DOWNSTREAM_REMOTE_ADDRESS%\",\n\t\"duration\": \"%DURATION%\",\n\t\"istio_policy_status\": \"%DYNAMIC_METADATA(istio.mixer:status)%\",\n\t\"method\": \"%REQ(:METHOD)%\",\n\t\"path\": \"%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%\",\n\t\"protocol\": \"%PROTOCOL%\",\n\t\"request_id\": \"%REQ(X-REQUEST-ID)%\",\n\t\"requested_server_name\": \"%REQUESTED_SERVER_NAME%\",\n\t\"response_code\": \"%RESPONSE_CODE%\",\n\t\"response_flags\": \"%RESPONSE_FLAGS%\",\n\t\"route_name\": \"%ROUTE_NAME%\",\n\t\"start_time\": \"%START_TIME%\",\n\t\"upstream_cluster\": \"%UPSTREAM_CLUSTER%\",\n\t\"upstream_host\": \"%UPSTREAM_HOST%\",\n\t\"upstream_local_address\": \"%UPSTREAM_LOCAL_ADDRESS%\",\n\t\"upstream_service_time\": \"%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%\",\n\t\"upstream_transport_failure_reason\": \"%UPSTREAM_TRANSPORT_FAILURE_REASON%\",\n\t\"user_agent\": \"%REQ(USER-AGENT)%\",\n\t\"x_forwarded_for\": \"%REQ(X-FORWARDED-FOR)%\"\n}\n"
    mesh_name     = tencentcloud_tcm_mesh.basic.id
    template      = "istio"

    cls {
        enable  = false
        # log_set = "SCF_logset_NLCsDxxx"
        # topic   = "SCF_logtopic_rPWZpxxx"
    }

    selected_range {
        all = true
    }
}

`
