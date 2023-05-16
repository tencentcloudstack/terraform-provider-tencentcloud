/*
Provides a resource to create a vpc normal_address_return

Example Usage

```hcl
resource "tencentcloud_eip_normal_address_return" "normal_address_return" {
  address_ips =
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudEipNormalAddressReturn() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipNormalAddressReturnCreate,
		Read:   resourceTencentCloudEipNormalAddressReturnRead,
		Delete: resourceTencentCloudEipNormalAddressReturnDelete,
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

func resourceTencentCloudEipNormalAddressReturnCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_normal_address_return.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = vpc.NewReturnNormalAddressesRequest()
		addressIps string
	)
	if v, ok := d.GetOk("address_ips"); ok {
		addressIpsSet := v.(*schema.Set).List()
		for i := range addressIpsSet {
			addressIp := addressIpsSet[i].(string)
			request.AddressIps = append(request.AddressIps, &addressIp)
			addressIps = addressIp + FILED_SP
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ReturnNormalAddresses(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc normalAddressReturn failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(addressIps)

	return resourceTencentCloudEipNormalAddressReturnRead(d, meta)
}

func resourceTencentCloudEipNormalAddressReturnRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_normal_address_return.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudEipNormalAddressReturnDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_normal_address_return.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
