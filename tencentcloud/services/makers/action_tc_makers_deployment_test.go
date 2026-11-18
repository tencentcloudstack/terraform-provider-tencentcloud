package makers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	makerssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/makers"
)

func ptrStringMakersDeployment(s string) *string {
	return &s
}

func ptrInt64MakersDeployment(i int64) *int64 {
	return &i
}

func makersDeploymentSchema(t *testing.T) action_schema.Schema {
	t.Helper()
	a := makers.NewMakersDeployment()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

func makersDeploymentTimeoutsObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"invoke": tftypes.String,
		},
	}
}

func nullMakersDeploymentTimeoutsValue() tftypes.Value {
	return tftypes.NewValue(makersDeploymentTimeoutsObjType(), nil)
}

// makersDeploymentAttributeTypes lists every attribute of the action schema:
// all CreateDeployment request fields plus the `timeouts` block.
func makersDeploymentAttributeTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"zone_id":          tftypes.String,
		"project_id":       tftypes.String,
		"via_meta":         tftypes.String,
		"env":              tftypes.String,
		"provider":         tftypes.String,
		"dist_type":        tftypes.String,
		"temp_bucket_path": tftypes.String,
		"repo_branch":      tftypes.String,
		"repo_commit_hash": tftypes.String,
		"repo_commit_msg":   tftypes.String,
		"build_from":       tftypes.String,
		"template_id":      tftypes.String,
		"ci_spec":          tftypes.String,
		"timeouts":         makersDeploymentTimeoutsObjType(),
	}
}

// newMakersDeploymentInvokeRequest builds an InvokeRequest whose config object
// carries the given attribute values; attributes absent from attrs are null.
func newMakersDeploymentInvokeRequest(t *testing.T, attrs map[string]tftypes.Value) action.InvokeRequest {
	return newMakersDeploymentInvokeRequestWithTimeouts(t, attrs, nullMakersDeploymentTimeoutsValue())
}

func newMakersDeploymentInvokeRequestWithTimeouts(t *testing.T, attrs map[string]tftypes.Value, timeoutsVal tftypes.Value) action.InvokeRequest {
	t.Helper()
	s := makersDeploymentSchema(t)

	attributeTypes := makersDeploymentAttributeTypes()
	values := make(map[string]tftypes.Value, len(attributeTypes))
	for name, typ := range attributeTypes {
		if v, ok := attrs[name]; ok {
			values[name] = v
			continue
		}
		values[name] = tftypes.NewValue(typ, nil)
	}
	values["timeouts"] = timeoutsVal

	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: attributeTypes}, values)

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

func baseMakersDeploymentAttrs() map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"zone_id":    tftypes.NewValue(tftypes.String, "zone-1"),
		"project_id": tftypes.NewValue(tftypes.String, "project-1"),
		"via_meta":   tftypes.NewValue(tftypes.String, "Api"),
		"env":        tftypes.NewValue(tftypes.String, "Production"),
		"provider":   tftypes.NewValue(tftypes.String, "Github"),
	}
}

func setMakersDeploymentClient(a *makers.MakersDeployment, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// shortenMakersDeploymentPollInterval shortens the package-level poll interval
// so the Wait-loop tests do not sleep for seconds.
func shortenMakersDeploymentPollInterval(t *testing.T) func() {
	t.Helper()
	orig := makers.MakersDeploymentPollIntervalForTest()
	makers.SetMakersDeploymentPollIntervalForTest(1 * time.Millisecond)
	return func() {
		makers.SetMakersDeploymentPollIntervalForTest(orig)
	}
}

// patchMakersClient installs a UseMakersClient patch returning a stub client.
func patchMakersClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) {
	patches.ApplyMethodFunc(client, "UseMakersClient", func() *makerssdk.Client {
		return &makerssdk.Client{}
	})
}

// patchCreateDeployment installs a CreateDeploymentWithContext patch.
func patchCreateDeployment(patches *gomonkey.Patches, fn func(ctx context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error)) {
	patches.ApplyMethodFunc(&makerssdk.Client{}, "CreateDeploymentWithContext", fn)
}

