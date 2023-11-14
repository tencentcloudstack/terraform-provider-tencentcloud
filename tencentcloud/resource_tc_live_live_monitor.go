/*
Provides a resource to create a live live_monitor

Example Usage

```hcl
resource "tencentcloud_live_live_monitor" "live_monitor" {
  monitor_id = ""
  audible_input_index_list =
}
```

Import

live live_monitor can be imported using the id, e.g.

```
terraform import tencentcloud_live_live_monitor.live_monitor live_monitor_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLiveLive_monitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveLive_monitorCreate,
		Read:   resourceTencentCloudLiveLive_monitorRead,
		Delete: resourceTencentCloudLiveLive_monitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Monitor ID。.",
			},

			"audible_input_index_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The input index for monitoring the screen audio, supports multiple input audio sources.The valid range for InputIndex is that it must already exist.If left blank, there will be no audio output by default.",
			},
		},
	}
}

func resourceTencentCloudLiveLive_monitorCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_live_monitor.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = live.NewStartLiveStreamMonitorRequest()
		response  = live.NewStartLiveStreamMonitorResponse()
		monitorId string
	)
	if v, ok := d.GetOk("monitor_id"); ok {
		monitorId = v.(string)
		request.MonitorId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("audible_input_index_list"); ok {
		audibleInputIndexListSet := v.(*schema.Set).List()
		for i := range audibleInputIndexListSet {
			audibleInputIndexList := audibleInputIndexListSet[i].(int)
			request.AudibleInputIndexList = append(request.AudibleInputIndexList, helper.IntUint64(audibleInputIndexList))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().StartLiveStreamMonitor(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate live live_monitor failed, reason:%+v", logId, err)
		return err
	}

	monitorId = *response.Response.MonitorId
	d.SetId(monitorId)

	return resourceTencentCloudLiveLive_monitorRead(d, meta)
}

func resourceTencentCloudLiveLive_monitorRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_live_monitor.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLiveLive_monitorDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_live_monitor.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
