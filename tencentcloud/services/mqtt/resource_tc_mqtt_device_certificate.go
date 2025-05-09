package mqtt

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMqttDeviceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMqttDeviceCertificateCreate,
		Read:   resourceTencentCloudMqttDeviceCertificateRead,
		Update: resourceTencentCloudMqttDeviceCertificateUpdate,
		Delete: resourceTencentCloudMqttDeviceCertificateDelete,
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

			"device_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Device certificate.",
			},

			"ca_sn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Associated CA certificate SN.",
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Client ID.",
			},

			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Certificate format, Default is PEM.",
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ACTIVE", "INACTIVE"}),
				Description:  "Certificate status, Default is ACTIVE.\\n  ACTIVE activation;\\n  INACTIVE not active.",
			},

			"device_certificate_sn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Equipment certificate serial number.",
			},

			"device_certificate_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate common name.",
			},

			"certificate_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate source.",
			},

			"created_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Certificate create time.",
			},

			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Certificate update time.",
			},

			"not_before_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Certificate effective start date.",
			},

			"not_after_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Certificate expiring date.",
			},
		},
	}
}

func resourceTencentCloudMqttDeviceCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_device_certificate.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request             = mqttv20240516.NewRegisterDeviceCertificateRequest()
		response            = mqttv20240516.NewRegisterDeviceCertificateResponse()
		instanceId          string
		deviceCertificateSn string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("device_certificate"); ok {
		request.DeviceCertificate = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ca_sn"); ok {
		request.CaSn = helper.String(v.(string))
	}

	if v, ok := d.GetOk("client_id"); ok {
		request.ClientId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().RegisterDeviceCertificateWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mqtt device certificate failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt device certificate failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.DeviceCertificateSn == nil {
		return fmt.Errorf("DeviceCertificateSn is nil.")
	}

	deviceCertificateSn = *response.Response.DeviceCertificateSn
	d.SetId(strings.Join([]string{instanceId, deviceCertificateSn}, tccommon.FILED_SP))

	return resourceTencentCloudMqttDeviceCertificateRead(d, meta)
}

func resourceTencentCloudMqttDeviceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_device_certificate.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	deviceCertificateSn := idSplit[1]

	respData, err := service.DescribeMqttDeviceCertificateById(ctx, instanceId, deviceCertificateSn)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt_device_certificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.DeviceCertificate != nil {
		_ = d.Set("device_certificate", respData.DeviceCertificate)
	}

	if respData.CaSn != nil {
		_ = d.Set("ca_sn", respData.CaSn)
	}

	if respData.ClientId != nil {
		_ = d.Set("client_id", respData.ClientId)
	}

	if respData.Format != nil {
		_ = d.Set("format", respData.Format)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.DeviceCertificateSn != nil {
		_ = d.Set("device_certificate_sn", respData.DeviceCertificateSn)
	}

	if respData.DeviceCertificateCn != nil {
		_ = d.Set("device_certificate_cn", respData.DeviceCertificateCn)
	}

	if respData.CertificateSource != nil {
		_ = d.Set("certificate_source", respData.CertificateSource)
	}

	if respData.CreatedTime != nil {
		_ = d.Set("created_time", respData.CreatedTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.NotBeforeTime != nil {
		_ = d.Set("not_before_time", respData.NotBeforeTime)
	}

	if respData.NotAfterTime != nil {
		_ = d.Set("not_after_time", respData.NotAfterTime)
	}

	return nil
}

func resourceTencentCloudMqttDeviceCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_device_certificate.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	deviceCertificateSn := idSplit[1]

	if d.HasChange("status") {
		var status string
		if v, ok := d.GetOk("status"); ok {
			status = v.(string)
		}

		if status == "ACTIVE" {
			request := mqttv20240516.NewActivateDeviceCertificateRequest()
			request.InstanceId = &instanceId
			request.DeviceCertificateSn = &deviceCertificateSn
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ActivateDeviceCertificateWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if reqErr != nil {
				log.Printf("[CRITAL]%s update mqtt device certificate activate failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else if status == "INACTIVE" {
			request := mqttv20240516.NewDeactivateDeviceCertificateRequest()
			request.InstanceId = &instanceId
			request.DeviceCertificateSn = &deviceCertificateSn
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeactivateDeviceCertificateWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update mqtt device certificate deactivate failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else {
			return fmt.Errorf("`status` only support `ACTIVE` and `INACTIVE`.")
		}
	}

	return resourceTencentCloudMqttDeviceCertificateRead(d, meta)
}

func resourceTencentCloudMqttDeviceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_device_certificate.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = mqttv20240516.NewDeleteDeviceCertificateRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	deviceCertificateSn := idSplit[1]

	respData, err := service.DescribeMqttDeviceCertificateById(ctx, instanceId, deviceCertificateSn)
	if err != nil {
		return err
	}

	if respData == nil {
		return nil
	}

	if respData.Status != nil {
		if *respData.Status == "ACTIVE" {
			DeactivateRequest := mqttv20240516.NewDeactivateDeviceCertificateRequest()
			DeactivateRequest.InstanceId = &instanceId
			DeactivateRequest.DeviceCertificateSn = &deviceCertificateSn
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeactivateDeviceCertificateWithContext(ctx, DeactivateRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, DeactivateRequest.GetAction(), DeactivateRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update mqtt device certificate deactivate failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			request.InstanceId = &instanceId
			request.DeviceCertificateSn = &deviceCertificateSn
			reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeleteDeviceCertificateWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s delete mqtt device certificate failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else if *respData.Status == "INACTIVE" {
			request.InstanceId = &instanceId
			request.DeviceCertificateSn = &deviceCertificateSn
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeleteDeviceCertificateWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s delete mqtt device certificate failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else {
			return fmt.Errorf("The current certificate status is %s and cannot be deleted.", *respData.Status)
		}

		return nil
	}

	return fmt.Errorf("Failed to obtain certificate status, unable to perform destruction operation.")
}
