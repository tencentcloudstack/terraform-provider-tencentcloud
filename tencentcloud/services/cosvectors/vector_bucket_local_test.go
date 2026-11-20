package cosvectors

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// cosVectorLocalProvider 只在测试中注入回环 SDK，不增加生产 endpoint 开关。
type cosVectorLocalProvider struct{ service *CosVectorService }

func (p *cosVectorLocalProvider) Metadata(_ context.Context, _ fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "tencentcloud"
	resp.Version = "local-test"
}
func (p *cosVectorLocalProvider) Schema(_ context.Context, _ fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}
func (p *cosVectorLocalProvider) Configure(_ context.Context, _ fwprovider.ConfigureRequest, _ *fwprovider.ConfigureResponse) {
}
func (p *cosVectorLocalProvider) Resources(_ context.Context) []func() fwresource.Resource {
	return []func() fwresource.Resource{
		func() fwresource.Resource { return &CosVectorBucketResource{service: p.service} },
		func() fwresource.Resource { return &CosVectorIndexResource{service: p.service} },
		func() fwresource.Resource { return &CosVectorBucketPolicyResource{service: p.service} },
	}
}
func (p *cosVectorLocalProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func cosVectorLocalSetup(t *testing.T) (*CosVectorService, *cosVectorBucketMock, map[string]func() (tfprotov6.ProviderServer, error)) {
	t.Helper()
	if os.Getenv("TENCENTCLOUD_COS_VECTOR_LOCAL_TEST") != "1" {
		t.Skip("set TENCENTCLOUD_COS_VECTOR_LOCAL_TEST=1 and TF_ACC_TERRAFORM_PATH for local Terraform tests")
	}
	service, mock := newCosVectorBucketTestService(t)
	factories := map[string]func() (tfprotov6.ProviderServer, error){"tencentcloud": providerserver.NewProtocol6WithError(&cosVectorLocalProvider{service: service})}
	return service, mock, factories
}

func cosVectorLocalConfig(name, extra string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_vector_bucket" "test" {
  vector_bucket_name = %q
  %s
}
`, name, extra)
}

func TestCosVectorBucketLocalLifecycle(t *testing.T) {
	_, mock, factories := cosVectorLocalSetup(t)
	mock.invisibleReads = 1
	name := "tf-local-1250000000"
	config := cosVectorLocalConfig(name, "")
	resource.Test(t, resource.TestCase{IsUnitTest: true, ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{
		{Config: config, Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("tencentcloud_cos_vector_bucket.test", "id", name),
			resource.TestCheckResourceAttr("tencentcloud_cos_vector_bucket.test", "encryption_configuration.0.sse_type", "AES256"),
			resource.TestCheckResourceAttrSet("tencentcloud_cos_vector_bucket.test", "vector_bucket_qcs"),
			resource.TestCheckResourceAttrSet("tencentcloud_cos_vector_bucket.test", "creation_time"))},
		{Config: config, PlanOnly: true},
		{ResourceName: "tencentcloud_cos_vector_bucket.test", ImportState: true, ImportStateVerify: true, ImportStateId: name},
		{ResourceName: "tencentcloud_cos_vector_bucket.test", ImportState: true, ImportStateVerify: true, ImportStateId: "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/" + name},
		{Config: config, PreConfig: func() { mock.drop(name) }, PlanOnly: true, ExpectNonEmptyPlan: true},
		{Config: config},
		{Config: cosVectorLocalConfig(name, `encryption_configuration { sse_type = "AES256" }`), PlanOnly: true},
		{Config: config, PlanOnly: true},
		{Config: cosVectorLocalConfig("tf-renamed-1250000000", `encryption_configuration { sse_type = "AES256" }`)},
		{Config: cosVectorLocalConfig("tf-renamed-1250000000", `encryption_configuration { sse_type = "AES256" }`), PlanOnly: true},
		{Config: cosVectorLocalConfig("tf-renamed-1250000000", ""), PlanOnly: true},
		{Config: cosVectorLocalConfig("tf-create-first-1250000000", "lifecycle { create_before_destroy = true }"), Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("tencentcloud_cos_vector_bucket.test", "id", "tf-create-first-1250000000"),
			resource.TestCheckResourceAttr("tencentcloud_cos_vector_bucket.test", "vector_bucket_qcs", "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-create-first-1250000000"),
			resource.TestCheckResourceAttr("tencentcloud_cos_vector_bucket.test", "creation_time", "1700000004"))},
		{Config: cosVectorLocalConfig("tf-create-first-1250000000", "lifecycle { create_before_destroy = true }"), PlanOnly: true},
	}})
	if mock.count() != 0 {
		t.Fatal("destroy left mock resources")
	}
	if mock.callCount("CreateVectorBucket") != 4 {
		t.Fatal("unexpected create count; rename must replace bucket")
	}
}

func TestCosVectorBucketLocalImportUnencrypted(t *testing.T) {
	service, mock, factories := cosVectorLocalSetup(t)
	name := "tf-local-import-1250000000"
	if _, err := service.CreateBucket(context.Background(), name, nil); err != nil {
		t.Fatal(err)
	}
	config := cosVectorLocalConfig(name, "")
	resource.Test(t, resource.TestCase{IsUnitTest: true, ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{
		{Config: config, ResourceName: "tencentcloud_cos_vector_bucket.test", ImportState: true, ImportStateId: name, ImportStatePersist: true},
		{Config: config, PlanOnly: true},
		{Config: cosVectorLocalConfig(name, `encryption_configuration { sse_type = "AES256" }`), Check: resource.TestCheckResourceAttr("tencentcloud_cos_vector_bucket.test", "encryption_configuration.0.sse_type", "AES256")},
	}})
	if mock.callCount("CreateVectorBucket") != 2 || mock.count() != 0 {
		t.Fatal("encryption change did not replace bucket or destroy failed")
	}
}

func TestCosVectorBucketLocalForbiddenRefresh(t *testing.T) {
	_, mock, factories := cosVectorLocalSetup(t)
	config := cosVectorLocalConfig("tf-local-denied-1250000000", "")
	resource.Test(t, resource.TestCase{IsUnitTest: true, ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{
		{Config: config},
		{Config: config, PreConfig: func() { mock.setFailure("GetVectorBucket", 403) }, ExpectError: regexp.MustCompile("Error reading COS vector bucket")},
		{Config: config, PreConfig: func() { mock.setFailure("GetVectorBucket", 0) }},
	}})
	if mock.callCount("CreateVectorBucket") != 1 || mock.count() != 0 {
		t.Fatal("forbidden refresh caused recreation or cleanup failed")
	}
}

func TestCosVectorBucketLocalInvalidConfiguration(t *testing.T) {
	_, _, factories := cosVectorLocalSetup(t)
	for _, tc := range []struct{ name, config, error string }{
		{"name", cosVectorLocalConfig("invalid-name", ""), "Invalid COS vector bucket name"},
		{"kms", cosVectorLocalConfig("tf-local-invalid-1250000000", `encryption_configuration { sse_type = "aws:kms" }`), "Unsupported COS vector encryption"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{IsUnitTest: true, ProtoV6ProviderFactories: factories, Steps: []resource.TestStep{{Config: tc.config, PlanOnly: true, ExpectError: regexp.MustCompile(tc.error)}}})
		})
	}
}
