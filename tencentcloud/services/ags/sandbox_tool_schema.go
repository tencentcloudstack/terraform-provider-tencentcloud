package ags

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"k8s.io/apimachinery/pkg/api/resource"
)

type semanticStringType struct {
	basetypes.StringType
	kind string
}
type semanticStringValue struct {
	basetypes.StringValue
	kind string
}

func (t semanticStringType) String() string { return "ags." + t.kind }
func (t semanticStringType) Equal(other attr.Type) bool {
	o, ok := other.(semanticStringType)
	return ok && t.kind == o.kind
}
func (t semanticStringType) ValueType(context.Context) attr.Value {
	return semanticStringValue{kind: t.kind}
}
func (t semanticStringType) ValueFromString(_ context.Context, v basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return semanticStringValue{StringValue: v, kind: t.kind}, nil
}
func (t semanticStringType) ValueFromTerraform(ctx context.Context, v tftypes.Value) (attr.Value, error) {
	raw, err := t.StringType.ValueFromTerraform(ctx, v)
	if err != nil {
		return nil, err
	}
	return semanticStringValue{StringValue: raw.(basetypes.StringValue), kind: t.kind}, nil
}
func (v semanticStringValue) Type(context.Context) attr.Type { return semanticStringType{kind: v.kind} }
func (v semanticStringValue) Equal(other attr.Value) bool {
	o, ok := other.(semanticStringValue)
	return ok && o.kind == v.kind && v.StringValue.Equal(o.StringValue)
}
func (v semanticStringValue) StringSemanticEquals(_ context.Context, other basetypes.StringValuable) (bool, diag.Diagnostics) {
	o, ok := other.(semanticStringValue)
	if !ok || o.kind != v.kind {
		return false, nil
	}
	return semanticallyEqual(v.kind, v.ValueString(), o.ValueString()), nil
}
func semanticallyEqual(kind, a, b string) bool {
	if a == b {
		return true
	}
	if kind == "duration" {
		x, ex := time.ParseDuration(a)
		y, ey := time.ParseDuration(b)
		return ex == nil && ey == nil && x == y
	}
	x, ex := resource.ParseQuantity(a)
	y, ey := resource.ParseQuantity(b)
	return ex == nil && ey == nil && x.Cmp(y) == 0
}

type replaceString struct {
	kind   string
	stable bool
}

func (m replaceString) Description(context.Context) string {
	return "Changing this configuration replaces the Tool; semantically equivalent values preserve state."
}
func (m replaceString) MarkdownDescription(ctx context.Context) string { return m.Description(ctx) }
func (m replaceString) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if m.stable && req.PlanValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}
	if req.PlanValue.IsUnknown() {
		return
	}
	if !req.StateValue.IsNull() && !req.PlanValue.IsNull() && m.kind != "" && semanticallyEqual(m.kind, req.StateValue.ValueString(), req.PlanValue.ValueString()) {
		// The custom string type also supplies semantic equality for refreshed
		// values. CLI plan tests must verify this normalization against Core.
		resp.PlanValue = req.StateValue
		return
	}
	resp.RequiresReplace = !req.StateValue.Equal(req.PlanValue)
}

type replaceBool struct{}

func (replaceBool) Description(context.Context) string {
	return "Changing this field replaces the Tool."
}
func (m replaceBool) MarkdownDescription(ctx context.Context) string { return m.Description(ctx) }
func (replaceBool) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if !req.State.Raw.IsNull() && !req.Plan.Raw.IsNull() && !req.PlanValue.IsUnknown() {
		resp.RequiresReplace = !req.StateValue.Equal(req.PlanValue)
	}
}

type replaceList struct{}

func (replaceList) Description(context.Context) string {
	return "Changing this block replaces the Tool."
}
func (m replaceList) MarkdownDescription(ctx context.Context) string { return m.Description(ctx) }
func (replaceList) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if !req.State.Raw.IsNull() && !req.Plan.Raw.IsNull() && !req.PlanValue.IsUnknown() {
		resp.RequiresReplace = !req.StateValue.Equal(req.PlanValue)
	}
}

