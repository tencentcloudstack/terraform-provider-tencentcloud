/*
Provides a resource to create a vpc normal_address_return

Example Usage

```hcl
resource "tencentcloud_vpc_normal_address_return" "normal_address_return" {
  address_ips =
}
```

Import

vpc normal_address_return can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_normal_address_return.normal_address_return normal_address_return_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"log"
)

func resourceTencentCloudVpcNormalAddressReturn() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcNormalAddressReturnCreate,
		Read:   resourceTencentCloudVpcNormalAddressReturnRead,
		Delete: resourceTencentCloudVpcNormalAddressReturnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"address_ips": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IP address of the EIP, example: 101.35.139.183.",
			},
		},
	}
}

func resourceTencentCloudVpcNormalAddressReturnCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_normal_address_return.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = vpc.NewReturnNormalAddressesRequest()
		response = vpc.NewReturnNormalAddressesResponse()
	)
	if v, ok := d.GetOk("address_ips"); ok {
		addressIpsSet := v.(*schema.Set).List()
		for i := range addressIpsSet {
			addressIps := addressIpsSet[i].(string)
			request.AddressIps = append(request.AddressIps, &addressIps)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ReturnNormalAddresses(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc normalAddressReturn failed, reason:%+v", logId, err)
		return err
	}

	addressIps = *response.Response.AddressIps
	d.SetId()

	return resourceTencentCloudVpcNormalAddressReturnRead(d, meta)
}

func resourceTencentCloudVpcNormalAddressReturnRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_normal_address_return.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcNormalAddressReturnDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_normal_address_return.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
