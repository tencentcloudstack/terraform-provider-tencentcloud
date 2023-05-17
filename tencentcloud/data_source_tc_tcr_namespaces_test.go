package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTCRNamespacesNameAll = "data.tencentcloud_tcr_namespaces.id_test"

func TestAccTencentCloudTcrNamespacesData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config:    testAccTencentCloudDataTCRNamespacesBasic,
				PreConfig: func() { testAccStepSetRegion(t, "ap-shanghai") },
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataTCRNamespacesNameAll, "namespace_list.0.name"),
					resource.TestCheckResourceAttrSet(testDataTCRNamespacesNameAll, "namespace_list.0.is_public"), // we only need to care whether the value is set or not, rather than the exact value itself, and this value of public cannot be confirmed when the e2e case parallel running
					resource.TestCheckResourceAttrSet(testDataTCRNamespacesNameAll, "namespace_list.0.id"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRNamespacesBasic = defaultTCRInstanceData + `
data "tencentcloud_tcr_namespaces" "id_test" {
  instance_id = local.tcr_id
}
`