// patchDescribeDeployments installs a DescribeDeploymentsWithContext patch.
func patchDescribeDeployments(patches *gomonkey.Patches, fn func(ctx context.Context, req *makerssdk.DescribeDeploymentsRequest) (*makerssdk.DescribeDeploymentsResponse, error)) {
	patches.ApplyMethodFunc(&makerssdk.Client{}, "DescribeDeploymentsWithContext", fn)
}

// go test ./tencentcloud/services/makers/ -run "TestMakersDeployment" -v -count=1 -gcflags="all=-l"

// TestMakersDeployment_Invoke_Success tests a successful build: CreateDeployment
// is accepted with an id, then polling observes the Process state followed by
// Success, and every request field is forwarded as expected.
func TestMakersDeployment_Invoke_Success(t *testing.T) {
	defer shortenMakersDeploymentPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchMakersClient(patches, client)

	var gotReq *makerssdk.CreateDeploymentRequest
	patchCreateDeployment(patches, func(_ context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error) {
		gotReq = req
		resp := makerssdk.NewCreateDeploymentResponse()
		resp.Response = &makerssdk.CreateDeploymentResponseParams{
			DeploymentId: ptrStringMakersDeployment("deploy-1"),
			RequestId:    ptrStringMakersDeployment("fake-request-id"),
		}
		return resp, nil
	})

	calls := 0
	patchDescribeDeployments(patches, func(_ context.Context, req *makerssdk.DescribeDeploymentsRequest) (*makerssdk.DescribeDeploymentsResponse, error) {
		calls++
		if req.ZoneId == nil || *req.ZoneId != "zone-1" {
			return nil, fmt.Errorf("expected ZoneId=zone-1 in status query, got %v", req.ZoneId)
		}
		if req.ProjectId == nil || *req.ProjectId != "project-1" {
			return nil, fmt.Errorf("expected ProjectId=project-1 in status query, got %v", req.ProjectId)
		}
		if len(req.DeploymentIds) != 1 || *req.DeploymentIds[0] != "deploy-1" {
			return nil, fmt.Errorf("expected DeploymentIds=[deploy-1] in status query, got %v", req.DeploymentIds)
		}
		status := makers.DeploymentStatusProcess
		if calls >= 2 {
			status = makers.DeploymentStatusSuccess
		}
		resp := makerssdk.NewDescribeDeploymentsResponse()
		resp.Response = &makerssdk.DescribeDeploymentsResponseParams{
			TotalCount:  ptrInt64MakersDeployment(1),
			Deployments: []*makerssdk.Deployment{{DeploymentId: ptrStringMakersDeployment("deploy-1"), Status: ptrStringMakersDeployment(status)}},
		}
		return resp, nil
	})

	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)
	setMakersDeploymentClient(a, client)

	attrs := baseMakersDeploymentAttrs()
	attrs["dist_type"] = tftypes.NewValue(tftypes.String, "Zip")
	attrs["repo_branch"] = tftypes.NewValue(tftypes.String, "main")
	attrs["repo_commit_hash"] = tftypes.NewValue(tftypes.String, "9a03d4e")
	attrs["repo_commit_msg"] = tftypes.NewValue(tftypes.String, "init")
	attrs["build_from"] = tftypes.NewValue(tftypes.String, "console")
	attrs["template_id"] = tftypes.NewValue(tftypes.String, "tpl-1")
	attrs["ci_spec"] = tftypes.NewValue(tftypes.String, "small")

	req := newMakersDeploymentInvokeRequest(t, attrs)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotReq)

	// Required attributes must be forwarded exactly as configured.
	assert.Equal(t, "zone-1", *gotReq.ZoneId)
	assert.Equal(t, "project-1", *gotReq.ProjectId)
	assert.Equal(t, "Api", *gotReq.ViaMeta)
	assert.Equal(t, "Production", *gotReq.Env)
	assert.Equal(t, "Github", *gotReq.Provider)

	// Optional attributes set in config must be forwarded.
	assert.NotNil(t, gotReq.DistType)
	assert.Equal(t, "Zip", *gotReq.DistType)
	assert.NotNil(t, gotReq.RepoBranch)
	assert.Equal(t, "main", *gotReq.RepoBranch)
	assert.NotNil(t, gotReq.RepoCommitHash)
	assert.Equal(t, "9a03d4e", *gotReq.RepoCommitHash)
	assert.NotNil(t, gotReq.RepoCommitMsg)
	assert.Equal(t, "init", *gotReq.RepoCommitMsg)
	assert.NotNil(t, gotReq.BuildFrom)
	assert.Equal(t, "console", *gotReq.BuildFrom)
	assert.NotNil(t, gotReq.TemplateId)
	assert.Equal(t, "tpl-1", *gotReq.TemplateId)
	assert.NotNil(t, gotReq.CiSpec)
	assert.Equal(t, "small", *gotReq.CiSpec)

	// Omitted optional attributes must stay absent from the request.
	assert.Nil(t, gotReq.TempBucketPath)
	assert.True(t, calls >= 2, "expected polling to run at least twice")
}

