package ags

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

func stringInput(fields map[string]attr.Value, key string) *string {
	v := fields[key]
	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil
	}
	s, ok := v.(interface{ ValueString() string })
	if !ok {
		return nil
	}
	return pointer(s.ValueString())
}
func boolInput(fields map[string]attr.Value, key string) *bool {
	v := fields[key]
	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil
	}
	b, ok := v.(types.Bool)
	if !ok {
		return nil
	}
	return pointer(b.ValueBool())
}
func intInput(fields map[string]attr.Value, key string) *int64 {
	v := fields[key]
	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil
	}
	n, ok := v.(types.Int64)
	if !ok {
		return nil
	}
	return pointer(n.ValueInt64())
}

func blockObjects(v attr.Value) ([]map[string]attr.Value, error) {
	out := []map[string]attr.Value{}
	if v == nil || v.IsNull() {
		return out, nil
	}
	if v.IsUnknown() {
		return nil, errors.New("block configuration is not known yet")
	}
	collection, ok := v.(interface{ Elements() []attr.Value })
	if !ok {
		return nil, errors.New("expected a block collection")
	}
	for _, item := range collection.Elements() {
		o, ok := item.(types.Object)
		if !ok || o.IsNull() || o.IsUnknown() {
			return nil, errors.New("block value is not a known object")
		}
		out = append(out, o.Attributes())
	}
	return out, nil
}
func oneBlock(fields map[string]attr.Value, key string, required bool) (map[string]attr.Value, error) {
	items, err := blockObjects(fields[key])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}
	if len(items) > 1 || required && len(items) != 1 {
		return nil, fmt.Errorf("%s must contain exactly %d block(s)", key, 1)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return items[0], nil
}
func stringElements(v attr.Value) ([]*string, error) {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	c, ok := v.(interface{ Elements() []attr.Value })
	if !ok {
		return nil, errors.New("expected a string collection")
	}
	out := make([]*string, 0, len(c.Elements()))
	for _, item := range c.Elements() {
		if item.IsUnknown() || item.IsNull() {
			return nil, errors.New("collection contains an unknown or null string")
		}
		s, ok := item.(interface{ ValueString() string })
		if !ok {
			return nil, errors.New("collection element is not a string")
		}
		out = append(out, pointer(s.ValueString()))
	}
	if _, ok := v.(types.Set); ok {
		sort.Slice(out, func(i, j int) bool { return *out[i] < *out[j] })
	}
	return out, nil
}
func stringMap(v attr.Value) (map[string]string, error) {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	m, ok := v.(types.Map)
	if !ok {
		return nil, errors.New("expected a string map")
	}
	out := map[string]string{}
	for key, item := range m.Elements() {
		s, ok := item.(types.String)
		if !ok || s.IsNull() || s.IsUnknown() {
			return nil, errors.New("map contains an unknown or null string")
		}
		out[key] = s.ValueString()
	}
	return out, nil
}

func expandNetwork(fields map[string]attr.Value) (*sdk.NetworkConfiguration, error) {
	n, err := oneBlock(fields, "network_configuration", true)
	if err != nil {
		return nil, err
	}
	mode := stringInput(n, "network_mode")
	if mode == nil {
		return nil, errors.New("network_mode is not known")
	}
	out := &sdk.NetworkConfiguration{NetworkMode: mode}
	vpc, err := oneBlock(n, "vpc_config", *mode == "VPC")
	if err != nil {
		return nil, err
	}
	if vpc != nil {
		if *mode != "VPC" {
			return nil, errors.New("vpc_config is only supported for VPC network mode")
		}
		subnets, err := stringElements(vpc["subnet_ids"])
		if err != nil {
			return nil, err
		}
		groups, err := stringElements(vpc["security_group_ids"])
		if err != nil {
			return nil, err
		}
		if len(subnets) == 0 || len(groups) == 0 {
			return nil, errors.New("VPC subnets and security groups must be nonempty")
		}
		out.VpcConfig = &sdk.VPCConfig{SubnetIds: subnets, SecurityGroupIds: groups}
	}
	return out, nil
}

