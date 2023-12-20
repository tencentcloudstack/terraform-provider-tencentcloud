package cdn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCdnDomainVerifyRecord_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainVerifyRecord,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "verify_result"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "file_verify_domains.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "file_verify_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "file_verify_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "record"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "record_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "sub_domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdn_domain_verifier.vr", "verify_result"),
				),
			},
		},
	})
}

const testAccCdnDomainVerifyRecord = `
data "tencentcloud_domains" "domains" {}

locals {
  domain1 = data.tencentcloud_domains.domains.list.0.domain_name
}

output "domain" {
  value = local.domain1
}

resource "tencentcloud_dnspod_record" "demo" {
  domain = local.domain1
  record_type = "A"
  record_line = "默认"
  value = "1.2.3.9"
  sub_domain="test"
}

data "tencentcloud_cdn_domain_verifier" "vr" {
  # test.xxxxxx.xxx
  domain = "${tencentcloud_dnspod_record.demo.sub_domain}.${local.domain1}"
  auto_verify = true
}
`
