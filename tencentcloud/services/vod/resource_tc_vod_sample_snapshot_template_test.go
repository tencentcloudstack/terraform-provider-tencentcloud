package vod_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

const SampleSnapshotTemplateResourceKey = "tencentcloud_vod_sample_snapshot_template.sample_snapshot_template"

func TestAccTencentCloudVodSampleSnapshotTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSampleSnapshotTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(SampleSnapshotTemplateResourceKey, "id"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "sample_type", "Percent"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "sample_interval", "10"),
					resource.TestCheckResourceAttrSet(SampleSnapshotTemplateResourceKey, "sub_app_id"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "name", "testSampleSnapshot"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "width", "500"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "height", "400"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "resolution_adaptive", "open"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "format", "jpg"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "comment", "test sample snopshot"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "fill_type", "black"),
				),
			},
			{
				Config: testAccVodSampleSnapshotTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(SampleSnapshotTemplateResourceKey, "id"),
					resource.TestCheckResourceAttrSet(SampleSnapshotTemplateResourceKey, "sub_app_id"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "width", "600"),
					resource.TestCheckResourceAttr(SampleSnapshotTemplateResourceKey, "height", "500"),
				),
			},
			{
				ResourceName:      SampleSnapshotTemplateResourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVodSampleSnapshotTemplate = `
resource  "tencentcloud_vod_sub_application" "snapshot_template_sub_application" {
	name = "snapshotTemplateSubApplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_sample_snapshot_template" "sample_snapshot_template" {
  sample_type = "Percent"
  sample_interval = 10
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.snapshot_template_sub_application.id)[1])
  name = "testSampleSnapshot"
  width = 500
  height = 400
  resolution_adaptive = "open"
  format = "jpg"
  comment = "test sample snopshot"
  fill_type = "black"
}
`

const testAccVodSampleSnapshotTemplateUpdate = `
resource  "tencentcloud_vod_sub_application" "snapshot_template_sub_application" {
	name = "snapshotTemplateSubApplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_sample_snapshot_template" "sample_snapshot_template" {
  sample_type = "Percent"
  sample_interval = 10
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.snapshot_template_sub_application.id)[1])
  name = "testSampleSnapshot"
  width = 600
  height = 500
  resolution_adaptive = "open"
  format = "jpg"
  comment = "update"
  fill_type = "black"
}
`