func expandCustom(fields map[string]attr.Value) (*sdk.CustomConfiguration, error) {
	c, err := oneBlock(fields, "custom_configuration", value(stringInput(fields, "tool_type")) == "custom")
	if err != nil || c == nil {
		return nil, err
	}
	resources, err := oneBlock(c, "resources", true)
	if err != nil {
		return nil, err
	}
	out := &sdk.CustomConfiguration{Image: stringInput(c, "image"), ImageRegistryType: stringInput(c, "image_registry_type"), Resources: &sdk.ResourceConfiguration{CPU: stringInput(resources, "cpu"), Memory: stringInput(resources, "memory")}}
	if out.Image == nil || out.ImageRegistryType == nil || out.Resources.CPU == nil || out.Resources.Memory == nil {
		return nil, errors.New("custom image, registry type and CPU/memory must be known")
	}
	out.Command, err = stringElements(c["command"])
	if err != nil {
		return nil, err
	}
	out.Args, err = stringElements(c["args"])
	if err != nil {
		return nil, err
	}
	env, err := stringMap(c["env"])
	if err != nil {
		return nil, err
	}
	if env != nil {
		out.Env = make([]*sdk.EnvVar, 0, len(env))
		keys := make([]string, 0, len(env))
		for key := range env {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			out.Env = append(out.Env, &sdk.EnvVar{Name: pointer(key), Value: pointer(env[key])})
		}
	}
	ports, err := blockObjects(c["ports"])
	if err != nil {
		return nil, err
	}
	if len(ports) > 0 {
		out.Ports = make([]*sdk.PortConfiguration, 0, len(ports))
		for _, p := range ports {
			out.Ports = append(out.Ports, &sdk.PortConfiguration{Name: stringInput(p, "name"), Port: intInput(p, "port"), Protocol: stringInput(p, "protocol")})
		}
	}
	probe, err := oneBlock(c, "probe", false)
	if err != nil {
		return nil, err
	}
	if probe != nil {
		h, err := oneBlock(probe, "http_get", true)
		if err != nil {
			return nil, err
		}
		out.Probe = &sdk.ProbeConfiguration{HttpGet: &sdk.HttpGetAction{Path: stringInput(h, "path"), Port: intInput(h, "port"), Scheme: stringInput(h, "scheme")}, ReadyTimeoutMs: intInput(probe, "ready_timeout_ms"), ProbeTimeoutMs: intInput(probe, "probe_timeout_ms"), ProbePeriodMs: intInput(probe, "probe_period_ms"), SuccessThreshold: intInput(probe, "success_threshold"), FailureThreshold: intInput(probe, "failure_threshold")}
	}
	return out, nil
}

func expandMounts(fields map[string]attr.Value) ([]*sdk.StorageMount, error) {
	items, err := blockObjects(fields["storage_mounts"])
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	out := make([]*sdk.StorageMount, 0, len(items))
	for _, item := range items {
		cos, err := oneBlock(item, "cos", false)
		if err != nil {
			return nil, err
		}
		cfs, err := oneBlock(item, "cfs", false)
		if err != nil {
			return nil, err
		}
		if (cos == nil) == (cfs == nil) {
			return nil, errors.New("each mount requires exactly one COS or CFS source")
		}
		source := &sdk.StorageSource{}
		if cos != nil {
			source.Cos = &sdk.CosStorageSource{Endpoint: stringInput(cos, "endpoint"), BucketName: stringInput(cos, "bucket_name"), BucketPath: stringInput(cos, "bucket_path")}
		}
		if cfs != nil {
			source.Cfs = &sdk.CfsStorageSource{FileSystemId: stringInput(cfs, "file_system_id"), Path: stringInput(cfs, "path")}
		}
		readOnly := boolInput(item, "read_only")
		if readOnly == nil {
			readOnly = pointer(false)
		}
		out = append(out, &sdk.StorageMount{Name: stringInput(item, "name"), MountPath: stringInput(item, "mount_path"), ReadOnly: readOnly, StorageSource: source})
	}
	sort.Slice(out, func(i, j int) bool { return value(out[i].Name) < value(out[j].Name) })
	return out, nil
}
func expandLogs(fields map[string]attr.Value) (*sdk.LogConfiguration, error) {
	logs, err := oneBlock(fields, "log_configuration", false)
	if err != nil || logs == nil {
		return nil, err
	}
	files, err := stringElements(logs["files"])
	if err != nil {
		return nil, err
	}
	return &sdk.LogConfiguration{CLSConfig: &sdk.CLSConfig{TopicId: stringInput(logs, "topic_id")}, LogSources: &sdk.LogSources{Files: files}}, nil
}

