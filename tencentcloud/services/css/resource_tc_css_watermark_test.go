package css_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccss "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/css"

	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_css_watermark", &resource.Sweeper{
		Name: "tencentcloud_css_watermark",
		F:    testSweepCSSWatermarkTask,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_css_watermark
func testSweepCSSWatermarkTask(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	cssService := svccss.NewCssService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

	marks, err := cssService.DescribeCssWatermarks(ctx)
	if err != nil {
		return err
	}
	if marks == nil {
		return fmt.Errorf("watermark instance not exists.")
	}

	for _, v := range marks {
		delName := v.WatermarkName
		delId := v.WatermarkId

		if strings.HasPrefix(*delName, tcacctest.DefaultCSSPrefix) {
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				err := cssService.DeleteCssWatermarkById(ctx, delId)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] sweeper watermark instance %s:%v failed! reason:[%s]", *delName, *delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudCssWatermarkResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCssWatermarkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssWatermark, tcacctest.DefaultCSSPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCssWatermarkExists("tencentcloud_css_watermark.watermark"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "id"),
					resource.TestMatchResourceAttr("tencentcloud_css_watermark.watermark", "picture_url", regexp.MustCompile("https://main.qcloudimg.com")),
					resource.TestMatchResourceAttr("tencentcloud_css_watermark.watermark", "watermark_name", regexp.MustCompile(tcacctest.DefaultCSSPrefix)),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "x_position"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "y_position"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "width"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "height"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCssWatermark_update, tcacctest.DefaultCSSPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCssWatermarkExists("tencentcloud_css_watermark.watermark"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "id"),
					resource.TestMatchResourceAttr("tencentcloud_css_watermark.watermark", "picture_url", regexp.MustCompile("changed")),
					resource.TestMatchResourceAttr("tencentcloud_css_watermark.watermark", "watermark_name", regexp.MustCompile("changed")),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "status"),
					resource.TestCheckResourceAttr("tencentcloud_css_watermark.watermark", "x_position", "5"),
					resource.TestCheckResourceAttr("tencentcloud_css_watermark.watermark", "y_position", "5"),
					resource.TestCheckResourceAttr("tencentcloud_css_watermark.watermark", "width", "100"),
					resource.TestCheckResourceAttr("tencentcloud_css_watermark.watermark", "height", "100"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_watermark.watermark",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCssWatermarkDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_css_watermark" {
			continue
		}

		watermark, err := cssService.DescribeCssWatermark(ctx, rs.Primary.ID)
		if err != nil {
			return nil
		}

		if watermark != nil {
			return fmt.Errorf("css watermark still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCssWatermarkExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("css watermark %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("css watermark id is not set")
		}

		cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		watermark, err := cssService.DescribeCssWatermark(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if watermark == nil {
			return fmt.Errorf("css watermark not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCssWatermark = `

resource "tencentcloud_css_watermark" "watermark" {
  picture_url = "https://main.qcloudimg.com/raw/c3e0cf113a5c5346b776ecbcfbdcfc72.svg"
  watermark_name = "%swm"
  x_position = 0
  y_position = 0
  width = 0
  height = 0
}

`

const testAccCssWatermark_update = `

resource "tencentcloud_css_watermark" "watermark" {
  picture_url = "https://main.qcloudimg.com/raw/changed.svg"
  watermark_name = "%schanged"
  x_position = 5
  y_position = 5
  width = 100
  height = 100
}

`
