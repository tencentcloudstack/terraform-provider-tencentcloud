package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTCRNamespacesNameAll = "data.tencentcloud_tcr_namespaces.id_test"

func TestAccTencentCloudDataTCRNamespaces(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRNamespacesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRNamespaceExists("tencentcloud_tcr_namespace.mytcr_namespace"),
					resource.TestCheckResourceAttrSet(testDataTCRNamespacesNameAll, "namespace_list.0.name"),
					resource.TestCheckResourceAttr(testDataTCRNamespacesNameAll, "namespace_list.0.is_public", "false"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRNamespacesBasic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "standard"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "mytcr_namespace" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  name        = "test"
  is_public   = false
}

data "tencentcloud_tcr_namespaces" "id_test" {
  instance_id = tencentcloud_tcr_namespace.mytcr_namespace.instance_id
}
`
