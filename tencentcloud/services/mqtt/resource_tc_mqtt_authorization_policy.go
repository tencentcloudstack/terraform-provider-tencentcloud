package mqtt

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMqttAuthorizationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMqttAuthorizationPolicyCreate,
		Read:   resourceTencentCloudMqttAuthorizationPolicyRead,
		Update: resourceTencentCloudMqttAuthorizationPolicyUpdate,
		Delete: resourceTencentCloudMqttAuthorizationPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "MQTT instance ID.",
			},

			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy name, cannot be empty, 3-64 characters, supports Chinese characters, letters, numbers, \"-\" and \"_\".",
			},

			"policy_version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy version, default is 1, currently only 1 is supported.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The strategy priority, the smaller the higher the priority, cannot be repeated.",
			},

			"effect": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Decision: allow/deny.",
			},

			"actions": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Operation - connect: connect; pub: publish; sub: subscribe.",
			},

			"retain": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Condition - Reserved message 1, match reserved message; 2, match unreserved message, 3. match reserved and unreserved message.",
			},

			"qos": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Condition: Quality of Service 0: At most once 1: At least once 2: Exactly once.",
			},

			"resources": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resources, requiring matching subscriptions.",
			},

			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Condition - Username.",
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Condition - Client ID, supports regular expressions.",
			},

			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Condition - Client IP address, supports IP or CIDR.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks, up to 128 characters.",
			},

			// computed
			"policy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Authorization policy rule id.",
			},
		},
	}
}

func resourceTencentCloudMqttAuthorizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_authorization_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = mqttv20240516.NewCreateAuthorizationPolicyRequest()
		response   = mqttv20240516.NewCreateAuthorizationPolicyResponse()
		instanceId string
		id         int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("policy_name"); ok {
		request.PolicyName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("policy_version"); ok {
		request.PolicyVersion = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("effect"); ok {
		request.Effect = helper.String(v.(string))
	}

	if v, ok := d.GetOk("actions"); ok {
		request.Actions = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("retain"); ok {
		request.Retain = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("qos"); ok {
		request.Qos = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resources"); ok {
		request.Resources = helper.String(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		request.Username = helper.String(v.(string))
	}

	if v, ok := d.GetOk("client_id"); ok {
		request.ClientId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip"); ok {
		request.Ip = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().CreateAuthorizationPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.RetryableError(fmt.Errorf("Create mqtt authorization policy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create mqtt authorization policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Id == nil {
		return fmt.Errorf("Id is nil.")
	}

	id = *response.Response.Id
	d.SetId(strings.Join([]string{instanceId, helper.Int64ToStr(id)}, tccommon.FILED_SP))

	return resourceTencentCloudMqttAuthorizationPolicyRead(d, meta)
}

func resourceTencentCloudMqttAuthorizationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_authorization_policy.read")()
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
	id := idSplit[1]

	respData, err := service.DescribeMqttAuthorizationPolicyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mqtt_authorization_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.Data != nil && len(respData.Data) > 0 {
		idInt, _ := strconv.ParseInt(id, 10, 64)
		for _, item := range respData.Data {
			if item.Id != nil && *item.Id == idInt {
				if item.InstanceId != nil {
					_ = d.Set("instance_id", item.InstanceId)
				}

				if item.PolicyName != nil {
					_ = d.Set("policy_name", item.PolicyName)
				}

				if item.Version != nil {
					_ = d.Set("policy_version", item.Version)
				}

				if item.Priority != nil {
					_ = d.Set("priority", item.Priority)
				}

				if item.Effect != nil {
					_ = d.Set("effect", item.Effect)
				}

				if item.Actions != nil {
					_ = d.Set("actions", item.Actions)
				}

				if item.Retain != nil {
					_ = d.Set("retain", item.Retain)
				}

				if item.Qos != nil {
					_ = d.Set("qos", item.Qos)
				}

				if item.Resources != nil {
					_ = d.Set("resources", item.Resources)
				}

				if item.Username != nil {
					_ = d.Set("username", item.Username)
				}

				if item.ClientId != nil {
					_ = d.Set("client_id", item.ClientId)
				}

				if item.Ip != nil {
					_ = d.Set("ip", item.Ip)
				}

				if item.Remark != nil {
					_ = d.Set("remark", item.Remark)
				}

				if item.Id != nil {
					_ = d.Set("policy_id", item.Id)
				}
			}
		}
	}

	return nil
}

func resourceTencentCloudMqttAuthorizationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_authorization_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = mqttv20240516.NewModifyAuthorizationPolicyRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	id := idSplit[1]

	if d.HasChange("policy_name") {
		if v, ok := d.GetOk("policy_name"); ok {
			request.PolicyName = helper.String(v.(string))
		}
	}

	if d.HasChange("policy_version") {
		if v, ok := d.GetOkExists("policy_version"); ok {
			request.PolicyVersion = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("priority") {
		if v, ok := d.GetOkExists("priority"); ok {
			request.Priority = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("effect") {
		if v, ok := d.GetOk("effect"); ok {
			request.Effect = helper.String(v.(string))
		}
	}

	if d.HasChange("actions") {
		if v, ok := d.GetOk("actions"); ok {
			request.Actions = helper.String(v.(string))
		}
	}

	if d.HasChange("retain") {
		if v, ok := d.GetOkExists("retain"); ok {
			request.Retain = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("qos") {
		if v, ok := d.GetOk("qos"); ok {
			request.Qos = helper.String(v.(string))
		}
	}

	if d.HasChange("resources") {
		if v, ok := d.GetOk("resources"); ok {
			request.Resources = helper.String(v.(string))
		}
	}

	if d.HasChange("username") {
		if v, ok := d.GetOk("username"); ok {
			request.Username = helper.String(v.(string))
		}
	}

	if d.HasChange("client_id") {
		if v, ok := d.GetOk("client_id"); ok {
			request.ClientId = helper.String(v.(string))
		}
	}

	if d.HasChange("ip") {
		if v, ok := d.GetOk("ip"); ok {
			request.Ip = helper.String(v.(string))
		}
	}

	if d.HasChange("remark") {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
	}

	request.Id = helper.StrToInt64Point(id)
	request.InstanceId = &instanceId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().ModifyAuthorizationPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update mqtt authorization policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudMqttAuthorizationPolicyRead(d, meta)
}

func resourceTencentCloudMqttAuthorizationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mqtt_authorization_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = mqttv20240516.NewDeleteAuthorizationPolicyRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	id := idSplit[1]

	request.InstanceId = &instanceId
	request.Id = helper.StrToInt64Point(id)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMqttV20240516Client().DeleteAuthorizationPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete mqtt authorization policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
