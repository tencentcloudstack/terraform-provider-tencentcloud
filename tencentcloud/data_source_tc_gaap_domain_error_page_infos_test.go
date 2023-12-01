package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapDomainErrorPageInfosDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfosDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_domain_error_page_infos.domain_error_page_infos"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_domain_error_page_infos.domain_error_page_infos", "error_page_set.#"),
				),
			},
		},
	})
}

const testAccGaapDomainErrorPageInfosDataSource = `
data "tencentcloud_gaap_domain_error_page_infos" "domain_error_page_infos" {
	error_page_ids = ["errorPage-mh4k07v5"]
}
`
