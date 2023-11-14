/*
Provides a resource to create a monitor enable_grafana_internet

Example Usage

```hcl
resource "tencentcloud_monitor_enable_grafana_internet" "enable_grafana_internet" {
  instance_i_d = "grafana-kleu3gt0"
  enable_internet = true
}
```

Import

monitor enable_grafana_internet can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_enable_grafana_internet.enable_grafana_internet enable_grafana_internet_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMonitorEnableGrafanaInternet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorEnableGrafanaInternetCreate,
		Read:   resourceTencentCloudMonitorEnableGrafanaInternetRead,
		Delete: resourceTencentCloudMonitorEnableGrafanaInternetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_i_d": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},

			"enable_internet": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the Internet access: true for enabling; false for disabling.",
			},
		},
	}
}

func resourceTencentCloudMonitorEnableGrafanaInternetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_enable_grafana_internet.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = monitor.NewEnableGrafanaInternetRequest()
		response   = monitor.NewEnableGrafanaInternetResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_i_d"); ok {
		request.InstanceID = helper.String(v.(string))
	}

	if v, _ := d.GetOk("enable_internet"); v != nil {
		request.EnableInternet = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().EnableGrafanaInternet(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate monitor enableGrafanaInternet failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMonitorEnableGrafanaInternetRead(d, meta)
}

func resourceTencentCloudMonitorEnableGrafanaInternetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_enable_grafana_internet.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMonitorEnableGrafanaInternetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_enable_grafana_internet.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
