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

func ResourceTencentCloudMqttCaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMqttCaCertificateCreate,
		Read:   resourceTencentCloudMqttCaCertificateRead,
		Update: resourceTencentCloudMqttCaCertificateUpdate,
		Delete: resourceTencentCloudMqttCaCertificateDelete,
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

			"ca_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CA certificate.",
			},

			"verification_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Verification certificate.",
			},

			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Certificate format, Default is PEM.",
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ACTIVE", "INACTIVE"}),
				Description:  "Certificate status, Default is ACTIVE.\n  ACTIVE activation;\n  INACTIVE not active.",
			},

			"ca_sn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate serial number.",
			},

			"ca_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate common name.",
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

func resourceTencentCloudMqttCaCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_ca_certificate.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewRegisterCaCertificateRequest()
		response   = mqttv20240516.NewRegisterCaCertificateResponse()
		instanceId string
		caSn       string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("ca_certificate"); ok {
		request.CaCertificate = helper.String(v.(string))
	}

	if v, ok := d.GetOk("verification_certificate"); ok {
		request.VerificationCertificate = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().RegisterCaCertificateWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mqtt ca certificate failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt ca certificate failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.CaSn == nil {
		return fmt.Errorf("CaSn is nil.")
	}

	caSn = *response.Response.CaSn
	d.SetId(strings.Join([]string{instanceId, caSn}, tccommon.FILED_SP))

	return resourceTencentCloudMqttCaCertificateRead(d, meta)
}

func resourceTencentCloudMqttCaCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_ca_certificate.read")()
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
	caSn := idSplit[1]

	respData1, err := service.DescribeMqttCaCertificateById(ctx, instanceId, caSn)
	if err != nil {
		return err
	}

	if respData1 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt_ca_certificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("ca_sn", caSn)

	if respData1.CaCertificate != nil {
		_ = d.Set("ca_certificate", respData1.CaCertificate)
	}

	if respData1.Format != nil {
		_ = d.Set("format", respData1.Format)
	}

	if respData1.Status != nil {
		_ = d.Set("status", respData1.Status)
	}

	if respData1.CreatedTime != nil {
		_ = d.Set("created_time", respData1.CreatedTime)
	}

	if respData1.UpdateTime != nil {
		_ = d.Set("update_time", respData1.UpdateTime)
	}

	if respData1.NotBeforeTime != nil {
		_ = d.Set("not_before_time", respData1.NotBeforeTime)
	}

	if respData1.NotAfterTime != nil {
		_ = d.Set("not_after_time", respData1.NotAfterTime)
	}

	respData2, err := service.DescribeMqttCaCertificatesById(ctx, instanceId, caSn)
	if err != nil {
		return err
	}

	if respData2 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt_ca_certificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData2.VerificationCertificate != nil {
		_ = d.Set("verification_certificate", respData2.VerificationCertificate)
	}

	return nil
}

func resourceTencentCloudMqttCaCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_ca_certificate.update")()
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
	caSn := idSplit[1]

	if d.HasChange("status") {
		var status string
		if v, ok := d.GetOk("status"); ok {
			status = v.(string)
		}

		if status == "ACTIVE" {
			request := mqttv20240516.NewActivateCaCertificateRequest()
			request.InstanceId = &instanceId
			request.CaSn = &caSn
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ActivateCaCertificateWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update mqtt ca certificate activate failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else if status == "INACTIVE" {
			request := mqttv20240516.NewDeactivateCaCertificateRequest()
			request.InstanceId = &instanceId
			request.CaSn = &caSn
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeactivateCaCertificateWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update mqtt ca certificate deactivate failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else {
			return fmt.Errorf("`status` only support `ACTIVE` and `INACTIVE`.")
		}
	}

	return resourceTencentCloudMqttCaCertificateRead(d, meta)
}

func resourceTencentCloudMqttCaCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_ca_certificate.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = mqttv20240516.NewDeleteCaCertificateRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	caSn := idSplit[1]

	respData, err := service.DescribeMqttCaCertificateById(ctx, instanceId, caSn)
	if err != nil {
		return err
	}

	if respData == nil {
		return nil
	}

	if respData.Status != nil {
		if *respData.Status == "ACTIVE" {
			DeactivateRequest := mqttv20240516.NewDeactivateCaCertificateRequest()
			DeactivateRequest.InstanceId = &instanceId
			DeactivateRequest.CaSn = &caSn
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeactivateCaCertificateWithContext(ctx, DeactivateRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, DeactivateRequest.GetAction(), DeactivateRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update mqtt ca certificate deactivate failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			request.InstanceId = &instanceId
			request.CaSn = &caSn
			reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeleteCaCertificateWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s delete mqtt ca certificate failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else if *respData.Status == "INACTIVE" {
			request.InstanceId = &instanceId
			request.CaSn = &caSn
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeleteCaCertificateWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s delete mqtt ca certificate failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		} else {
			return fmt.Errorf("The current certificate status is %s and cannot be deleted.", *respData.Status)
		}

		return nil
	}

	return fmt.Errorf("Failed to obtain certificate status, unable to perform destruction operation.")
}