type replaceSet struct{}

func (replaceSet) Description(context.Context) string {
	return "Changing these bindings replaces the Tool."
}
func (m replaceSet) MarkdownDescription(ctx context.Context) string { return m.Description(ctx) }
func (replaceSet) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.State.Raw.IsNull() && !req.Plan.Raw.IsNull() && !req.PlanValue.IsUnknown() {
		resp.RequiresReplace = !req.StateValue.Equal(req.PlanValue)
	}
}

func requiredString(description string) schema.StringAttribute {
	return schema.StringAttribute{Required: true, Description: description}
}
func optionalString(description string) schema.StringAttribute {
	return schema.StringAttribute{Optional: true, Computed: true, Description: description}
}
func singleton(description string, attributes map[string]schema.Attribute, blocks map[string]schema.Block) schema.ListNestedBlock {
	return schema.ListNestedBlock{Description: description, NestedObject: schema.NestedBlockObject{Attributes: attributes, Blocks: blocks}}
}

func toolSchema() schema.Schema {
	name := requiredString("Unique Tool name, 1-50 letters, digits, underscores or hyphens. Changes replace the Tool.")
	name.PlanModifiers = []planmodifier.String{replaceString{}}
	kind := requiredString("V1 Tool type: browser, code-interpreter or custom. Changes replace the Tool.")
	kind.PlanModifiers = []planmodifier.String{replaceString{}}
	duration := optionalString("Default instance lifetime, expressed as a positive whole-second duration up to 24h. Omission accepts the remote default. Changes replace the Tool; equivalent durations do not.")
	duration.CustomType = semanticStringType{kind: "duration"}
	duration.PlanModifiers = []planmodifier.String{replaceString{kind: "duration"}}
	role := optionalString("Existing CAM data-access role ARN, distinct from Provider credentials. Omission accepts the current role. Changes replace the Tool.")
	role.PlanModifiers = []planmodifier.String{replaceString{}}
	vpc := singleton("Required exactly once for VPC mode; forbidden otherwise.", map[string]schema.Attribute{
		"subnet_ids":         schema.SetAttribute{Required: true, ElementType: types.StringType, Description: "Nonempty set of existing subnets in one VPC. Reordering is not a change."},
		"security_group_ids": schema.SetAttribute{Required: true, ElementType: types.StringType, Description: "Nonempty set of existing security groups compatible with the selected VPC."},
	}, nil)
	network := singleton("Exactly one network block. V1 conservatively replaces the Tool for network changes, including VPC mode and bindings.", map[string]schema.Attribute{"network_mode": requiredString("PUBLIC, SANDBOX or VPC.")}, map[string]schema.Block{"vpc_config": vpc})
	network.PlanModifiers = []planmodifier.List{replaceList{}}
	quantity := func(description string) schema.StringAttribute {
		a := requiredString(description)
		a.CustomType = semanticStringType{kind: "quantity"}
		return a
	}
	resources := singleton("Exactly one CPU/memory configuration for custom Tools.", map[string]schema.Attribute{"cpu": quantity("Positive supported CPU quantity, for example 1 or 1000m; values must satisfy the regional support matrix."), "memory": quantity("Positive supported memory quantity, for example 1Gi; values must satisfy the regional support matrix.")}, nil)
	httpGet := singleton("Exactly one HTTP GET action when a probe is configured.", map[string]schema.Attribute{"path": requiredString("Absolute HTTP request path."), "port": schema.Int64Attribute{Required: true, Description: "Probe port, 1-65535."}, "scheme": requiredString("HTTP or HTTPS.")}, nil)
	probeAttrs := map[string]schema.Attribute{}
	for _, key := range []string{"ready_timeout_ms", "probe_timeout_ms", "probe_period_ms", "success_threshold", "failure_threshold"} {
		unit := "positive count"
		if len(key) > 3 && key[len(key)-3:] == "_ms" {
			unit = "positive milliseconds"
		}
		probeAttrs[key] = schema.Int64Attribute{Optional: true, Computed: true, Description: fmt.Sprintf("Probe %s (%s). Omission accepts the service default.", key, unit)}
	}
	probe := singleton("Optional readiness probe. Value updates are in place; removing the whole block replaces the Tool until a public reset contract exists.", probeAttrs, map[string]schema.Block{"http_get": httpGet})
	ports := schema.SetNestedBlock{Description: "Declared business ports; removing entries updates in place. Order is insignificant. Does not itself close listening processes or establish a firewall.", NestedObject: schema.NestedBlockObject{Attributes: map[string]schema.Attribute{"name": requiredString("Unique port name."), "port": schema.Int64Attribute{Required: true, Description: "Business port, 1-65535, excluding platform-reserved ports."}, "protocol": requiredString("Supported port protocol.")}}}
	custom := singleton("Exactly one block for custom Tools; forbidden for built-in types. Uses prebuilt, AGS-compatible enterprise/personal registry images. Changes affect newly created instances according to the service contract.", map[string]schema.Attribute{
		"image":               requiredString("Prebuilt image reference. In-place update; immutable references are recommended."),
		"image_registry_type": requiredString("enterprise for TCR or personal for CCR. Third-party credential configuration is not supported."),
		"command":             schema.ListAttribute{Optional: true, Computed: true, ElementType: types.StringType, Description: "Ordered command override. Omission preserves the remote setting; [] explicitly removes the override when supported."},
		"args":                schema.ListAttribute{Optional: true, Computed: true, ElementType: types.StringType, Description: "Ordered command arguments. Omission preserves the remote setting; [] explicitly clears the override."},
		"env":                 schema.MapAttribute{Optional: true, Computed: true, ElementType: types.StringType, Description: "Non-sensitive environment variables. Omission preserves remote values; {} clears user overrides. Values are stored in Terraform state."},
	}, map[string]schema.Block{"resources": resources, "probe": probe, "ports": ports})
	cos := singleton("Choose exactly one COS or CFS source per mount.", map[string]schema.Attribute{"endpoint": requiredString("Existing COS endpoint without embedded credentials."), "bucket_name": requiredString("Existing bucket name."), "bucket_path": requiredString("Absolute path within the bucket.")}, nil)
	cfs := singleton("Choose exactly one COS or CFS source per mount.", map[string]schema.Attribute{"file_system_id": requiredString("Existing CFS filesystem ID."), "path": requiredString("CFS source path.")}, nil)
	mounts := schema.SetNestedBlock{Description: "Existing storage mounts. Changes replace the Tool; the Tool never creates or deletes the storage itself.", PlanModifiers: []planmodifier.Set{replaceSet{}}, NestedObject: schema.NestedBlockObject{Attributes: map[string]schema.Attribute{"name": requiredString("Unique mount name."), "mount_path": requiredString("Absolute path inside the instance."), "read_only": schema.BoolAttribute{Optional: true, Computed: true, Description: "Read-only mount; defaults to false when omitted in a declared mount."}}, Blocks: map[string]schema.Block{"cos": cos, "cfs": cfs}}}
	logs := singleton("Optional CLS log delivery; changes replace the Tool. Existing Topic and role permissions must be provided.", map[string]schema.Attribute{"topic_id": requiredString("Existing CLS Topic ID."), "files": schema.SetAttribute{Required: true, ElementType: types.StringType, Description: "One to ten files directly under /logs/; subdirectories are unsupported."}}, nil)
	logs.PlanModifiers = []planmodifier.List{replaceList{}}
	timeoutAttrs := map[string]schema.Attribute{}
	for _, operation := range []string{"create", "update", "delete", "read"} {
		d := "30m"
		if operation == "read" {
			d = "5m"
		}
		timeoutAttrs[operation] = schema.StringAttribute{Optional: true, Description: operation + " operation waiting budget; defaults to " + d + ". Distinct from instance lifetime."}
	}
	return schema.Schema{Version: 0, Description: "Manages an AGS SandboxTool configuration. Instances remain owned by the runtime application. Active or paused instances can block deletion. Fixed-field changes replace the Tool; same-name replacement may require downtime. V1 API contract gates must pass before production release.", Attributes: map[string]schema.Attribute{
		"id":        schema.StringAttribute{Computed: true, Description: "Tool ID consumed by the runtime SDK and used for import.", PlanModifiers: []planmodifier.String{replaceString{stable: true}}},
		"tool_name": name, "tool_type": kind, "description": optionalString("Tool description, up to 200 characters. Empty string requests clearing; omitted values accept the current remote value. The API must support explicit clearing."), "default_timeout": duration,
		"persistent":              schema.BoolAttribute{Optional: true, Computed: true, Description: "Default persistent mode for future instances. Does not start an instance. Type support must be verified; changes replace the Tool.", PlanModifiers: []planmodifier.Bool{replaceBool{}}},
		"role_arn":                role,
		"tags":                    schema.MapAttribute{Optional: true, Computed: true, ElementType: types.StringType, Description: "User-managed tags. {} removes all user tags. The documented immutable qcs:tag:createdBy system tag is read-only and remains visible in effective_configuration.tags."},
		"status":                  schema.StringAttribute{Computed: true, Description: "Actual CREATING, ACTIVE, DELETING or FAILED status."},
		"status_reason":           schema.StringAttribute{Computed: true, Description: "Service-reported status reason."},
		"image_digest":            schema.StringAttribute{Computed: true, Description: "Resolved image digest; recalculated when image changes."},
		"effective_configuration": schema.ObjectAttribute{Computed: true, AttributeTypes: configurationTypes(), Description: "Supported actual configuration, including remote defaults and system tags. Does not grant configuration ownership to omitted blocks."},
	}, Blocks: map[string]schema.Block{"network_configuration": network, "custom_configuration": custom, "storage_mounts": mounts, "log_configuration": logs, "timeouts": singleton("Optional operation waiting budgets; at most one block.", timeoutAttrs, nil)}}
}

