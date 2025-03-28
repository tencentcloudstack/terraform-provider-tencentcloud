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

func ResourceTencentCloudMqttJwksAuthenticator() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMqttJwksAuthenticatorCreate,
		Read:   resourceTencentCloudMqttJwksAuthenticatorRead,
		Update: resourceTencentCloudMqttJwksAuthenticatorUpdate,
		Delete: resourceTencentCloudMqttJwksAuthenticatorDelete,
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

			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JWKS endpoint.",
			},

			"refresh_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "JWKS refresh interval. unit: s.",
			},

			"text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JWKS text.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark.",
			},

			"from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pass the key of JWT when connecting the device; Username - passed using the username field; Password - Pass using password field.",
			},
		},
	}
}

func resourceTencentCloudMqttJwksAuthenticatorCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_jwks_authenticator.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewCreateJWKSAuthenticatorRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("endpoint"); ok {
		request.Endpoint = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("refresh_interval"); ok {
		request.RefreshInterval = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("text"); ok {
		request.Text = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("from"); ok {
		request.From = helper.String(v.(string))
	}

	request.Status = helper.String("open")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().CreateJWKSAuthenticatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mqtt jwks authenticator failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt jwks authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(instanceId)

	return resourceTencentCloudMqttJwksAuthenticatorRead(d, meta)
}

func resourceTencentCloudMqttJwksAuthenticatorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_jwks_authenticator.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeMqttJwksAuthenticatorById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt_jwks_authenticator` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.Config != nil {
		var configMap map[string]interface{}
		e := json.Unmarshal([]byte(*respData.Config), &configMap)
		if e != nil {
			return fmt.Errorf("Failed to parse config content: %s", e.Error())
		}

		if v, ok := configMap["endpoint"].(string); ok && v != "" {
			_ = d.Set("endpoint", v)
		}

		if v, ok := configMap["refreshInterval"].(float64); ok {
			_ = d.Set("refresh_interval", int(v))
		}

		if v, ok := configMap["text"].(string); ok && v != "" {
			_ = d.Set("text", v)
		}

		if v, ok := configMap["from"].(string); ok && v != "" {
			_ = d.Set("from", v)
		}
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	return nil
}

func resourceTencentCloudMqttJwksAuthenticatorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_jwks_authenticator.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	request := mqttv20240516.NewModifyJWKSAuthenticatorRequest()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("endpoint"); ok {
		request.Endpoint = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("refresh_interval"); ok {
		request.RefreshInterval = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("text"); ok {
		request.Text = helper.String(v.(string))
	}

	if v, ok := d.GetOk("from"); ok {
		request.From = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	request.Status = helper.String("open")
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ModifyJWKSAuthenticatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update mqtt jwks authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudMqttJwksAuthenticatorRead(d, meta)
}

func resourceTencentCloudMqttJwksAuthenticatorDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_jwks_authenticator.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewDeleteAuthenticatorRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	request.Type = helper.String("JWKS")

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
		log.Printf("[CRITAL]%s delete mqtt jwks authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
