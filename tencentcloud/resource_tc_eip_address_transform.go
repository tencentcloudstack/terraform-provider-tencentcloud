/*
Provides a resource to create a eip address_transform

Example Usage

```hcl
resource "tencentcloud_eip_address_transform" "address_transform" {
  instance_id = ""
}
```

Import

eip address_transform can be imported using the id, e.g.

```
terraform import tencentcloud_eip_address_transform.address_transform address_transform_id
```
*/
package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eip "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudEipAddressTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipAddressTransformCreate,
		Read:   resourceTencentCloudEipAddressTransformRead,
		Delete: resourceTencentCloudEipAddressTransformDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "the instance ID of a normal public network IP to be operated. eg:ins-23mk45jn.",
			},
		},
	}
}

func resourceTencentCloudEipAddressTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_address_transform.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = eip.NewTransformAddressRequest()
		response   = eip.NewTransformAddressResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().TransformAddress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate eip addressTransform failed, reason:%+v", logId, err)
		return err
	}

	taskId := *response.Response.TaskId
	d.SetId(instanceId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 1*readRetryTimeout, time.Second, service.VpcIpv6AddressStateRefreshFunc(helper.UInt64ToStr(taskId), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudEipAddressTransformRead(d, meta)
}

func resourceTencentCloudEipAddressTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_address_transform.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudEipAddressTransformDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_address_transform.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
