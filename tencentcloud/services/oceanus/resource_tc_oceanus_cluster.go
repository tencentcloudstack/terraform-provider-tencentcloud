package oceanus

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var (
	_ resource.ResourceWithConfigure   = &OceanusClusterResource{}
	_ resource.ResourceWithImportState = &OceanusClusterResource{}
)

// Polling tuning knobs. Tests in the same package shrink them so the mocked
// async wait converges immediately.
var (
	oceanusClusterPollDelay      = 10 * time.Second
	oceanusClusterPollMinTimeout = 15 * time.Second
	oceanusClusterPollTimeout    = 30 * time.Minute
)

// Oceanus cluster status values reported by DescribeClusters. Only 1
// (not initialized), 3 (initializing) and 2 (running) are documented by the
// cloud API; any other value is treated as an unknown transition state and
// keeps the caller waiting until the poll timeout.
const (
	oceanusClusterStatusNotInitialized = 1
	oceanusClusterStatusInitializing   = 3
	oceanusClusterStatusRunning        = 2

	oceanusClusterScaleModeUp   = "ScaleUp"
	oceanusClusterScaleModeDown = "ScaleDown"

	oceanusClusterTaskSuccess = "success"
)

// NewOceanusClusterResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewOceanusClusterResource() resource.Resource {
	return &OceanusClusterResource{}
}

// OceanusClusterResource implements resource.Resource for
// tencentcloud_oceanus_oceanus_cluster.
//
// CreateOceanusCluster is asynchronous: after the API returns the ClusterId
// the state is persisted immediately and the creation is awaited by polling
// DescribeClusters until the cluster reports status 2 (running).
// ScaleOceanusCluster (CU resize) and RenewOceanusCluster (prepaid renewal)
// are the supported update paths; every other attribute change forces
// replacement because no cloud API can modify it in place.
type OceanusClusterResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// OceanusClusterVpcDescriptionModel maps one element of the
// vpc_descriptions / slave_vpc_descriptions lists.
type OceanusClusterVpcDescriptionModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
	AppId    types.Int64  `tfsdk:"app_id"`
	OwnerUin types.String `tfsdk:"owner_uin"`
}

// OceanusClusterVpcDescriptionAttrTypes returns the object attribute types of
// a single vpc description element.
func OceanusClusterVpcDescriptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vpc_id":    attr.Type(types.StringType),
		"subnet_id": attr.Type(types.StringType),
		"app_id":    attr.Type(types.Int64Type),
		"owner_uin": attr.Type(types.StringType),
	}
}

// OceanusClusterResourceModel maps the schema attributes.
type OceanusClusterResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	ClusterName                 types.String `tfsdk:"cluster_name"`
	RegionId                    types.Int64  `tfsdk:"region_id"`
	ZoneId                      types.Int64  `tfsdk:"zone_id"`
	LoginPassword               types.String `tfsdk:"login_password"`
	VpcDescriptions             types.List   `tfsdk:"vpc_descriptions"`
	DefaultCosBucket            types.String `tfsdk:"default_cos_bucket"`
	Cu                          types.Int64  `tfsdk:"cu"`
	Remark                      types.String `tfsdk:"remark"`
	Period                      types.Int64  `tfsdk:"period"`
	InstanceChargeType          types.String `tfsdk:"instance_charge_type"`
	ClusterType                 types.String `tfsdk:"cluster_type"`
	RenewFlag                   types.String `tfsdk:"renew_flag"`
	FlinkWebUiNetworkAccessType types.String `tfsdk:"flink_web_ui_network_access_type"`
	SlaveVpcDescriptions        types.List   `tfsdk:"slave_vpc_descriptions"`
	CuMemory                    types.Int64  `tfsdk:"cu_memory"`
}

func (r *OceanusClusterResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_oceanus_oceanus_cluster"
}

