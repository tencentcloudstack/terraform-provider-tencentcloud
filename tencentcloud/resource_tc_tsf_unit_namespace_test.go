package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfUnitNamespaceResource_basic -v
func TestAccTencentCloudTsfUnitNamespaceResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfUnitNamespace,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfUnitNamespaceExists("tencentcloud_tsf_unit_namespace.unit_namespace"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_unit_namespace.unit_namespace", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_namespace.unit_namespace", "gateway_instance_id", defaultTsfGateway),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_namespace.unit_namespace", "namespace_id", defaultTsfGWNamespaceId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_unit_namespace.unit_namespace", "namespace_name", "keep-terraform-cls"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_unit_namespace.unit_namespace", "created_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_unit_namespace.unit_namespace", "updated_time"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_unit_namespace.unit_namespace",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfUnitNamespaceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_unit_namespace" {
			continue
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		gatewayInstanceId := ids[0]
		namespaceId := ids[1]

		res, err := service.DescribeTsfUnitNamespaceById(ctx, gatewayInstanceId, namespaceId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf unitNamespace %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfUnitNamespaceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		gatewayInstanceId := ids[0]
		namespaceId := ids[1]

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfUnitNamespaceById(ctx, gatewayInstanceId, namespaceId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf unitNamespace %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfUnitNamespaceVar = `
variable "gateway_instance_id" {
	default = "` + defaultTsfGateway + `"
}

variable "namespace_id" {
	default = "` + defaultTsfGWNamespaceId + `"
}
`

const testAccTsfUnitNamespace = testAccTsfUnitNamespaceVar + `

resource "tencentcloud_tsf_unit_namespace" "unit_namespace" {
	gateway_instance_id = var.gateway_instance_id
	namespace_id = var.namespace_id
	namespace_name = "keep-terraform-cls"
}

`
