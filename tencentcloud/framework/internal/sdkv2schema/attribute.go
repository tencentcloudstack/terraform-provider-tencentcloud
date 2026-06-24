// Package sdkv2schema contains compatibility helpers for mirroring SDKv2
// provider schema semantics in terraform-plugin-framework schemas.
package sdkv2schema

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

var _ schema.Attribute = StringAttributeWithEnvDefault{}

// StringAttributeWithEnvDefault mirrors an SDKv2 required string field using
// schema.EnvDefaultFunc(envName, nil).
//
// SDKv2 does not serialize DefaultFunc itself into the protocol schema, but
// during schema conversion it treats a Required field with a non-nil default
// value as Optional. schema.EnvDefaultFunc only returns a non-nil default when
// the environment variable is non-empty, so this attribute makes the same
// runtime Required/Optional decision.
//
// The embedded schema.StringAttribute carries the standard framework string
// attribute metadata. Its Required/Optional fields are intentionally ignored;
// this wrapper owns those two protocol flags.
type StringAttributeWithEnvDefault struct {
	EnvName string
	schema.StringAttribute
}

// IsRequired returns true when the SDKv2 EnvDefaultFunc would not provide a
// default value, matching SDKv2's protocol schema conversion behavior.
func (a StringAttributeWithEnvDefault) IsRequired() bool {
	return !a.IsOptional()
}

// IsOptional returns true when the SDKv2 EnvDefaultFunc would provide a default
// value, matching SDKv2's protocol schema conversion behavior.
func (a StringAttributeWithEnvDefault) IsOptional() bool {
	return os.Getenv(a.EnvName) != ""
}
