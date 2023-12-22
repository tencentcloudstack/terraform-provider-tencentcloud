package tcm_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTcmTracingConfigResource_basic -v
func TestAccTencentCloudTcmTracingConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmTracingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcmTracingConfigExists("tencentcloud_tcm_tracing_config.tracing_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_tracing_config.tracing_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_tracing_config.tracing_config", "mesh_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_tracing_config.tracing_config", "sampling", "30"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_tracing_config.tracing_config", "apm.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_tracing_config.tracing_config", "apm.0.enable", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_tracing_config.tracing_config", "zipkin.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_tracing_config.tracing_config", "zipkin.0.address", "10.10.10.10:9411"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_tracing_config.tracing_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTcmTracingConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctcm.NewTcmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if mesh.Mesh.Config.Tracing == nil {
			return fmt.Errorf("tcm tracing %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTcmTracingConfig = `

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

resource "tencentcloud_tcm_tracing_config" "tracing_config" {
	mesh_id = tencentcloud_tcm_mesh.basic.id
	enable = true
	apm {
	  enable = false
	  # region = "ap-guangzhou"
	  # instance_id = "apm-kSy0jYxxx"
	}
	sampling = 30
	zipkin {
		address = "10.10.10.10:9411"
	}
}

`
