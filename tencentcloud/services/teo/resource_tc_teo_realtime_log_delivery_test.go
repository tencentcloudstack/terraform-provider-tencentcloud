package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoRealtimeLogDeliveryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRealtimeLogDelivery,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "delivery_status", "enabled"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "entity_list.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "fields.#"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "log_type", "application"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "sample", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "task_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "task_type", "s3"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "log_format.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.access_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.access_key"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.compress_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.endpoint"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.region"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoRealtimeLogDeliveryUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "delivery_status", "disabled"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "entity_list.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "fields.#"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "log_type", "application"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "sample", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "task_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "task_type", "s3"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "log_format.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.access_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.access_key"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.compress_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.endpoint"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_realtime_log_delivery.teo_realtime_log_delivery", "s3.0.region"),
				),
			},
		},
	})
}

const testAccTeoRealtimeLogDelivery = `

resource "tencentcloud_teo_realtime_log_delivery" "teo_realtime_log_delivery" {
    area            = "overseas"
    delivery_status = "enabled"
    entity_list     = [
        "sid-2yvhjw98uaco",
    ]
    fields          = [
        "ServiceID",
        "ConnectTimeStamp",
        "DisconnetTimeStamp",
        "DisconnetReason",
        "ClientRealIP",
        "ClientRegion",
        "EdgeIP",
        "ForwardProtocol",
        "ForwardPort",
        "SentBytes",
        "ReceivedBytes",
        "LogTimeStamp",
    ]
    log_type        = "application"
    sample          = 0
    task_name       = "test"
    task_type       = "s3"
    zone_id         = "zone-2qtuhspy7cr6"

    log_format {
        field_delimiter  = ","
        format_type      = "json"
        record_delimiter = "\n"
        record_prefix    = "{"
        record_suffix    = "}"
    }

    s3 {
        access_id     = "xxxxxxxxxx"
        access_key    = "xxxxxxxxxx"
        bucket        = "test-1253833068"
        compress_type = "gzip"
        endpoint      = "https://test-1253833068.cos.ap-nanjing.myqcloud.com"
        region        = "ap-nanjing"
    }
}
`

const testAccTeoRealtimeLogDeliveryUp = `

resource "tencentcloud_teo_realtime_log_delivery" "teo_realtime_log_delivery" {
    area            = "overseas"
    delivery_status = "disabled"
    entity_list     = [
        "sid-2yvhjw98uaco",
    ]
    fields          = [
        "ServiceID",
        "ConnectTimeStamp",
        "DisconnetTimeStamp",
        "DisconnetReason",
        "ClientRealIP",
        "ClientRegion",
        "EdgeIP",
        "ForwardProtocol",
        "ForwardPort",
        "SentBytes",
        "ReceivedBytes",
        "LogTimeStamp",
    ]
    log_type        = "application"
    sample          = 0
    task_name       = "test"
    task_type       = "s3"
    zone_id         = "zone-2qtuhspy7cr6"

    log_format {
        field_delimiter  = ","
        format_type      = "json"
        record_delimiter = "\n"
        record_prefix    = "{"
        record_suffix    = "}"
    }

    s3 {
        access_id     = "xxxxxxxxxx"
        access_key    = "xxxxxxxxxx"
        bucket        = "test-1253833068"
        compress_type = "gzip"
        endpoint      = "https://test-1253833068.cos.ap-nanjing.myqcloud.com"
        region        = "ap-nanjing"
    }
}
`
