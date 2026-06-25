package ga2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGa2ForwardingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2ForwardingPolicyCreate,
		Read:   resourceTencentCloudGa2ForwardingPolicyRead,
		Update: resourceTencentCloudGa2ForwardingPolicyUpdate,
		Delete: resourceTencentCloudGa2ForwardingPolicyDelete,
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
				Description: "Global accelerator instance ID this forwarding policy belongs to.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Listener ID this forwarding policy belongs to.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The domain/host for the forwarding policy.",
			},

			// Computed
			"forwarding_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Forwarding policy ID.",
			},
			"default_host_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is the default host policy for the listener.",
			},
		},
	}
}

func resourceTencentCloudGa2ForwardingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = ga2v20250115.NewCreateForwardingPolicyRequest()
		response   = ga2v20250115.NewCreateForwardingPolicyResponse()
		gaId       string
		listenerId string
		policyId   string
		taskId     string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		gaId = v.(string)
		request.GlobalAcceleratorId = helper.String(gaId)
	}

	if v, ok := d.GetOk("listener_id"); ok {
		listenerId = v.(string)
		request.ListenerId = helper.String(listenerId)
	}

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateForwardingPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 forwarding policy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 forwarding policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ForwardingPolicyId == nil || *response.Response.ForwardingPolicyId == "" {
		return fmt.Errorf("create ga2 forwarding policy failed, ForwardingPolicyId is nil or empty")
	}
	policyId = *response.Response.ForwardingPolicyId

	if response.Response.TaskId == nil {
		return fmt.Errorf("create ga2 forwarding policy failed, TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	d.SetId(strings.Join([]string{gaId, listenerId, policyId}, tccommon.FILED_SP))

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	return resourceTencentCloudGa2ForwardingPolicyRead(d, meta)
}

func resourceTencentCloudGa2ForwardingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, listenerId, policyId, err := parseGa2ForwardingPolicyId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeGa2ForwardingPolicyById(ctx, gaId, listenerId, policyId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_forwarding_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.GlobalAcceleratorId != nil {
		_ = d.Set("global_accelerator_id", respData.GlobalAcceleratorId)
	}

	if respData.ListenerId != nil {
		_ = d.Set("listener_id", respData.ListenerId)
	}

	if respData.ForwardingPolicyId != nil {
		_ = d.Set("forwarding_policy_id", respData.ForwardingPolicyId)
	}

	if respData.Host != nil {
		_ = d.Set("host", respData.Host)
	}

	if respData.DefaultHostFlag != nil {
		_ = d.Set("default_host_flag", respData.DefaultHostFlag)
	}

	return nil
}

func resourceTencentCloudGa2ForwardingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	gaId, listenerId, policyId, err := parseGa2ForwardingPolicyId(d.Id())
	if err != nil {
		return err
	}

	// Body fields supported by ModifyForwardingPolicy (the ID fields are ForceNew, so they
	// cannot trigger Update directly).
	if !d.HasChange("host") {
		return resourceTencentCloudGa2ForwardingPolicyRead(d, meta)
	}

	request := ga2v20250115.NewModifyForwardingPolicyRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.ForwardingPolicyId = helper.String(policyId)

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyForwardingPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2 forwarding policy failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2 forwarding policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceTencentCloudGa2ForwardingPolicyRead(d, meta)
}

func resourceTencentCloudGa2ForwardingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ga2v20250115.NewDeleteForwardingPolicyRequest()
	)

	gaId, listenerId, policyId, err := parseGa2ForwardingPolicyId(d.Id())
	if err != nil {
		return err
	}

	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.ForwardingPolicyId = helper.String(policyId)

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteForwardingPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 forwarding policy failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 forwarding policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

// parseGa2ForwardingPolicyId splits the composite resource ID into its three components.
func parseGa2ForwardingPolicyId(id string) (gaId, listenerId, policyId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<listener_id>%s<forwarding_policy_id>", id, tccommon.FILED_SP, tccommon.FILED_SP)
		return
	}
	gaId, listenerId, policyId = parts[0], parts[1], parts[2]
	if gaId == "" || listenerId == "" || policyId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}
