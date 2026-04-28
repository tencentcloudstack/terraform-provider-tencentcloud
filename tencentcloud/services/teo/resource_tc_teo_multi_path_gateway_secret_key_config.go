package teo

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigCreate,
		Read:   resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigRead,
		Update: resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigUpdate,
		Delete: resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Multi-path gateway secret key, base64 string, the string length before encoding is 32-48 characters.",
			},
		},
	}
}

func resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_secret_key_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	zoneId := d.Get("zone_id").(string)
	d.SetId(zoneId)

	return resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigUpdate(d, meta)
}

func resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_secret_key_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  = d.Id()
	)

	secretKey, err := service.DescribeTeoMultiPathGatewaySecretKeyById(ctx, zoneId)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFound") {
			log.Printf("[WARN]%s resource `tencentcloud_teo_multi_path_gateway_secret_key_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if secretKey == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_multi_path_gateway_secret_key_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("secret_key", *secretKey)

	return nil
}

func resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_secret_key_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  = d.Id()
	)

	// First, call Describe API to check if the secret key already exists
	existingKey, err := service.DescribeTeoMultiPathGatewaySecretKeyById(ctx, zoneId)
	if err != nil {
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return err
		}
		// ResourceNotFound means the key does not exist, treat as nil
		existingKey = nil
	}

	if existingKey != nil {
		// Key exists: call CreateMultiPathGatewaySecretKey API to replace the key
		request := teov20220901.NewCreateMultiPathGatewaySecretKeyRequest()
		request.ZoneId = helper.String(zoneId)

		if v, ok := d.GetOk("secret_key"); ok {
			request.SecretKey = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateMultiPathGatewaySecretKeyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s create teo multi path gateway secret key failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	} else {
		// Key does not exist: call ModifyMultiPathGatewaySecretKey API to set the key
		request := teov20220901.NewModifyMultiPathGatewaySecretKeyRequest()
		request.ZoneId = helper.String(zoneId)

		if v, ok := d.GetOk("secret_key"); ok {
			request.SecretKey = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyMultiPathGatewaySecretKeyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo multi path gateway secret key failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigRead(d, meta)
}

func resourceTencentCloudTeoMultiPathGatewaySecretKeyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_secret_key_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
