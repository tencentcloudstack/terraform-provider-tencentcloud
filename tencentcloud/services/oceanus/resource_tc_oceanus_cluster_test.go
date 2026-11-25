package oceanus

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	oceanussdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	regionv20220627 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const (
	oceanusTestClusterId   = "cluster-tf123456"
	oceanusTestClusterName = "tf-example-cluster"
	oceanusTestCosBucket   = "tf-example-bucket-1250000000"
)

// shortenPolling makes the async wait helpers converge immediately under
// mocked API responses. It returns a restore function to be deferred.
func shortenPolling() func() {
	origDelay, origMinTimeout, origTimeout := oceanusClusterPollDelay, oceanusClusterPollMinTimeout, oceanusClusterPollTimeout
	oceanusClusterPollDelay = 0
	oceanusClusterPollMinTimeout = 0
	oceanusClusterPollTimeout = 10 * time.Second
	return func() {
		oceanusClusterPollDelay, oceanusClusterPollMinTimeout, oceanusClusterPollTimeout = origDelay, origMinTimeout, origTimeout
	}
}

// oceanusClusterTestContext wires the resource to a patched SDK client.
type oceanusClusterTestContext struct {
	resource resource.Resource
	client   *oceanussdk.Client
	schema   schema.Schema
}

func newOceanusClusterTestContext(t *testing.T, patches *gomonkey.Patches) *oceanusClusterTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	oceanusClient := &oceanussdk.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseOceanusClient", oceanusClient)

	res := NewOceanusClusterResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}
	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp,
	)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	return &oceanusClusterTestContext{resource: res, client: oceanusClient, schema: schemaResp.Schema}
}

func (c *oceanusClusterTestContext) vpcList(t *testing.T) types.List {
	t.Helper()
	elemType := c.schema.Attributes["vpc_descriptions"].GetType().(attr.TypeWithElementType).ElementType()
	list, diags := types.ListValueFrom(context.Background(), elemType, []OceanusClusterVpcDescriptionModel{
		{
			VpcId:    types.StringValue("vpc-tf1234"),
			SubnetId: types.StringValue("subnet-tf1234"),
		},
	})
	if diags.HasError() {
		t.Fatalf("failed to build vpc list: %v", diags.Errors())
	}
	return list
}

// config mirrors the plan construction for the Config side of a create or
// update request: write-only attribute values (login_password) are only
// provided through the configuration, never through the plan.
func (c *oceanusClusterTestContext) config(t *testing.T, model OceanusClusterResourceModel) tfsdk.Config {
	t.Helper()
	plan := c.plan(t, model)
	return tfsdk.Config{Schema: c.schema, Raw: plan.Raw}
}

