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

func ResourceTencentCloudMqttHttpAuthenticator() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMqttHttpAuthenticatorCreate,
		Read:   resourceTencentCloudMqttHttpAuthenticatorRead,
		Update: resourceTencentCloudMqttHttpAuthenticatorUpdate,
		Delete: resourceTencentCloudMqttHttpAuthenticatorDelete,
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
				Required:    true,
				Description: "JWKS endpoint.",
			},

			"concurrency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum concurrent connections, default 8, range: 1-20.",
			},

			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network request method GET or POST, default POST.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Is the authenticator enabled: open enable; Close close.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark.",
			},

			"connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Connection timeout, unit: seconds, range: 1-30.",
			},

			"read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Request timeout, unit: seconds, range: 1-30.",
			},

			"header": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Forwarding request header.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Header key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Header value.",
						},
					},
				},
			},

			"body": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Forwarding request body.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Body key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Body key.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMqttHttpAuthenticatorCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_http_authenticator.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewCreateHttpAuthenticatorRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("endpoint"); ok {
		request.Endpoint = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("concurrency"); ok {
		request.Concurrency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("method"); ok {
		request.Method = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("connect_timeout"); ok {
		request.ConnectTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("read_timeout"); ok {
		request.ReadTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("header"); ok {
		for _, item := range v.([]interface{}) {
			headerMap := item.(map[string]interface{})
			headerItem := mqttv20240516.HeaderItem{}
			if v, ok := headerMap["key"].(string); ok && v != "" {
				headerItem.Key = helper.String(v)
			}

			if v, ok := headerMap["value"].(string); ok && v != "" {
				headerItem.Value = helper.String(v)
			}

			request.Header = append(request.Header, &headerItem)
		}
	}

	if v, ok := d.GetOk("body"); ok {
		for _, item := range v.([]interface{}) {
			bodyMap := item.(map[string]interface{})
			bodyItem := mqttv20240516.BodyItem{}
			if v, ok := bodyMap["key"].(string); ok && v != "" {
				bodyItem.Key = helper.String(v)
			}

			if v, ok := bodyMap["value"].(string); ok && v != "" {
				bodyItem.Value = helper.String(v)
			}

			request.Body = append(request.Body, &bodyItem)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().CreateHttpAuthenticatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mqtt http authenticator failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt http authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(instanceId)

	return resourceTencentCloudMqttHttpAuthenticatorRead(d, meta)
}

func resourceTencentCloudMqttHttpAuthenticatorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_http_authenticator.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeMqttHttpAuthenticatorById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt_http_authenticator` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

		if v, ok := configMap["concurrency"].(float64); ok {
			_ = d.Set("concurrency", int(v))
		}

		if v, ok := configMap["method"].(string); ok && v != "" {
			_ = d.Set("method", v)
		}

		if v, ok := configMap["connectTimeout"].(float64); ok {
			_ = d.Set("connect_timeout", int(v))
		}

		if v, ok := configMap["readTimeout"].(float64); ok {
			_ = d.Set("read_timeout", int(v))
		}

		if v, ok := configMap["headers"].([]interface{}); ok {
			tmpList := make([]map[string]interface{}, 0)
			for _, item := range v {
				bodyMap := item.(map[string]interface{})
				dMap := map[string]interface{}{}
				if v, ok := bodyMap["key"].(string); ok && v != "" {
					dMap["key"] = v
				}

				if v, ok := bodyMap["value"].(string); ok && v != "" {
					dMap["value"] = v
				}

				tmpList = append(tmpList, dMap)
			}

			_ = d.Set("header", tmpList)
		}

		if v, ok := configMap["body"].([]interface{}); ok {
			tmpList := make([]map[string]interface{}, 0)
			for _, item := range v {
				bodyMap := item.(map[string]interface{})
				dMap := map[string]interface{}{}
				if v, ok := bodyMap["key"].(string); ok && v != "" {
					dMap["key"] = v
				}

				if v, ok := bodyMap["value"].(string); ok && v != "" {
					dMap["value"] = v
				}

				tmpList = append(tmpList, dMap)
			}

			_ = d.Set("body", tmpList)
		}
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	return nil
}

func resourceTencentCloudMqttHttpAuthenticatorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_http_authenticator.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	request := mqttv20240516.NewModifyHttpAuthenticatorRequest()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("endpoint"); ok {
		request.Endpoint = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("concurrency"); ok {
		request.Concurrency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("connect_timeout"); ok {
		request.ConnectTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("read_timeout"); ok {
		request.ReadTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("method"); ok {
		request.Method = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOk("header"); ok {
		for _, item := range v.([]interface{}) {
			headerMap := item.(map[string]interface{})
			headerItem := mqttv20240516.HeaderItem{}
			if v, ok := headerMap["key"].(string); ok && v != "" {
				headerItem.Key = helper.String(v)
			}

			if v, ok := headerMap["value"].(string); ok && v != "" {
				headerItem.Value = helper.String(v)
			}

			request.Header = append(request.Header, &headerItem)
		}
	}

	if v, ok := d.GetOk("body"); ok {
		for _, item := range v.([]interface{}) {
			bodyMap := item.(map[string]interface{})
			bodyItem := mqttv20240516.BodyItem{}
			if v, ok := bodyMap["key"].(string); ok && v != "" {
				bodyItem.Key = helper.String(v)
			}

			if v, ok := bodyMap["value"].(string); ok && v != "" {
				bodyItem.Value = helper.String(v)
			}

			request.Body = append(request.Body, &bodyItem)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ModifyHttpAuthenticatorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update mqtt http authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudMqttHttpAuthenticatorRead(d, meta)
}

func resourceTencentCloudMqttHttpAuthenticatorDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_http_authenticator.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewDeleteAuthenticatorRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	request.Type = helper.String("HTTP")

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
		log.Printf("[CRITAL]%s delete mqtt http authenticator failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
