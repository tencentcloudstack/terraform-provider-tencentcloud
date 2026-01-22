package billing_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBillingAllocationTagResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingAllocationTag,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_billing_allocation_tag.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_billing_allocation_tag.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBillingAllocationTag = `
resource "tencentcloud_tag" "example" {
  tag_key   = "tagKey"
  tag_value = "tagValue"
}

resource "tencentcloud_billing_allocation_tag" "example" {
  tag_key = tencentcloud_tag.example.tag_key
}
`