func (c *oceanusClusterTestContext) plan(t *testing.T, model OceanusClusterResourceModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *oceanusClusterTestContext) state(t *testing.T, model OceanusClusterResourceModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

// nullTimeouts builds a null timeouts.Value matching the resource schema, so
// model (de)serialization through tfsdk.Plan/State behaves exactly like a
// configuration that does not set the timeouts block.
func (c *oceanusClusterTestContext) nullTimeouts() resourcetimeouts.Value {
	timeoutsBlock, ok := c.schema.Blocks["timeouts"]
	if !ok {
		panic("timeouts block not found in schema")
	}
	timeoutsType, ok := timeoutsBlock.Type().(resourcetimeouts.Type)
	if !ok {
		panic("unexpected timeouts block type")
	}
	return resourcetimeouts.Value{
		Object: types.ObjectNull(timeoutsType.AttributeTypes()),
	}
}

// oceanusClusterTestModel builds a full resource model. When withPassword is
// true, login_password carries a plaintext value (plan side); otherwise it is
// null (state side, write-only). When id is empty the id is unknown so the
// plan behaves like a create.
func (c *oceanusClusterTestContext) oceanusClusterTestModel(t *testing.T, id string, withPassword bool, cu int64, period *int64) OceanusClusterResourceModel {
	t.Helper()
	idValue := types.StringNull()
	if id != "" {
		idValue = types.StringValue(id)
	}
	periodValue := types.Int64Null()
	if period != nil {
		periodValue = types.Int64Value(*period)
	}
	return OceanusClusterResourceModel{
		Id:                          idValue,
		ClusterName:                 types.StringValue(oceanusTestClusterName),
		Region:                      types.StringValue("ap-guangzhou"),
		Zone:                        types.StringValue("ap-guangzhou-1"),
		LoginPassword:               types.StringPointerValue(strOrNil(withPassword, "Password@123")),
		VpcDescriptions:             c.vpcList(t),
		DefaultCosBucket:            types.StringValue(oceanusTestCosBucket),
		Cu:                          types.Int64Value(cu),
		Remark:                      types.StringValue("managed by terraform"),
		Period:                      periodValue,
		InstanceChargeType:          types.StringValue("POSTPAID_BY_SECOND"),
		ClusterType:                 types.StringNull(),
		RenewFlag:                   types.StringNull(),
		FlinkWebUiNetworkAccessType: types.StringNull(),
		SlaveVpcDescriptions:        types.ListNull(c.vpcListElemType(t)),
		CuMemory:                    types.Int64Null(),
		Timeouts:                    c.nullTimeouts(),
	}
}

func (c *oceanusClusterTestContext) vpcListElemType(t *testing.T) attr.Type {
	t.Helper()
	return c.schema.Attributes["vpc_descriptions"].GetType().(attr.TypeWithElementType).ElementType()
}

func strOrNil(set bool, v string) *string {
	if !set {
		return nil
	}
	return &v
}

func oceanusClusterTestItem(status int64) *oceanussdk.Cluster {
	return &oceanussdk.Cluster{
		ClusterId:        strOrNil(true, oceanusTestClusterId),
		Name:             strOrNil(true, oceanusTestClusterName),
		Region:           strOrNil(true, "ap-guangzhou"),
		Status:           int64OrNil(status),
		CuNum:            uint64OrNil(12),
		Remark:           strOrNil(true, "managed by terraform"),
		DefaultCOSBucket: strOrNil(true, oceanusTestCosBucket),
		Zone:             strOrNil(true, "ap-guangzhou-1"),
	}
}

func int64OrNil(v int64) *int64 { return &v }

func uint64OrNil(v uint64) *uint64 { return &v }

func oceanusTestDescribeResponse(clusters ...*oceanussdk.Cluster) *oceanussdk.DescribeClustersResponse {
	return &oceanussdk.DescribeClustersResponse{
		Response: &oceanussdk.DescribeClustersResponseParams{
			TotalCount: int64OrNil(int64(len(clusters))),
			ClusterSet: clusters,
		},
	}
}

// patchOceanusRegionResolution stubs the region management Describe calls so
// the Create flows can run without the real region service. The name -> id
// translation logic itself keeps running: ap-guangzhou resolves to region id
// 1 and ap-guangzhou-1 to zone id 100001. The stubs target the Describe
// methods (large retry functions, never inlined) instead of the small
// translation helpers, which the compiler inlines into the Create call site.
func patchOceanusRegionResolution(patches *gomonkey.Patches) {
	patches.ApplyMethodFunc(&OceanusService{}, "DescribeOceanusRegions",
		func(_ context.Context) ([]*regionv20220627.RegionInfo, error) {
			return []*regionv20220627.RegionInfo{
				{Region: strOrNil(true, "ap-guangzhou"), RegionIdMC: strOrNil(true, "1")},
			}, nil
		})
	patches.ApplyMethodFunc(&OceanusService{}, "DescribeOceanusZones",
		func(_ context.Context) ([]*regionv20220627.ZoneInfo, error) {
			return []*regionv20220627.ZoneInfo{
				{Zone: strOrNil(true, "ap-guangzhou-1"), ZoneId: strOrNil(true, "100001")},
			}, nil
		})
}

// TestOceanusClusterCreateSuccess verifies that Create sends all configured
// fields, persists the returned ClusterId before polling, and refreshes the
// read-back fields once the cluster is running.
func TestOceanusClusterCreateSuccess(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patchOceanusRegionResolution(patches)

	var captured *oceanussdk.CreateOceanusClusterRequest
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "CreateOceanusClusterWithContext",
		func(_ context.Context, req *oceanussdk.CreateOceanusClusterRequest) (*oceanussdk.CreateOceanusClusterResponse, error) {
			captured = req
			return &oceanussdk.CreateOceanusClusterResponse{
				Response: &oceanussdk.CreateOceanusClusterResponseParams{
					ClusterId: strOrNil(true, oceanusTestClusterId),
				},
			}, nil
		})
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, req *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			if len(req.ClusterIds) != 1 || *req.ClusterIds[0] != oceanusTestClusterId {
				t.Errorf("DescribeClusters request ids = %v, want [%s]", req.ClusterIds, oceanusTestClusterId)
			}
			return oceanusTestDescribeResponse(oceanusClusterTestItem(2)), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	model := ctx.oceanusClusterTestModel(t, "", true, 12, nil)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model), Config: ctx.config(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateOceanusCluster request is nil")
	}
	if captured.ClusterName == nil || *captured.ClusterName != oceanusTestClusterName {
		t.Errorf("request.ClusterName = %v, want %q", captured.ClusterName, oceanusTestClusterName)
	}
	if captured.RegionId == nil || *captured.RegionId != 1 {
		t.Errorf("request.RegionId = %v, want 1", captured.RegionId)
	}
	if captured.ZoneId == nil || *captured.ZoneId != 100001 {
		t.Errorf("request.ZoneId = %v, want 100001", captured.ZoneId)
	}
	if captured.LoginPassword == nil || *captured.LoginPassword != "Password@123" {
		t.Errorf("request.LoginPassword = %v, want the configured password", captured.LoginPassword)
	}
	if captured.DefaultCOSBucket == nil || *captured.DefaultCOSBucket != oceanusTestCosBucket {
		t.Errorf("request.DefaultCOSBucket = %v, want %q", captured.DefaultCOSBucket, oceanusTestCosBucket)
	}
	if captured.CU == nil || *captured.CU != 12 {
		t.Errorf("request.CU = %v, want 12", captured.CU)
	}
	if len(captured.VpcDescriptions) != 1 || captured.VpcDescriptions[0].VpcId == nil || *captured.VpcDescriptions[0].VpcId != "vpc-tf1234" {
		t.Errorf("request.VpcDescriptions = %v, want one entry with vpc-tf1234", captured.VpcDescriptions)
	}
	if captured.Period != nil {
		t.Errorf("request.Period = %v, want nil (unset)", captured.Period)
	}

	var got OceanusClusterResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != oceanusTestClusterId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), oceanusTestClusterId)
	}
	if !got.LoginPassword.IsNull() {
		t.Errorf("login_password = %v, want null (write-only)", got.LoginPassword)
	}
	if got.ClusterName.ValueString() != oceanusTestClusterName {
		t.Errorf("cluster_name = %q, want %q", got.ClusterName.ValueString(), oceanusTestClusterName)
	}
	if got.DefaultCosBucket.ValueString() != oceanusTestCosBucket {
		t.Errorf("default_cos_bucket = %q, want %q", got.DefaultCosBucket.ValueString(), oceanusTestCosBucket)
	}
	if got.Cu.ValueInt64() != 12 {
		t.Errorf("cu = %d, want 12", got.Cu.ValueInt64())
	}
	if got.Region.ValueString() != "ap-guangzhou" {
		t.Errorf("region = %q, want ap-guangzhou", got.Region.ValueString())
	}
	var vpcs []OceanusClusterVpcDescriptionModel
	if diags := got.VpcDescriptions.ElementsAs(context.Background(), &vpcs, false); diags.HasError() {
		t.Fatalf("failed to read back vpc_descriptions: %v", diags.Errors())
	}
	if len(vpcs) != 1 || vpcs[0].VpcId.ValueString() != "vpc-tf1234" || vpcs[0].SubnetId.ValueString() != "subnet-tf1234" {
		t.Errorf("vpc_descriptions = %v, want one entry vpc-tf1234/subnet-tf1234", vpcs)
	}
}

