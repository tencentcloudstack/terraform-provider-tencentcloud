package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodNeedFixModifyDomainOwnerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodModifyDomainOwner,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_modify_domain_owner_operation.modify_domain_owner", "domain", "terraform.com"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_modify_domain_owner_operation.modify_domain_owner", "account", "xxxxxx"),
					// resource.TestCheckResourceAttr("tencentcloud_dnspod_modify_domain_owner_operation.modify_domain_owner", "domain_id", "123"),
				),
			},
		},
	})
}

const testAccDnspodModifyDomainOwner = `

resource "tencentcloud_dnspod_modify_domain_owner_operation" "modify_domain_owner" {
  domain = "terraform.com"
  account = "xxxxxx"
  # domain_id = "123"
}

`
