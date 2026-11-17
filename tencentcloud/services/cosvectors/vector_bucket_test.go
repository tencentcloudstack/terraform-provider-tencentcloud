package cosvectors

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

func cosVectorTestState(t *testing.T, r *CosVectorBucketResource, data CosVectorBucketResourceModel) tfsdk.State {
	t.Helper()
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	state := tfsdk.State{Schema: schemaResp.Schema}
	if diags := state.Set(context.Background(), &data); diags.HasError() {
		t.Fatal(diags)
	}
	return state
}

func TestCosVectorBucketCleanupFailurePreservesState(t *testing.T) {
	for _, failure := range []string{"deadline", "permission", "listNotFound"} {
		t.Run(failure, func(t *testing.T) {
			service, mock := newCosVectorBucketTestService(t)
			name := "tf-gc-state-1250000000"
			if _, err := service.CreateBucket(context.Background(), name, nil); err != nil {
				t.Fatal(err)
			}
			bucket, err := service.GetBucket(context.Background(), name)
			if err != nil {
				t.Fatal(err)
			}
			r := &CosVectorBucketResource{service: service}
			data := CosVectorBucketResourceModel{Region: types.StringValue(service.region)}
			data.setBucket(bucket)
			state := cosVectorTestState(t, r, data)
			mock.mu.Lock()
			mock.cleanupPending = 100
			mock.mu.Unlock()
			duration := 3 * time.Second
			if failure == "deadline" {
				duration = 30 * time.Millisecond
			} else if failure == "permission" {
				mock.setFailure("ListIndexes", 403)
			} else {
				mock.setFailure("ListIndexes", 404)
			}
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			resp := resource.DeleteResponse{State: state}
			r.Delete(ctx, resource.DeleteRequest{State: state}, &resp)
			if !resp.Diagnostics.HasError() || !resp.State.Raw.Equal(state.Raw) || mock.count() != 1 || mock.callCount("DeleteVectorBucket") != 1 {
				t.Fatal("failed cleanup lost state or retried deletion")
			}
		})
	}
}

func TestCosVectorBucketMetadataAndConfigure(t *testing.T) {
	r := NewCosVectorBucketResource().(*CosVectorBucketResource)
	var meta resource.MetadataResponse
	r.Metadata(context.Background(), resource.MetadataRequest{}, &meta)
	if meta.TypeName != "tencentcloud_cos_vector_bucket" {
		t.Fatal("wrong resource type")
	}
	var config resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{}, &config)
	if config.Diagnostics.HasError() {
		t.Fatal("nil configure should be a no-op")
	}
	var imported resource.ImportStateResponse
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "tf-vector-1250000000"}, &imported)
	if !imported.Diagnostics.HasError() {
		t.Fatal("import without configured service accepted")
	}
	client := &connectivity.TencentCloudClient{Region: "ap-shanghai", Credential: common.NewTokenCredential("shared-id", "shared-key", "shared-token")}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}}, &config)
	if config.Diagnostics.HasError() || r.Client() != client || r.service == nil || r.service.region != client.Region || r.service.secrets[2] != "shared-token" || r.service.client.BaseURL.VectorURL.String() != "https://vectors.ap-shanghai.coslake.com" {
		t.Fatal("resource did not use the configured shared provider client")
	}
}

func TestCosVectorBucketStateErrorsAndIdentity(t *testing.T) {
	ctx := context.Background()
	service, mock := newCosVectorBucketTestService(t)
	name := "tf-state-1250000000"
	if _, err := service.CreateBucket(ctx, name, &cosSDK.VectorEncryptionConfig{SseType: "AES256"}); err != nil {
		t.Fatal(err)
	}
	bucket, err := service.GetBucket(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	r := &CosVectorBucketResource{service: service}
	data := CosVectorBucketResourceModel{Region: types.StringValue(service.region)}
	data.setBucket(bucket)
	state := cosVectorTestState(t, r, data)
	for _, status := range []int{403, 409, 500} {
		mock.setFailure("GetVectorBucket", status)
		resp := resource.ReadResponse{State: state}
		r.Read(ctx, resource.ReadRequest{State: state}, &resp)
		if !resp.Diagnostics.HasError() || !resp.State.Raw.Equal(state.Raw) {
			t.Fatal("read error lost state")
		}
	}
	mock.setFailure("GetVectorBucket", 0)
	mock.mu.Lock()
	mock.invalidRead = true
	mock.mu.Unlock()
	readResp := resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, &readResp)
	if !readResp.Diagnostics.HasError() || !readResp.State.Raw.Equal(state.Raw) {
		t.Fatal("invalid response lost state")
	}
	mock.mu.Lock()
	mock.invalidRead = false
	mock.nonEmpty = true
	mock.mu.Unlock()
	var deleteResp resource.DeleteResponse
	r.Delete(ctx, resource.DeleteRequest{State: state}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() || mock.count() != 1 {
		t.Fatal("nonempty bucket was deleted")
	}
	data.Region = types.StringValue("ap-shanghai")
	wrongState := cosVectorTestState(t, r, data)
	calls := mock.callCount("DeleteVectorBucket")
	deleteResp = resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: wrongState}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() || mock.callCount("DeleteVectorBucket") != calls {
		t.Fatal("region mismatch made a remote request")
	}
	mock.drop(name)
	mock.mu.Lock()
	mock.malformedNotFound = true
	mock.mu.Unlock()
	readResp = resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, &readResp)
	if !readResp.Diagnostics.HasError() || !readResp.State.Raw.Equal(state.Raw) {
		t.Fatal("malformed 404 lost state")
	}
	mock.mu.Lock()
	mock.malformedNotFound = false
	mock.mu.Unlock()
	readResp = resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, &readResp)
	if readResp.Diagnostics.HasError() || !readResp.State.Raw.IsNull() {
		t.Fatal("missing bucket state was not removed")
	}
	deleteResp = resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: state}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatal("delete missing bucket failed")
	}
}

