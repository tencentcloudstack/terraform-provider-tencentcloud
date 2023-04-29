/*
Provides a resource to create a mariadb encrypt_attributes

Example Usage

```hcl
resource "tencentcloud_mariadb_encrypt_attributes" "encrypt_attributes" {
  instance_id = "tdsql-ow728lmc"
  encrypt_enabled = 1
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "instance id.",
			},

			"encrypt_enabled": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "whether to enable data encryption, it is not supported to turn it off after it is turned on. The optional values: 0-disable, 1-enable.",
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

	instanceId := d.Id()

	encryptAttributes, err := service.DescribeDBEncryptAttributes(ctx, instanceId)

	if err != nil {
		return err
	}

	if encryptAttributes == nil {
		d.SetId("")
		return fmt.Errorf("resource `encryptAttributes` %s does not exist", instanceId)
	}

	_ = d.Set("instance_id", instanceId)

	if encryptAttributes.EncryptStatus != nil {
		_ = d.Set("encrypt_enabled", encryptAttributes.EncryptStatus)
	}

	return nil
}

func resourceTencentCloudMariadbEncryptAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_encrypt_attributes.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mariadb.NewModifyDBEncryptAttributesRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, _ := d.GetOk("encrypt_enabled"); v != nil {
		request.EncryptEnabled = helper.IntInt64(v.(int))
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
