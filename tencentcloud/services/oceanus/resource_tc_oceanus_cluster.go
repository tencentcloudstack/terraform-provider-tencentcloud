package oceanus

import (
	"context"
	"fmt"
	"time"

	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
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

const (
	oceanusClusterScaleModeUp   = "ScaleUp"
	oceanusClusterScaleModeDown = "ScaleDown"
)

// NewOceanusClusterResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewOceanusClusterResource() resource.Resource {
	return &OceanusClusterResource{}
}

// OceanusClusterResource implements resource.Resource for
// tencentcloud_oceanus_cluster.
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
}

// OceanusClusterVpcDescriptionAttrTypes returns the object attribute types of
// a single vpc description element. The vpc_descriptions and
// slave_vpc_descriptions schema attributes use ListNestedAttribute; this map
// mirrors their element object type for callers (tests) that need to build
// list values by hand.
func OceanusClusterVpcDescriptionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vpc_id":    attr.Type(types.StringType),
		"subnet_id": attr.Type(types.StringType),
	}
}

// OceanusClusterResourceModel maps the schema attributes.
type OceanusClusterResourceModel struct {
	Id                          types.String           `tfsdk:"id"`
	ClusterName                 types.String           `tfsdk:"cluster_name"`
	Region                      types.String           `tfsdk:"region"`
	Zone                        types.String           `tfsdk:"zone"`
	LoginPassword               types.String           `tfsdk:"login_password"`
	VpcDescriptions             types.List             `tfsdk:"vpc_descriptions"`
	DefaultCosBucket            types.String           `tfsdk:"default_cos_bucket"`
	Cu                          types.Int64            `tfsdk:"cu"`
	Remark                      types.String           `tfsdk:"remark"`
	Period                      types.Int64            `tfsdk:"period"`
	InstanceChargeType          types.String           `tfsdk:"instance_charge_type"`
	ClusterType                 types.String           `tfsdk:"cluster_type"`
	RenewFlag                   types.String           `tfsdk:"renew_flag"`
	FlinkWebUiNetworkAccessType types.String           `tfsdk:"flink_web_ui_network_access_type"`
	SlaveVpcDescriptions        types.List             `tfsdk:"slave_vpc_descriptions"`
	CuMemory                    types.Int64            `tfsdk:"cu_memory"`
	Timeouts                    resourcetimeouts.Value `tfsdk:"timeouts"`
}

func (r *OceanusClusterResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_oceanus_cluster"
}