func (r *OceanusClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provide an Oceanus cluster resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the cluster, in the format of cluster-xxx.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Required:    true,
				Description: "Cluster name. Supported characters: 1-50 English letters, Chinese characters, digits, hyphens (-) or underscores (_). Changing the name forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.Int64Attribute{
				Required:    true,
				Description: "Region ID, which can be queried via the DescribeRegions API with the Product parameter set to oceanus. Changing it forces recreation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.Int64Attribute{
				Required:    true,
				Description: "Availability zone ID, which can be queried via the DescribeZones API with the Product parameter set to oceanus. Changing it forces recreation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"login_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				WriteOnly:   true,
				Description: "Flink UI access password. The cluster user name is admin. This value is write-only and never persisted in state. Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_descriptions": schema.ListAttribute{
				Required:    true,
				ElementType: types.ObjectType{AttrTypes: OceanusClusterVpcDescriptionAttrTypes()},
				Description: "VPC used by the stream computing cluster to access other cloud resources in the same region via elastic network cards. Make sure the selected subnet has enough free IP addresses. Changing it forces recreation.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"default_cos_bucket": schema.StringAttribute{
				Required:    true,
				Description: "Default COS bucket used to store checkpoints, jar packages, or delivered logs. Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cu": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Number of cluster CUs. 12 CU is the minimal computing and billing unit; 1 CU contains 1 CPU and 4 GB memory. " +
					"Management node fees are waived when the CU number is greater than or equal to 48. " +
					"The cluster CU must satisfy 12 + 7*n (n >= 0). Resizing an existing cluster is applied in place via ScaleOceanusCluster.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster description. Supported characters: 1-50 English letters, Chinese characters, digits, hyphens (-) or underscores (_). Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Description: "Purchase duration in months. Only effective for prepaid (PREPAID) clusters. Create accepts [1, 48]; renewal accepts [1, 36]. Changed in place via RenewOceanusCluster.",
			},
			"instance_charge_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster billing type. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_SECOND` (pay as you go). Default value: `POSTPAID_BY_SECOND`. Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster type. Valid values: `MULTI_AZ_CLUSTER` (multi availability zone cluster). Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"renew_flag": schema.StringAttribute{
				Optional:    true,
				Description: "Auto renewal flag, only valid for prepaid clusters. Valid values: `NOTIFY_AND_MANUAL_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `DISABLE_NOTIFY_AND_MANUAL_RENEW`. Default value: `NOTIFY_AND_MANUAL_RENEW`. Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flink_web_ui_network_access_type": schema.StringAttribute{
				Optional:    true,
				Description: "Flink UI network access type. Valid values: `NetworkAccess_INTERNAL` (internal network), `NetworkAccess_EXTERNAL` (public network). Default value: `NetworkAccess_EXTERNAL`. Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"slave_vpc_descriptions": schema.ListAttribute{
				Optional:    true,
				ElementType: types.ObjectType{AttrTypes: OceanusClusterVpcDescriptionAttrTypes()},
				Description: "VPC of the slave availability zone, only effective for multi availability zone clusters. Changing it forces recreation.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"cu_memory": schema.Int64Attribute{
				Optional:    true,
				Description: "CU to memory ratio. Valid values: 0, 2, 4, 8. Changing it forces recreation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources. The returned ClusterId
// is persisted into state before the asynchronous creation is awaited, so
// `terraform destroy` can still clean up the half-created cluster.
func (r *OceanusClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OceanusClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var vpcs, slaves []OceanusClusterVpcDescriptionModel
	resp.Diagnostics.Append(plan.VpcDescriptions.ElementsAs(ctx, &vpcs, false)...)
	resp.Diagnostics.Append(plan.SlaveVpcDescriptions.ElementsAs(ctx, &slaves, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := oceanus.NewCreateOceanusClusterRequest()
	request.ClusterName = plan.ClusterName.ValueStringPointer()
	request.RegionId = oceanusClusterInt64Ptr(plan.RegionId.ValueInt64())
	request.ZoneId = oceanusClusterInt64Ptr(plan.ZoneId.ValueInt64())
	request.LoginPassword = plan.LoginPassword.ValueStringPointer()
	request.VpcDescriptions = oceanusClusterVpcItems(vpcs)
	request.DefaultCOSBucket = plan.DefaultCosBucket.ValueStringPointer()
	if !plan.Cu.IsNull() && !plan.Cu.IsUnknown() {
		request.CU = oceanusClusterInt64Ptr(plan.Cu.ValueInt64())
	}
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		request.Remark = plan.Remark.ValueStringPointer()
	}
	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		request.Period = oceanusClusterInt64Ptr(plan.Period.ValueInt64())
	}
	if !plan.InstanceChargeType.IsNull() && !plan.InstanceChargeType.IsUnknown() {
		request.InstanceChargeType = plan.InstanceChargeType.ValueStringPointer()
	}
	if !plan.ClusterType.IsNull() && !plan.ClusterType.IsUnknown() {
		request.ClusterType = plan.ClusterType.ValueStringPointer()
	}
	if !plan.RenewFlag.IsNull() && !plan.RenewFlag.IsUnknown() {
		request.RenewFlag = plan.RenewFlag.ValueStringPointer()
	}
	if !plan.FlinkWebUiNetworkAccessType.IsNull() && !plan.FlinkWebUiNetworkAccessType.IsUnknown() {
		request.FlinkWebUINetworkAccessType = plan.FlinkWebUiNetworkAccessType.ValueStringPointer()
	}
	if len(slaves) > 0 {
		request.SlaveVpcDescriptions = oceanusClusterSlaveItems(slaves)
	}
	if !plan.CuMemory.IsNull() && !plan.CuMemory.IsUnknown() {
		request.CUMemory = oceanusClusterInt64Ptr(plan.CuMemory.ValueInt64())
	}

	response, err := r.Client().UseOceanusClient().CreateOceanusClusterWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Create Oceanus cluster failed", "API CreateOceanusCluster returned error: "+err.Error())
		return
	}
	if response == nil || response.Response == nil || response.Response.ClusterId == nil || *response.Response.ClusterId == "" {
		resp.Diagnostics.AddError("Create Oceanus cluster failed", "API CreateOceanusCluster returned no ClusterId")
		return
	}

	plan.Id = types.StringValue(*response.Response.ClusterId)
	// login_password is write-only, but it is normalised here as well so the
	// same model can be reused for the final state write.
	plan.LoginPassword = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.waitForClusterReady(ctx, plan.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Wait for Oceanus cluster ready failed", err.Error())
		return
	}

	final, err := r.refreshClusterState(ctx, plan, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Oceanus cluster failed", err.Error())
		return
	}
	if final == nil {
		resp.Diagnostics.AddError(
			"Read Oceanus cluster failed",
			fmt.Sprintf("cluster %s was created but is no longer returned by DescribeClusters", plan.Id.ValueString()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, final)...)
}

// Read refreshes the resource state from DescribeClusters. The whole model is
// rebuilt once and written with a single State.Set call.
func (r *OceanusClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	final, err := r.refreshClusterState(ctx, state, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Oceanus cluster failed", err.Error())
		return
	}
	if final == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, final)...)
}

