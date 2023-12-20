package cvm_test

import (
	"context"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cvm_sync_image
	resource.AddTestSweepers("tencentcloud_cvm_sync_image", &resource.Sweeper{
		Name: "tencentcloud_cvm_sync_image",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			request := cvm.NewDescribeImagesRequest()
			response, err := client.UseCvmClient().DescribeImages(request)
			if err != nil {
				return err
			}
			for _, image := range response.Response.ImageSet {
				imageName := *image.ImageName
				imageId := *image.ImageId

				now := time.Now()

				createTime := tccommon.StringToTime(*image.CreatedTime)
				interval := now.Sub(createTime).Minutes()
				if strings.HasPrefix(imageName, tcacctest.KeepResource) || strings.HasPrefix(imageName, tcacctest.DefaultResource) {
					continue
				}
				// less than 30 minute, not delete
				if tccommon.NeedProtect == 1 && int64(interval) < 30 {
					continue
				}

				service := svccvm.NewCvmService(client)
				if err := service.DeleteImage(ctx, imageId); err != nil {
					continue
				}
				return nil
			}

			return nil
		},
	})
}

func TestAccTencentCloudCvmSyncImageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmSyncImage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_sync_image.sync_image", "id")),
			},
		},
	})
}

const testAccCvmSyncImage = `
resource "tencentcloud_cvm_sync_image" "sync_image" {
	image_id = "img-k4h0m5la" 
	destination_regions = ["ap-shanghai"]
}
`
