package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTCRNamespacesNameAll = "data.tencentcloud_tcr_namespaces.id_test"

func TestAccTencentCloudTCRNamespacesData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRNamespacesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataTCRNamespacesNameAll, "namespace_list.0.name"),
					resource.TestCheckResourceAttrSet(testDataTCRNamespacesNameAll, "namespace_list.0.is_public"),
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