// TestMakersDeployment_Invoke_OptionalOmitted verifies that optional attributes
// omitted from the configuration are not sent to the API.
func TestMakersDeployment_Invoke_OptionalOmitted(t *testing.T) {
	defer shortenMakersDeploymentPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchMakersClient(patches, client)

	var gotReq *makerssdk.CreateDeploymentRequest
	patchCreateDeployment(patches, func(_ context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error) {
		gotReq = req
		resp := makerssdk.NewCreateDeploymentResponse()
		resp.Response = &makerssdk.CreateDeploymentResponseParams{
			DeploymentId: ptrStringMakersDeployment("deploy-2"),
		}
		return resp, nil
	})

	patchDescribeDeployments(patches, func(_ context.Context, req *makerssdk.DescribeDeploymentsRequest) (*makerssdk.DescribeDeploymentsResponse, error) {
		resp := makerssdk.NewDescribeDeploymentsResponse()
		resp.Response = &makerssdk.DescribeDeploymentsResponseParams{
			Deployments: []*makerssdk.Deployment{{
				DeploymentId: ptrStringMakersDeployment(*req.DeploymentIds[0]),
				Status:       ptrStringMakersDeployment(makers.DeploymentStatusSuccess),
			}},
		}
		return resp, nil
	})

	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)
	setMakersDeploymentClient(a, client)

	req := newMakersDeploymentInvokeRequest(t, baseMakersDeploymentAttrs())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotReq)
	assert.Nil(t, gotReq.DistType)
	assert.Nil(t, gotReq.TempBucketPath)
	assert.Nil(t, gotReq.RepoBranch)
	assert.Nil(t, gotReq.RepoCommitHash)
	assert.Nil(t, gotReq.RepoCommitMsg)
	assert.Nil(t, gotReq.BuildFrom)
	assert.Nil(t, gotReq.TemplateId)
	assert.Nil(t, gotReq.CiSpec)
}

// TestMakersDeployment_Invoke_CreateError tests an API error from
// CreateDeployment.
func TestMakersDeployment_Invoke_CreateError(t *testing.T) {
	defer shortenMakersDeploymentPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchMakersClient(patches, client)

	patchCreateDeployment(patches, func(_ context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=bad input")
	})

	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)
	setMakersDeploymentClient(a, client)

	req := newMakersDeploymentInvokeRequest(t, baseMakersDeploymentAttrs())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "triggering Makers deployment")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "InvalidParameter")
}

// TestMakersDeployment_Invoke_NilDeploymentId tests that a CreateDeployment
// success without a deployment id is reported as a failure, because the build
// outcome could not be confirmed.
func TestMakersDeployment_Invoke_NilDeploymentId(t *testing.T) {
	defer shortenMakersDeploymentPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchMakersClient(patches, client)

	patchCreateDeployment(patches, func(_ context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error) {
		resp := makerssdk.NewCreateDeploymentResponse()
		resp.Response = &makerssdk.CreateDeploymentResponseParams{}
		return resp, nil
	})

	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)
	setMakersDeploymentClient(a, client)

	req := newMakersDeploymentInvokeRequest(t, baseMakersDeploymentAttrs())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "no deployment id")
}