// refreshClusterState reads the cluster and returns a complete model. A nil
// model with a nil error means the cluster no longer exists.
//
// Only the attributes with a confirmed read-back path are refreshed:
//   - cluster_name        <- Cluster.Name
//   - cu                  <- Cluster.CuNum
//   - remark              <- Cluster.Remark (only when configured)
//   - default_cos_bucket  <- Cluster.DefaultCOSBucket (only when configured)
//
// Every other attribute keeps the value carried by base. Those attributes are
// either write-only (login_password) or have no single-valued read-back path
// in DescribeClusters (region_id, zone_id, period, instance_charge_type,
// cluster_type, renew_flag, flink_web_ui_network_access_type, cu_memory and the
// VPC lists), so drift on them cannot be detected by the provider.
func (r *OceanusClusterResource) refreshClusterState(
	ctx context.Context,
	base OceanusClusterResourceModel,
	clusterId string,
) (*OceanusClusterResourceModel, error) {
	cluster, err := r.describeCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, nil
	}

	result := base
	result.Id = types.StringValue(clusterId)
	if cluster.ClusterId != nil && *cluster.ClusterId != "" {
		result.Id = types.StringValue(*cluster.ClusterId)
	}
	// A write-only attribute is always absent from state.
	result.LoginPassword = types.StringNull()

	if cluster.Name != nil {
		result.ClusterName = types.StringValue(*cluster.Name)
	}
	if cluster.CuNum != nil {
		result.Cu = types.Int64Value(int64(*cluster.CuNum))
	}
	// remark and default_cos_bucket are only refreshed when the model holds a
	// configured value, so an empty cloud echo cannot turn a null attribute
	// into a non-null state value.
	if cluster.Remark != nil && !result.Remark.IsNull() && !result.Remark.IsUnknown() {
		result.Remark = types.StringValue(*cluster.Remark)
	}
	if cluster.DefaultCOSBucket != nil && !result.DefaultCosBucket.IsNull() && !result.DefaultCosBucket.IsUnknown() {
		result.DefaultCosBucket = types.StringValue(*cluster.DefaultCOSBucket)
	}

	return &result, nil
}

