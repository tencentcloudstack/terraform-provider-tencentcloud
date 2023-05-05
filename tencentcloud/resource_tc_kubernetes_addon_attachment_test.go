package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func TestAccTencentCloudKubernetesAddonAttachmentResource(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_attachment.cos", "version"),
				),
			},
		},
	})
}

func testAccTkeAddonAttachment() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_kubernetes_addon_attachment" "cos" {
  cluster_id = local.cluster_id
  name = "%s"
}
`, TkeDataSource, defaultAddonName)
}
