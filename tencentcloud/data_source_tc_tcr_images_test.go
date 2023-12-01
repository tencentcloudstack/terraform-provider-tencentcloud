package tencentcloud

import (
	"fmt"
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
				Config: fmt.Sprintf(testAccTcrImagesDataSource_id, defaultTCRInstanceId, defaultTCRNamespace, defaultTCRRepoName),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttr(testObjectName, "registry_id", defaultTCRInstanceId),
					resource.TestCheckResourceAttr(testObjectName, "namespace_name", defaultTCRNamespace),
					resource.TestCheckResourceAttr(testObjectName, "repository_name", defaultTCRRepoName),
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
				Config: fmt.Sprintf(testAccTcrImagesDataSource_exact, defaultTCRInstanceId, defaultTCRNamespace, defaultTCRRepoName),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttr(testObjectName, "registry_id", defaultTCRInstanceId),
					resource.TestCheckResourceAttr(testObjectName, "namespace_name", defaultTCRNamespace),
					resource.TestCheckResourceAttr(testObjectName, "repository_name", defaultTCRRepoName),
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
				Config: fmt.Sprintf(testAccTcrImagesDataSource_exact_version, defaultTCRInstanceId, defaultTCRNamespace, defaultTCRRepoName),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttr(testObjectName, "registry_id", defaultTCRInstanceId),
					resource.TestCheckResourceAttr(testObjectName, "namespace_name", defaultTCRNamespace),
					resource.TestCheckResourceAttr(testObjectName, "repository_name", defaultTCRRepoName),
					resource.TestCheckResourceAttr(testObjectName, "image_version", "v1"),
					resource.TestCheckResourceAttr(testObjectName, "exact_match", "true"),
					resource.TestCheckResourceAttrSet(testObjectName, "image_info_list.#"),
				),
			},
		},
	})
}

const testAccTcrImagesDataSource_id = `

data "tencentcloud_tcr_images" "images" {
  registry_id = "%s"
  namespace_name = "%s" 
  repository_name = "%s"
  }

`

const testAccTcrImagesDataSource_exact = `

data "tencentcloud_tcr_images" "images" {
  registry_id = "%s"
  namespace_name = "%s" 
  repository_name = "%s"
  exact_match = true
  }

`

const testAccTcrImagesDataSource_exact_version = `

data "tencentcloud_tcr_images" "images" {
  registry_id = "%s"
  namespace_name = "%s" 
  repository_name = "%s"
  image_version = "v1"
  exact_match = true
  }

`