// TestOceanusClusterCreateApiError verifies Create surfaces the API error.
func TestOceanusClusterCreateApiError(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patchOceanusRegionResolution(patches)
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "CreateOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.CreateOceanusClusterRequest) (*oceanussdk.CreateOceanusClusterResponse, error) {
			return nil, fmt.Errorf("create failed")
		})

	ctx := newOceanusClusterTestContext(t, patches)
	model := ctx.oceanusClusterTestModel(t, "", true, 12, nil)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model), Config: ctx.config(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an error")
	}
}

// TestOceanusClusterCreateEmptyId verifies Create fails loudly when the API
// returns no ClusterId instead of writing an unmanageable resource.
func TestOceanusClusterCreateEmptyId(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patchOceanusRegionResolution(patches)
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "CreateOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.CreateOceanusClusterRequest) (*oceanussdk.CreateOceanusClusterResponse, error) {
			return &oceanussdk.CreateOceanusClusterResponse{
				Response: &oceanussdk.CreateOceanusClusterResponseParams{},
			}, nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	model := ctx.oceanusClusterTestModel(t, "", true, 12, nil)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model), Config: ctx.config(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns no ClusterId")
	}
}

// TestOceanusClusterReadSuccess verifies Read refreshes the read-back fields
// and preserves the write-only and no-read-back fields.
func TestOceanusClusterReadSuccess(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, req *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			if len(req.ClusterIds) != 1 || *req.ClusterIds[0] != oceanusTestClusterId {
				t.Errorf("DescribeClusters request ids = %v, want [%s]", req.ClusterIds, oceanusTestClusterId)
			}
			return oceanusTestDescribeResponse(oceanusClusterTestItem(2)), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, stateModel)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, stateModel)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got OceanusClusterResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != oceanusTestClusterId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), oceanusTestClusterId)
	}
	if got.ClusterName.ValueString() != oceanusTestClusterName {
		t.Errorf("cluster_name = %q, want %q", got.ClusterName.ValueString(), oceanusTestClusterName)
	}
	if got.Cu.ValueInt64() != 12 {
		t.Errorf("cu = %d, want 12", got.Cu.ValueInt64())
	}
	if got.Region.ValueString() != "ap-guangzhou" {
		t.Errorf("region = %q, want preserved ap-guangzhou", got.Region.ValueString())
	}
	if got.Zone.ValueString() != "ap-guangzhou-1" {
		t.Errorf("zone = %q, want preserved ap-guangzhou-1", got.Zone.ValueString())
	}
	if got.InstanceChargeType.ValueString() != "POSTPAID_BY_SECOND" {
		t.Errorf("instance_charge_type = %q, want preserved POSTPAID_BY_SECOND", got.InstanceChargeType.ValueString())
	}
	if !got.LoginPassword.IsNull() {
		t.Errorf("login_password = %v, want null (write-only)", got.LoginPassword)
	}
	if !got.Period.IsNull() {
		t.Errorf("period = %v, want preserved null", got.Period)
	}
}

