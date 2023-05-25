/*
Provides a resource to create a dc internet_address_config

Example Usage

```hcl
resource "tencentcloud_dc_internet_address" "internet_address" {
  mask_len = 30
  addr_type = 2
  addr_proto = 0
}

resource "tencentcloud_dc_internet_address_config" "internet_address_config" {
  instance_id = tencentcloud_dc_internet_address.internet_address.id
  enable = false
}
```

Import

dc internet_address_config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_internet_address_config.internet_address_config internet_address_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
)

func resourceTencentCloudDcInternetAddressConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcInternetAddressConfigCreate,
		Read:   resourceTencentCloudDcInternetAddressConfigRead,
		Update: resourceTencentCloudDcInternetAddressConfigUpdate,
		Delete: resourceTencentCloudDcInternetAddressConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "internet public address id.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "whether enable internet address.",
			},
		},
	}
}

func resourceTencentCloudDcInternetAddressConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address_config.create")()
	defer inconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)

	d.SetId(instanceId)

	return resourceTencentCloudDcInternetAddressConfigUpdate(d, meta)
}

func resourceTencentCloudDcInternetAddressConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	internetAddressConfig, err := service.DescribeDcInternetAddressById(ctx, instanceId)
	if err != nil {
		return err
	}

	if internetAddressConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcInternetAddressConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if internetAddressConfig.InstanceId != nil {
		_ = d.Set("instance_id", internetAddressConfig.InstanceId)
	}

	if *internetAddressConfig.Status == 0 {
		_ = d.Set("enable", true)
	} else {
		_ = d.Set("enable", false)
	}

	return nil
}

func resourceTencentCloudDcInternetAddressConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		enable         bool
		enableRequest  = dc.NewEnableInternetAddressRequest()
		disableRequest = dc.NewDisableInternetAddressRequest()
	)

	instanceId := d.Id()

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	}

	if enable {
		enableRequest.InstanceId = &instanceId

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().EnableInternetAddress(enableRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dc internetAddressConfig failed, reason:%+v", logId, err)
			return err
		}
	} else {
		disableRequest.InstanceId = &instanceId

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().DisableInternetAddress(disableRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dc internetAddressConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDcInternetAddressConfigRead(d, meta)
}

func resourceTencentCloudDcInternetAddressConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
