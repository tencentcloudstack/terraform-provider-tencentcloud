package mqtt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMqttInstanceDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMqttInstanceDetailRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID.",
			},

			// computed
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance type. BASIC- Basic Edition; PRO- professional edition; PLATINUM- Platinum version.",
			},

			"topic_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Topic num.",
			},

			"topic_num_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of instance topics.",
			},

			"tps_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Elastic TPS current limit value.",
			},

			"created_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation time, millisecond timestamp.",
			},

			"remark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Remark.",
			},

			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance status. RUNNING- In operation; MAINTAINING- Under Maintenance; ABNORMAL- abnormal; OVERDUE- Arrears of fees; DESTROYED- Deleted; CREATING- Creating in progress; MODIFYING- In the process of transformation; CREATE_FAILURE- Creation failed; MODIFY_FAILURE- Transformation failed; DELETING- deleting.",
			},

			"sku_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Product specifications.",
			},

			"max_subscription_per_client": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of subscriptions per client.",
			},

			"authorization_policy_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Limit on the number of authorization rules.",
			},

			"client_num_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of client connections online.",
			},

			"device_certificate_provision_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Client certificate registration method: JITP: Automatic Registration; API: Manually register through API.",
			},

			"automatic_activation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is it automatically activated when registering device certificates automatically.",
			},

			"renew_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether to renew automatically. Only the annual and monthly package cluster is effective. 1: Automatic renewal; 0: Non automatic renewal.",
			},

			"pay_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Billing mode, POSTPAID, pay as you go PREPAID, annual and monthly package.",
			},

			"expiry_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Expiration time, millisecond level timestamp.",
			},

			"destroy_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Pre destruction time, millisecond timestamp.",
			},

			"x509_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TLS, Unidirectional authentication mTLS, bidirectional authentication BYOC; One machine, one certificate.",
			},

			"max_ca_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum Ca quota.",
			},

			"registration_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate registration code.",
			},

			"max_subscription": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of subscriptions in the cluster.",
			},

			"authorization_policy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Authorization Policy Switch.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMqttInstanceDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mqtt_instance_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = MqttService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	var respData *mqttv20240516.DescribeInstanceResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMqttInstanceDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if respData.InstanceType != nil {
		_ = d.Set("instance_type", respData.InstanceType)
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
	}

	if respData.TopicNum != nil {
		_ = d.Set("topic_num", respData.TopicNum)
	}

	if respData.TopicNumLimit != nil {
		_ = d.Set("topic_num_limit", respData.TopicNumLimit)
	}

	if respData.TpsLimit != nil {
		_ = d.Set("tps_limit", respData.TpsLimit)
	}

	if respData.CreatedTime != nil {
		_ = d.Set("created_time", respData.CreatedTime)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.InstanceStatus != nil {
		_ = d.Set("instance_status", respData.InstanceStatus)
	}

	if respData.SkuCode != nil {
		_ = d.Set("sku_code", respData.SkuCode)
	}

	if respData.MaxSubscriptionPerClient != nil {
		_ = d.Set("max_subscription_per_client", respData.MaxSubscriptionPerClient)
	}

	if respData.AuthorizationPolicyLimit != nil {
		_ = d.Set("authorization_policy_limit", respData.AuthorizationPolicyLimit)
	}

	if respData.ClientNumLimit != nil {
		_ = d.Set("client_num_limit", respData.ClientNumLimit)
	}

	if respData.DeviceCertificateProvisionType != nil {
		_ = d.Set("device_certificate_provision_type", respData.DeviceCertificateProvisionType)
	}

	if respData.AutomaticActivation != nil {
		_ = d.Set("automatic_activation", respData.AutomaticActivation)
	}

	if respData.RenewFlag != nil {
		_ = d.Set("renew_flag", respData.RenewFlag)
	}

	if respData.PayMode != nil {
		_ = d.Set("pay_mode", respData.PayMode)
	}

	if respData.ExpiryTime != nil {
		_ = d.Set("expiry_time", respData.ExpiryTime)
	}

	if respData.DestroyTime != nil {
		_ = d.Set("destroy_time", respData.DestroyTime)
	}

	if respData.X509Mode != nil {
		_ = d.Set("x509_mode", respData.X509Mode)
	}

	if respData.MaxCaNum != nil {
		_ = d.Set("max_ca_num", respData.MaxCaNum)
	}

	if respData.RegistrationCode != nil {
		_ = d.Set("registration_code", respData.RegistrationCode)
	}

	if respData.MaxSubscription != nil {
		_ = d.Set("max_subscription", respData.MaxSubscription)
	}

	if respData.AuthorizationPolicy != nil {
		_ = d.Set("authorization_policy", respData.AuthorizationPolicy)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