// TestOceanusClusterReadNotFound verifies Read removes the state when the
// cluster no longer exists.
func TestOceanusClusterReadNotFound(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, _ *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			return oceanusTestDescribeResponse(), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil)
	readResp := resource.ReadResponse{State: ctx.state(t, stateModel)}
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, stateModel)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}
	// RemoveResource rewrites the state raw value to null.
	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state should be removed for a missing cluster, got %s", readResp.State.Raw.String())
	}
}

// TestOceanusClusterUpdateScaleAndRenew verifies that a CU change goes through
// ScaleOceanusCluster with the right ScaleMode and a period change goes
// through RenewOceanusCluster.
func TestOceanusClusterUpdateScaleAndRenew(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var scaleReq *oceanussdk.ScaleOceanusClusterRequest
	var renewReq *oceanussdk.RenewOceanusClusterRequest
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "ScaleOceanusClusterWithContext",
		func(_ context.Context, req *oceanussdk.ScaleOceanusClusterRequest) (*oceanussdk.ScaleOceanusClusterResponse, error) {
			scaleReq = req
			return &oceanussdk.ScaleOceanusClusterResponse{
				Response: &oceanussdk.ScaleOceanusClusterResponseParams{
					TaskExecResult: strOrNil(true, "success"),
				},
			}, nil
		})
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "RenewOceanusClusterWithContext",
		func(_ context.Context, req *oceanussdk.RenewOceanusClusterRequest) (*oceanussdk.RenewOceanusClusterResponse, error) {
			renewReq = req
			return &oceanussdk.RenewOceanusClusterResponse{
				Response: &oceanussdk.RenewOceanusClusterResponseParams{
					TaskExecResult: strOrNil(true, "success"),
				},
			}, nil
		})
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, _ *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			// The describe happens after the resize converged: report the
			// scaled-up cluster.
			item := oceanusClusterTestItem(2)
			item.CuNum = uint64OrNil(19)
			return oceanusTestDescribeResponse(item), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	period := int64(12)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil)
	planModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, true, 19, &period)
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(), resource.UpdateRequest{Plan: ctx.plan(t, planModel), State: ctx.state(t, stateModel)}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if scaleReq == nil {
		t.Fatal("ScaleOceanusCluster was not called for a CU change")
	}
	if scaleReq.ClusterId == nil || *scaleReq.ClusterId != oceanusTestClusterId {
		t.Errorf("ScaleOceanusCluster request ClusterId = %v, want %s", scaleReq.ClusterId, oceanusTestClusterId)
	}
	if scaleReq.NewCU == nil || *scaleReq.NewCU != 19 {
		t.Errorf("ScaleOceanusCluster request NewCU = %v, want 19", scaleReq.NewCU)
	}
	if scaleReq.ScaleMode == nil || *scaleReq.ScaleMode != "ScaleUp" {
		t.Errorf("ScaleOceanusCluster request ScaleMode = %v, want ScaleUp", scaleReq.ScaleMode)
	}
	if renewReq == nil {
		t.Fatal("RenewOceanusCluster was not called for a period change")
	}
	if renewReq.Period == nil || *renewReq.Period != 12 {
		t.Errorf("RenewOceanusCluster request Period = %v, want 12", renewReq.Period)
	}

	var got OceanusClusterResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Cu.ValueInt64() != 19 {
		t.Errorf("cu = %d, want 19", got.Cu.ValueInt64())
	}
	if got.Period.ValueInt64() != 12 {
		t.Errorf("period = %d, want 12", got.Period.ValueInt64())
	}
	if !got.LoginPassword.IsNull() {
		t.Errorf("login_password = %v, want null (write-only)", got.LoginPassword)
	}
}

