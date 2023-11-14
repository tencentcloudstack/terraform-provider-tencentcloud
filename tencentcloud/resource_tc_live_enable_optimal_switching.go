/*
Provides a resource to create a live enable_optimal_switching

Example Usage

```hcl
resource "tencentcloud_live_enable_optimal_switching" "enable_optimal_switching" {
  stream_name = ""
  enable_switch =
  host_group_name = ""
}
```

Import

live enable_optimal_switching can be imported using the id, e.g.

```
terraform import tencentcloud_live_enable_optimal_switching.enable_optimal_switching enable_optimal_switching_id
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

func resourceTencentCloudLiveEnableOptimalSwitching() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveEnableOptimalSwitchingCreate,
		Read:   resourceTencentCloudLiveEnableOptimalSwitchingRead,
		Delete: resourceTencentCloudLiveEnableOptimalSwitchingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"stream_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Stream id.",
			},

			"enable_switch": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "1:enable.",
			},

			"host_group_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group name.",
			},
		},
	}
}

func resourceTencentCloudLiveEnableOptimalSwitchingCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_enable_optimal_switching.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewEnableOptimalSwitchingRequest()
		response   = live.NewEnableOptimalSwitchingResponse()
		streamName string
	)
	if v, ok := d.GetOk("stream_name"); ok {
		streamName = v.(string)
		request.StreamName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("enable_switch"); v != nil {
		request.EnableSwitch = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("host_group_name"); ok {
		request.HostGroupName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().EnableOptimalSwitching(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate live enableOptimalSwitching failed, reason:%+v", logId, err)
		return err
	}

	streamName = *response.Response.StreamName
	d.SetId(streamName)

	return resourceTencentCloudLiveEnableOptimalSwitchingRead(d, meta)
}

func resourceTencentCloudLiveEnableOptimalSwitchingRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_enable_optimal_switching.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLiveEnableOptimalSwitchingDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_enable_optimal_switching.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
