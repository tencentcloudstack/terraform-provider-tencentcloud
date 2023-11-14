/*
Provides a resource to create a lighthouse reboot_instance

Example Usage

```hcl
resource "tencentcloud_lighthouse_reboot_instance" "reboot_instance" {
  instance_ids =
}
```

Import

lighthouse reboot_instance can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_reboot_instance.reboot_instance reboot_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"log"
	"time"
)

func resourceTencentCloudLighthouseRebootInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseRebootInstanceCreate,
		Read:   resourceTencentCloudLighthouseRebootInstanceRead,
		Delete: resourceTencentCloudLighthouseRebootInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list.",
			},
		},
	}
}

func resourceTencentCloudLighthouseRebootInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_reboot_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewRebootInstancesRequest()
		response   = lighthouse.NewRebootInstancesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().RebootInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse rebootInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseRebootInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseRebootInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseRebootInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_reboot_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseRebootInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_reboot_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
