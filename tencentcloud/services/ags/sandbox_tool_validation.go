package ags

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"k8s.io/apimachinery/pkg/api/resource"
)

func validateToolConfig(config types.Object) diag.Diagnostics {
	var ds diag.Diagnostics
	f := config.Attributes()
	add := func(key, why string) { ds.AddError("Invalid AGS Tool configuration", key+": "+why) }
	name := stringInput(f, "tool_name")
	if name != nil && !regexp.MustCompile(`^[A-Za-z0-9_-]{1,50}$`).MatchString(*name) {
		add("tool_name", "must contain 1-50 letters, digits, underscores or hyphens")
	}
	kind := stringInput(f, "tool_type")
	if kind != nil && *kind != "custom" && *kind != "code-interpreter" && *kind != "browser" {
		add("tool_type", "V1 supports custom, code-interpreter and browser")
	}
	if text := stringInput(f, "description"); text != nil && utf8.RuneCountInString(*text) > 200 {
		add("description", "maximum length is 200 characters")
	}
	if text := stringInput(f, "default_timeout"); text != nil {
		d, e := time.ParseDuration(*text)
		if e != nil || d <= 0 || d > 24*time.Hour || d%time.Second != 0 {
			add("default_timeout", "use a positive whole-second duration no greater than 24h")
		}
	}
	if persistent := boolInput(f, "persistent"); persistent != nil && *persistent && kind != nil && *kind != "custom" {
		add("persistent", "among V1 types, only custom supports persistent=true")
	}
	if role := stringInput(f, "role_arn"); role != nil && (*role == "" || !strings.HasPrefix(*role, "qcs::cam::")) {
		add("role_arn", "must be a nonempty CAM role ARN; omit to accept the remote role")
	}
	if tags, err := stringMap(f["tags"]); err == nil {
		for key := range tags {
			if key == "" || strings.HasPrefix(key, "qcs:") {
				add("tags", "keys must be nonempty user tags, excluding reserved qcs: system tags")
			}
		}
	}
	for _, key := range []string{"network_configuration", "custom_configuration", "log_configuration", "timeouts"} {
		if v := f[key]; v != nil && !v.IsUnknown() {
			blocks, e := blockObjects(v)
			if e != nil {
				add(key, e.Error())
			} else if len(blocks) > 1 || key == "network_configuration" && len(blocks) != 1 {
				add(key, "unexpected number of singleton blocks")
			}
		}
	}
	custom, e := oneBlock(f, "custom_configuration", false)
	if e == nil && kind != nil {
		if *kind == "custom" && custom == nil {
			add("custom_configuration", "required for custom Tools")
		}
		if *kind != "custom" && custom != nil {
			add("custom_configuration", "forbidden for built-in Tool types")
		}
	}
	network, e := oneBlock(f, "network_configuration", false)
	if e == nil && network != nil {
		mode := stringInput(network, "network_mode")
		if mode != nil {
			if *mode != "PUBLIC" && *mode != "SANDBOX" && *mode != "VPC" {
				add("network_configuration.network_mode", "expected PUBLIC, SANDBOX or VPC")
			}
			vpc, e := oneBlock(network, "vpc_config", *mode == "VPC")
			if e != nil {
				if known(network["vpc_config"]) {
					add("network_configuration.vpc_config", e.Error())
				}
			} else if vpc != nil {
				if *mode != "VPC" {
					add("network_configuration.vpc_config", "forbidden outside VPC mode")
				}
				for _, key := range []string{"subnet_ids", "security_group_ids"} {
					if known(vpc[key]) {
						items, pending, err := validationStringElements(vpc[key])
						if err != nil || !pending && len(items) == 0 {
							add("network_configuration.vpc_config."+key, "must contain nonempty IDs once dependency references resolve")
						}
						for _, item := range items {
							if strings.TrimSpace(*item) == "" {
								add("network_configuration.vpc_config."+key, "IDs may not be empty")
							}
						}
					}
				}
			}
			if *mode == "SANDBOX" && hasBlocks(f["log_configuration"]) {
				add("log_configuration", "CLS delivery is not supported with SANDBOX networking")
			}
		}
	}
	if custom != nil {
		if registry := stringInput(custom, "image_registry_type"); registry != nil && *registry != "enterprise" && *registry != "personal" {
			add("custom_configuration.image_registry_type", "expected enterprise or personal")
		}
		if image := stringInput(custom, "image"); image != nil && (*image == "" || strings.ContainsAny(*image, " \t\n") || strings.Contains(*image, "://") || strings.Contains(*image, "@") && !strings.Contains(*image, "@sha256:")) {
			add("custom_configuration.image", "expected a credential-free image reference, optionally using @sha256:")
		}
		if resources, e := oneBlock(custom, "resources", true); e != nil {
			if known(custom["resources"]) {
				add("custom_configuration.resources", e.Error())
			}
		} else {
			for _, key := range []string{"cpu", "memory"} {
				if text := stringInput(resources, key); text != nil {
					q, e := resource.ParseQuantity(*text)
					if e != nil || q.Sign() <= 0 {
						add("custom_configuration.resources."+key, "expected a positive supported resource quantity")
					}
				}
			}
		}
		ports, e := blockObjects(custom["ports"])
		if e == nil {
			names := map[string]bool{}
			bindings := map[string]bool{}
			for _, port := range ports {
				if name := stringInput(port, "name"); name != nil {
					if *name == "" || names[*name] {
						add("custom_configuration.ports.name", "must be nonempty and unique")
					}
					names[*name] = true
				}
				if n := intInput(port, "port"); n != nil {
					if *n < 1 || *n > 65535 || *n == 49983 || *n == 49999 {
						add("custom_configuration.ports.port", "must be a valid business port, excluding platform ports 49983/49999")
					}
					key := fmt.Sprintf("%d/%s", *n, value(stringInput(port, "protocol")))
					if bindings[key] {
						add("custom_configuration.ports", "duplicate port/protocol binding")
					}
					bindings[key] = true
				}
				if p := stringInput(port, "protocol"); p != nil && *p != "TCP" && *p != "UDP" {
					add("custom_configuration.ports.protocol", "expected TCP or UDP")
				}
			}
		}
		if probe, e := oneBlock(custom, "probe", false); e != nil {
			if known(custom["probe"]) {
				add("custom_configuration.probe", e.Error())
			}
		} else if probe != nil {
			for _, key := range []string{"ready_timeout_ms", "probe_timeout_ms", "probe_period_ms", "success_threshold", "failure_threshold"} {
				if n := intInput(probe, key); n != nil && *n <= 0 {
					add("custom_configuration.probe."+key, "must be positive; time values use milliseconds")
				}
			}
			if h, e := oneBlock(probe, "http_get", true); e != nil {
				if known(probe["http_get"]) {
					add("custom_configuration.probe.http_get", e.Error())
				}
			} else {
				if p := stringInput(h, "path"); p != nil && !strings.HasPrefix(*p, "/") {
					add("custom_configuration.probe.http_get.path", "must be absolute")
				}
				if p := intInput(h, "port"); p != nil && (*p < 1 || *p > 65535) {
					add("custom_configuration.probe.http_get.port", "must be 1-65535")
				}
				if p := stringInput(h, "scheme"); p != nil && *p != "HTTP" && *p != "HTTPS" {
					add("custom_configuration.probe.http_get.scheme", "expected HTTP or HTTPS")
				}
			}
		}
	}
	if mounts, e := blockObjects(f["storage_mounts"]); e == nil {
		names := map[string]bool{}
		paths := map[string]bool{}
		for _, mount := range mounts {
			if name := stringInput(mount, "name"); name != nil {
				if *name == "" || names[*name] {
					add("storage_mounts.name", "must be nonempty and unique")
				}
				names[*name] = true
			}
			if target := stringInput(mount, "mount_path"); target != nil {
				if !strings.HasPrefix(*target, "/") || path.Clean(*target) != *target || paths[*target] {
					add("storage_mounts.mount_path", "must be a unique clean absolute path")
				}
				paths[*target] = true
			}
			cos, ce := oneBlock(mount, "cos", false)
			cfs, fe := oneBlock(mount, "cfs", false)
			if ce != nil || fe != nil {
				add("storage_mounts", "each source block is a singleton")
			} else if known(mount["cos"]) && known(mount["cfs"]) && (cos == nil) == (cfs == nil) {
				add("storage_mounts", "choose exactly one COS or CFS source")
			}
			if cos != nil {
				if p := stringInput(cos, "bucket_path"); p != nil && !strings.HasPrefix(*p, "/") {
					add("storage_mounts.cos.bucket_path", "must be absolute")
				}
				if endpoint := stringInput(cos, "endpoint"); endpoint != nil && strings.Contains(*endpoint, "@") {
					add("storage_mounts.cos.endpoint", "embedded credentials are forbidden")
				}
			}
		}
	}
	if logs, e := oneBlock(f, "log_configuration", false); e == nil && logs != nil && known(logs["files"]) {
		files, pending, e := validationStringElements(logs["files"])
		if e != nil || !pending && (len(files) == 0 || len(files) > 10) {
			add("log_configuration.files", "must contain 1-10 files")
		}
		for _, file := range files {
			if !strings.HasPrefix(*file, "/logs/") || strings.Contains(strings.TrimPrefix(*file, "/logs/"), "/") || strings.ContainsAny(*file, "*?[]") || path.Clean(*file) != *file {
				add("log_configuration.files", "files must be directly under /logs/ without wildcard or subdirectory")
			}
		}
	}
	if timeouts, e := oneBlock(f, "timeouts", false); e == nil {
		for _, key := range []string{"create", "read", "update", "delete"} {
			if text := stringInput(timeouts, key); text != nil {
				d, e := time.ParseDuration(*text)
				if e != nil || d <= 0 {
					add("timeouts."+key, "must be a positive duration")
				}
			}
		}
	}
	return ds
}
func known(v attr.Value) bool { return v != nil && !v.IsUnknown() }

// Collection.IsUnknown describes the collection itself, not its elements.
// A set of dependency references is known during plan even while each ID is
// unknown. Validate known elements now and defer the rest to apply; request
// expansion deliberately keeps stringElements' strict unresolved-value check.
func validationStringElements(v attr.Value) (values []*string, pending bool, err error) {
	if v == nil || v.IsNull() {
		return nil, false, nil
	}
	if v.IsUnknown() {
		return nil, true, nil
	}
	collection, ok := v.(interface{ Elements() []attr.Value })
	if !ok {
		return nil, false, errors.New("expected a string collection")
	}
	for _, item := range collection.Elements() {
		if item.IsUnknown() {
			pending = true
			continue
		}
		value, ok := item.(types.String)
		if !ok || value.IsNull() {
			return values, pending, errors.New("collection contains a null or non-string value")
		}
		values = append(values, pointer(value.ValueString()))
	}
	return values, pending, nil
}