func objectType(fields map[string]attr.Type) attr.Type { return types.ObjectType{AttrTypes: fields} }
func configurationTypes() map[string]attr.Type {
	stringSet := types.SetType{ElemType: types.StringType}
	stringList := types.ListType{ElemType: types.StringType}
	port := objectType(map[string]attr.Type{"name": types.StringType, "port": types.Int64Type, "protocol": types.StringType})
	probe := map[string]attr.Type{"http_get": objectType(map[string]attr.Type{"path": types.StringType, "port": types.Int64Type, "scheme": types.StringType})}
	for _, key := range []string{"ready_timeout_ms", "probe_timeout_ms", "probe_period_ms", "success_threshold", "failure_threshold"} {
		probe[key] = types.Int64Type
	}
	custom := objectType(map[string]attr.Type{"image": types.StringType, "image_registry_type": types.StringType, "command": stringList, "args": stringList, "env": types.MapType{ElemType: types.StringType}, "ports": types.SetType{ElemType: port}, "resources": objectType(map[string]attr.Type{"cpu": types.StringType, "memory": types.StringType}), "probe": objectType(probe)})
	mount := objectType(map[string]attr.Type{"name": types.StringType, "mount_path": types.StringType, "read_only": types.BoolType, "cos": objectType(map[string]attr.Type{"endpoint": types.StringType, "bucket_name": types.StringType, "bucket_path": types.StringType}), "cfs": objectType(map[string]attr.Type{"file_system_id": types.StringType, "path": types.StringType})})
	return map[string]attr.Type{"tool_name": types.StringType, "tool_type": types.StringType, "description": types.StringType, "default_timeout": types.StringType, "persistent": types.BoolType, "role_arn": types.StringType, "tags": types.MapType{ElemType: types.StringType}, "network_configuration": objectType(map[string]attr.Type{"network_mode": types.StringType, "vpc_config": objectType(map[string]attr.Type{"subnet_ids": stringSet, "security_group_ids": stringSet})}), "custom_configuration": custom, "storage_mounts": types.SetType{ElemType: mount}, "log_configuration": objectType(map[string]attr.Type{"topic_id": types.StringType, "files": stringSet})}
}
