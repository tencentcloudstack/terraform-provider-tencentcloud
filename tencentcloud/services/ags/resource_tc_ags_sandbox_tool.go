package ags

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

type SandboxToolResource struct {
	fw.ResourceWithConfigure
	testService *Service
}

var (
	_ resource.ResourceWithConfigure      = &SandboxToolResource{}
	_ resource.ResourceWithImportState    = &SandboxToolResource{}
	_ resource.ResourceWithValidateConfig = &SandboxToolResource{}
	_ resource.ResourceWithModifyPlan     = &SandboxToolResource{}
)

func NewAgsSandboxToolResource() resource.Resource { return &SandboxToolResource{} }
func (r *SandboxToolResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_tool"
}
func (r *SandboxToolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = toolSchema()
}
func (r *SandboxToolResource) service() (*Service, error) {
	if r.testService != nil {
		return r.testService, nil
	}
	client, err := r.Client().UseAgsClient()
	if err != nil {
		return nil, err
	}
	return NewService(client), nil
}

func (r *SandboxToolResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config types.Object
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || config.IsNull() || config.IsUnknown() {
		return
	}
	resp.Diagnostics.Append(validateToolConfig(config)...)
}

func (r *SandboxToolResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}
	var prior, plan, config types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Normalize only the schema's explicitly typed duration/quantity values.
	// In particular custom CPU 2 -> 2000m and memory 2Gi -> 2048Mi must not
	// produce a local update when the API configuration is unchanged.
	plan = preserveEquivalent(prior, plan).(types.Object)
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	// Framework invalidates computed fields before attribute modifiers run.
	// Once typed normalization proves every managed input unchanged, keep the
	// readback metadata too; otherwise an equivalent 5m -> 300s edit alone
	// creates a spurious update. Any configured unknown prevents this path.
	unchanged := !containsUnknown(config)
	inputs := []string{"tool_name", "tool_type", "description", "default_timeout", "persistent", "role_arn", "tags", "network_configuration", "custom_configuration", "storage_mounts", "log_configuration"}
	comparison := plan
	if unchanged {
		// Framework can mark an omitted Optional/Computed nested input unknown
		// when a caller relinquishes ownership. Compare that value with the
		// refreshed state. Commit this comparison only if *all* managed inputs
		// are unchanged: another update may change the remote default, in which
		// case its future value must remain unknown until the write is read back.
		fields := plan.Attributes()
		for _, key := range inputs {
			fields[key] = unownedComparisonValue(ctx, prior.Attributes()[key], fields[key], config.Attributes()[key])
		}
		var ds diag.Diagnostics
		comparison, ds = types.ObjectValue(inputAttributeTypes(), fields)
		resp.Diagnostics.Append(ds...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	for _, key := range inputs {
		unchanged = unchanged && plannedEqual(comparison.Attributes()[key], prior.Attributes()[key])
	}
	if unchanged {
		plan = comparison
		fields := plan.Attributes()
		for key, old := range prior.Attributes() {
			if fields[key].IsUnknown() {
				fields[key] = old
			}
		}
		out, ds := types.ObjectValue(inputAttributeTypes(), fields)
		resp.Diagnostics.Append(ds...)
		if !resp.Diagnostics.HasError() {
			resp.Diagnostics.Append(resp.Plan.Set(ctx, out)...)
		}
	}
	oldCustom, err := oneBlock(prior.Attributes(), "custom_configuration", false)
	if err != nil {
		return
	}
	newCustom, err := oneBlock(plan.Attributes(), "custom_configuration", false)
	if err != nil {
		return
	}
	if newCustom != nil && hasBlocks(oldCustom["probe"]) && !hasBlocks(newCustom["probe"]) {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("custom_configuration").AtListIndex(0).AtName("probe"))
	}
}