func createRequest(plan types.Object) (*sdk.CreateSandboxToolRequest, error) {
	f := plan.Attributes()
	req := sdk.NewCreateSandboxToolRequest()
	req.ToolName, req.ToolType = stringInput(f, "tool_name"), stringInput(f, "tool_type")
	if req.ToolName == nil || req.ToolType == nil {
		return nil, errors.New("tool_name and tool_type must be known")
	}
	req.Description, req.DefaultTimeout, req.RoleArn, req.Persistent = stringInput(f, "description"), stringInput(f, "default_timeout"), stringInput(f, "role_arn"), boolInput(f, "persistent")
	var err error
	req.NetworkConfiguration, err = expandNetwork(f)
	if err != nil {
		return nil, err
	}
	req.CustomConfiguration, err = expandCustom(f)
	if err != nil {
		return nil, err
	}
	req.StorageMounts, err = expandMounts(f)
	if err != nil {
		return nil, err
	}
	req.LogConfiguration, err = expandLogs(f)
	if err != nil {
		return nil, err
	}
	tags, err := stringMap(f["tags"])
	if err != nil {
		return nil, err
	}
	if tags != nil {
		req.Tags = expandTags(tags)
	}
	return req, nil
}

// updateRequest sends only changed input fields. It does not flatten the remote
// Tool into an Update request, and does not invoke a mutation for a local-only change.
func updateRequest(prior, plan types.Object) (*sdk.UpdateSandboxToolRequest, bool, error) {
	before, err := createRequest(prior)
	if err != nil {
		return nil, false, err
	}
	after, err := createRequest(plan)
	if err != nil {
		return nil, false, err
	}
	req := sdk.NewUpdateSandboxToolRequest()
	req.ToolId = stringInput(prior.Attributes(), "id")
	changed := false
	for _, key := range []string{"tool_name", "tool_type", "default_timeout", "persistent", "role_arn", "network_configuration", "storage_mounts", "log_configuration"} {
		a, b := prior.Attributes()[key], plan.Attributes()[key]
		if !plannedEqual(a, b) && !b.IsUnknown() {
			return nil, false, fmt.Errorf("%s changed without a replacement plan", key)
		}
	}
	if after.Description != nil && !reflect.DeepEqual(before.Description, after.Description) {
		req.Description = after.Description
		changed = true
	}
	if after.Tags != nil && !reflect.DeepEqual(before.Tags, after.Tags) {
		req.Tags = after.Tags
		changed = true
	}
	if after.CustomConfiguration != nil {
		a, b := before.CustomConfiguration, after.CustomConfiguration
		if a == nil {
			return nil, false, errors.New("custom block was introduced without replacement")
		}
		patch := &sdk.CustomConfiguration{}
		any := false
		if !reflect.DeepEqual(a.Image, b.Image) {
			patch.Image = b.Image
			any = true
		}
		if !reflect.DeepEqual(a.ImageRegistryType, b.ImageRegistryType) {
			patch.ImageRegistryType = b.ImageRegistryType
			any = true
		}
		if b.Command != nil && !reflect.DeepEqual(a.Command, b.Command) {
			patch.Command = b.Command
			any = true
		}
		if b.Args != nil && !reflect.DeepEqual(a.Args, b.Args) {
			patch.Args = b.Args
			any = true
		}
		if b.Env != nil && !reflect.DeepEqual(a.Env, b.Env) {
			patch.Env = b.Env
			any = true
		}
		if !reflect.DeepEqual(a.Ports, b.Ports) {
			patch.Ports = b.Ports
			if patch.Ports == nil {
				patch.Ports = []*sdk.PortConfiguration{}
			}
			any = true
		}
		if !quantitiesEqual(a.Resources, b.Resources) {
			patch.Resources = b.Resources
			any = true
		}
		if !reflect.DeepEqual(a.Probe, b.Probe) {
			if b.Probe == nil {
				return nil, false, errors.New("removing the whole probe requires replacement")
			}
			patch.Probe = b.Probe
			any = true
		}
		if any {
			req.CustomConfiguration = patch
			changed = true
		}
	}
	return req, changed, nil
}
func quantitiesEqual(a, b *sdk.ResourceConfiguration) bool {
	if a == nil || b == nil {
		return a == b
	}
	return semanticallyEqual("quantity", value(a.CPU), value(b.CPU)) && semanticallyEqual("quantity", value(a.Memory), value(b.Memory))
}

func plainStrings(v []*string) []string {
	out := make([]string, 0, len(v))
	for _, s := range v {
		out = append(out, value(s))
	}
	return out
}
func optionalScalar[T any](v *T) any {
	if v == nil {
		return nil
	}
	return *v
}

