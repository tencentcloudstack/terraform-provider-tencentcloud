/*
Provides a resource to create a dc internet_address

Example Usage

```hcl
resource "tencentcloud_dc_internet_address" "internet_address" {
  mask_len = 30
  addr_type = 2
  addr_proto = 0
}
```

Import

dc internet_address can be imported using the id, e.g.

```
terraform import tencentcloud_dc_internet_address.internet_address internet_address_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcInternetAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcInternetAddressCreate,
		Read:   resourceTencentCloudDcInternetAddressRead,
		Delete: resourceTencentCloudDcInternetAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mask_len": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "CIDR address mask.",
			},

			"addr_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "0: BGP, 1: china telecom, 2: china mobile, 3: china unicom.",
			},

			"addr_proto": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "0: IPv4, 1: IPv6.",
			},
		},
	}
}

func resourceTencentCloudDcInternetAddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dc.NewApplyInternetAddressRequest()
		response   = dc.NewApplyInternetAddressResponse()
		instanceId string
	)
	if v, ok := d.GetOkExists("mask_len"); ok {
		request.MaskLen = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("addr_type"); ok {
		request.AddrType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("addr_proto"); ok {
		request.AddrProto = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().ApplyInternetAddress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dc internetAddress failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcInternetAddressRead(d, meta)
}

func resourceTencentCloudDcInternetAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	internetAddress, err := service.DescribeDcInternetAddressById(ctx, instanceId)
	if err != nil {
		return err
	}

	if internetAddress == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcInternetAddress` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if internetAddress.MaskLen != nil {
		_ = d.Set("mask_len", internetAddress.MaskLen)
	}

	if internetAddress.AddrType != nil {
		_ = d.Set("addr_type", internetAddress.AddrType)
	}

	if internetAddress.AddrProto != nil {
		_ = d.Set("addr_proto", internetAddress.AddrProto)
	}

	return nil
}

func resourceTencentCloudDcInternetAddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.DeleteDcInternetAddressById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
