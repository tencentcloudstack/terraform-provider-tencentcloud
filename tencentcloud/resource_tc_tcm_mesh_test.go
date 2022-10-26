package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTcmMesh_basic -v
func TestAccTencentCloudTcmMesh_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMeshDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmMesh,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMeshExists("tencentcloud_tcm_mesh.basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_mesh.basic", "display_name", "test_mesh"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_mesh.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMeshDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcm_mesh" {
			continue
		}

		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if err != nil {
			if isExpectError(err, []string{"ResourceNotFound"}) {
				return nil
			}
		}
		if mesh != nil {
			return fmt.Errorf("tcm mesh %v still exists", *mesh.Mesh.State)
		}
	}
	return nil
}

func testAccCheckMeshExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		mesh, err := service.DescribeTcmMesh(ctx, rs.Primary.ID)
		if mesh.Mesh == nil {
			return fmt.Errorf("tcm mesh %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTcmMesh = `

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
  }
  tag_list {
    key = "key"
    value = "value"
    passthrough = false
  }
}

`
