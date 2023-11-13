/*
Provides a resource to create a mariadb encrypt_attributes

Example Usage

```hcl
resource "tencentcloud_mariadb_encrypt_attributes" "encrypt_attributes" {
  instance_id = "tdsql-e9tklsgz"
  encrypt_enabled =
}
```

Import

mariadb encrypt_attributes can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_encrypt_attributes.encrypt_attributes encrypt_attributes_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"log"
)

func resourceTencentCloudMariadbEncryptAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbEncryptAttributesCreate,
		Read:   resourceTencentCloudMariadbEncryptAttributesRead,
		Update: resourceTencentCloudMariadbEncryptAttributesUpdate,
		Delete: resourceTencentCloudMariadbEncryptAttributesDelete,
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

func resourceTencentCloudMariadbEncryptAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_encrypt_attributes.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbEncryptAttributesUpdate(d, meta)
}

func resourceTencentCloudMariadbEncryptAttributesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_encrypt_attributes.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	encryptAttributesId := d.Id()

	encryptAttributes, err := service.DescribeMariadbEncryptAttributesById(ctx, instanceId)
	if err != nil {
		return err
	}

	if encryptAttributes == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbEncryptAttributes` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudMariadbEncryptAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_encrypt_attributes.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mariadb.NewModifyDBEncryptAttributesRequest()

	encryptAttributesId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "encrypt_enabled"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBEncryptAttributes(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mariadb encryptAttributes failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbEncryptAttributesRead(d, meta)
}

func resourceTencentCloudMariadbEncryptAttributesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_encrypt_attributes.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
