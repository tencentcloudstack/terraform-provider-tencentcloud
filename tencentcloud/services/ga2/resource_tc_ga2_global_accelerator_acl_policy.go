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

func ResourceTencentCloudGa2GlobalAcceleratorAclPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2GlobalAcceleratorAclPolicyCreate,
		Read:   resourceTencentCloudGa2GlobalAcceleratorAclPolicyRead,
		Update: resourceTencentCloudGa2GlobalAcceleratorAclPolicyUpdate,
		Delete: resourceTencentCloudGa2GlobalAcceleratorAclPolicyDelete,
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
				Description: "Global accelerator instance ID this ACL policy belongs to.",
			},
			"default_action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Default traffic action. Enumerated values: `ACCEPT` (allow all traffic by default), `DROP` (deny all traffic by default).",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ACL policy state. Enumerated values: `OPEN` (enabled), `CLOSE` (disabled).",
			},

			// Computed
			"global_accelerator_acl_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ACL policy ID.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Async task ID returned by the most recent write call (Create / Modify / Delete).",
			},
		},
	}
}

func resourceTencentCloudGa2GlobalAcceleratorAclPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = ga2v20250115.NewCreateGlobalAcceleratorAclPolicyRequest()
		response = ga2v20250115.NewCreateGlobalAcceleratorAclPolicyResponse()
		gaId     string
		policyId string
		taskId   string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		gaId = v.(string)
		request.GlobalAcceleratorId = helper.String(gaId)
	}

	if v, ok := d.GetOk("default_action"); ok {
		request.DefaultAction = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateGlobalAcceleratorAclPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2_global_accelerator_acl_policy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2_global_accelerator_acl_policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.GlobalAcceleratorAclPolicyId == nil || *response.Response.GlobalAcceleratorAclPolicyId == "" {
		log.Printf("[CRITAL]%s create ga2_global_accelerator_acl_policy failed, id=%s, GlobalAcceleratorAclPolicyId is nil or empty", logId, d.Id())
		return fmt.Errorf("create ga2_global_accelerator_acl_policy failed, GlobalAcceleratorAclPolicyId is nil or empty")
	}
	policyId = *response.Response.GlobalAcceleratorAclPolicyId

	if response.Response.TaskId == nil {
		return fmt.Errorf("create ga2_global_accelerator_acl_policy failed, TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(gaId + tccommon.FILED_SP + policyId)
	_ = d.Set("task_id", taskId)

	return resourceTencentCloudGa2GlobalAcceleratorAclPolicyRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorAclPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, policyId, err := parseGa2AclPolicyCompositeId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeGa2GlobalAcceleratorAclPolicyById(ctx, gaId, policyId)
	if err != nil {
		if !d.IsNewResource() && IsGa2ResourceNotFoundError(err) {
			log.Printf("[CRUD] ga2_global_accelerator_acl_policy id=%s", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] ga2_global_accelerator_acl_policy id=%s", d.Id())
		if d.IsNewResource() {
			return fmt.Errorf("ga2_global_accelerator_acl_policy [%s] not found after creation", d.Id())
		}
		d.SetId("")
		return nil
	}

	if respData.GlobalAcceleratorAclPolicyId != nil {
		_ = d.Set("global_accelerator_acl_policy_id", respData.GlobalAcceleratorAclPolicyId)
	}

	if respData.DefaultAction != nil {
		_ = d.Set("default_action", respData.DefaultAction)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	// global_accelerator_id is ForceNew and recovered from the composite ID.
	_ = d.Set("global_accelerator_id", gaId)

	return nil
}

func resourceTencentCloudGa2GlobalAcceleratorAclPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	gaId, policyId, err := parseGa2AclPolicyCompositeId(d.Id())
	if err != nil {
		return err
	}

	// default_action and global_accelerator_id are ForceNew and must NOT be sent on Modify.
	if !d.HasChange("status") {
		return resourceTencentCloudGa2GlobalAcceleratorAclPolicyRead(d, meta)
	}

	request := ga2v20250115.NewModifyGlobalAcceleratorAclPolicyRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.GlobalAcceleratorAclPolicyId = helper.String(policyId)

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyGlobalAcceleratorAclPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2_global_accelerator_acl_policy failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2_global_accelerator_acl_policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	_ = d.Set("task_id", taskId)

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceTencentCloudGa2GlobalAcceleratorAclPolicyRead(d, meta)
}

func resourceTencentCloudGa2GlobalAcceleratorAclPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ga2v20250115.NewDeleteGlobalAcceleratorAclPolicyRequest()
	)

	gaId, policyId, err := parseGa2AclPolicyCompositeId(d.Id())
	if err != nil {
		return err
	}

	request.GlobalAcceleratorId = helper.String(gaId)
	request.GlobalAcceleratorAclPolicyId = helper.String(policyId)

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteGlobalAcceleratorAclPolicyWithContext(ctx, request)
		if e != nil {
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound" {
					return nil
				}
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2_global_accelerator_acl_policy failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2_global_accelerator_acl_policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// ResourceNotFound short-circuit returned nil inside the retry block; when that happens
	// taskId stays empty and there is no task to poll.
	if taskId != "" {
		_ = d.Set("task_id", taskId)
		service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	return nil
}

// parseGa2AclPolicyCompositeId splits the composite resource ID into its two components.
func parseGa2AclPolicyCompositeId(id string) (gaId, policyId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<global_accelerator_acl_policy_id>", id, tccommon.FILED_SP)
		return
	}
	gaId, policyId = parts[0], parts[1]
	if gaId == "" || policyId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}
