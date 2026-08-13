package ssm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"

	svcssm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"
)

// TestSsmSecretVersionEphemeralResource_Metadata validates that
// Metadata.TypeName is set to "tencentcloud_ssm_secret_version".
func TestSsmSecretVersionEphemeralResource_Metadata(t *testing.T) {
	e := &svcssm.SsmSecretVersionEphemeralResource{}
	req := ephemeral.MetadataRequest{}
	resp := &ephemeral.MetadataResponse{}

	e.Metadata(t.Context(), req, resp)

	if resp.TypeName != "tencentcloud_ssm_secret_version" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_ssm_secret_version")
	}
}

// TestSsmSecretVersionEphemeralResource_Schema validates that the schema
// includes the required attributes (secret_name, version_id) and computed
// attributes (secret_binary, secret_string) with the right flags.
func TestSsmSecretVersionEphemeralResource_Schema(t *testing.T) {
	e := &svcssm.SsmSecretVersionEphemeralResource{}
	req := ephemeral.SchemaRequest{}
	resp := &ephemeral.SchemaResponse{}

	e.Schema(t.Context(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	attrKeys := make(map[string]bool)
	for k := range resp.Schema.Attributes {
		attrKeys[k] = true
	}

	for _, required := range []string{"secret_name", "version_id"} {
		if !attrKeys[required] {
			t.Errorf("schema missing required attribute: %s", required)
			continue
		}
		schemaAttr, ok := resp.Schema.Attributes[required].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %s is not a StringAttribute", required)
			continue
		}
		if !schemaAttr.Required {
			t.Errorf("attribute %s should be Required", required)
		}
	}

	for _, computed := range []string{"secret_binary", "secret_string"} {
		if !attrKeys[computed] {
			t.Errorf("schema missing computed attribute: %s", computed)
			continue
		}
		schemaAttr, ok := resp.Schema.Attributes[computed].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %s is not a StringAttribute", computed)
			continue
		}
		if !schemaAttr.Computed {
			t.Errorf("attribute %s should be Computed", computed)
		}
		if !schemaAttr.Sensitive {
			t.Errorf("attribute %s should be Sensitive", computed)
		}
	}
}

// TestSsmSecretVersionEphemeralResource_ModelTypes validates the model struct
// fields are accessible (compile-time guard against accidental renames).
func TestSsmSecretVersionEphemeralResource_ModelTypes(t *testing.T) {
	model := svcssm.SsmSecretVersionEphemeralResourceModel{}
	_ = model.SecretName
	_ = model.VersionId
	_ = model.SecretBinary
	_ = model.SecretString
}

// TestSsmSecretVersionEphemeralResource_ImplementsInterface validates that
// SsmSecretVersionEphemeralResource implements ephemeral.EphemeralResource
// and ephemeral.EphemeralResourceWithConfigure.
func TestSsmSecretVersionEphemeralResource_ImplementsInterface(t *testing.T) {
	var _ ephemeral.EphemeralResource = &svcssm.SsmSecretVersionEphemeralResource{}
	var _ ephemeral.EphemeralResourceWithConfigure = &svcssm.SsmSecretVersionEphemeralResource{}
}

// TestSsmSecretVersionEphemeralResource_Factory validates that the registry
// factory returns a non-nil EphemeralResource.
func TestSsmSecretVersionEphemeralResource_Factory(t *testing.T) {
	got := svcssm.NewSsmSecretVersionEphemeralResource()
	if got == nil {
		t.Fatal("NewSsmSecretVersionEphemeralResource returned nil")
	}
}

// TestSsmSecretVersionEphemeralResource_ConfigureNilProviderData validates
// that Configure handles nil ProviderData gracefully (this is the normal
// pre-Configure invocation pattern used by the framework runtime).
func TestSsmSecretVersionEphemeralResource_ConfigureNilProviderData(t *testing.T) {
	e := &svcssm.SsmSecretVersionEphemeralResource{}
	req := ephemeral.ConfigureRequest{}
	resp := &ephemeral.ConfigureResponse{}

	e.Configure(t.Context(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestSsmSecretVersionEphemeralResource_ConfigureWrongType validates that
// Configure adds an error diagnostic when ProviderData is of an unexpected
// type (defensive guard against future refactors that change the meta type).
func TestSsmSecretVersionEphemeralResource_ConfigureWrongType(t *testing.T) {
	e := &svcssm.SsmSecretVersionEphemeralResource{}
	req := ephemeral.ConfigureRequest{
		ProviderData: "wrong_type",
	}
	resp := &ephemeral.ConfigureResponse{}

	e.Configure(t.Context(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Error("expected error diagnostic for wrong ProviderData type")
	}
}
