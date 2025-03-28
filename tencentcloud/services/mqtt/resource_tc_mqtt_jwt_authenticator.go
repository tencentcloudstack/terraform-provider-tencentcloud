package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMqttJwtAuthenticator() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMqttJwtAuthenticatorCreate,
		Read:   resourceTencentCloudMqttJwtAuthenticatorRead,
		Update: resourceTencentCloudMqttJwtAuthenticatorUpdate,
		Delete: resourceTencentCloudMqttJwtAuthenticatorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Algorithm. hmac-based, public-key.",
			},

			"from": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pass the key of JWT when connecting the device; Username - passed using the username field; Password - Pass using password field.",
			},

			"secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secret.",
			},

			"public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Public key.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark.",
			},
		},
	}
}

func resourceTencentCloudMqttJwtAuthenticatorCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_jwt_authenticator.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewCreateJWTAuthenticatorRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("algorithm"); ok {
		request.Algorithm = helper.String(v.(string))
	}

	if v, ok := d.GetOk("from"); ok {
		request.From = helper.String(v.(string))
	}

	if v, ok := d.GetOk("secret"); ok {
		request.Secret = helper.String(v.(string))
	}

	if v, ok := d.GetOk("public_key"); ok {
		request.PublicKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	request.Status = helper.String("open")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().CreateJWTAuthenticatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mqtt jwt authenticator failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt jwt authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(instanceId)

	return resourceTencentCloudMqttJwtAuthenticatorRead(d, meta)
}

func resourceTencentCloudMqttJwtAuthenticatorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_jwt_authenticator.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeMqttJwtAuthenticatorById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt_jwt_authenticator` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.Config != nil {
		var configMap map[string]interface{}
		e := json.Unmarshal([]byte(*respData.Config), &configMap)
		if e != nil {
			return fmt.Errorf("Failed to parse config content: %s", e.Error())
		}

		if v, ok := configMap["algorithm"].(string); ok && v != "" {
			_ = d.Set("algorithm", v)
		}

		if v, ok := configMap["from"].(string); ok && v != "" {
			_ = d.Set("from", v)
		}

		if v, ok := configMap["secret"].(string); ok && v != "" {
			_ = d.Set("secret", v)
		}

		if v, ok := configMap["public_key"].(string); ok && v != "" {
			_ = d.Set("publicKey", v)
		}
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	return nil
}

func resourceTencentCloudMqttJwtAuthenticatorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_jwt_authenticator.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	request := mqttv20240516.NewModifyJWTAuthenticatorRequest()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("algorithm"); ok {
		request.Algorithm = helper.String(v.(string))
	}

	if v, ok := d.GetOk("from"); ok {
		request.From = helper.String(v.(string))
	}

	if v, ok := d.GetOk("secret"); ok {
		request.Secret = helper.String(v.(string))
	}

	if v, ok := d.GetOk("public_key"); ok {
		request.PublicKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("text"); ok {
		request.Text = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ModifyJWTAuthenticatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update mqtt jwt authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudMqttJwtAuthenticatorRead(d, meta)
}

func resourceTencentCloudMqttJwtAuthenticatorDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_jwt_authenticator.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewDeleteAuthenticatorRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	request.Type = helper.String("JWT")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeleteAuthenticatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete mqtt jwt authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
