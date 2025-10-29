package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataLineageAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataLineageAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_lineage_attachment.example", "id"),
				),
			},
		},
	})
}

const testAccWedataLineageAttachment = `
resource "tencentcloud_wedata_lineage_attachment" "example" {
  relations {
    source {
      resource_unique_id = "2s5veseIo2AXGOHJkKjBvQ"
      resource_type      = "TABLE"
      platform           = "WEDATA"
      resource_name      = "smoke_testing_db_no_delete.tdsqlc2dlc_sink_01_bj"
      description        = "DLC"
    }

    target {
      resource_unique_id = "fM8OgzE-AM2h4aaJmdXoPg"
      resource_type      = "TABLE"
      platform           = "WEDATA"
      resource_name      = "smoke_testing_db_no_delete.dlc_source_4col"
      description        = "DLC"
    }

    processes {
      process_id       = "20241107221758402"
      process_type     = "SCHEDULE_TASK"
      platform         = "WEDATA"
      process_sub_type = "SQL_TASK"
    }
  }
}
`