// TestOceanusClusterUpdateScaleDown verifies the ScaleMode is ScaleDown when
// the target CU is smaller than the current one.
func TestOceanusClusterUpdateScaleDown(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var scaleReq *oceanussdk.ScaleOceanusClusterRequest
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "ScaleOceanusClusterWithContext",
		func(_ context.Context, req *oceanussdk.ScaleOceanusClusterRequest) (*oceanussdk.ScaleOceanusClusterResponse, error) {
			scaleReq = req
			return &oceanussdk.ScaleOceanusClusterResponse{
				Response: &oceanussdk.ScaleOceanusClusterResponseParams{
					TaskExecResult: strOrNil(true, "success"),
				},
			}, nil
		})
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, _ *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			return oceanusTestDescribeResponse(oceanusClusterTestItem(2)), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 19, nil)
	planModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, true, 12, nil)
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(), resource.UpdateRequest{Plan: ctx.plan(t, planModel), State: ctx.state(t, stateModel)}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if scaleReq == nil || scaleReq.ScaleMode == nil || *scaleReq.ScaleMode != "ScaleDown" {
		t.Fatalf("ScaleOceanusCluster request ScaleMode = %v, want ScaleDown", scaleReq)
	}
}

// TestOceanusClusterDeleteSuccess verifies Delete reports success and waits
// until DescribeClusters no longer returns the cluster.
func TestOceanusClusterDeleteSuccess(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *oceanussdk.DeleteOceanusClusterRequest
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DeleteOceanusClusterWithContext",
		func(_ context.Context, req *oceanussdk.DeleteOceanusClusterRequest) (*oceanussdk.DeleteOceanusClusterResponse, error) {
			captured = req
			return &oceanussdk.DeleteOceanusClusterResponse{
				Response: &oceanussdk.DeleteOceanusClusterResponseParams{
					TaskExecResult: strOrNil(true, "success"),
				},
			}, nil
		})
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, req *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			if len(req.ClusterIds) != 1 || *req.ClusterIds[0] != oceanusTestClusterId {
				t.Errorf("DescribeClusters request ids = %v, want [%s]", req.ClusterIds, oceanusTestClusterId)
			}
			return oceanusTestDescribeResponse(), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	deleteResp.State = ctx.state(t, ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil))
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: deleteResp.State}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
	if captured == nil || captured.ClusterId == nil || *captured.ClusterId != oceanusTestClusterId {
		t.Errorf("DeleteOceanusCluster request ClusterId = %v, want %s", captured, oceanusTestClusterId)
	}
}

