package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmImageDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"), resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"), resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_id")),
			},
		},
	})
}

const testAccCvmImageDataSource_BasicCreate = `

data "tencentcloud_image" "public_image" {
}

`

func TestAccTencentCloudCvmImageDataSource_Filter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageDataSource_FilterCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"), resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"), resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_id")),
			},
		},
	})
}

const testAccCvmImageDataSource_FilterCreate = `

data "tencentcloud_image" "public_image" {
    
    filter {
        values = ["PUBLIC_IMAGE"]
        name = "image-type"
    }
}

`

func TestAccTencentCloudCvmImageDataSource_WithOsName(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageDataSource_WithOsNameCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"), resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"), resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_id")),
			},
		},
	})
}

const testAccCvmImageDataSource_WithOsNameCreate = `

data "tencentcloud_image" "public_image" {
    os_name = "TencentOS Server 3.2"
}

`

func TestAccTencentCloudCvmImageDataSource_WithImageNameRegex(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageDataSource_WithImageNameRegexCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"), resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"), resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_id")),
			},
		},
	})
}

const testAccCvmImageDataSource_WithImageNameRegexCreate = `

data "tencentcloud_image" "public_image" {
    image_name_regex = "^Windows\\s.*$"
}

`

// go test -i; go test -test.run TestAccTencentCloudDataSourceImageBase_basic -v
// func TestAccTencentCloudDataSourceImageBase_basic(t *testing.T) {
// 	t.Parallel()
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { tcacctest.AccPreCheck(t) },
// 		Providers: tcacctest.AccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccTencentCloudDataSourceImageBase,
// 				Check: resource.ComposeTestCheckFunc(
// 					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"),
// 					resource.TestMatchResourceAttr("data.tencentcloud_image.public_image", "image_id", regexp.MustCompile("^img-")),
// 					resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"),
// 				),
// 			},
// 			{
// 				Config: testAccTencentCloudDataSourceImageBaseWithFilter,
// 				Check: resource.ComposeTestCheckFunc(
// 					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"),
// 					resource.TestMatchResourceAttr("data.tencentcloud_image.public_image", "image_id", regexp.MustCompile("^img-")),
// 					resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"),
// 				),
// 			},
// 			{
// 				Config: testAccTencentCloudDataSourceImageBaseWithOsName,
// 				Check: resource.ComposeTestCheckFunc(
// 					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"),
// 					resource.TestMatchResourceAttr("data.tencentcloud_image.public_image", "image_id", regexp.MustCompile("^img-")),
// 					resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"),
// 				),
// 			},
// 			{
// 				Config: testAccTencentCloudDataSourceImageBaseWithImageNameRegex,
// 				Check: resource.ComposeTestCheckFunc(
// 					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_image.public_image"),
// 					resource.TestMatchResourceAttr("data.tencentcloud_image.public_image", "image_id", regexp.MustCompile("^img-")),
// 					resource.TestCheckResourceAttrSet("data.tencentcloud_image.public_image", "image_name"),
// 				),
// 			},
// 		},
// 	})
// }

// const testAccTencentCloudDataSourceImageBase = `
// data "tencentcloud_image" "public_image" {}
// `

// const testAccTencentCloudDataSourceImageBaseWithFilter = `
// data "tencentcloud_image" "public_image" {
//   filter {
//     name   = "image-type"
//     values = ["PUBLIC_IMAGE"]
//   }
// }
// `

// const testAccTencentCloudDataSourceImageBaseWithOsName = `
// data "tencentcloud_image" "public_image" {
//   os_name = "TencentOS Server 3.2"

//   filter {
//     name   = "image-type"
//     values = ["PUBLIC_IMAGE"]
//   }
// }
// `

// const testAccTencentCloudDataSourceImageBaseWithImageNameRegex = `
// data "tencentcloud_image" "public_image" {
//   image_name_regex = "^Windows\\s.*$"

//   filter {
//     name   = "image-type"
//     values = ["PUBLIC_IMAGE"]
//   }
// }
// `