// TestMakersDeployment_Invoke_PollFailed tests that polling to a Failed
// terminal state surfaces an error carrying the last observed status snapshot.
func TestMakersDeployment_Invoke_PollFailed(t *testing.T) {
	defer shortenMakersDeploymentPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchMakersClient(patches, client)

	patchCreateDeployment(patches, func(_ context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error) {
		resp := makerssdk.NewCreateDeploymentResponse()
		resp.Response = &makerssdk.CreateDeploymentResponseParams{
			DeploymentId: ptrStringMakersDeployment("deploy-3"),
		}
		return resp, nil
	})

	patchDescribeDeployments(patches, func(_ context.Context, req *makerssdk.DescribeDeploymentsRequest) (*makerssdk.DescribeDeploymentsResponse, error) {
		resp := makerssdk.NewDescribeDeploymentsResponse()
		resp.Response = &makerssdk.DescribeDeploymentsResponseParams{
			Deployments: []*makerssdk.Deployment{{
				DeploymentId: ptrStringMakersDeployment("deploy-3"),
				Status:       ptrStringMakersDeployment(makers.DeploymentStatusFailed),
				Code:         ptrInt64MakersDeployment(1),
			}},
		}
		return resp, nil
	})

	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)
	setMakersDeploymentClient(a, client)

	req := newMakersDeploymentInvokeRequest(t, baseMakersDeploymentAttrs())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for Makers deployment")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), makers.DeploymentStatusFailed)
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "last observed: status=Failed, code=1")
}

// TestMakersDeployment_Invoke_PollTimeout tests that a build which never
// reaches a terminal state fails with a timeout diagnostic carrying the last
// observed snapshot.
func TestMakersDeployment_Invoke_PollTimeout(t *testing.T) {
	defer shortenMakersDeploymentPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchMakersClient(patches, client)

	patchCreateDeployment(patches, func(_ context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error) {
		resp := makerssdk.NewCreateDeploymentResponse()
		resp.Response = &makerssdk.CreateDeploymentResponseParams{
			DeploymentId: ptrStringMakersDeployment("deploy-4"),
		}
		return resp, nil
	})

	patchDescribeDeployments(patches, func(_ context.Context, req *makerssdk.DescribeDeploymentsRequest) (*makerssdk.DescribeDeploymentsResponse, error) {
		resp := makerssdk.NewDescribeDeploymentsResponse()
		resp.Response = &makerssdk.DescribeDeploymentsResponseParams{
			Deployments: []*makerssdk.Deployment{{
				DeploymentId: ptrStringMakersDeployment("deploy-4"),
				Status:       ptrStringMakersDeployment(makers.DeploymentStatusProcess),
			}},
		}
		return resp, nil
	})

	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)
	setMakersDeploymentClient(a, client)

	timeoutVal := tftypes.NewValue(makersDeploymentTimeoutsObjType(), map[string]tftypes.Value{
		"invoke": tftypes.NewValue(tftypes.String, "50ms"),
	})
	req := newMakersDeploymentInvokeRequestWithTimeouts(t, baseMakersDeploymentAttrs(), timeoutVal)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for Makers deployment")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "did not converge")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "last observed: status=Process")
}

// TestMakersDeployment_Invoke_PollQueryError tests an error from the
// DescribeDeployments poll query.
func TestMakersDeployment_Invoke_PollQueryError(t *testing.T) {
	defer shortenMakersDeploymentPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchMakersClient(patches, client)

	patchCreateDeployment(patches, func(_ context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error) {
		resp := makerssdk.NewCreateDeploymentResponse()
		resp.Response = &makerssdk.CreateDeploymentResponseParams{
			DeploymentId: ptrStringMakersDeployment("deploy-5"),
		}
		return resp, nil
	})

	patchDescribeDeployments(patches, func(_ context.Context, req *makerssdk.DescribeDeploymentsRequest) (*makerssdk.DescribeDeploymentsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=boom")
	})

	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)
	setMakersDeploymentClient(a, client)

	req := newMakersDeploymentInvokeRequest(t, baseMakersDeploymentAttrs())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for Makers deployment")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "InternalError")
}

