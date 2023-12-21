package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapDomainErrorPageInfosDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapDomainErrorPageInfosDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_domain_error_page_infos.domain_error_page_infos"),
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
