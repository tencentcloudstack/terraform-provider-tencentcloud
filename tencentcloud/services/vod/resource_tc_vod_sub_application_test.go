package vod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVodSubApplicationResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSubApplication,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vod_sub_application.foo", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.foo", "name", "foo"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.foo", "status", "On"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.foo", "description", "this is sub application"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_sub_application.foo", "create_time"),
				),
			},
			{
				ResourceName:            "tencentcloud_vod_sub_application.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func TestAccTencentCloudVodSubApplicationResource_complete(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSubApplicationComplete,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vod_sub_application.complete", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "name", "tf-test-complete"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "status", "On"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "description", "Complete sub application with all parameters"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "type", "Professional"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "mode", "fileid+path"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "storage_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.%", "3"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.team", "media"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.environment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.project", "terraform"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_sub_application.complete", "create_time"),
				),
			},
			// Test tags update
			{
				Config: testAccVodSubApplicationCompleteUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.%", "2"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.team", "media-updated"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.environment", "production"),
				),
			},
			{
				ResourceName:            "tencentcloud_vod_sub_application.complete",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

func TestAccTencentCloudVodSubApplicationResource_professional(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSubApplicationProfessional,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vod_sub_application.professional", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.professional", "name", "tf-test-professional"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.professional", "type", "Professional"),
					resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.professional", "mode", "fileid+path"),
				),
			},
		},
	})
}

const testAccVodSubApplication = `
resource  "tencentcloud_vod_sub_application" "foo" {
	name = "foo"
	status = "On"
	description = "this is sub application"
  }
`

const testAccVodSubApplicationComplete = `
resource "tencentcloud_vod_sub_application" "complete" {
  name           = "tf-test-complete"
  status         = "On"
  description    = "Complete sub application with all parameters"
  type           = "Professional"
  mode           = "fileid+path"
  storage_region = "ap-guangzhou"
  
  tags = {
    "team"        = "media"
    "environment" = "test"
    "project"     = "terraform"
  }
}
`

const testAccVodSubApplicationCompleteUpdate = `
resource "tencentcloud_vod_sub_application" "complete" {
  name           = "tf-test-complete"
  status         = "On"
  description    = "Complete sub application with all parameters"
  type           = "Professional"
  mode           = "fileid+path"
  storage_region = "ap-guangzhou"
  
  tags = {
    "team"        = "media-updated"
    "environment" = "production"
  }
}
`

const testAccVodSubApplicationProfessional = `
resource "tencentcloud_vod_sub_application" "professional" {
  name           = "tf-test-professional"
  status         = "On"
  description    = "Professional sub application"
  type           = "Professional"
  mode           = "fileid+path"
  storage_region = "ap-beijing"
}
`
