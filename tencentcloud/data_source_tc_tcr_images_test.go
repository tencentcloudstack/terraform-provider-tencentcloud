package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testObjectName = "data.tencentcloud_tcr_images.images"

func TestAccTencentCloudTcrImagesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImagesDataSource_id,
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttrSet(testObjectName, "registry_id"),
					resource.TestCheckResourceAttrSet(testObjectName, "namespace_name"),
					resource.TestCheckResourceAttrSet(testObjectName, "repository_name"),
					resource.TestCheckResourceAttrSet(testObjectName, "image_info_list.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudTcrImagesDataSource_exact(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImagesDataSource_exact,
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttrSet(testObjectName, "registry_id"),
					resource.TestCheckResourceAttrSet(testObjectName, "namespace_name"),
					resource.TestCheckResourceAttrSet(testObjectName, "repository_name"),
					resource.TestCheckResourceAttr(testObjectName, "exact_match", "true"),
					resource.TestCheckResourceAttrSet(testObjectName, "image_info_list.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudTcrImagesDataSource_exact_version(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImagesDataSource_exact_version,
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttrSet(testObjectName, "registry_id"),
					resource.TestCheckResourceAttrSet(testObjectName, "namespace_name"),
					resource.TestCheckResourceAttrSet(testObjectName, "repository_name"),
					resource.TestCheckResourceAttr(testObjectName, "image_version", "v1"),
					resource.TestCheckResourceAttr(testObjectName, "exact_match", "true"),
					resource.TestCheckResourceAttrSet(testObjectName, "image_info_list.#"),
				),
			},
		},
	})
}

const testAccTcrImagesDataSource_id = TCRDataSource + `

data "tencentcloud_tcr_images" "images" {
  registry_id = local.tcr_id
  namespace_name = local.tcr_ns_name
  repository_name = local.tcr_repo
  }

`

const testAccTcrImagesDataSource_exact = TCRDataSource + `

data "tencentcloud_tcr_images" "images" {
  registry_id = local.tcr_id
  namespace_name = local.tcr_ns_name
  repository_name = local.tcr_repo
  exact_match = true
  }

`

const testAccTcrImagesDataSource_exact_version = TCRDataSource + `

data "tencentcloud_tcr_images" "images" {
  registry_id = local.tcr_id
  namespace_name = local.tcr_ns_name
  repository_name = local.tcr_repo
  image_version = "vv1"
  exact_match = true
  }

`