func effectiveProjection(tool *sdk.SandboxTool) (map[string]any, error) {
	if tool == nil {
		return nil, errors.New("missing Tool")
	}
	if tool.Persistent == nil {
		return nil, errors.New("Tool read omitted persistent mode")
	}
	if !*tool.Persistent && tool.DefaultTimeoutSeconds == nil {
		return nil, errors.New("nonpersistent Tool read omitted default timeout")
	}
	allTags := map[string]string{}
	for _, tag := range tool.Tags {
		if tag == nil || tag.Key == nil || tag.Value == nil {
			return nil, errors.New("incomplete tag")
		}
		if _, exists := allTags[*tag.Key]; exists {
			return nil, errors.New("duplicate tag")
		}
		allTags[*tag.Key] = *tag.Value
	}
	out := map[string]any{"tool_name": value(tool.ToolName), "tool_type": value(tool.ToolType), "description": value(tool.Description), "default_timeout": nil, "persistent": *tool.Persistent, "role_arn": optionalScalar(tool.RoleArn), "tags": allTags, "network_configuration": nil, "custom_configuration": nil, "storage_mounts": []map[string]any{}, "log_configuration": nil}
	if tool.DefaultTimeoutSeconds != nil {
		out["default_timeout"] = fmt.Sprintf("%ds", *tool.DefaultTimeoutSeconds)
	}
	if tool.NetworkConfiguration == nil {
		return nil, errors.New("Tool read omitted network configuration")
	}
	network := map[string]any{"network_mode": value(tool.NetworkConfiguration.NetworkMode), "vpc_config": nil}
	if v := tool.NetworkConfiguration.VpcConfig; v != nil {
		network["vpc_config"] = map[string]any{"subnet_ids": plainStrings(v.SubnetIds), "security_group_ids": plainStrings(v.SecurityGroupIds)}
	}
	out["network_configuration"] = network
	if c := tool.CustomConfiguration; c != nil {
		env := map[string]string{}
		for _, entry := range c.Env {
			if entry == nil || entry.Name == nil || entry.Value == nil {
				return nil, errors.New("incomplete environment variable")
			}
			if _, exists := env[*entry.Name]; exists {
				return nil, errors.New("duplicate environment variable name")
			}
			env[*entry.Name] = *entry.Value
		}
		ports := []map[string]any{}
		for _, port := range c.Ports {
			if port == nil || port.Port == nil {
				return nil, errors.New("incomplete port")
			}
			ports = append(ports, map[string]any{"name": value(port.Name), "port": *port.Port, "protocol": value(port.Protocol)})
		}
		custom := map[string]any{"image": value(c.Image), "image_registry_type": value(c.ImageRegistryType), "command": plainStrings(c.Command), "args": plainStrings(c.Args), "env": env, "ports": ports, "resources": nil, "probe": nil}
		if r := c.Resources; r != nil {
			custom["resources"] = map[string]any{"cpu": value(r.CPU), "memory": value(r.Memory)}
		}
		if p := c.Probe; p != nil {
			probe := map[string]any{"http_get": nil, "ready_timeout_ms": optionalScalar(p.ReadyTimeoutMs), "probe_timeout_ms": optionalScalar(p.ProbeTimeoutMs), "probe_period_ms": optionalScalar(p.ProbePeriodMs), "success_threshold": optionalScalar(p.SuccessThreshold), "failure_threshold": optionalScalar(p.FailureThreshold)}
			if h := p.HttpGet; h != nil {
				probe["http_get"] = map[string]any{"path": value(h.Path), "port": optionalScalar(h.Port), "scheme": value(h.Scheme)}
			}
			custom["probe"] = probe
		}
		out["custom_configuration"] = custom
	}
	mounts := []map[string]any{}
	for _, mount := range tool.StorageMounts {
		if mount == nil || mount.Name == nil || mount.StorageSource == nil {
			return nil, errors.New("incomplete mount or a mount outside the V1 COS/CFS contract")
		}
		m := map[string]any{"name": *mount.Name, "mount_path": value(mount.MountPath), "read_only": value(mount.ReadOnly), "cos": nil, "cfs": nil}
		if c := mount.StorageSource.Cos; c != nil {
			m["cos"] = map[string]any{"endpoint": value(c.Endpoint), "bucket_name": value(c.BucketName), "bucket_path": value(c.BucketPath)}
		}
		if c := mount.StorageSource.Cfs; c != nil {
			m["cfs"] = map[string]any{"file_system_id": value(c.FileSystemId), "path": value(c.Path)}
		}
		if m["cos"] == nil && m["cfs"] == nil {
			return nil, errors.New("mount uses a source outside the V1 COS/CFS contract")
		}
		mounts = append(mounts, m)
	}
	out["storage_mounts"] = mounts
	if l := tool.LogConfiguration; l != nil {
		if l.CLSConfig == nil || l.LogSources == nil {
			return nil, errors.New("incomplete log configuration")
		}
		out["log_configuration"] = map[string]any{"topic_id": value(l.CLSConfig.TopicId), "files": plainStrings(l.LogSources.Files)}
	}
	return out, nil
}

