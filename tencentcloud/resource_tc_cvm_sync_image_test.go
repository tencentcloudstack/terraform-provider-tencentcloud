package tencentcloud

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cvm_sync_image
	resource.AddTestSweepers("tencentcloud_cvm_sync_image", &resource.Sweeper{
		Name: "tencentcloud_cvm_sync_image",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			request := cvm.NewDescribeImagesRequest()
			response, err := client.UseCvmClient().DescribeImages(request)
			if err != nil {
				return err
			}
			for _, image := range response.Response.ImageSet {
				imageName := *image.ImageName
				imageId := *image.ImageId

				now := time.Now()

				createTime := stringTotime(*image.CreatedTime)
				interval := now.Sub(createTime).Minutes()
				if strings.HasPrefix(imageName, keepResource) || strings.HasPrefix(imageName, defaultResource) {
					continue
				}
				// less than 30 minute, not delete
				if needProtect == 1 && int64(interval) < 30 {
					continue
				}

				service := CvmService{client}
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
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
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
