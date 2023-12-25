package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfUsableUnitNamespacesDataSource_basic -v
func TestAccTencentCloudTsfUsableUnitNamespacesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfUsableUnitNamespacesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_usable_unit_namespaces.usable_unit_namespaces"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_usable_unit_namespaces.usable_unit_namespaces", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_usable_unit_namespaces.usable_unit_namespaces", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_usable_unit_namespaces.usable_unit_namespaces", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_usable_unit_namespaces.usable_unit_namespaces", "result.0.content.0.namespace_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_usable_unit_namespaces.usable_unit_namespaces", "result.0.content.0.namespace_name"),
				),
			},
		},
	})
}

const testAccTsfUsableUnitNamespacesDataSource = `

data "tencentcloud_tsf_usable_unit_namespaces" "usable_unit_namespaces" {
  search_word = "test"
}

`