func reflectedValue(ctx context.Context, t attr.Type, raw any) (attr.Value, error) {
	if raw == nil {
		return t.ValueFromTerraform(ctx, tftypes.NewValue(t.TerraformType(ctx), nil))
	}
	if existing, ok := raw.(attr.Value); ok {
		if existing.Type(ctx).Equal(t) {
			return existing, nil
		}
		encoded, err := existing.ToTerraformValue(ctx)
		if err != nil {
			return nil, err
		}
		return t.ValueFromTerraform(ctx, encoded)
	}
	v := reflect.ValueOf(raw)
	switch target := t.(type) {
	case types.ObjectType:
		if v.Kind() != reflect.Map || v.Type().Key().Kind() != reflect.String {
			return nil, errors.New("expected an object with string field names")
		}
		fields := map[string]attr.Value{}
		for key, kind := range target.AttrTypes {
			var input any
			if item := v.MapIndex(reflect.ValueOf(key)); item.IsValid() {
				input = item.Interface()
			}
			converted, err := reflectedValue(ctx, kind, input)
			if err != nil {
				return nil, fmt.Errorf("field %s: %w", key, err)
			}
			fields[key] = converted
		}
		out, ds := types.ObjectValue(target.AttrTypes, fields)
		if ds.HasError() {
			return nil, fmt.Errorf("object conversion: %v", ds.Errors())
		}
		return out, nil
	case types.ListType:
		items, err := reflectedElements(ctx, target.ElemType, v)
		if err != nil {
			return nil, err
		}
		out, ds := types.ListValue(target.ElemType, items)
		if ds.HasError() {
			return nil, fmt.Errorf("list conversion: %v", ds.Errors())
		}
		return out, nil
	case types.SetType:
		items, err := reflectedElements(ctx, target.ElemType, v)
		if err != nil {
			return nil, err
		}
		out, ds := types.SetValue(target.ElemType, items)
		if ds.HasError() {
			return nil, fmt.Errorf("set conversion: %v", ds.Errors())
		}
		return out, nil
	case types.MapType:
		if v.Kind() != reflect.Map || v.Type().Key().Kind() != reflect.String {
			return nil, errors.New("expected a string-keyed map")
		}
		fields := map[string]attr.Value{}
		iter := v.MapRange()
		for iter.Next() {
			converted, err := reflectedValue(ctx, target.ElemType, iter.Value().Interface())
			if err != nil {
				return nil, err
			}
			fields[iter.Key().String()] = converted
		}
		out, ds := types.MapValue(target.ElemType, fields)
		if ds.HasError() {
			return nil, fmt.Errorf("map conversion: %v", ds.Errors())
		}
		return out, nil
	case semanticStringType:
		text, ok := raw.(string)
		if !ok {
			return nil, errors.New("expected semantic string")
		}
		return semanticStringValue{StringValue: types.StringValue(text), kind: target.kind}, nil
	case basetypes.StringType:
		text, ok := raw.(string)
		if !ok {
			return nil, errors.New("expected string")
		}
		return types.StringValue(text), nil
	case basetypes.BoolType:
		b, ok := raw.(bool)
		if !ok {
			return nil, errors.New("expected boolean")
		}
		return types.BoolValue(b), nil
	case basetypes.Int64Type:
		if v.Kind() < reflect.Int || v.Kind() > reflect.Int64 {
			return nil, errors.New("expected signed integer")
		}
		return types.Int64Value(v.Int()), nil
	}
	return nil, fmt.Errorf("unsupported configuration type %T", t)
}