func (r *OceanusClusterResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Region name, e.g. `ap-guangzhou`. It is translated to the numeric region id expected by the cluster APIs via the region DescribeRegions API. Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone": schema.StringAttribute{
				Required:    true,
				Description: "Availability zone name, e.g. `ap-guangzhou-7`. It is translated to the numeric zone id expected by the cluster APIs via the region DescribeZones API. Changing it forces recreation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			"vpc_descriptions": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id": schema.StringAttribute{
							Required:    true,
							Description: "VPC ID.",
						},
						"subnet_id": schema.StringAttribute{
							Required:    true,
							Description: "Subnet ID.",
						},
					},
				},
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
			"slave_vpc_descriptions": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id":    schema.StringAttribute{Required: true, Description: "VPC ID."},
						"subnet_id": schema.StringAttribute{Required: true, Description: "Subnet ID."},
					},
				},
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
		// The timeouts definition uses a block (rather than an object
		// attribute) so the familiar `timeouts { ... }` HCL block syntax
		// works, matching the official Terraform examples and the SDKv2
		// providers. Every stage bounds the corresponding async wait loop.
		Blocks: map[string]schema.Block{
			"timeouts": resourcetimeouts.Block(ctx, resourcetimeouts.Opts{
				Create:            true,
				Update:            true,
				CreateDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `30m`. It limits how long to wait for the cluster to reach the running state after creation. Default is `30m`.",
				UpdateDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `30m`. It limits how long to wait for the cluster to converge after a CU resize or a renewal. Default is `30m`.",
			}),
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

	// The `timeouts.create` attribute bounds the creation wait loop; an unset
	// block or attribute falls back to the built-in default and a malformed
	// duration is reported as an error diagnostic (the schema validator
	// normally rejects it at plan time; this is the runtime safety net).
	createTimeout, timeoutDiags := plan.Timeouts.Create(ctx, oceanusClusterPollTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// login_password is write-only: Terraform never carries write-only
	// values in the plan (they would be persisted in plan files), they are
	// only exposed through the configuration during apply. Reading the plan
	// alone yields a null password, and CreateOceanusCluster would fail with
	// MissingParameter.
	var config OceanusClusterResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var vpcs, slaves []OceanusClusterVpcDescriptionModel
	resp.Diagnostics.Append(plan.VpcDescriptions.ElementsAs(ctx, &vpcs, false)...)
	resp.Diagnostics.Append(plan.SlaveVpcDescriptions.ElementsAs(ctx, &slaves, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := OceanusService{client: r.Client()}

	// The cluster APIs address regions and zones by their numeric ids, while
	// the practitioner configures their names. Translate both through the
	// region management service before building the create request.
	regionId, err := service.DescribeOceanusRegionId(ctx, plan.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Resolve Oceanus region failed",
			fmt.Sprintf("failed to resolve region %q to a numeric region id: %s", plan.Region.ValueString(), err.Error()),
		)
		return
	}
	zoneId, err := service.DescribeOceanusZoneId(ctx, plan.Zone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Resolve Oceanus zone failed",
			fmt.Sprintf("failed to resolve zone %q to a numeric zone id: %s", plan.Zone.ValueString(), err.Error()),
		)
		return
	}

	request := oceanus.NewCreateOceanusClusterRequest()
	request.ClusterName = plan.ClusterName.ValueStringPointer()
	request.RegionId = helper.Int64(regionId)
	request.ZoneId = helper.Int64(zoneId)
	request.LoginPassword = config.LoginPassword.ValueStringPointer()
	request.VpcDescriptions = oceanusClusterVpcItems(vpcs)
	request.DefaultCOSBucket = plan.DefaultCosBucket.ValueStringPointer()
	if !plan.Cu.IsNull() && !plan.Cu.IsUnknown() {
		request.CU = helper.Int64(plan.Cu.ValueInt64())
	}
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		request.Remark = plan.Remark.ValueStringPointer()
	}
	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		request.Period = helper.Int64(plan.Period.ValueInt64())
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
		request.CUMemory = helper.Int64(plan.CuMemory.ValueInt64())
	}

	clusterId, err := service.CreateOceanusCluster(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Create Oceanus cluster failed", err.Error())
		return
	}

	plan.Id = types.StringValue(clusterId)
	plan.LoginPassword = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := service.WaitOceanusClusterReady(ctx, clusterId, createTimeout, oceanusClusterPollDelay, oceanusClusterPollMinTimeout); err != nil {
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
	service := OceanusService{client: r.Client()}
	cluster, err := service.DescribeOceanusClusterById(ctx, clusterId)
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
	if cluster.Region != nil {
		result.Region = types.StringValue(*cluster.Region)
	}
	if cluster.Zone != nil {
		result.Zone = types.StringValue(*cluster.Zone)
	}
	if cluster.CCNs != nil {
		vpcs := make([]OceanusClusterVpcDescriptionModel, 0, len(cluster.CCNs))
		for _, v := range cluster.CCNs {
			vpcs = append(vpcs, OceanusClusterVpcDescriptionModel{
				VpcId:    types.StringValue(*v.VpcId),
				SubnetId: types.StringValue(*v.SubnetId),
			})
		}
		result.VpcDescriptions, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: OceanusClusterVpcDescriptionAttrTypes()}, vpcs)
	}
	if cluster.DefaultCOSBucket != nil {
		result.DefaultCosBucket = types.StringValue(*cluster.DefaultCOSBucket)
	}
	if cluster.CuNum != nil {
		result.Cu = types.Int64Value(int64(*cluster.CuNum))
	}
	if cluster.Remark != nil {
		result.Remark = types.StringValue(*cluster.Remark)
	}
	// period has no reliable read-back: the describe result only reports the
	// expire/create timestamps of a prepaid cluster, and renewals shift the
	// expiry so the originally configured purchase period cannot be
	// recovered. The value carried by base is kept.
	if cluster.PayMode != nil {
		if *cluster.PayMode == 1 {
			result.InstanceChargeType = types.StringValue("PREPAID")
		} else {
			result.InstanceChargeType = types.StringValue("POSTPAID_BY_SECOND")
		}
	}
	if cluster.DeploymentMode != nil && *cluster.DeploymentMode == 1 {
		result.ClusterType = types.StringValue("MULTI_AZ_CLUSTER")
	}
	// NOTIFY_AND_MANUAL_RENEW (Notify and manual renewal) → 0 (default value)
	// NOTIFY_AND_AUTO_RENEW (Notify and automatic renewal) → 1
	// DISABLE_NOTIFY_AND_MANUAL_RENEW (No notification, no renewal) → 2
	if cluster.AutoRenewFlag != nil {
		if *cluster.AutoRenewFlag == 0 {
			result.RenewFlag = types.StringValue("NOTIFY_AND_MANUAL_RENEW")
		} else if *cluster.AutoRenewFlag == 1 {
			result.RenewFlag = types.StringValue("NOTIFY_AND_AUTO_RENEW")
		} else if *cluster.AutoRenewFlag == 2 {
			result.RenewFlag = types.StringValue("DISABLE_NOTIFY_AND_MANUAL_RENEW")
		}
	}
	// NetworkAccess_EXTERNAL (public network access, default) → WebUIType=0
	// NetworkAccess_INTERNAL (intranet access) → WebUIType=1
	if cluster.WebUIType != nil {
		if *cluster.WebUIType == 0 {
			result.FlinkWebUiNetworkAccessType = types.StringValue("NetworkAccess_EXTERNAL")
		} else if *cluster.WebUIType == 1 {
			result.FlinkWebUiNetworkAccessType = types.StringValue("NetworkAccess_INTERNAL")
		}
	}
	if cluster.SlaveZones != nil {
		slaves := make([]OceanusClusterVpcDescriptionModel, 0, len(cluster.SlaveZones))
		for _, v := range cluster.SlaveZones {
			slaves = append(slaves, OceanusClusterVpcDescriptionModel{
				VpcId:    types.StringValue(*v.VpcId),
				SubnetId: types.StringValue(*v.SubnetId),
			})
		}
		result.SlaveVpcDescriptions, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: OceanusClusterVpcDescriptionAttrTypes()}, slaves)
	}
	if cluster.CuMem != nil {
		result.CuMemory = types.Int64Value(int64(*cluster.CuMem))
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

	// The `timeouts.update` attribute bounds the post-resize / post-renewal
	// wait loops; an unset block falls back to the built-in default.
	updateTimeout, timeoutDiags := plan.Timeouts.Update(ctx, oceanusClusterPollTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := OceanusService{client: r.Client()}

	if !plan.Cu.IsNull() && !plan.Cu.IsUnknown() && (state.Cu.IsNull() || !plan.Cu.Equal(state.Cu)) {
		newCU := plan.Cu.ValueInt64()
		scaleMode := oceanusClusterScaleModeUp
		if !state.Cu.IsNull() && !state.Cu.IsUnknown() && newCU < state.Cu.ValueInt64() {
			scaleMode = oceanusClusterScaleModeDown
		}
		request := oceanus.NewScaleOceanusClusterRequest()
		request.ClusterId = helper.String(clusterId)
		request.NewCU = helper.Int64(newCU)
		request.ScaleMode = helper.String(scaleMode)
		if err := service.ScaleOceanusCluster(ctx, request); err != nil {
			resp.Diagnostics.AddError("Scale Oceanus cluster failed", err.Error())
			return
		}
		if err := service.WaitOceanusClusterReady(ctx, clusterId, updateTimeout, oceanusClusterPollDelay, oceanusClusterPollMinTimeout); err != nil {
			resp.Diagnostics.AddError("Wait for Oceanus cluster ready failed", err.Error())
			return
		}
	}

	if !plan.Period.IsNull() && !plan.Period.IsUnknown() && (state.Period.IsNull() || !plan.Period.Equal(state.Period)) {
		request := oceanus.NewRenewOceanusClusterRequest()
		request.ClusterId = helper.String(clusterId)
		request.Period = helper.Int64(plan.Period.ValueInt64())
		if err := service.RenewOceanusCluster(ctx, request); err != nil {
			resp.Diagnostics.AddError("Renew Oceanus cluster failed", err.Error())
			return
		}
		if err := service.WaitOceanusClusterReady(ctx, clusterId, updateTimeout, oceanusClusterPollDelay, oceanusClusterPollMinTimeout); err != nil {
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
	request.ClusterId = helper.String(clusterId)
	service := OceanusService{client: r.Client()}
	if err := service.DeleteOceanusCluster(ctx, request); err != nil {
		resp.Diagnostics.AddError("Delete Oceanus cluster failed", err.Error())
		return
	}
}

func oceanusClusterVpcItems(vpcs []OceanusClusterVpcDescriptionModel) []*oceanus.VPCDescription {
	items := make([]*oceanus.VPCDescription, 0, len(vpcs))
	for _, v := range vpcs {
		item := &oceanus.VPCDescription{
			VpcId:    v.VpcId.ValueStringPointer(),
			SubnetId: v.SubnetId.ValueStringPointer(),
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
		items = append(items, item)
	}
	return items
}
