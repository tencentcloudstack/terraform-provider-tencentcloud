package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsSampleSnapshotTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsSampleSnapshotTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_sample_snapshot_template.sample_snapshot_template", "id")),
			},
			{
				Config: testAccMpsSampleSnapshotTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_sample_snapshot_template.sample_snapshot_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_sample_snapshot_template.sample_snapshot_template", "name", "terraform-for-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_mps_sample_snapshot_template.sample_snapshot_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsSampleSnapshotTemplate = `

resource "tencentcloud_mps_sample_snapshot_template" "sample_snapshot_template" {
  fill_type           = "stretch"
  format              = "jpg"
  height              = 128
  name                = "terraform-test"
  resolution_adaptive = "open"
  sample_interval     = 10
  sample_type         = "Percent"
  width               = 140
}

`

const testAccMpsSampleSnapshotTemplateUpdate = `

resource "tencentcloud_mps_sample_snapshot_template" "sample_snapshot_template" {
  fill_type           = "stretch"
  format              = "jpg"
  height              = 128
  name                = "terraform-for-test"
  resolution_adaptive = "open"
  sample_interval     = 10
  sample_type         = "Percent"
  width               = 140
}

`