// unownedComparisonValue is for an otherwise unchanged plan, not an attribute
// UseStateForUnknown modifier. It never replaces configured values and does
// not pair unordered set elements. The caller discards the entire comparison
// whenever any managed input changes or any configured value is unknown.
func unownedComparisonValue(ctx context.Context, prior, plan, config attr.Value) attr.Value {
	if prior == nil || plan == nil || config == nil || prior.IsUnknown() || config.IsUnknown() {
		return plan
	}
	if config.IsNull() && plan.IsUnknown() {
		return prior
	}
	if prior.IsNull() || plan.IsNull() || plan.IsUnknown() || config.IsNull() {
		return plan
	}
	if planned, ok := plan.(types.Object); ok {
		before, beforeOK := prior.(types.Object)
		configured, configOK := config.(types.Object)
		if beforeOK && configOK {
			fields := planned.Attributes()
			for key, value := range fields {
				fields[key] = unownedComparisonValue(ctx, before.Attributes()[key], value, configured.Attributes()[key])
			}
			if result, ds := types.ObjectValue(planned.AttributeTypes(ctx), fields); !ds.HasError() {
				return result
			}
		}
	}
	if planned, ok := plan.(types.List); ok {
		before, beforeOK := prior.(types.List)
		configured, configOK := config.(types.List)
		if beforeOK && configOK && len(planned.Elements()) == len(before.Elements()) && len(planned.Elements()) == len(configured.Elements()) {
			items := planned.Elements()
			for index, value := range items {
				items[index] = unownedComparisonValue(ctx, before.Elements()[index], value, configured.Elements()[index])
			}
			if result, ds := types.ListValue(planned.ElementType(ctx), items); !ds.HasError() {
				return result
			}
		}
	}
	return plan
}

func containsUnknown(v attr.Value) bool {
	if v == nil || v.IsNull() {
		return false
	}
	if v.IsUnknown() {
		return true
	}
	if object, ok := v.(interface{ Attributes() map[string]attr.Value }); ok {
		for _, child := range object.Attributes() {
			if containsUnknown(child) {
				return true
			}
		}
	}
	if collection, ok := v.(interface{ Elements() []attr.Value }); ok {
		for _, child := range collection.Elements() {
			if containsUnknown(child) {
				return true
			}
		}
	}
	if collection, ok := v.(interface{ Elements() map[string]attr.Value }); ok {
		for _, child := range collection.Elements() {
			if containsUnknown(child) {
				return true
			}
		}
	}
	return false
}

