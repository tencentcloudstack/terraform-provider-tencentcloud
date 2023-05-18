/*
Provides a resource to create a lighthouse start_instance

Example Usage

```hcl
resource "tencentcloud_lighthouse_start_instance" "start_instance" {
  instance_id = "lhins-xxxxxx"
}
```
*/
package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

func resourceTencentCloudLighthouseStartInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseStartInstanceCreate,
		Read:   resourceTencentCloudLighthouseStartInstanceRead,
		Delete: resourceTencentCloudLighthouseStartInstanceDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudLighthouseStartInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_start_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := lighthouse.NewStartInstancesRequest()
	instanceId := d.Get("instance_id").(string)
	request.InstanceIds = []*string{&instanceId}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().StartInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse startInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseStartInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseStartInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_start_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseStartInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_start_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
