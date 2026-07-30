package dbdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDbdcNodeToDbCustomClusterAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbdcNodeToDbCustomClusterAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example", "cluster_id", "dbcc-xxxxxxxx"),
					resource.TestCheckResourceAttr("tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example", "node_id", "dbcn-xxxxxxxx"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example", "status"),
				),
			},
			{
				ResourceName:            "tencentcloud_dbdc_node_to_db_custom_cluster_attachment.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_id", "login_settings"},
			},
		},
	})
}

const testAccDbdcNodeToDbCustomClusterAttachment = `
resource "tencentcloud_dbdc_node_to_db_custom_cluster_attachment" "example" {
  cluster_id = "dbcc-xxxxxxxx"
  node_id    = "dbcn-xxxxxxxx"
  image_id   = "img-xxxxxxxx"

  login_settings {
    password = "Passw0rd@2024"
  }
}
`