// Update supports two in-place operations only:
//   - CU resize via ScaleOceanusCluster (with ScaleMode derived from the
//     direction of the change);
//   - prepaid renewal via RenewOceanusCluster (period change).
//
// All other attribute changes are handled by replacement through the
// RequiresReplace plan modifiers.
func (r *OceanusClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state OceanusClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.Id.ValueString()
	if clusterId == "" {
		resp.Diagnostics.AddError("Update Oceanus cluster failed", "cluster id is missing from state")
		return
	}

	if !plan.Cu.IsNull() && !plan.Cu.IsUnknown() && (state.Cu.IsNull() || !plan.Cu.Equal(state.Cu)) {
		newCU := plan.Cu.ValueInt64()
		scaleMode := oceanusClusterScaleModeUp
		if !state.Cu.IsNull() && !state.Cu.IsUnknown() && newCU < state.Cu.ValueInt64() {
			scaleMode = oceanusClusterScaleModeDown
		}
		request := oceanus.NewScaleOceanusClusterRequest()
		request.ClusterId = oceanusClusterStrPtr(clusterId)
		request.NewCU = oceanusClusterInt64Ptr(newCU)
		request.ScaleMode = oceanusClusterStrPtr(scaleMode)
		response, err := r.Client().UseOceanusClient().ScaleOceanusClusterWithContext(ctx, request)
		if err != nil {
			resp.Diagnostics.AddError("Scale Oceanus cluster failed", "API ScaleOceanusCluster returned error: "+err.Error())
			return
		}
		if response == nil || response.Response == nil || response.Response.TaskExecResult == nil ||
			*response.Response.TaskExecResult != oceanusClusterTaskSuccess {
			resp.Diagnostics.AddError(
				"Scale Oceanus cluster failed",
				"API ScaleOceanusCluster did not report a successful task execution result",
			)
			return
		}
		if err := r.waitForClusterReady(ctx, clusterId); err != nil {
			resp.Diagnostics.AddError("Wait for Oceanus cluster ready failed", err.Error())
			return
		}
	}

	if !plan.Period.IsNull() && !plan.Period.IsUnknown() && (state.Period.IsNull() || !plan.Period.Equal(state.Period)) {
		request := oceanus.NewRenewOceanusClusterRequest()
		request.ClusterId = oceanusClusterStrPtr(clusterId)
		request.Period = oceanusClusterInt64Ptr(plan.Period.ValueInt64())
		if _, err := r.Client().UseOceanusClient().RenewOceanusClusterWithContext(ctx, request); err != nil {
			resp.Diagnostics.AddError("Renew Oceanus cluster failed", "API RenewOceanusCluster returned error: "+err.Error())
			return
		}
		if err := r.waitForClusterReady(ctx, clusterId); err != nil {
			resp.Diagnostics.AddError("Wait for Oceanus cluster ready failed", err.Error())
			return
		}
	}

	// login_password is write-only, but it is normalised here as well so the
	// same model can be reused for the final state write.
	plan.LoginPassword = types.StringNull()

	final, err := r.refreshClusterState(ctx, plan, clusterId)
	if err != nil {
		resp.Diagnostics.AddError("Read Oceanus cluster failed", err.Error())
		return
	}
	if final == nil {
		resp.Diagnostics.AddError(
			"Read Oceanus cluster failed",
			fmt.Sprintf("cluster %s is no longer returned by DescribeClusters", clusterId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, final)...)
}

// Delete destroys the cluster via DeleteOceanusCluster and waits until
// DescribeClusters no longer returns the cluster.
func (r *OceanusClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.Id.ValueString()
	if clusterId == "" {
		resp.Diagnostics.AddError("Delete Oceanus cluster failed", "cluster id is missing from state")
		return
	}

	request := oceanus.NewDeleteOceanusClusterRequest()
	request.ClusterId = oceanusClusterStrPtr(clusterId)
	response, err := r.Client().UseOceanusClient().DeleteOceanusClusterWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Delete Oceanus cluster failed", "API DeleteOceanusCluster returned error: "+err.Error())
		return
	}
	if response == nil || response.Response == nil || response.Response.TaskExecResult == nil ||
		*response.Response.TaskExecResult != oceanusClusterTaskSuccess {
		resp.Diagnostics.AddError(
			"Delete Oceanus cluster failed",
			"API DeleteOceanusCluster did not report a successful task execution result",
		)
		return
	}

	if err := r.waitForClusterDeleted(ctx, clusterId); err != nil {
		resp.Diagnostics.AddError("Wait for Oceanus cluster deletion failed", err.Error())
	}
}