// TestOceanusClusterDeleteApiError verifies Delete surfaces the API error.
func TestOceanusClusterDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DeleteOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.DeleteOceanusClusterRequest) (*oceanussdk.DeleteOceanusClusterResponse, error) {
			return nil, fmt.Errorf("delete failed")
		})

	ctx := newOceanusClusterTestContext(t, patches)
	var deleteResp resource.DeleteResponse
	deleteResp.State = ctx.state(t, ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil))
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: deleteResp.State}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail when the API returns an error")
	}
}

// TestOceanusClusterCreateTimeoutPassedToWait verifies that the
// `timeouts.create` attribute bounds the creation wait loop: a cluster stuck
// in the initializing state must fail after the configured short timeout
// instead of the built-in default.
func TestOceanusClusterCreateTimeoutPassedToWait(t *testing.T) {
	restore := shortenPolling()
	defer restore()

	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patchOceanusRegionResolution(patches)

	patches.ApplyMethodFunc(&oceanussdk.Client{}, "CreateOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.CreateOceanusClusterRequest) (*oceanussdk.CreateOceanusClusterResponse, error) {
			return &oceanussdk.CreateOceanusClusterResponse{
				Response: &oceanussdk.CreateOceanusClusterResponseParams{
					ClusterId: strOrNil(true, oceanusTestClusterId),
				},
			}, nil
		})
	// The cluster never leaves the initializing state.
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, _ *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			return oceanusTestDescribeResponse(oceanusClusterTestItem(3)), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	model := ctx.oceanusClusterTestModel(t, "", true, 12, nil)
	model.Timeouts = ctx.timeoutsWith(t, "create", "1s")

	createResp := resource.CreateResponse{State: ctx.state(t, model)}
	ctx.resource.Create(context.Background(), resource.CreateRequest{
		Plan:   ctx.plan(t, model),
		Config: ctx.config(t, model),
	}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the cluster never becomes ready")
	}

	summary := createResp.Diagnostics.Errors()[0].Summary()
	if !strings.Contains(summary, "Wait for Oceanus cluster ready failed") {
		t.Fatalf("unexpected error summary: %s", summary)
	}
}

// timeoutsWith builds a timeouts.Value with the given stage set to a duration
// string and every other stage null.
func (c *oceanusClusterTestContext) timeoutsWith(t *testing.T, stage, value string) resourcetimeouts.Value {
	t.Helper()
	timeoutsBlock, ok := c.schema.Blocks["timeouts"]
	if !ok {
		t.Fatalf("timeouts block not found in schema")
	}
	timeoutsType, ok := timeoutsBlock.Type().(resourcetimeouts.Type)
	if !ok {
		t.Fatalf("timeouts does not have the framework timeouts type: %T", timeoutsBlock.Type())
	}

	attrTypes := timeoutsType.AttributeTypes()
	attributes := map[string]attr.Value{}
	for name := range attrTypes {
		if name == stage {
			attributes[name] = types.StringValue(value)
		} else {
			attributes[name] = types.StringNull()
		}
	}
	object, diags := types.ObjectValue(attrTypes, attributes)
	if diags.HasError() {
		t.Fatalf("failed to build timeouts object: %v", diags.Errors())
	}
	timeoutsValue, diags := timeoutsType.ValueFromObject(context.Background(), object)
	if diags.HasError() {
		t.Fatalf("failed to build timeouts value: %v", diags.Errors())
	}
	result, ok := timeoutsValue.(resourcetimeouts.Value)
	if !ok {
		t.Fatalf("timeouts value has unexpected type %T", timeoutsValue)
	}
	return result
}