func (r *SandboxToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.State.Schema = req.Plan.Schema
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validateToolConfig(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	payload, err := createRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AGS Tool configuration", err.Error())
		return
	}
	ctx, cancel, err := operationContext(ctx, plan, "create")
	if err != nil {
		resp.Diagnostics.AddError("Invalid operation timeout", err.Error())
		return
	}
	defer cancel()
	svc, err := r.service()
	if err != nil {
		resp.Diagnostics.AddError("AGS client unavailable", err.Error())
		return
	}
	token, err := NewClientToken()
	if err != nil {
		resp.Diagnostics.AddError("Cannot create AGS idempotency token", err.Error())
		return
	}
	payload.ClientToken = pointer(token)
	id, createErr := svc.Create(ctx, payload)
	if id != "" {
		provisional, err := provisionalState(ctx, plan, id)
		if err != nil {
			resp.Diagnostics.AddError("Cannot preserve accepted Tool identity", fmt.Sprintf("Tool %s was accepted; reconcile its ID before retrying: %v", id, err))
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, provisional)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if createErr != nil {
		resp.Diagnostics.AddError("AGS Tool creation was not completed", fmt.Sprintf("%v; known Tool ID=%q. If no ID was received, reconcile owned name %q in this account/region and import; do not blindly create again.", createErr, id, value(payload.ToolName)))
		return
	}
	var pending []string
	tool, err := svc.WaitReady(ctx, id, true, func(tool *sdk.SandboxTool) (bool, error) {
		var matchErr error
		pending, matchErr = pendingPlanFields(ctx, plan, plan, tool)
		return len(pending) == 0, matchErr
	})
	if tool != nil {
		state, flattenErr := flattenTool(ctx, plan, tool, false)
		if flattenErr != nil {
			resp.Diagnostics.AddError("Incomplete AGS Tool readback", flattenErr.Error())
		} else {
			resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		}
	}
	if err != nil {
		resp.Diagnostics.AddError("AGS Tool did not become ready", fmt.Sprintf("%v; pending fields=%v; accepted Tool ID=%s is retained for read/import recovery. Creation errors can taint state; verify remote readiness before untainting.", err, pending, id))
	}
}

func (r *SandboxToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.State = req.State
	var prior types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := value(stringInput(prior.Attributes(), "id"))
	if id == "" {
		resp.Diagnostics.AddError("Missing AGS Tool identity", "State has no Tool ID; refusing an unconstrained read.")
		return
	}
	ctx, cancel, err := operationContext(ctx, prior, "read")
	if err != nil {
		resp.Diagnostics.AddError("Invalid read timeout", err.Error())
		return
	}
	defer cancel()
	svc, err := r.service()
	if err != nil {
		resp.Diagnostics.AddError("AGS client unavailable", err.Error())
		return
	}
	tool, err := svc.ReadManaged(ctx, id)
	if errors.Is(err, ErrConfirmedAbsent) && value(stringInput(prior.Attributes(), "status")) == "CREATING" {
		// A failed create can return its accepted ID before list visibility.
		// Do not discard that provisional identity as an external deletion.
		tool, err = svc.Lookup(ctx, id)
	}
	if errors.Is(err, ErrConfirmedAbsent) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Cannot reliably read AGS Tool", err.Error())
		return
	}
	state, err := flattenTool(ctx, prior, tool, false)
	if err != nil {
		resp.Diagnostics.AddError("Incomplete AGS Tool configuration", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if value(tool.Status) == "FAILED" {
		resp.Diagnostics.AddWarning("AGS Tool is in FAILED state", fmt.Sprintf("Tool %s: %s. The state is retained and destroy remains available.", id, value(tool.StatusReason)))
	}
}

func (r *SandboxToolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var prior, plan types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validateToolConfig(plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	payload, changed, err := updateRequest(prior, plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid AGS update plan", err.Error())
		return
	}
	ctx, cancel, err := operationContext(ctx, plan, "update")
	if err != nil {
		resp.Diagnostics.AddError("Invalid operation timeout", err.Error())
		return
	}
	defer cancel()
	svc, err := r.service()
	if err != nil {
		resp.Diagnostics.AddError("AGS client unavailable", err.Error())
		return
	}
	id := value(stringInput(prior.Attributes(), "id"))
	shape, err := ownershipForUpdate(ctx, prior, plan)
	if err != nil {
		resp.Diagnostics.AddError("Cannot establish field ownership", err.Error())
		return
	}
	var operationErr error
	if changed {
		operationErr = svc.Update(ctx, payload)
	}
	var tool *sdk.SandboxTool
	if operationErr == nil {
		tool, operationErr = svc.WaitReady(ctx, id, false, func(tool *sdk.SandboxTool) (bool, error) { return matchesOwnedPlan(ctx, plan, shape, tool) })
	} else if ctx.Err() == nil {
		tool, _ = svc.Lookup(ctx, id)
	}
	if tool != nil {
		state, flattenErr := flattenTool(ctx, shape, tool, false)
		if flattenErr != nil {
			resp.Diagnostics.AddError("Incomplete AGS update readback", flattenErr.Error())
		} else {
			resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		}
	}
	if operationErr != nil {
		resp.Diagnostics.AddError("AGS Tool update did not converge", fmt.Sprintf("%v; Tool ID %s is retained. Successfully read changes are recorded; reconcile the remaining plan. No automatic rollback or instance action was attempted.", operationErr, id))
	}
}

func (r *SandboxToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State = req.State
	var prior types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := value(stringInput(prior.Attributes(), "id"))
	if id == "" {
		resp.Diagnostics.AddError("Missing AGS Tool identity", "Refusing deletion without an exact Tool ID.")
		return
	}
	ctx, cancel, err := operationContext(ctx, prior, "delete")
	if err != nil {
		resp.Diagnostics.AddError("Invalid delete timeout", err.Error())
		return
	}
	defer cancel()
	svc, err := r.service()
	if err != nil {
		resp.Diagnostics.AddError("AGS client unavailable", err.Error())
		return
	}
	if err := svc.DeleteAndWait(ctx, id); err != nil {
		resp.Diagnostics.AddError("Cannot confirm AGS Tool deletion", fmt.Sprintf("%v; ID %s is retained. Completion requires explicit NotFound or an acknowledged Delete followed by a complete empty point query. Active or paused instances must be released by their runtime owner; the Provider never stops instances.", err, id))
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *SandboxToolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !regexp.MustCompile(`^sdt-[a-z0-9]+$`).MatchString(req.ID) {
		resp.Diagnostics.AddError("Invalid AGS Tool import ID", "Use the actual sdt-... Tool ID with the correct account and region; names cannot be imported.")
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	svc, err := r.service()
	if err != nil {
		resp.Diagnostics.AddError("AGS client unavailable", err.Error())
		return
	}
	tool, err := svc.Lookup(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Cannot import AGS Tool", err.Error())
		return
	}
	if err := checkImportScope(tool); err != nil {
		resp.Diagnostics.AddError("Tool configuration is outside V1 import support", err.Error())
		return
	}
	state, err := flattenTool(ctx, types.ObjectNull(inputAttributeTypes()), tool, true)
	if err != nil {
		resp.Diagnostics.AddError("Cannot import incomplete configuration", err.Error())
		return
	}
	if diagnostics := validateToolConfig(state); diagnostics.HasError() {
		resp.Diagnostics.AddError("Tool configuration is outside V1 import support", "Remote settings cannot be faithfully represented by the V1 configuration constraints; no state ownership was established.")
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	resp.State.Schema = toolSchema()
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func checkImportScope(tool *sdk.SandboxTool) error {
	if tool.ComputerConfiguration != nil {
		return errors.New("computer configuration is not managed by this V1 schema")
	}
	if c := tool.CustomConfiguration; c != nil {
		if c.DNSConfig != nil || c.SecurityContext != nil {
			return errors.New("DNS/security context cannot be faithfully represented by this V1 schema; do not silently discard it")
		}
		if c.Resources != nil && value(c.Resources.Storage) != "" {
			return errors.New("custom disk storage cannot be faithfully represented by this V1 schema; do not silently discard it")
		}
	}
	for _, mount := range tool.StorageMounts {
		if mount == nil || mount.Volume != nil || mount.VolumeTemplate != nil || mount.StorageSource == nil || mount.StorageSource.Image != nil || mount.StorageSource.AgentBucket != nil {
			return errors.New("mount contains a source outside the V1 COS/CFS schema")
		}
	}
	return nil
}

func provisionalState(ctx context.Context, plan types.Object, id string) (types.Object, error) {
	root := knownOrNull(ctx, plan).(types.Object)
	fields := root.Attributes()
	fields["id"] = types.StringValue(id)
	fields["status"] = types.StringValue("CREATING")
	fields["status_reason"] = types.StringNull()
	fields["image_digest"] = types.StringNull()
	fields["effective_configuration"] = types.ObjectNull(configurationTypes())
	out, ds := types.ObjectValue(inputAttributeTypes(), fields)
	if ds.HasError() {
		return types.Object{}, fmt.Errorf("%v", ds.Errors())
	}
	return out, nil
}
func knownOrNull(ctx context.Context, v attr.Value) attr.Value {
	if v.IsUnknown() {
		null, err := v.Type(ctx).ValueFromTerraform(ctx, tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil))
		if err == nil {
			return null
		}
		return v
	}
	if v.IsNull() {
		return v
	}
	switch x := v.(type) {
	case types.Object:
		fields := x.Attributes()
		for k, item := range fields {
			fields[k] = knownOrNull(ctx, item)
		}
		out, ds := types.ObjectValue(x.AttributeTypes(ctx), fields)
		if !ds.HasError() {
			return out
		}
	case types.List:
		items := x.Elements()
		for i, item := range items {
			items[i] = knownOrNull(ctx, item)
		}
		out, ds := types.ListValue(x.ElementType(ctx), items)
		if !ds.HasError() {
			return out
		}
	case types.Set:
		items := x.Elements()
		for i, item := range items {
			items[i] = knownOrNull(ctx, item)
		}
		out, ds := types.SetValue(x.ElementType(ctx), items)
		if !ds.HasError() {
			return out
		}
	case types.Map:
		items := x.Elements()
		for k, item := range items {
			items[k] = knownOrNull(ctx, item)
		}
		out, ds := types.MapValue(x.ElementType(ctx), items)
		if !ds.HasError() {
			return out
		}
	}
	return v
}
