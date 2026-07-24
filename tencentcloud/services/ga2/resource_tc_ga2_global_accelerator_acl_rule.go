package ga2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGa2GlobalAcceleratorAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2GlobalAcceleratorAclRuleCreate,
		Read:   resourceTencentCloudGa2GlobalAcceleratorAclRuleRead,
		Update: resourceTencentCloudGa2GlobalAcceleratorAclRuleUpdate,
		Delete: resourceTencentCloudGa2GlobalAcceleratorAclRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID.",
			},
			"global_accelerator_acl_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL policy ID.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol. Valid values: `TCP`, `UDP`, `ALL`.",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port.",
			},
			"source_cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source CIDR block.",
			},
			"policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Action. Valid values: `ACCEPT` (allow), `DROP` (deny).",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description. Maximum length is 100 bytes.",
			},
			"global_accelerator_acl_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ACL rule ID.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Async task ID for the last operation on this resource.",
			},
		},
	}
}

func resourceTencentCloudGa2GlobalAcceleratorAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = ga2v20250115.NewCreateGlobalAcceleratorAclRuleRequest()
		response = ga2v20250115.NewCreateGlobalAcceleratorAclRuleResponse()
		gaId     string
		policyId string
		ruleId   string
		taskId   string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		gaId = v.(string)
		request.GlobalAcceleratorId = helper.String(gaId)
	}

	if v, ok := d.GetOk("global_accelerator_acl_policy_id"); ok {
		policyId = v.(string)
		request.GlobalAcceleratorAclPolicyId = helper.String(policyId)
	}

	request.AclEntries = []*ga2v20250115.AclEntries{
		{
			Protocol:        helper.String(d.Get("protocol").(string)),
			Port:            helper.String(d.Get("port").(string)),
			SourceCidrBlock: helper.String(d.Get("source_cidr_block").(string)),
			Policy:          helper.String(d.Get("policy").(string)),
			Description:     helper.String(d.Get("description").(string)),
		},
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateGlobalAcceleratorAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 global_accelerator_acl_rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 global_accelerator_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create ga2 global_accelerator_acl_rule, d.Id()=%s", logId, d.Id())

	if response.Response.GlobalAcceleratorAclRuleIds == nil || len(response.Response.GlobalAcceleratorAclRuleIds) == 0 {
		return fmt.Errorf("GlobalAcceleratorAclRuleIds is nil or empty.")
	}
	ruleId = *response.Response.GlobalAcceleratorAclRuleIds[0]

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	_ = d.Set("task_id", taskId)

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{gaId, policyId, ruleId}, tccommon.FILED_SP))
	return resourceTencentCloudGa2GlobalAcceleratorAclRuleRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, policyId, ruleId, err := parseGa2GlobalAcceleratorAclRuleId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeGa2GlobalAcceleratorAclRuleById(ctx, policyId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_global_accelerator_acl_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("global_accelerator_id", gaId)
	_ = d.Set("global_accelerator_acl_policy_id", policyId)

	if respData.GlobalAcceleratorAclRuleId != nil {
		_ = d.Set("global_accelerator_acl_rule_id", respData.GlobalAcceleratorAclRuleId)
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
	}

	if respData.SourceCidrBlock != nil {
		_ = d.Set("source_cidr_block", respData.SourceCidrBlock)
	}

	if respData.Policy != nil {
		_ = d.Set("policy", respData.Policy)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	return nil
}

func resourceTencentCloudGa2GlobalAcceleratorAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	gaId, policyId, ruleId, err := parseGa2GlobalAcceleratorAclRuleId(d.Id())
	if err != nil {
		return err
	}

	request := ga2v20250115.NewModifyGlobalAcceleratorAclRuleRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.GlobalAcceleratorAclPolicyId = helper.String(policyId)
	request.GlobalAcceleratorAclRuleId = helper.String(ruleId)

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		request.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_cidr_block"); ok {
		request.SourceCidrBlock = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy"); ok {
		request.Policy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyGlobalAcceleratorAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2 global_accelerator_acl_rule failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2 global_accelerator_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	_ = d.Set("task_id", taskId)

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceTencentCloudGa2GlobalAcceleratorAclRuleRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ga2v20250115.NewDeleteGlobalAcceleratorAclRuleRequest()
	)

	gaId, policyId, ruleId, err := parseGa2GlobalAcceleratorAclRuleId(d.Id())
	if err != nil {
		return err
	}

	request.GlobalAcceleratorId = helper.String(gaId)
	request.GlobalAcceleratorAclPolicyId = helper.String(policyId)
	request.GlobalAcceleratorAclRuleIds = []*string{helper.String(ruleId)}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteGlobalAcceleratorAclRuleWithContext(ctx, request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound" {
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 global_accelerator_acl_rule failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 global_accelerator_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

func parseGa2GlobalAcceleratorAclRuleId(id string) (gaId, policyId, ruleId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<global_accelerator_acl_policy_id>%s<global_accelerator_acl_rule_id>", id, tccommon.FILED_SP, tccommon.FILED_SP)
		return
	}
	gaId, policyId, ruleId = parts[0], parts[1], parts[2]
	if gaId == "" || policyId == "" || ruleId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}
