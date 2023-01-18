package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfNamespaceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfNamespace,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_namespace.namespace", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_namespace.namespace",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfNamespace = `

resource "tencentcloud_tsf_namespace" "namespace" {
  namespace_name = ""
  cluster_id = ""
  namespace_desc = ""
  namespace_resource_type = ""
  namespace_type = ""
  namespace_id = ""
  is_ha_enable = ""
  program_id = ""
    program_id_list = 
              }

`
