/*
Provides a resource to create a vpc address_transform

Example Usage

```hcl
resource "tencentcloud_vpc_address_transform" "address_transform" {
  instance_id = ""
}
```

Import

vpc address_transform can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_address_transform.address_transform address_transform_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudVpcAddressTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcAddressTransformCreate,
		Read:   resourceTencentCloudVpcAddressTransformRead,
		Delete: resourceTencentCloudVpcAddressTransformDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance ID of a normal public network IP to be operated. eg:ins-23mk45jn.",
			},
		},
	}
}

func resourceTencentCloudVpcAddressTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_address_transform.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = vpc.NewTransformAddressRequest()
		response   = vpc.NewTransformAddressResponse()
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
		log.Printf("[CRITAL]%s operate vpc addressTransform failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 1*readRetryTimeout, time.Second, service.VpcAddressTransformStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudVpcAddressTransformRead(d, meta)
}

func resourceTencentCloudVpcAddressTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_address_transform.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcAddressTransformDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_address_transform.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
