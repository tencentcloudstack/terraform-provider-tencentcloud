/*
Provides a resource to create a dasb device_account

Example Usage

```hcl
resource "tencentcloud_dasb_device_account" "example" {
  device_id = 100
  account   = "root"
}
```

Import

dasb device_account can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device_account.example 11
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDasbDeviceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbDeviceAccountCreate,
		Read:   resourceTencentCloudDasbDeviceAccountRead,
		Delete: resourceTencentCloudDasbDeviceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Device ID.",
			},
			"account": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Device account.",
			},
		},
	}
}

func resourceTencentCloudDasbDeviceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_account.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		request         = dasb.NewCreateDeviceAccountRequest()
		response        = dasb.NewCreateDeviceAccountResponse()
		deviceAccountId string
	)

	if v, ok := d.GetOkExists("device_id"); ok {
		request.DeviceId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("account"); ok {
		request.Account = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDasbClient().CreateDeviceAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response.Id != nil {
			e = fmt.Errorf("dasb DeviceAccount not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb DeviceAccount failed, reason:%+v", logId, err)
		return err
	}

	deviceAccountIdInt := *response.Response.Id
	deviceAccountId = strconv.FormatUint(deviceAccountIdInt, 10)
	d.SetId(deviceAccountId)

	return resourceTencentCloudDasbDeviceAccountRead(d, meta)
}

func resourceTencentCloudDasbDeviceAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_account.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		service         = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
		deviceAccountId = d.Id()
	)

	DeviceAccount, err := service.DescribeDasbDeviceAccountById(ctx, deviceAccountId)
	if err != nil {
		return err
	}

	if DeviceAccount == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbDeviceAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if DeviceAccount.DeviceId != nil {
		_ = d.Set("device_id", DeviceAccount.DeviceId)
	}

	if DeviceAccount.Account != nil {
		_ = d.Set("account", DeviceAccount.Account)
	}

	return nil
}

func resourceTencentCloudDasbDeviceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dasb_device_account.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		service         = DasbService{client: meta.(*TencentCloudClient).apiV3Conn}
		deviceAccountId = d.Id()
	)

	if err := service.DeleteDasbDeviceAccountById(ctx, deviceAccountId); err != nil {
		return err
	}

	return nil
}
