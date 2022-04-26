package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const defaultAddonName = "cos"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_kubernetes_addon_attachment
	resource.AddTestSweepers("tencentcloud_kubernetes_addon_attachment", &resource.Sweeper{
		Name: "tencentcloud_kubernetes_addon_attachment",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := TkeService{client}

			cls, err := service.DescribeClusters(ctx, "", "keep")
			if err != nil {
				return err
			}

			if len(cls) == 0 {
				return fmt.Errorf("no persistent cluster")
			}

			clusterId := cls[0].ClusterId

			if err = service.DeleteExtensionAddon(ctx, clusterId, defaultAddonName); err != nil {
				return err
			}

			return nil
		},
	})
}

func TestAccTencentCloudTkeAddonAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeAddonAttachment(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_attachment.cos", "response_body"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_attachment.cos", "name", "cos"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_attachment.cos", "version", "1.0.0"),
				),
			},
		},
	})
}

func testAccTkeAddonAttachment() string {
	return fmt.Sprintf(`

data "tencentcloud_kubernetes_clusters" "cls" {
  cluster_name = "keep"
}

resource "tencentcloud_kubernetes_addon_attachment" "cos" {
  cluster_id = data.tencentcloud_kubernetes_clusters.cls.list.0.cluster_id // tencentcloud_kubernetes_cluster.managed_cluster.id
  name = %s
  version = "1.0.0"
}
`, defaultAddonName)
}