// describeCluster queries the cluster by its ID. A nil result with no error
// means the cluster no longer exists.
func (r *OceanusClusterResource) describeCluster(ctx context.Context, clusterId string) (*oceanus.Cluster, error) {
	request := oceanus.NewDescribeClustersRequest()
	request.ClusterIds = []*string{oceanusClusterStrPtr(clusterId)}
	response, err := r.Client().UseOceanusClient().DescribeClustersWithContext(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("API DescribeClusters returned error: %w", err)
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("API DescribeClusters returned nil response")
	}
	if len(response.Response.ClusterSet) == 0 {
		return nil, nil
	}
	for _, cluster := range response.Response.ClusterSet {
		if cluster == nil {
			continue
		}
		if cluster.ClusterId == nil || *cluster.ClusterId == clusterId {
			return cluster, nil
		}
	}
	return nil, nil
}

// waitForClusterReady polls DescribeClusters until the cluster reports
// status 2 (running).
func (r *OceanusClusterResource) waitForClusterReady(ctx context.Context, clusterId string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"INITIALIZING"},
		Target:  []string{"RUNNING"},
		Refresh: func() (interface{}, string, error) {
			cluster, err := r.describeCluster(ctx, clusterId)
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return nil, "MISSING", fmt.Errorf("cluster %s not found", clusterId)
			}
			if cluster.Status == nil {
				return cluster, "INITIALIZING", nil
			}
			switch *cluster.Status {
			case oceanusClusterStatusRunning:
				return cluster, "RUNNING", nil
			case oceanusClusterStatusNotInitialized, oceanusClusterStatusInitializing:
				return cluster, "INITIALIZING", nil
			default:
				// Undocumented transition state: keep waiting instead of
				// guessing an extra enum value.
				return cluster, "INITIALIZING", nil
			}
		},
		Timeout:    oceanusClusterPollTimeout,
		Delay:      oceanusClusterPollDelay,
		MinTimeout: oceanusClusterPollMinTimeout,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// waitForClusterDeleted polls DescribeClusters until the cluster disappears.
func (r *OceanusClusterResource) waitForClusterDeleted(ctx context.Context, clusterId string) error {
	stateConf := &sdkresource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			cluster, err := r.describeCluster(ctx, clusterId)
			if err != nil {
				return nil, "", err
			}
			if cluster == nil {
				return nil, "DELETED", nil
			}
			return cluster, "DELETING", nil
		},
		Timeout:    oceanusClusterPollTimeout,
		Delay:      oceanusClusterPollDelay,
		MinTimeout: oceanusClusterPollMinTimeout,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func oceanusClusterInt64Ptr(v int64) *int64 {
	return &v
}

func oceanusClusterStrPtr(v string) *string {
	return &v
}

func oceanusClusterVpcItems(vpcs []OceanusClusterVpcDescriptionModel) []*oceanus.VPCDescription {
	items := make([]*oceanus.VPCDescription, 0, len(vpcs))
	for _, v := range vpcs {
		item := &oceanus.VPCDescription{
			VpcId:    v.VpcId.ValueStringPointer(),
			SubnetId: v.SubnetId.ValueStringPointer(),
		}
		if !v.AppId.IsNull() && !v.AppId.IsUnknown() {
			item.AppId = oceanusClusterInt64Ptr(v.AppId.ValueInt64())
		}
		if !v.OwnerUin.IsNull() && !v.OwnerUin.IsUnknown() {
			item.OwnerUin = v.OwnerUin.ValueStringPointer()
		}
		items = append(items, item)
	}
	return items
}

func oceanusClusterSlaveItems(slaves []OceanusClusterVpcDescriptionModel) []*oceanus.SlaveVpcDescriptions {
	items := make([]*oceanus.SlaveVpcDescriptions, 0, len(slaves))
	for _, v := range slaves {
		item := &oceanus.SlaveVpcDescriptions{
			VpcId:    v.VpcId.ValueStringPointer(),
			SubnetId: v.SubnetId.ValueStringPointer(),
		}
		if !v.AppId.IsNull() && !v.AppId.IsUnknown() {
			item.AppId = oceanusClusterInt64Ptr(v.AppId.ValueInt64())
		}
		if !v.OwnerUin.IsNull() && !v.OwnerUin.IsUnknown() {
			item.OwnerUin = v.OwnerUin.ValueStringPointer()
		}
		items = append(items, item)
	}
	return items
}
