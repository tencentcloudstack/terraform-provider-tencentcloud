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
			//// NOTE this test case is dependent in the account which already created some private images
			//{
			//	Config: testAccTencentCloudImagesDataSourceConfigFilterWithPrivateImage,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckTencentCloudDataSourceID("data.tencentcloud_image.private_image"),
			//		resource.TestMatchResourceAttr("data.tencentcloud_image.private_image", "image_id", regexp.MustCompile("^img-")),
			//	),
			//},
			//// NOTE this test case is dependent in the account which already created some private images
			//{
			//	Config: testAccTencentCloudImagesDataSourceConfigFilterWithShareImageAndOsName,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheckTencentCloudDataSourceID("data.tencentcloud_image.private_image"),
			//		resource.TestMatchResourceAttr("data.tencentcloud_image.private_image", "image_id", regexp.MustCompile("^img-")),
			//	),
			//},
		},
	})
}

const testAccTencentCloudImagesDataSourceConfigFilter = `
data "tencentcloud_image" "public_image" {
  filter {
    name = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
`

const testAccTencentCloudImagesDataSourceConfigFilterWithOsName = `
data "tencentcloud_image" "public_image" {
  os_name = "CentOS 7.5 64位"
  filter {
    name = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
`

const testAccTencentCloudImagesDataSourceConfigFilterWithImageNameRegex = `
data "tencentcloud_image" "public_image" {
  image_name_regex = "^CentOS\\s+7\\.5\\s+64\\w*"
  filter {
    name = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
`

const testAccTencentCloudImagesDataSourceConfigFilterWithPrivateImage = `
data "tencentcloud_image" "private_image" {
  image_name_regex = "^batch-tensorflow"
  filter {
    name = "image-type"
    values = ["PRIVATE_IMAGE"]
  }
}
`

const testAccTencentCloudImagesDataSourceConfigFilterWithShareImageAndOsName = `
data "tencentcloud_image" "private_image" {
  os_name = "ubuntu"
  filter {
    name = "image-type"
    values = ["SHARED_IMAGE"]
  }
}
`
