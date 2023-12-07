package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCfwVpcPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwVpcPolicyCreate,
		Read:   resourceTencentCloudCfwVpcPolicyRead,
		Update: resourceTencentCloudCfwVpcPolicyUpdate,
		Delete: resourceTencentCloudCfwVpcPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access source examplnet:IP/CIDR(192.168.0.2).",
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access source type, the type can be: net, template.",
			},
			"dest_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access purpose example: net:IP/CIDR(192.168.0.2) domain:domain rule, for example*.qq.com.",
			},
			"dest_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access purpose type, the type can be: net, template.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol, optional value:TCP, UDP, ICMP, ANY, HTTP, HTTPS, HTTP/HTTPS, SMTP, SMTPS, SMTP/SMTPS, FTP, DNS, TLS/SSL.",
			},
			"rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(POLICY_RULE_ACTION),
				Description:  "How traffic set in the access control policy passes through the cloud firewall. Value: accept:accept, drop:drop, log:log.",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The port for the access control policy. Value: -1/-1: All ports; 80: port 80.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Describe.",
			},
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      POLICY_ENABLE_TRUE,
				ValidateFunc: validateAllowedStringValue(POLICY_ENABLE),
				Description:  "Rule status, true means enabled, false means disabled. Default is true.",
			},
			"fw_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALL",
				Description: "Firewall instance ID where the rule takes effect. Default is ALL.",
			},
			"uuid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique id corresponding to the rule.",
			},
			"internal_uuid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Uuid used internally, this field is generally not used.",
			},
			"fw_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Firewall name.",
			},
		},
	}
}

func resourceTencentCloudCfwVpcPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_policy.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = cfw.NewAddVpcAcRuleRequest()
		response    = cfw.NewAddVpcAcRuleResponse()
		vpcRuleItem = cfw.VpcRuleItem{}
		uuid        string
	)

	if v, ok := d.GetOk("source_content"); ok {
		vpcRuleItem.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		vpcRuleItem.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_content"); ok {
		vpcRuleItem.DestContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_type"); ok {
		vpcRuleItem.DestType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		vpcRuleItem.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		vpcRuleItem.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		vpcRuleItem.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		vpcRuleItem.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("enable"); ok {
		vpcRuleItem.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fw_group_id"); ok {
		vpcRuleItem.FwGroupId = helper.String(v.(string))
	}

	vpcRuleItem.EdgeId = helper.String("all")

	request.Rules = append(request.Rules, &vpcRuleItem)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().AddVpcAcRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result.Response == nil || len(result.Response.RuleUuids) != 1 {
			e = fmt.Errorf("create cfw vpcPolicy failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw vpcPolicy failed, reason:%+v", logId, err)
		return err
	}

	ruleUuid := *response.Response.RuleUuids[0]
	uuid = strconv.FormatInt(ruleUuid, 10)
	d.SetId(uuid)

	return resourceTencentCloudCfwVpcPolicyRead(d, meta)
}

func resourceTencentCloudCfwVpcPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_policy.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		uuid    = d.Id()
	)

	vpcPolicy, err := service.DescribeCfwVpcPolicyById(ctx, uuid)
	if err != nil {
		return err
	}

	if vpcPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwVpcPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if vpcPolicy.SourceContent != nil {
		_ = d.Set("source_content", vpcPolicy.SourceContent)
	}

	if vpcPolicy.SourceType != nil {
		_ = d.Set("source_type", vpcPolicy.SourceType)
	}

	if vpcPolicy.DestContent != nil {
		_ = d.Set("dest_content", vpcPolicy.DestContent)
	}

	if vpcPolicy.DestType != nil {
		_ = d.Set("dest_type", vpcPolicy.DestType)
	}

	if vpcPolicy.Protocol != nil {
		_ = d.Set("protocol", vpcPolicy.Protocol)
	}

	if vpcPolicy.RuleAction != nil {
		_ = d.Set("rule_action", vpcPolicy.RuleAction)
	}

	if vpcPolicy.Port != nil {
		_ = d.Set("port", vpcPolicy.Port)
	}

	if vpcPolicy.Description != nil {
		_ = d.Set("description", vpcPolicy.Description)
	}

	if vpcPolicy.Enable != nil {
		_ = d.Set("enable", vpcPolicy.Enable)
	}

	if vpcPolicy.FwGroupId != nil {
		_ = d.Set("fw_group_id", vpcPolicy.FwGroupId)
	}

	if vpcPolicy.Uuid != nil {
		_ = d.Set("uuid", vpcPolicy.Uuid)
	}

	if vpcPolicy.InternalUuid != nil {
		_ = d.Set("internal_uuid", vpcPolicy.InternalUuid)
	}

	if vpcPolicy.FwGroupName != nil {
		_ = d.Set("fw_group_name", vpcPolicy.FwGroupName)
	}

	return nil
}

func resourceTencentCloudCfwVpcPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_policy.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		request     = cfw.NewModifyVpcAcRuleRequest()
		vpcRuleItem = cfw.VpcRuleItem{}
		uuid        = d.Id()
	)

	uuidInt, _ := strconv.ParseInt(uuid, 10, 64)
	vpcRuleItem.Uuid = &uuidInt

	if v, ok := d.GetOk("source_content"); ok {
		vpcRuleItem.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		vpcRuleItem.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_content"); ok {
		vpcRuleItem.DestContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dest_type"); ok {
		vpcRuleItem.DestType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		vpcRuleItem.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		vpcRuleItem.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		vpcRuleItem.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		vpcRuleItem.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("enable"); ok {
		vpcRuleItem.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fw_group_id"); ok {
		vpcRuleItem.FwGroupId = helper.String(v.(string))
	}

	vpcRuleItem.EdgeId = helper.String("all")

	request.Rules = append(request.Rules, &vpcRuleItem)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyVpcAcRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw vpcPolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwVpcPolicyRead(d, meta)
}

func resourceTencentCloudCfwVpcPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_policy.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		uuid    = d.Id()
	)

	if err := service.DeleteCfwVpcPolicyById(ctx, uuid); err != nil {
		return err
	}

	return nil
}
