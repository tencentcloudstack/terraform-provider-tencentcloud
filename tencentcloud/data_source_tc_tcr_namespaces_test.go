package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTCRNamespacesNameAll = "data.tencentcloud_tcr_namespaces.id_test"

func TestAccTencentCloudDataTCRNamespaces(t *testing.T) {
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
					resource.TestCheckResourceAttr(testDataTCRNamespacesNameAll, "namespace_list.0.is_public", "false"),
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
