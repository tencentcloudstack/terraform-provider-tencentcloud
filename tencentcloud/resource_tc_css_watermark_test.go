package tencentcloud

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_css_watermark", &resource.Sweeper{
		Name: "tencentcloud_css_watermark",
		F:    testSweepCSSWatermarkTask,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_css_watermark
func testSweepCSSWatermarkTask(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	cssService := CssService{client: cli.(*TencentCloudClient).apiV3Conn}

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

		if strings.HasPrefix(*delName, defaultCSSPrefix) {
			err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := cssService.DeleteCssWatermarkById(ctx, delId)
				if err != nil {
					return retryError(err)
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

func TestAccTencentCloudCSSWatermarkResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCssWatermarkDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssWatermark, defaultCSSPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCssWatermarkExists("tencentcloud_css_watermark.watermark"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "id"),
					resource.TestMatchResourceAttr("tencentcloud_css_watermark.watermark", "picture_url", regexp.MustCompile("https://main.qcloudimg.com")),
					resource.TestMatchResourceAttr("tencentcloud_css_watermark.watermark", "watermark_name", regexp.MustCompile(defaultCSSPrefix)),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "x_position"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "y_position"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "width"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "height"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCssWatermark_update, defaultCSSPrefix),
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cssService := CssService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("css watermark %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("css watermark id is not set")
		}

		cssService := CssService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
