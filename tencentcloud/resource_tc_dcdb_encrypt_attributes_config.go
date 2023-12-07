package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcdbEncryptAttributesConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbEncryptAttributesConfigCreate,
		Read:   resourceTencentCloudDcdbEncryptAttributesConfigRead,
		Update: resourceTencentCloudDcdbEncryptAttributesConfigUpdate,
		Delete: resourceTencentCloudDcdbEncryptAttributesConfigDelete,
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
				Description: "whether to enable data encryption. Notice: it is not supported to turn it off after it is turned on. The optional values: 0-disable, 1-enable.",
			},
		},
	}
}

func resourceTencentCloudDcdbEncryptAttributesConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	d.SetId(instanceId)

	return resourceTencentCloudDcdbEncryptAttributesConfigUpdate(d, meta)
}

func resourceTencentCloudDcdbEncryptAttributesConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	encryptAttributesConfig, err := service.DescribeDcdbEncryptAttributesConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if encryptAttributesConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbEncryptAttributesConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if encryptAttributesConfig.EncryptStatus != nil {
		_ = d.Set("encrypt_enabled", encryptAttributesConfig.EncryptStatus)
	}

	return nil
}

func resourceTencentCloudDcdbEncryptAttributesConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyDBEncryptAttributesRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("encrypt_enabled") {
		if v, ok := d.GetOkExists("encrypt_enabled"); ok {
			request.EncryptEnabled = helper.IntInt64(v.(int))
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
		log.Printf("[CRITAL]%s update dcdb encryptAttributesConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbEncryptAttributesConfigRead(d, meta)
}

func resourceTencentCloudDcdbEncryptAttributesConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
