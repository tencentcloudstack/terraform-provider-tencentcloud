/*
Provides a resource to create a css start_stream_monitor

Example Usage

```hcl
resource "tencentcloud_css_start_stream_monitor" "start_stream_monitor" {
  monitor_id               = "3d5738dd-1ca2-4601-a6e9-004c5ec75c0b"
  audible_input_index_list = [1]
}
```

Import

css start_stream_monitor can be imported using the id, e.g.

```
terraform import tencentcloud_css_start_stream_monitor.start_stream_monitor start_stream_monitor_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssStartStreamMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssStartStreamMonitorCreate,
		Read:   resourceTencentCloudCssStartStreamMonitorRead,
		Delete: resourceTencentCloudCssStartStreamMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Monitor id.",
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

func resourceTencentCloudCssStartStreamMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_start_stream_monitor.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = css.NewStartLiveStreamMonitorRequest()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().StartLiveStreamMonitor(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css StartStreamMonitor failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(monitorId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"1"}, 6*readRetryTimeout, time.Second, service.CssStartStreamMonitorStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCssStartStreamMonitorRead(d, meta)
}

func resourceTencentCloudCssStartStreamMonitorRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_start_stream_monitor.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	monitorId := d.Id()

	StartStreamMonitor, err := service.DescribeCssStreamMonitorById(ctx, monitorId)
	if err != nil {
		return err
	}

	if StartStreamMonitor == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssStartStreamMonitor` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if StartStreamMonitor.MonitorId != nil {
		_ = d.Set("monitor_id", StartStreamMonitor.MonitorId)
	}

	if StartStreamMonitor.AudibleInputIndexList != nil {
		_ = d.Set("audible_input_index_list", StartStreamMonitor.AudibleInputIndexList)
	}

	return nil
}

func resourceTencentCloudCssStartStreamMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_start_stream_monitor.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	monitorId := d.Id()

	if err := service.DeleteCssStartStreamMonitorById(ctx, monitorId); err != nil {
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"0"}, 6*readRetryTimeout, time.Second, service.CssStartStreamMonitorStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
