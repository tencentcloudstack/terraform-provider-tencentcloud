/*
Provides a resource to create a mariadb encrypt_attributes_operation

Example Usage

```hcl
resource "tencentcloud_mariadb_encrypt_attributes_operation" "encrypt_attributes_operation" {
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

func resourceTencentCloudMariadbEncryptAttributesOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbEncryptAttributesOperationCreate,
		Read:   resourceTencentCloudMariadbEncryptAttributesOperationRead,
		Delete: resourceTencentCloudMariadbEncryptAttributesOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id, in the form of: tdsql-ow728lmc.",
			},

			"encrypt_enabled": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "whether to enable data encryption, it is not supported to turn it off after it is turned on. The optional values: 0-disable, 1-enable.",
			},
		},
	}
}

func resourceTencentCloudMariadbEncryptAttributesOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_encrypt_attributes_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewModifyDBEncryptAttributesRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

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
		log.Printf("[CRITAL]%s operate mariadb encryptAttributesOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbEncryptAttributesOperationRead(d, meta)
}

func resourceTencentCloudMariadbEncryptAttributesOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_encrypt_attributes_operation.read")()
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

func resourceTencentCloudMariadbEncryptAttributesOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_encrypt_attributes_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
