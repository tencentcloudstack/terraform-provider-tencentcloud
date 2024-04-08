package tke_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctke "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

const DefaultAddonName = "cos"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_kubernetes_addon_attachment
	resource.AddTestSweepers("tencentcloud_kubernetes_addon_attachment", &resource.Sweeper{
		Name: "tencentcloud_kubernetes_addon_attachment",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svctke.NewTkeService(client)

			cls, err := service.DescribeClusters(ctx, "", "keep")
			if err != nil {
				return err
			}

			if len(cls) == 0 {
				log.Println("no persistent cluster")
				return nil
			}

			for _, c := range cls {
				clusterId := c.ClusterId
				if err = service.DeleteExtensionAddon(ctx, clusterId, DefaultAddonName); err != nil {
					if e, ok := err.(*errors.TencentCloudSDKError); ok {
						// suppress the not found error when cos doesn't exist
						if strings.Contains(e.GetMessage(), "application cos not found") {
							continue
						}
					}
					return err
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudKubernetesAddonAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
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

func TestAccTencentCloudKubernetesAddonAttachmentResource_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeAddonAttachmentCos_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_attachment.cos", "response_body"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_attachment.cos", "name", "cos"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_attachment.cos", "version"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_attachment.cos", "request_body"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_attachment.cos", "version", "1.0.2"),
				),
			},
			{
				Config: testAccTkeAddonAttachmentCos_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_attachment.cos", "response_body"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_attachment.cos", "name", "cos"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_attachment.cos", "version", "1.0.3"),
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
`, tcacctest.TkeDataSource, DefaultAddonName)
}

func testAccTkeAddonAttachmentCos_basic() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_kubernetes_addon_attachment" "cos" {
  cluster_id = local.cluster_id
  name = "%s"
  request_body = jsonencode({
	kind = "App"
	spec = {
	  chart = {
		chartName    = "cos"
		chartVersion = "1.0.2"
	  }
	  values = {
		values        = []
		rawValues     = "e30="
		rawValuesType = "json"
	  }
	}
  })
}
`, tcacctest.TkeDataSource, DefaultAddonName)
}

func testAccTkeAddonAttachmentCos_update() string {
	return fmt.Sprintf(`
%s

resource "tencentcloud_kubernetes_addon_attachment" "cos" {
  cluster_id = local.cluster_id
  name = "%s"
  version = "1.0.3"
}
`, tcacctest.TkeDataSource, DefaultAddonName)
}