// TestOceanusClusterImportState verifies that the resource supports import via
// the id attribute passthrough.
func TestOceanusClusterImportState(t *testing.T) {
	ctx := newOceanusClusterTestContext(t, gomonkey.NewPatches())
	importer, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}
	var importResp resource.ImportStateResponse
	importResp.State.Schema = ctx.schema
	// The framework server initializes the import state raw value to a fully
	// null object of the schema type before the provider ImportState runs;
	// mirror that so SetAttribute can store into it.
	importResp.State.Raw = tftypes.NewValue(ctx.schema.Type().TerraformType(context.Background()), nil)
	importer.ImportState(context.Background(), resource.ImportStateRequest{
		ID: oceanusTestClusterId,
	}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}
	var got OceanusClusterResourceModel
	if diags := importResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back imported state: %v", diags.Errors())
	}
	if got.Id.ValueString() != oceanusTestClusterId {
		t.Errorf("imported id = %q, want %q", got.Id.ValueString(), oceanusTestClusterId)
	}
}

// TestOceanusClusterReadRefreshesReadableFields verifies that Read refreshes
// the attributes with a confirmed read-back path and preserves every attribute
// that DescribeClusters cannot echo back.
func TestOceanusClusterReadRefreshesReadableFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cluster := &oceanussdk.Cluster{
		ClusterId:        strOrNil(true, oceanusTestClusterId),
		Name:             strOrNil(true, "cloud-side-name"),
		Status:           int64OrNil(2),
		CuNum:            uint64OrNil(19),
		Remark:           strOrNil(true, "cloud-side-remark"),
		DefaultCOSBucket: strOrNil(true, "cloud-side-bucket-1250000000"),
	}
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, _ *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			return oceanusTestDescribeResponse(cluster), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	period := int64(12)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, &period)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, stateModel)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, stateModel)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got OceanusClusterResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != oceanusTestClusterId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), oceanusTestClusterId)
	}
	if got.ClusterName.ValueString() != "cloud-side-name" {
		t.Errorf("cluster_name = %q, want cloud-side-name", got.ClusterName.ValueString())
	}
	if got.Cu.ValueInt64() != 19 {
		t.Errorf("cu = %d, want 19", got.Cu.ValueInt64())
	}
	if got.Remark.ValueString() != "cloud-side-remark" {
		t.Errorf("remark = %q, want cloud-side-remark", got.Remark.ValueString())
	}
	if got.DefaultCosBucket.ValueString() != "cloud-side-bucket-1250000000" {
		t.Errorf("default_cos_bucket = %q, want cloud-side-bucket-1250000000", got.DefaultCosBucket.ValueString())
	}
	if !got.LoginPassword.IsNull() {
		t.Errorf("login_password = %v, want null (write-only)", got.LoginPassword)
	}
	if got.Region.ValueString() != "ap-guangzhou" {
		t.Errorf("region = %q, want ap-guangzhou (preserved from state)", got.Region.ValueString())
	}
	if got.Zone.ValueString() != "ap-guangzhou-1" {
		t.Errorf("zone = %q, want ap-guangzhou-1 (preserved from state)", got.Zone.ValueString())
	}
	if got.Period.ValueInt64() != 12 {
		t.Errorf("period = %d, want 12 (preserved from state)", got.Period.ValueInt64())
	}
	if got.InstanceChargeType.ValueString() != "POSTPAID_BY_SECOND" {
		t.Errorf("instance_charge_type = %q, want POSTPAID_BY_SECOND (preserved from state)", got.InstanceChargeType.ValueString())
	}
}

// TestOceanusClusterReadApiError verifies Read surfaces the API error.
func TestOceanusClusterReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, _ *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			return nil, fmt.Errorf("describe failed")
		})

	ctx := newOceanusClusterTestContext(t, patches)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, stateModel)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, stateModel)}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail when the API returns an error")
	}
}

