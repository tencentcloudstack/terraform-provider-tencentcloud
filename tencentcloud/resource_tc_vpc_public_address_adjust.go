/*
Provides a resource to create a vpc public_address_adjust

Example Usage

```hcl
resource "tencentcloud_vpc_public_address_adjust" "public_address_adjust" {
  instance_id = "ins-osckfnm7"
  address_id = "eip-erft45fu"
}
```

Import

vpc public_address_adjust can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_public_address_adjust.public_address_adjust public_address_adjust_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudVpcPublicAddressAdjust() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPublicAddressAdjustCreate,
		Read:   resourceTencentCloudVpcPublicAddressAdjustRead,
		Delete: resourceTencentCloudVpcPublicAddressAdjustDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A unique ID that identifies the CVM instance. The unique ID of CVM is in the form:&amp;amp;#39;ins-osckfnm7&amp;amp;#39;.",
			},

			"address_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A unique ID that identifies an EIP instance. The unique ID of EIP is in the form:&amp;amp;#39;eip-erft45fu&amp;amp;#39;.",
			},
		},
	}
}

func resourceTencentCloudVpcPublicAddressAdjustCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_public_address_adjust.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = vpc.NewAdjustPublicAddressRequest()
		response   = vpc.NewAdjustPublicAddressResponse()
		instanceId string
		addressId  string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_id"); ok {
		addressId = v.(string)
		request.AddressId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AdjustPublicAddress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc publicAddressAdjust failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{instanceId, addressId}, FILED_SP))

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 1*readRetryTimeout, time.Second, service.VpcPublicAddressAdjustStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudVpcPublicAddressAdjustRead(d, meta)
}

func resourceTencentCloudVpcPublicAddressAdjustRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_public_address_adjust.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcPublicAddressAdjustDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_public_address_adjust.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