func TestCosVectorBucketAcceptedCreateTracksState(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	service, mock := newCosVectorBucketTestService(t)
	mock.malformedCreate = true
	mock.setFailure("GetVectorBucket", 500)
	r := &CosVectorBucketResource{service: service}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	data := CosVectorBucketResourceModel{Id: types.StringUnknown(), VectorBucketName: types.StringValue("tf-accepted-1250000000"), VectorBucketQcs: types.StringUnknown(), Region: types.StringUnknown(), CreationTime: types.Int64Unknown(), EncryptionConfiguration: types.ListUnknown(types.ObjectType{AttrTypes: cosVectorEncryptionTypes})}
	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := plan.Set(ctx, &data); diags.HasError() {
		t.Fatal(diags)
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
	if !resp.Diagnostics.HasError() || resp.State.Raw.IsNull() {
		t.Fatal("accepted create lost resource identity")
	}
	var saved CosVectorBucketResourceModel
	if diags := resp.State.Get(ctx, &saved); diags.HasError() {
		t.Fatal(diags)
	}
	if saved.Id.ValueString() != data.VectorBucketName.ValueString() || saved.Region.ValueString() != service.region || mock.callCount("CreateVectorBucket") != 1 {
		t.Fatal("missing partial state or create retried")
	}
}

func TestCosVectorBucketCreateEncryptionNotConfirmed(t *testing.T) {
	ctx := context.Background()
	service, mock := newCosVectorBucketTestService(t)
	mock.omitEncryption = true
	r := &CosVectorBucketResource{service: service}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	data := CosVectorBucketResourceModel{Id: types.StringUnknown(), VectorBucketName: types.StringValue("tf-encryption-1250000000"), VectorBucketQcs: types.StringUnknown(), Region: types.StringUnknown(), CreationTime: types.Int64Unknown(), EncryptionConfiguration: cosVectorEncryptionValue("AES256")}
	plan := tfsdk.Plan{Schema: schemaResp.Schema}
	if diags := plan.Set(ctx, &data); diags.HasError() {
		t.Fatal(diags)
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
	if !resp.Diagnostics.HasError() || resp.State.Raw.IsNull() {
		t.Fatal("missing encryption was silently accepted")
	}
	var saved CosVectorBucketResourceModel
	if diags := resp.State.Get(ctx, &saved); diags.HasError() {
		t.Fatal(diags)
	}
	if !saved.EncryptionConfiguration.Equal(data.EncryptionConfiguration) || saved.Id.ValueString() != data.VectorBucketName.ValueString() {
		t.Fatal("requested encryption or bucket identity lost")
	}
}

func TestCosVectorBucketEncryptionPlanDefault(t *testing.T) {
	resourceType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String}}
	for _, tc := range []struct {
		name            string
		create, destroy bool
		config          types.List
		wantDefault     bool
	}{
		{"newOmitted", true, false, types.ListNull(types.ObjectType{AttrTypes: cosVectorEncryptionTypes}), true},
		{"existingOmitted", false, false, types.ListNull(types.ObjectType{AttrTypes: cosVectorEncryptionTypes}), false},
		{"newUnknown", true, false, types.ListUnknown(types.ObjectType{AttrTypes: cosVectorEncryptionTypes}), false},
		{"newConfigured", true, false, cosVectorEncryptionValue("AES256"), false},
		{"destroy", false, true, types.ListNull(types.ObjectType{AttrTypes: cosVectorEncryptionTypes}), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var state, plan any
			if !tc.create {
				state = map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, "existing")}
			}
			if !tc.destroy {
				plan = map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, tftypes.UnknownValue)}
			}
			req := planmodifier.ListRequest{State: tfsdk.State{Raw: tftypes.NewValue(resourceType, state)}, Plan: tfsdk.Plan{Raw: tftypes.NewValue(resourceType, plan)}, ConfigValue: tc.config, PlanValue: types.ListUnknown(types.ObjectType{AttrTypes: cosVectorEncryptionTypes})}
			resp := planmodifier.ListResponse{PlanValue: req.PlanValue}
			cosVectorEncryptionDefault{}.PlanModifyList(context.Background(), req, &resp)
			if resp.Diagnostics.HasError() || resp.PlanValue.Equal(cosVectorEncryptionValue("AES256")) != tc.wantDefault {
				t.Fatal("default did not respect create/config/import semantics")
			}
		})
	}
}

func TestCosVectorBucketImportUnencrypted(t *testing.T) {
	ctx := context.Background()
	service, _ := newCosVectorBucketTestService(t)
	name := "tf-import-1250000000"
	qcs, err := service.CreateBucket(ctx, name, nil)
	if err != nil {
		t.Fatal(err)
	}
	r := &CosVectorBucketResource{service: service}
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	imported := resource.ImportStateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}
	r.ImportState(ctx, resource.ImportStateRequest{ID: qcs}, &imported)
	if imported.Diagnostics.HasError() {
		t.Fatal(imported.Diagnostics)
	}
	readResp := resource.ReadResponse{State: imported.State}
	r.Read(ctx, resource.ReadRequest{State: imported.State}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatal(readResp.Diagnostics)
	}
	var data CosVectorBucketResourceModel
	if diags := readResp.State.Get(ctx, &data); diags.HasError() {
		t.Fatal(diags)
	}
	if !data.EncryptionConfiguration.IsNull() || data.Id.ValueString() != name {
		t.Fatal("import applied new bucket defaults")
	}
}
