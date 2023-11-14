/*
Provides a resource to create a dcdb encrypt_attributes

Example Usage

```hcl
resource "tencentcloud_dcdb_encrypt_attributes" "encrypt_attributes" {
  instance_id = ""
  encrypt_enabled =
}
```

Import

dcdb encrypt_attributes can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_encrypt_attributes.encrypt_attributes encrypt_attributes_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"log"
)

func resourceTencentCloudDcdbEncryptAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbEncryptAttributesCreate,
		Read:   resourceTencentCloudDcdbEncryptAttributesRead,
		Update: resourceTencentCloudDcdbEncryptAttributesUpdate,
		Delete: resourceTencentCloudDcdbEncryptAttributesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"encrypt_enabled": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable data encryption, it is not supported to turn it off after it is turned on. The optional values: 0-disable, 1-enable.",
			},
		},
	}
}

func resourceTencentCloudDcdbEncryptAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcdbEncryptAttributesUpdate(d, meta)
}

func resourceTencentCloudDcdbEncryptAttributesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	encryptAttributesId := d.Id()

	encryptAttributes, err := service.DescribeDcdbEncryptAttributesById(ctx, instanceId)
	if err != nil {
		return err
	}

	if encryptAttributes == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbEncryptAttributes` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if encryptAttributes.InstanceId != nil {
		_ = d.Set("instance_id", encryptAttributes.InstanceId)
	}

	if encryptAttributes.EncryptEnabled != nil {
		_ = d.Set("encrypt_enabled", encryptAttributes.EncryptEnabled)
	}

	return nil
}

func resourceTencentCloudDcdbEncryptAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyDBEncryptAttributesRequest()

	encryptAttributesId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "encrypt_enabled"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBEncryptAttributes(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb encryptAttributes failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbEncryptAttributesRead(d, meta)
}

func resourceTencentCloudDcdbEncryptAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
