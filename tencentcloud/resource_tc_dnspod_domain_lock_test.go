package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodDomainLockResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDomainLock,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_lock.domain_lock", "id"),
				),
			},
		},
	})
}

const testAccDnspodDomainLock = `

resource "tencentcloud_dnspod_domain_lock" "domain_lock" {
  domain = "iac-tf.cloud"
  lock_days = 30
}

`
