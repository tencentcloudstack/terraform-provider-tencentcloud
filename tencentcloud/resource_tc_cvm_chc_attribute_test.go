package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCvmChcAttributeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcAttribute,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_attribute.chc_attribute", "id")),
			},
		},
	})
}

const testAccCvmChcAttribute = `

resource "tencentcloud_cvm_chc_attribute" "chc_attribute" {
	chc_id = "chc-0brmw3wl"
	instance_name = "test"
}
`
