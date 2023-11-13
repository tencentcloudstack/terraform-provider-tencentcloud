package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfUnitNamespaceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfUnitNamespace,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_unit_namespace.unit_namespace", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_unit_namespace.unit_namespace",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfUnitNamespace = `

resource "tencentcloud_tsf_unit_namespace" "unit_namespace" {
  gateway_instance_id = ""
  unit_namespace_list {
		namespace_id = ""
		namespace_name = ""
		id = ""
		gateway_instance_id = ""
		created_time = ""
		updated_time = ""

  }
}

`
