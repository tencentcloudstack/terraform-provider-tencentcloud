package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudImagesDataSource_filter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudImagesDataSourceConfigBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"),
					resource.TestMatchResourceAttr("data.tencentcloud_image.public_image", "image_id", regexp.MustCompile("^img-")),
					resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"),
				),
			},
			{
				Config: testAccTencentCloudImagesDataSourceConfigFilter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"),
					resource.TestMatchResourceAttr("data.tencentcloud_image.public_image", "image_id", regexp.MustCompile("^img-")),
					resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"),
				),
			},
			{
				Config: testAccTencentCloudImagesDataSourceConfigFilterWithOsName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"),
					resource.TestMatchResourceAttr("data.tencentcloud_image.public_image", "image_id", regexp.MustCompile("^img-")),
				),
			},
			{
				Config: testAccTencentCloudImagesDataSourceConfigFilterWithImageNameRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"),
					resource.TestMatchResourceAttr("data.tencentcloud_image.public_image", "image_id", regexp.MustCompile("^img-")),
				),
			},
		},
	})
}

const testAccTencentCloudImagesDataSourceConfigBase = `
data "tencentcloud_image" "public_image" {
}
`

const testAccTencentCloudImagesDataSourceConfigFilter = `
data "tencentcloud_image" "public_image" {
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
`

const testAccTencentCloudImagesDataSourceConfigFilterWithOsName = `
data "tencentcloud_image" "public_image" {
  os_name = "CentOS 7.5"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
`

const testAccTencentCloudImagesDataSourceConfigFilterWithImageNameRegex = `
data "tencentcloud_image" "public_image" {
  image_name_regex = "^CentOS\\s+7\\.5\\s+64\\w*"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
`
