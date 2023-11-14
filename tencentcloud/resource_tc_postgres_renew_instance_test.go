package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresRenewInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresRenewInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_renew_instance.renew_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_renew_instance.renew_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresRenewInstance = `

resource "tencentcloud_postgres_renew_instance" "renew_instance" {
  d_b_instance_id = "postgres-6fego161"
  period = 12
  auto_voucher = 0
  voucher_ids = &lt;nil&gt;
  tags = {
    "createdBy" = "terraform"
  }
}

`
