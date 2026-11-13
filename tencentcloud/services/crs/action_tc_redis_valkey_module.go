package crs

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helperResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"

	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
)

// _ action.Action is intentionally omitted: the factory function
// NewRedisValkeyModule returns action.Action, which already performs the
// implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here.
var _ action.ActionWithConfigure = &RedisValkeyModule{}

// NewRedisValkeyModule is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewRedisValkeyModule() action.Action {
	return &RedisValkeyModule{}
}

// RedisValkeyModule implements action.Action for
// tencentcloud_redis_valkey_module.
type RedisValkeyModule struct {
	fw.ActionWithConfigure
}

// RedisValkeyModuleModel maps the schema attributes.
type RedisValkeyModuleModel struct {
	InstanceId types.String `tfsdk:"instance_id"`
	Name       types.String `tfsdk:"name"`
	Status     types.String `tfsdk:"status"`
}

func (a *RedisValkeyModule) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_redis_valkey_module"
}

func (a *RedisValkeyModule) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Modifies (enables) an extension module of a Valkey (Redis) instance via the ModifyValkeyModule API. " +
			"This is a one-time operation; it does not create, update or destroy any Terraform-managed resource, " +
			"and it exposes no attributes referenceable from other resources. " +
			"The cloud-side change (module enabled) persists after the action completes; " +
			"destroying this action does not disable the module. " +
			"Repeating the invocation repeats the cloud-side side effect. " +
			"Triggering this action through lifecycle.action_trigger runs it automatically during apply.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the instance to modify, e.g. `crs-xjhsdj****`. Please log in to the Redis console and copy the instance ID from the instance list.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The extension module to modify. Valid values: `ValkeyBloom`, `ValkeyJSON`, `ValkeySearch`.",
			},
			"status": schema.StringAttribute{
				Required:    true,
				Description: "Sets the enabled status of the extension module. Only `ON` (enable the module) is supported.",
			},
		},
	}
}

// Invoke is called to run the logic of the action.
func (a *RedisValkeyModule) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RedisValkeyModuleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.InstanceId.IsNull() || data.InstanceId.IsUnknown() || data.InstanceId.ValueString() == "" {
		resp.Diagnostics.AddError("Missing instance_id", "instance_id is required and cannot be empty.")
		return
	}
	if data.Name.IsNull() || data.Name.IsUnknown() || data.Name.ValueString() == "" {
		resp.Diagnostics.AddError("Missing name", "name is required and cannot be empty.")
		return
	}
	if data.Status.IsNull() || data.Status.IsUnknown() || data.Status.ValueString() == "" {
		resp.Diagnostics.AddError("Missing status", "status is required and cannot be empty.")
		return
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	progress := func(msg string) {
		if resp.SendProgress != nil {
			resp.SendProgress(action.InvokeProgressEvent{Message: msg})
		}
	}

	instanceId := data.InstanceId.ValueString()
	moduleName := data.Name.ValueString()
	client := a.Client().UseRedisClient()

	request := redis.NewModifyValkeyModuleRequest()
	request.InstanceId = data.InstanceId.ValueStringPointer()
	request.Name = data.Name.ValueStringPointer()
	request.Status = data.Status.ValueStringPointer()

	progress(fmt.Sprintf("invoking ModifyValkeyModule on %s for module %s", instanceId, moduleName))

	response, err := client.ModifyValkeyModuleWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("modifying Valkey module (%s) of instance (%s)", moduleName, instanceId),
			err.Error(),
		)
		return
	}
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("modifying Valkey module (%s) of instance (%s)", moduleName, instanceId),
			"empty response returned by the ModifyValkeyModule API",
		)
		return
	}
	if response.Response.TaskId == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("modifying Valkey module (%s) of instance (%s)", moduleName, instanceId),
			"no task id returned by the ModifyValkeyModule API",
		)
		return
	}

	// ModifyValkeyModule is asynchronous: poll the shared RedisService
	// DescribeTaskInfo helper with resource.Retry, following the same pattern
	// as other redis operations in service_tencentcloud_redis.go.
	taskId := *response.Response.TaskId
	progress(fmt.Sprintf("ModifyValkeyModule submitted, task id %d; waiting for completion", taskId))

	service := NewRedisService(a.Client())
	err = helperResource.Retry(tccommon.WriteRetryTimeout, func() *helperResource.RetryError {
		ok, e := service.DescribeTaskInfo(ctx, instanceId, taskId)
		if e != nil {
			if _, isSdkErr := e.(*sdkErrors.TencentCloudSDKError); isSdkErr {
				return helperResource.NonRetryableError(e)
			}
			return helperResource.RetryableError(e)
		}
		if ok {
			return nil
		}
		progress(fmt.Sprintf("task %d is processing", taskId))
		return helperResource.RetryableError(fmt.Errorf("ModifyValkeyModule task %d is processing", taskId))
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("waiting for Valkey module (%s) task (%d) of instance (%s)", moduleName, taskId, instanceId),
			err.Error(),
		)
		return
	}

	progress(fmt.Sprintf("ModifyValkeyModule on %s for module %s completed", instanceId, moduleName))
	// No state to write: actions have no terraform state and no output attributes.
}