func reflectedElements(ctx context.Context, element attr.Type, v reflect.Value) ([]attr.Value, error) {
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return nil, errors.New("expected collection")
	}
	out := make([]attr.Value, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		value, err := reflectedValue(ctx, element, v.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}
func inputAttributeTypes() map[string]attr.Type {
	return toolSchema().Type().(types.ObjectType).AttrTypes
}
func asBlocks(v any) []map[string]any {
	if v == nil {
		return []map[string]any{}
	}
	return []map[string]any{v.(map[string]any)}
}
func cloneMap(in map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
func hasBlocks(v attr.Value) bool { items, err := blockObjects(v); return err == nil && len(items) > 0 }

func flattenTool(ctx context.Context, prior types.Object, tool *sdk.SandboxTool, importing bool) (types.Object, error) {
	projection, err := effectiveProjection(tool)
	if err != nil {
		return types.Object{}, err
	}
	if value(tool.ToolType) != "custom" && value(tool.ToolType) != "code-interpreter" && value(tool.ToolType) != "browser" {
		return types.Object{}, errors.New("Tool type is outside the V1 support matrix")
	}
	typesByKey := inputAttributeTypes()
	fields := prior.Attributes()
	if fields == nil {
		fields = map[string]attr.Value{}
	}
	set := func(key string, raw any) error {
		v, err := reflectedValue(ctx, typesByKey[key], raw)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
		if old := fields[key]; old != nil {
			v = preserveEquivalent(old, v)
		}
		fields[key] = v
		return nil
	}
	for _, key := range []string{"tool_name", "tool_type", "description", "default_timeout", "persistent", "role_arn"} {
		if err := set(key, projection[key]); err != nil {
			return types.Object{}, err
		}
	}
	tags, err := UserTags(tool.Tags)
	if err != nil {
		return types.Object{}, err
	}
	if err := set("tags", tags); err != nil {
		return types.Object{}, err
	}
	network := cloneMap(projection["network_configuration"].(map[string]any))
	network["vpc_config"] = asBlocks(network["vpc_config"])
	if err := set("network_configuration", asBlocks(network)); err != nil {
		return types.Object{}, err
	}
	customBlocks := []map[string]any{}
	if value(tool.ToolType) == "custom" && (importing || hasBlocks(fields["custom_configuration"])) {
		if projection["custom_configuration"] == nil {
			return types.Object{}, errors.New("custom Tool read omitted custom configuration")
		}
		custom := cloneMap(projection["custom_configuration"].(map[string]any))
		if custom["resources"] == nil {
			return types.Object{}, errors.New("custom Tool read omitted CPU/memory")
		}
		custom["resources"] = asBlocks(custom["resources"])
		oldCustom, _ := oneBlock(fields, "custom_configuration", false)
		if !importing && !hasBlocks(oldCustom["ports"]) {
			custom["ports"] = []map[string]any{}
		}
		if !importing && !hasBlocks(oldCustom["probe"]) {
			custom["probe"] = []map[string]any{}
		} else {
			p := custom["probe"]
			if p != nil {
				probe := cloneMap(p.(map[string]any))
				probe["http_get"] = asBlocks(probe["http_get"])
				custom["probe"] = asBlocks(probe)
			} else {
				custom["probe"] = []map[string]any{}
			}
		}
		customBlocks = append(customBlocks, custom)
	}
	if err := set("custom_configuration", customBlocks); err != nil {
		return types.Object{}, err
	}
	mountBlocks := []map[string]any{}
	if importing || hasBlocks(fields["storage_mounts"]) {
		for _, raw := range projection["storage_mounts"].([]map[string]any) {
			mount := cloneMap(raw)
			mount["cos"], mount["cfs"] = asBlocks(mount["cos"]), asBlocks(mount["cfs"])
			mountBlocks = append(mountBlocks, mount)
		}
	}
	if err := set("storage_mounts", mountBlocks); err != nil {
		return types.Object{}, err
	}
	logBlocks := []map[string]any{}
	if importing || hasBlocks(fields["log_configuration"]) {
		logBlocks = asBlocks(projection["log_configuration"])
	}
	if err := set("log_configuration", logBlocks); err != nil {
		return types.Object{}, err
	}
	for key, raw := range map[string]any{"id": value(tool.ToolId), "status": value(tool.Status), "status_reason": value(tool.StatusReason), "image_digest": nil, "effective_configuration": projection} {
		if key == "image_digest" && tool.CustomConfiguration != nil {
			raw = optionalScalar(tool.CustomConfiguration.ImageDigest)
		}
		if err := set(key, raw); err != nil {
			return types.Object{}, err
		}
	}
	if fields["timeouts"] == nil {
		if err := set("timeouts", []map[string]any{}); err != nil {
			return types.Object{}, err
		}
	}
	out, diags := types.ObjectValue(typesByKey, fields)
	if diags.HasError() {
		return types.Object{}, fmt.Errorf("Tool state conversion failed: %s", diags.Errors())
	}
	return out, nil
}

// preserveEquivalent retains typed duration/quantity spelling and the RPC
// representation of absent nested blocks (null versus zero elements). An
// actual nonempty block or changed configured value is never normalized away.
func preserveEquivalent(before, after attr.Value) attr.Value {
	if emptyBlockValue(before) && emptyBlockValue(after) {
		return before
	}
	if a, ok := before.(semanticStringValue); ok {
		if b, ok := after.(semanticStringValue); ok && !a.IsNull() && !a.IsUnknown() && !b.IsNull() && !b.IsUnknown() && a.kind == b.kind && semanticallyEqual(a.kind, a.ValueString(), b.ValueString()) {
			return a
		}
	}
	if a, ok := before.(types.Object); ok && !a.IsUnknown() && !a.IsNull() {
		if b, ok := after.(types.Object); ok && !b.IsNull() && !b.IsUnknown() {
			fields := b.Attributes()
			for k, v := range fields {
				if old := a.Attributes()[k]; old != nil {
					fields[k] = preserveEquivalent(old, v)
				}
			}
			out, d := types.ObjectValue(b.AttributeTypes(context.Background()), fields)
			if !d.HasError() {
				return out
			}
		}
	}
	if a, ok := before.(types.List); ok && !a.IsUnknown() && !a.IsNull() {
		if b, ok := after.(types.List); ok && !b.IsUnknown() && !b.IsNull() && len(a.Elements()) == len(b.Elements()) {
			items := b.Elements()
			for i, v := range items {
				items[i] = preserveEquivalent(a.Elements()[i], v)
			}
			out, d := types.ListValue(b.ElementType(context.Background()), items)
			if !d.HasError() {
				return out
			}
		}
	}
	return after
}

func emptyBlockValue(v attr.Value) bool {
	if v == nil || v.IsUnknown() {
		return false
	}
	collection, ok := v.(interface {
		ElementType(context.Context) attr.Type
		Elements() []attr.Value
	})
	if !ok {
		return false
	}
	_, objectElements := collection.ElementType(context.Background()).(types.ObjectType)
	return objectElements && (v.IsNull() || len(collection.Elements()) == 0)
}

func plannedEqual(expected, actual attr.Value) bool {
	if expected == nil || expected.IsUnknown() {
		return true
	}
	if actual == nil {
		return false
	}
	if emptyBlockValue(expected) && emptyBlockValue(actual) {
		return true
	}
	if expected.IsNull() {
		return actual.IsNull()
	}
	if a, ok := expected.(semanticStringValue); ok {
		if b, ok := actual.(semanticStringValue); ok && !b.IsNull() && !b.IsUnknown() {
			return a.kind == b.kind && semanticallyEqual(a.kind, a.ValueString(), b.ValueString())
		}
	}
	if a, ok := expected.(interface{ Attributes() map[string]attr.Value }); ok {
		b, ok := actual.(interface{ Attributes() map[string]attr.Value })
		if !ok || actual.IsNull() || actual.IsUnknown() || len(a.Attributes()) != len(b.Attributes()) {
			return false
		}
		for k, v := range a.Attributes() {
			if !plannedEqual(v, b.Attributes()[k]) {
				return false
			}
		}
		return true
	}
	if a, ok := expected.(interface{ Elements() []attr.Value }); ok {
		b, ok := actual.(interface{ Elements() []attr.Value })
		if !ok || actual.IsNull() || actual.IsUnknown() || len(a.Elements()) != len(b.Elements()) {
			return false
		}
		if _, set := expected.(types.Set); set {
			matched := make([]bool, len(b.Elements()))
			for _, item := range a.Elements() {
				found := false
				for i, candidate := range b.Elements() {
					if !matched[i] && plannedEqual(item, candidate) {
						matched[i], found = true, true
						break
					}
				}
				if !found {
					return false
				}
			}
			return true
		}
		for i, v := range a.Elements() {
			if !plannedEqual(v, b.Elements()[i]) {
				return false
			}
		}
		return true
	}
	if a, ok := expected.(interface{ Elements() map[string]attr.Value }); ok {
		b, ok := actual.(interface{ Elements() map[string]attr.Value })
		if !ok || actual.IsNull() || actual.IsUnknown() || len(a.Elements()) != len(b.Elements()) {
			return false
		}
		for key, item := range a.Elements() {
			if !plannedEqual(item, b.Elements()[key]) {
				return false
			}
		}
		return true
	}
	if expected.Equal(actual) {
		return true
	}
	// Framework may retain different collection/custom-type metadata at the
	// RPC boundary. Equal protocol values still represent exactly the same
	// configuration; unequal values are never suppressed here.
	a, aErr := expected.ToTerraformValue(context.Background())
	b, bErr := actual.ToTerraformValue(context.Background())
	return aErr == nil && bErr == nil && a.Equal(b)
}

func matchesPlan(ctx context.Context, plan types.Object, tool *sdk.SandboxTool) (bool, error) {
	return matchesOwnedPlan(ctx, plan, plan, tool)
}

func matchesOwnedPlan(ctx context.Context, plan, ownership types.Object, tool *sdk.SandboxTool) (bool, error) {
	fields, err := pendingPlanFields(ctx, plan, ownership, tool)
	return len(fields) == 0, err
}

func pendingPlanFields(ctx context.Context, plan, ownership types.Object, tool *sdk.SandboxTool) ([]string, error) {
	actual, err := flattenTool(ctx, ownership, tool, false)
	if err != nil {
		return nil, err
	}
	pending := []string{}
	for _, key := range []string{"tool_name", "tool_type", "description", "default_timeout", "persistent", "role_arn", "tags", "network_configuration", "custom_configuration", "storage_mounts", "log_configuration"} {
		if !plannedEqual(plan.Attributes()[key], actual.Attributes()[key]) {
			pending = append(pending, fmt.Sprintf("%s(plan=%T,null=%v;actual=%T,null=%v)", key, plan.Attributes()[key], plan.Attributes()[key].IsNull(), actual.Attributes()[key], actual.Attributes()[key].IsNull()))
		}
	}
	return pending, nil
}

// While clearing ports, keep the previous block's ownership until the read
// proves it was cleared. Otherwise an unchanged server value could be hidden
// simply because the new configuration no longer contains the block.
func ownershipForUpdate(ctx context.Context, prior, plan types.Object) (types.Object, error) {
	oldCustom, err := oneBlock(prior.Attributes(), "custom_configuration", false)
	if err != nil {
		return types.Object{}, err
	}
	newCustom, err := oneBlock(plan.Attributes(), "custom_configuration", false)
	if err != nil {
		return types.Object{}, err
	}
	if newCustom == nil || !hasBlocks(oldCustom["ports"]) || hasBlocks(newCustom["ports"]) {
		return plan, nil
	}
	newCustom["ports"] = oldCustom["ports"]
	root := plan.Attributes()
	block := root["custom_configuration"].(types.List)
	item, ds := types.ObjectValue(block.ElementType(ctx).(types.ObjectType).AttrTypes, newCustom)
	if ds.HasError() {
		return types.Object{}, fmt.Errorf("port ownership object: %v", ds.Errors())
	}
	list, ds := types.ListValue(block.ElementType(ctx), []attr.Value{item})
	if ds.HasError() {
		return types.Object{}, fmt.Errorf("port ownership block: %v", ds.Errors())
	}
	root["custom_configuration"] = list
	out, ds := types.ObjectValue(inputAttributeTypes(), root)
	if ds.HasError() {
		return types.Object{}, fmt.Errorf("port ownership state: %v", ds.Errors())
	}
	return out, nil
}

func operationContext(ctx context.Context, config types.Object, operation string) (context.Context, context.CancelFunc, error) {
	duration := 30 * time.Minute
	if operation == "read" {
		duration = 5 * time.Minute
	}
	if settings, err := oneBlock(config.Attributes(), "timeouts", false); err != nil {
		return nil, nil, err
	} else if text := stringInput(settings, operation); text != nil {
		parsed, err := time.ParseDuration(*text)
		if err != nil || parsed <= 0 {
			return nil, nil, fmt.Errorf("invalid %s operation timeout", operation)
		}
		duration = parsed
	}
	ctx, cancel := context.WithTimeout(ctx, duration)
	return ctx, cancel, nil
}

func stringValue(ctx context.Context, v attr.Value) (types.String, error) {
	s, ok := v.(basetypes.StringValuable)
	if !ok {
		return types.String{}, errors.New("expected string value")
	}
	out, d := s.ToStringValue(ctx)
	if d.HasError() {
		return types.String{}, errors.New("invalid string value")
	}
	return out, nil
}