// TestMakersDeployment_Invoke_PollNotVisibleYet tests that a deployment which
// is not yet visible in DescribeDeployments keeps the poll loop waiting until
// it shows up and reaches a terminal state.
func TestMakersDeployment_Invoke_PollNotVisibleYet(t *testing.T) {
	defer shortenMakersDeploymentPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchMakersClient(patches, client)

	patchCreateDeployment(patches, func(_ context.Context, req *makerssdk.CreateDeploymentRequest) (*makerssdk.CreateDeploymentResponse, error) {
		resp := makerssdk.NewCreateDeploymentResponse()
		resp.Response = &makerssdk.CreateDeploymentResponseParams{
			DeploymentId: ptrStringMakersDeployment("deploy-6"),
		}
		return resp, nil
	})

	calls := 0
	patchDescribeDeployments(patches, func(_ context.Context, req *makerssdk.DescribeDeploymentsRequest) (*makerssdk.DescribeDeploymentsResponse, error) {
		calls++
		resp := makerssdk.NewDescribeDeploymentsResponse()
		// The first poll returns no matching deployment (eventual consistency
		// window); the second one reports the build as already finished.
		if calls == 1 {
			resp.Response = &makerssdk.DescribeDeploymentsResponseParams{
				TotalCount: ptrInt64MakersDeployment(0),
			}
			return resp, nil
		}
		resp.Response = &makerssdk.DescribeDeploymentsResponseParams{
			TotalCount: ptrInt64MakersDeployment(1),
			Deployments: []*makerssdk.Deployment{{
				DeploymentId: ptrStringMakersDeployment("deploy-6"),
				Status:       ptrStringMakersDeployment(makers.DeploymentStatusSuccess),
			}},
		}
		return resp, nil
	})

	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)
	setMakersDeploymentClient(a, client)

	req := newMakersDeploymentInvokeRequest(t, baseMakersDeploymentAttrs())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.True(t, calls >= 2, "expected the poll loop to keep waiting after a not-found result")
}

// TestMakersDeployment_Invoke_ClientNotConfigured tests the nil-client guard.
func TestMakersDeployment_Invoke_ClientNotConfigured(t *testing.T) {
	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)

	req := newMakersDeploymentInvokeRequest(t, baseMakersDeploymentAttrs())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestMakersDeployment_MetadataAndSchema validates metadata and schema
// definition: the type name, the exact attribute set required by the
// CreateDeployment API, and the timeouts block.
func TestMakersDeployment_MetadataAndSchema(t *testing.T) {
	aImpl := makers.NewMakersDeployment()
	a := aImpl.(*makers.MakersDeployment)

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_makers_deployment", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	// The CreateDeployment API carries exactly these fields; no other business
	// attribute is allowed to leak into the schema.
	assert.Len(t, attrs, 13)

	required := []string{"zone_id", "project_id", "via_meta", "env", "provider"}
	for _, name := range required {
		assert.Contains(t, attrs, name)
		assert.True(t, attrs[name].IsRequired(), "%s must be required", name)
		assert.False(t, attrs[name].IsOptional(), "%s must not be optional", name)
	}

	optional := []string{
		"dist_type", "temp_bucket_path", "repo_branch", "repo_commit_hash",
		"repo_commit_msg", "build_from", "template_id", "ci_spec",
	}
	for _, name := range optional {
		assert.Contains(t, attrs, name)
		assert.True(t, attrs[name].IsOptional(), "%s must be optional", name)
		assert.False(t, attrs[name].IsRequired(), "%s must not be required", name)
	}

	blocks := schemaResp.Schema.Blocks
	assert.Contains(t, blocks, "timeouts")
	timeoutsBlock, ok := blocks["timeouts"].(action_schema.SingleNestedBlock)
	assert.True(t, ok)
	assert.IsType(t, actiontimeouts.Type{}, timeoutsBlock.CustomType)
	invokeAttr, ok := timeoutsBlock.Attributes["invoke"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, invokeAttr.Optional)
	assert.NotEmpty(t, invokeAttr.Validators)
}
