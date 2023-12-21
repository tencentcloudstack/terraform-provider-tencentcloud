package mariadb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMariadbEncryptAttributes() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_encrypt_attributes.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbEncryptAttributesUpdate(d, meta)
}

func resourceTencentCloudMariadbEncryptAttributesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_encrypt_attributes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_encrypt_attributes.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mariadb.NewModifyDBEncryptAttributesRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, _ := d.GetOk("encrypt_enabled"); v != nil {
		request.EncryptEnabled = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ModifyDBEncryptAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_encrypt_attributes.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
