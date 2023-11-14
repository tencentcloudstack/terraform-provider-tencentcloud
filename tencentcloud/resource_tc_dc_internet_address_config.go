/*
Provides a resource to create a dc internet_address_config

Example Usage

```hcl
resource "tencentcloud_dc_internet_address_config" "internet_address_config" {
  instance_id = "ipv4-ljm17pbl"
  enable = true
}
```

Import

dc internet_address_config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_internet_address_config.internet_address_config internet_address_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	"log"
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
				Description: "Internet public address id.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether enable internet address.",
			},
		},
	}
}

func resourceTencentCloudDcInternetAddressConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcInternetAddressConfigUpdate(d, meta)
}

func resourceTencentCloudDcInternetAddressConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	internetAddressConfigId := d.Id()

	internetAddressConfig, err := service.DescribeDcInternetAddressConfigById(ctx, instanceId)
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

	if internetAddressConfig.Enable != nil {
		_ = d.Set("enable", internetAddressConfig.Enable)
	}

	return nil
}

func resourceTencentCloudDcInternetAddressConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		enableInternetAddressRequest  = dc.NewEnableInternetAddressRequest()
		enableInternetAddressResponse = dc.NewEnableInternetAddressResponse()
	)

	internetAddressConfigId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "enable"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().EnableInternetAddress(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dc internetAddressConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcInternetAddressConfigRead(d, meta)
}

func resourceTencentCloudDcInternetAddressConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_internet_address_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
