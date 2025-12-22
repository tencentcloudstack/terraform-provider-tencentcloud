package tcr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTcrReplicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTcrReplication,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_tcr_replication.example", "id"),
			),
		}},
	})
}

const testAccTcrReplication = `
resource "tencentcloud_tcr_replication" "example" {
  source_registry_id      = "tcr-9q9h1nof"
  destination_registry_id = "tcr-jtih9ngc"
  rule {
    name           = "tf-example"
    dest_namespace = ""
    override       = true
    deletion       = true
    filters {
      type  = "name"
      value = "tf-example/**"
    }
  }

  destination_region_id = 1
  description           = "remark."
}
`