// TestOceanusClusterUpdateNoChanges verifies that an update which changes
// neither cu nor period does not call any mutation API.
func TestOceanusClusterUpdateNoChanges(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	scaleCalled := false
	renewCalled := false
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "ScaleOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.ScaleOceanusClusterRequest) (*oceanussdk.ScaleOceanusClusterResponse, error) {
			scaleCalled = true
			return nil, fmt.Errorf("ScaleOceanusCluster must not be called")
		})
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "RenewOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.RenewOceanusClusterRequest) (*oceanussdk.RenewOceanusClusterResponse, error) {
			renewCalled = true
			return nil, fmt.Errorf("RenewOceanusCluster must not be called")
		})
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DescribeClustersWithContext",
		func(_ context.Context, _ *oceanussdk.DescribeClustersRequest) (*oceanussdk.DescribeClustersResponse, error) {
			return oceanusTestDescribeResponse(oceanusClusterTestItem(2)), nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil)
	planModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, true, 12, nil)
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(), resource.UpdateRequest{
		Plan:  ctx.plan(t, planModel),
		State: ctx.state(t, stateModel),
	}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if scaleCalled {
		t.Error("ScaleOceanusCluster should not be called when cu is unchanged")
	}
	if renewCalled {
		t.Error("RenewOceanusCluster should not be called when period is unchanged")
	}

	var got OceanusClusterResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if !got.LoginPassword.IsNull() {
		t.Errorf("login_password = %v, want null (write-only)", got.LoginPassword)
	}
	if got.Cu.ValueInt64() != 12 {
		t.Errorf("cu = %d, want 12", got.Cu.ValueInt64())
	}
}

// TestOceanusClusterUpdateScaleTaskFailure verifies that a non-successful
// TaskExecResult fails the update instead of silently continuing.
func TestOceanusClusterUpdateScaleTaskFailure(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "ScaleOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.ScaleOceanusClusterRequest) (*oceanussdk.ScaleOceanusClusterResponse, error) {
			return &oceanussdk.ScaleOceanusClusterResponse{
				Response: &oceanussdk.ScaleOceanusClusterResponseParams{
					TaskExecResult: strOrNil(true, "failed"),
				},
			}, nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil)
	planModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, true, 19, nil)
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(), resource.UpdateRequest{
		Plan:  ctx.plan(t, planModel),
		State: ctx.state(t, stateModel),
	}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail when ScaleOceanusCluster does not report success")
	}
}

// TestOceanusClusterCreateNilResponse verifies Create fails when the API
// returns no usable payload.
func TestOceanusClusterCreateNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patchOceanusRegionResolution(patches)
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "CreateOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.CreateOceanusClusterRequest) (*oceanussdk.CreateOceanusClusterResponse, error) {
			return &oceanussdk.CreateOceanusClusterResponse{}, nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	model := ctx.oceanusClusterTestModel(t, "", true, 12, nil)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model), Config: ctx.config(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns no ClusterId")
	}
}

// TestOceanusClusterDeleteTaskFailure verifies that a non-successful
// TaskExecResult fails the deletion instead of treating it as done.
func TestOceanusClusterDeleteTaskFailure(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyMethodFunc(&oceanussdk.Client{}, "DeleteOceanusClusterWithContext",
		func(_ context.Context, _ *oceanussdk.DeleteOceanusClusterRequest) (*oceanussdk.DeleteOceanusClusterResponse, error) {
			return &oceanussdk.DeleteOceanusClusterResponse{
				Response: &oceanussdk.DeleteOceanusClusterResponseParams{
					TaskExecResult: strOrNil(true, "failed"),
				},
			}, nil
		})

	ctx := newOceanusClusterTestContext(t, patches)
	stateModel := ctx.oceanusClusterTestModel(t, oceanusTestClusterId, false, 12, nil)
	var deleteResp resource.DeleteResponse
	deleteResp.State = ctx.state(t, stateModel)
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, stateModel)}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail when DeleteOceanusCluster does not report success")
	}
}
