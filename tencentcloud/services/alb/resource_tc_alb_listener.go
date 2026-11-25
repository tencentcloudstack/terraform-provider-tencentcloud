package alb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions (repo convention: only assert the
// extension interfaces; the base interface is checked implicitly by the
// factory return value).
var (
	_ resource.ResourceWithConfigure   = &AlbListenerResource{}
	_ resource.ResourceWithImportState = &AlbListenerResource{}
)

// NewAlbListenerResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAlbListenerResource() resource.Resource {
	return &AlbListenerResource{}
}

// AlbListenerResource implements resource.Resource for
// tencentcloud_alb_listener.
type AlbListenerResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AlbListenerDefaultActionModel maps one element of the `default_actions`
// attribute.
type AlbListenerDefaultActionModel struct {
	Type              types.String                       `tfsdk:"type"`
	TargetGroupConfig *AlbListenerTargetGroupConfigModel `tfsdk:"target_group_config"`
}

// AlbListenerTargetGroupConfigModel maps the `target_group_config` attribute.
type AlbListenerTargetGroupConfigModel struct {
	TargetGroups             []AlbListenerTargetGroupModel      `tfsdk:"target_groups"`
	TargetGroupStickySession *AlbListenerTargetGroupStickyModel `tfsdk:"target_group_sticky_session"`
}

// AlbListenerTargetGroupModel maps one element of the `target_groups`
// attribute.
type AlbListenerTargetGroupModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

// AlbListenerTargetGroupStickyModel maps the `target_group_sticky_session`
// attribute.
type AlbListenerTargetGroupStickyModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

// AlbListenerTagModel maps one element of the `tags` attribute.
type AlbListenerTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// AlbListenerXForwardedForConfigModel maps the `x_forwarded_for_config`
// attribute.
type AlbListenerXForwardedForConfigModel struct {
	XForwardedForAlbIdEnabled         types.Bool   `tfsdk:"x_forwarded_for_alb_id_enabled"`
	XForwardedForClientSrcPortEnabled types.Bool   `tfsdk:"x_forwarded_for_client_src_port_enabled"`
	XForwardedForHostEnabled          types.Bool   `tfsdk:"x_forwarded_for_host_enabled"`
	XForwardedForMode                 types.String `tfsdk:"x_forwarded_for_mode"`
	XForwardedForPortEnabled          types.Bool   `tfsdk:"x_forwarded_for_port_enabled"`
	XForwardedForProtoEnabled         types.Bool   `tfsdk:"x_forwarded_for_proto_enabled"`
	XTencentClientIDNEnabled          types.Bool   `tfsdk:"x_tencent_client_idn_enabled"`
	XTencentClientSDNEnabled          types.Bool   `tfsdk:"x_tencent_client_sdn_enabled"`
	XTencentClientSerialEnabled       types.Bool   `tfsdk:"x_tencent_client_serial_enabled"`
	XTencentClientVerifyEnabled       types.Bool   `tfsdk:"x_tencent_client_verify_enabled"`
}

// AlbListenerResourceModel maps the schema attributes.
//
// The describe API (DescribeListeners) does not echo back
// ca_certificate_ids / certificate_ids / client_token / default_actions, so
// Read keeps whatever values are already stored in state for them.
type AlbListenerResourceModel struct {
	Id                  types.String                         `tfsdk:"id"`
	LoadBalancerId      types.String                         `tfsdk:"load_balancer_id"`
	ListenerPort        types.Int64                          `tfsdk:"listener_port"`
	ListenerProtocol    types.String                         `tfsdk:"listener_protocol"`
	DefaultActions      []AlbListenerDefaultActionModel      `tfsdk:"default_actions"`
	CaCertificateIds    types.List                           `tfsdk:"ca_certificate_ids"`
	CaEnabled           types.Bool                           `tfsdk:"ca_enabled"`
	CertificateIds      types.List                           `tfsdk:"certificate_ids"`
	ClientToken         types.String                         `tfsdk:"client_token"`
	GzipEnabled         types.Bool                           `tfsdk:"gzip_enabled"`
	Http2Enabled        types.Bool                           `tfsdk:"http2_enabled"`
	IdleTimeout         types.Int64                          `tfsdk:"idle_timeout"`
	ListenerName        types.String                         `tfsdk:"listener_name"`
	RequestTimeout      types.Int64                          `tfsdk:"request_timeout"`
	SecurityPolicyId    types.String                         `tfsdk:"security_policy_id"`
	Tags                []AlbListenerTagModel                `tfsdk:"tags"`
	XForwardedForConfig *AlbListenerXForwardedForConfigModel `tfsdk:"x_forwarded_for_config"`
}

func (r *AlbListenerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listener"
}

func (r *AlbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage an ALB listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, the listener ID, in the format of `lst-` followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer instance ID, in the format of `alb-` followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_port": schema.Int64Attribute{
				Required:    true,
				Description: "Port used by the load balancer instance frontend. Value range: 1~65535.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"listener_protocol": schema.StringAttribute{
				Required:    true,
				Description: "Listener protocol. Valid values: `HTTP`, `HTTPS`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"default_actions": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of default forwarding rule actions. Currently, only 1 default forwarding rule action can be added to a listener.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Forwarding action type. When creating a listener, the default forwarding action type only supports forwarding to a target group.",
						},
						"target_group_config": schema.SingleNestedAttribute{
							Required:    true,
							Description: "Target group configuration of the forwarding action. When creating a listener, the target group configuration in the forwarding action only supports a single target group.",
							Attributes: map[string]schema.Attribute{
								"target_groups": schema.ListNestedAttribute{
									Required:    true,
									Description: "List of target groups.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"target_group_id": schema.StringAttribute{
												Required:    true,
												Description: "Target group ID, in the format of `lbtg-` followed by 8 alphanumeric characters.",
											},
											"weight": schema.Int64Attribute{
												Optional:    true,
												Description: "Weight. Value range: [0, 100]. Default value: 10.",
											},
										},
									},
								},
								"target_group_sticky_session": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "Sticky session between target groups.",
									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Required:    true,
											Description: "Whether to enable sticky session. Disabled by default.",
										},
										"timeout": schema.Int64Attribute{
											Optional:    true,
											Description: "Timeout in seconds. Value range: 1~86400. Default value: 1000.",
										},
									},
								},
							},
						},
					},
				},
			},
			"ca_certificate_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of CA certificate IDs configured for the listener. Currently, only 1 CA certificate can be added to a listener. Required when `ca_enabled` is set to `true`.",
			},
			"ca_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to enable mutual authentication. Valid values: `true`: enabled; `false` (default): disabled. Read back from the cloud as `CaEnable`.",
			},
			"certificate_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of server certificate IDs.",
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "Client token used to guarantee request idempotency. Generate a parameter value from your client ensuring the parameter value is unique across different requests. Only ASCII characters are supported.",
			},
			"gzip_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to enable Gzip compression. Valid values: `true` (default): yes; `false`: no.",
			},
			"http2_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to enable the HTTP/2 feature. `false` by default for the HTTP protocol and `true` by default for the HTTPS protocol. Only the HTTPS protocol supports this parameter. Read back from the cloud as `Http2Enable`.",
			},
			"idle_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Connection idle timeout in seconds. Value range: 1~600. Default value: 15. If no access request arrives within the timeout period, the load balancer disconnects the current connection and creates a new one when the next request arrives.",
			},
			"listener_name": schema.StringAttribute{
				Optional:    true,
				Description: "Custom listener name. 1~255 characters long, must be harmless characters, and can contain Chinese characters, letters, digits, hyphen (-), forward slash (/), dot (.) and underscore (_).",
			},
			"request_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Connection request timeout in seconds. Value range: 1~600. Default value: 60. If the backend server does not return a response within the timeout period, the load balancer gives up waiting and returns an HTTP 504 error to the client.",
			},
			"security_policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "Security policy ID, in the format of `tls-` followed by 8 alphanumeric characters. Read back from the cloud as `TlsSecurityPolicyId`.",
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of tags. Up to 20 supported.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag_key": schema.StringAttribute{
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": schema.StringAttribute{
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"x_forwarded_for_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "X-Forwarded-For configuration.",
				Attributes: map[string]schema.Attribute{
					"x_forwarded_for_alb_id_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the load balancer instance ID through the ALB-ID header field. `true`: yes; `false`: no.",
					},
					"x_forwarded_for_client_src_port_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the port of the client accessing the load balancer instance through the X-Forwarded-Client-srcport header field. `true`: yes; `false`: no.",
					},
					"x_forwarded_for_host_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the domain name of the client accessing the load balancer instance through the X-Forwarded-Host header field. `true`: yes; `false`: no.",
					},
					"x_forwarded_for_mode": schema.StringAttribute{
						Optional:    true,
						Description: "Specifies how to process the X-Forwarded-For (XFF) HTTP header field. Valid values: `append`: append mode (default), which appends the client's real IP address to the end of the X-Forwarded-For header, retaining the original XFF chain information; `remove`: remove mode, which removes the X-Forwarded-For header and does not pass it to the backend server; `passthrough`: pass-through mode, which keeps the X-Forwarded-For header unchanged and passes it to the backend server without any modification.",
					},
					"x_forwarded_for_port_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the listener port of the load balancer instance through the X-Forwarded-Port header field. `true`: yes; `false`: no.",
					},
					"x_forwarded_for_proto_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the listener protocol of the load balancer instance through the X-Forwarded-Proto header field. `true`: yes; `false`: no.",
					},
					"x_tencent_client_idn_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the issuer of the client certificate, $ssl_client_i_dn, through the X-Tencent-Client-IDN header. `true`: yes; `false`: no.",
					},
					"x_tencent_client_sdn_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the subject of the client certificate, $ssl_client_s_dn, through the X-Tencent-Client-SDN header. `true`: yes; `false`: no.",
					},
					"x_tencent_client_serial_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the serial number of the client certificate, $ssl_client_serial, through the X-Tencent-Client-Serial header. `true`: yes; `false`: no.",
					},
					"x_tencent_client_verify_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the verification result of the client certificate, $ssl_client_verify, through the X-Tencent-Client-Verify header. `true`: yes; `false`: no.",
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AlbListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlbListenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createListenerRequest := albv20251030.NewCreateListenerRequest()
	createListenerRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	createListenerRequest.ListenerPort = fwHelper.Uint64PointerFromValue(plan.ListenerPort)
	createListenerRequest.ListenerProtocol = plan.ListenerProtocol.ValueStringPointer()
	createListenerRequest.DefaultActions = expandAlbListenerDefaultActions(plan.DefaultActions)
	if !plan.CaCertificateIds.IsNull() && !plan.CaCertificateIds.IsUnknown() {
		createListenerRequest.CaCertificateIds = fwHelper.PointersFromList[string](ctx, plan.CaCertificateIds, &resp.Diagnostics)
	}
	if !plan.CaEnabled.IsNull() && !plan.CaEnabled.IsUnknown() {
		createListenerRequest.CaEnabled = plan.CaEnabled.ValueBoolPointer()
	}
	if !plan.CertificateIds.IsNull() && !plan.CertificateIds.IsUnknown() {
		createListenerRequest.CertificateIds = fwHelper.PointersFromList[string](ctx, plan.CertificateIds, &resp.Diagnostics)
	}
	if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() {
		createListenerRequest.ClientToken = plan.ClientToken.ValueStringPointer()
	}
	if !plan.GzipEnabled.IsNull() && !plan.GzipEnabled.IsUnknown() {
		createListenerRequest.GzipEnabled = plan.GzipEnabled.ValueBoolPointer()
	}
	if !plan.Http2Enabled.IsNull() && !plan.Http2Enabled.IsUnknown() {
		createListenerRequest.Http2Enabled = plan.Http2Enabled.ValueBoolPointer()
	}
	if !plan.IdleTimeout.IsNull() && !plan.IdleTimeout.IsUnknown() {
		createListenerRequest.IdleTimeout = fwHelper.Uint64PointerFromValue(plan.IdleTimeout)
	}
	if !plan.ListenerName.IsNull() && !plan.ListenerName.IsUnknown() {
		createListenerRequest.ListenerName = plan.ListenerName.ValueStringPointer()
	}
	if !plan.RequestTimeout.IsNull() && !plan.RequestTimeout.IsUnknown() {
		createListenerRequest.RequestTimeout = fwHelper.Uint64PointerFromValue(plan.RequestTimeout)
	}
	if !plan.SecurityPolicyId.IsNull() && !plan.SecurityPolicyId.IsUnknown() {
		createListenerRequest.SecurityPolicyId = plan.SecurityPolicyId.ValueStringPointer()
	}
	createListenerRequest.Tags = expandAlbListenerTags(plan.Tags)
	if plan.XForwardedForConfig != nil {
		createListenerRequest.XForwardedForConfig = expandAlbListenerXForwardedForConfig(plan.XForwardedForConfig)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAlbService(r.Client())

	listenerId, err := service.CreateListener(ctx, createListenerRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			fmt.Sprintf("CreateListener failed: %v", err),
		)

		return
	}

	// Never write an empty id into state.
	if listenerId == "" {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			"CreateListener returned an empty listener id.",
		)

		return
	}

	plan.Id = types.StringValue(listenerId)

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	describeListenersRequest.ListenerIds = []*string{fwHelper.String(listenerId)}
	listeners, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener after create", err.Error())
		return
	}

	item := findAlbListener(listeners, listenerId)
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB listener not readable after create",
			fmt.Sprintf("listener %s was created but cannot be read back", plan.Id.ValueString()),
		)

		return
	}

	refreshAlbListenerModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AlbListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlbListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId := state.Id.ValueString()

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = state.LoadBalancerId.ValueStringPointer()
	describeListenersRequest.ListenerIds = []*string{fwHelper.String(listenerId)}

	service := NewAlbService(r.Client())

	listeners, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener", err.Error())
		return
	}

	item := findAlbListener(listeners, listenerId)

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ALB listener %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	refreshAlbListenerModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *AlbListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AlbListenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AlbListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId := state.Id.ValueString()

	service := NewAlbService(r.Client())

	// ModifyListenerAttributes covers every mutable attribute except
	// listener_port / listener_protocol / load_balancer_id / tags, which are
	// marked RequiresReplace in the schema. The write API is only called when
	// one of the attributes it covers actually changed.
	if albListenerAttributesChanged(&plan, &state) {
		// TODO(human): the clear semantics of the optional fields of
		// ModifyListenerAttributes (ca_certificate_ids, certificate_ids,
		// default_actions, x_forwarded_for_config, ...) is not documented by
		// the cloud API: removing an attribute from the configuration sends
		// an omitted field, which the cloud may treat as "keep the previous
		// value", producing a permanent diff. A human should confirm whether
		// the cloud supports explicit clearing and how.
		modifyListenerRequest := albv20251030.NewModifyListenerAttributesRequest()
		modifyListenerRequest.ListenerId = fwHelper.String(listenerId)
		modifyListenerRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
		if !plan.CaCertificateIds.IsNull() && !plan.CaCertificateIds.IsUnknown() {
			modifyListenerRequest.CaCertificateIds = fwHelper.PointersFromList[string](ctx, plan.CaCertificateIds, &resp.Diagnostics)
		}
		if !plan.CaEnabled.IsNull() && !plan.CaEnabled.IsUnknown() {
			modifyListenerRequest.CaEnabled = plan.CaEnabled.ValueBoolPointer()
		}
		if !plan.CertificateIds.IsNull() && !plan.CertificateIds.IsUnknown() {
			modifyListenerRequest.CertificateIds = fwHelper.PointersFromList[string](ctx, plan.CertificateIds, &resp.Diagnostics)
		}
		if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() {
			modifyListenerRequest.ClientToken = plan.ClientToken.ValueStringPointer()
		}
		modifyListenerRequest.DefaultActions = expandAlbListenerDefaultActions(plan.DefaultActions)
		if !plan.GzipEnabled.IsNull() && !plan.GzipEnabled.IsUnknown() {
			modifyListenerRequest.GzipEnabled = plan.GzipEnabled.ValueBoolPointer()
		}
		if !plan.Http2Enabled.IsNull() && !plan.Http2Enabled.IsUnknown() {
			modifyListenerRequest.Http2Enabled = plan.Http2Enabled.ValueBoolPointer()
		}
		if !plan.IdleTimeout.IsNull() && !plan.IdleTimeout.IsUnknown() {
			modifyListenerRequest.IdleTimeout = fwHelper.Uint64PointerFromValue(plan.IdleTimeout)
		}
		if !plan.ListenerName.IsNull() && !plan.ListenerName.IsUnknown() {
			modifyListenerRequest.ListenerName = plan.ListenerName.ValueStringPointer()
		}
		if !plan.RequestTimeout.IsNull() && !plan.RequestTimeout.IsUnknown() {
			modifyListenerRequest.RequestTimeout = fwHelper.Uint64PointerFromValue(plan.RequestTimeout)
		}
		if !plan.SecurityPolicyId.IsNull() && !plan.SecurityPolicyId.IsUnknown() {
			modifyListenerRequest.SecurityPolicyId = plan.SecurityPolicyId.ValueStringPointer()
		}
		if plan.XForwardedForConfig != nil {
			modifyListenerRequest.XForwardedForConfig = expandAlbListenerXForwardedForConfig(plan.XForwardedForConfig)
		}
		if resp.Diagnostics.HasError() {
			return
		}

		if err := service.ModifyListenerAttributes(ctx, modifyListenerRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating ALB listener",
				fmt.Sprintf("ModifyListenerAttributes failed: %v", err),
			)

			return
		}
	}

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	describeListenersRequest.ListenerIds = []*string{fwHelper.String(listenerId)}
	listeners, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener after update", err.Error())
		return
	}

	item := findAlbListener(listeners, listenerId)
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB listener not readable after update",
			fmt.Sprintf("listener %s was updated but cannot be read back", state.Id.ValueString()),
		)

		return
	}

	plan.Id = state.Id
	refreshAlbListenerModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *AlbListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlbListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listenerId := state.Id.ValueString()

	deleteListenerRequest := albv20251030.NewDeleteListenerRequest()
	deleteListenerRequest.ListenerIds = []*string{fwHelper.String(listenerId)}
	deleteListenerRequest.LoadBalancerId = state.LoadBalancerId.ValueStringPointer()

	service := NewAlbService(r.Client())

	if err := service.DeleteListener(ctx, deleteListenerRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ALB listener",
			fmt.Sprintf("DeleteListener failed: %v", err),
		)

		return
	}
}

// refreshAlbListenerModel copies cloud-owned fields into the model.
//
// The describe API does not echo back ca_certificate_ids, certificate_ids,
// client_token or default_actions, so those keep the values already held in
// the model. Fields whose create/update request name differs from the
// response name are mapped: CaEnabled<-CaEnable, Http2Enabled<-Http2Enable,
// SecurityPolicyId<-TlsSecurityPolicyId.
func refreshAlbListenerModel(m *AlbListenerResourceModel, item *albv20251030.ListenerOutput) {
	if item == nil {
		return
	}

	m.ListenerPort = fwHelper.Int64ValueFromUintOrPrev(item.ListenerPort, m.ListenerPort)
	m.ListenerProtocol = fwHelper.StringValueOrPrev(item.ListenerProtocol, m.ListenerProtocol)
	m.ListenerName = fwHelper.StringValueOrPrev(item.ListenerName, m.ListenerName)
	m.CaEnabled = fwHelper.BoolValueOrPrev(item.CaEnable, m.CaEnabled)
	m.GzipEnabled = fwHelper.BoolValueOrPrev(item.GzipEnabled, m.GzipEnabled)
	m.Http2Enabled = fwHelper.BoolValueOrPrev(item.Http2Enable, m.Http2Enabled)
	m.IdleTimeout = fwHelper.Int64ValueFromUintOrPrev(item.IdleTimeout, m.IdleTimeout)
	m.RequestTimeout = fwHelper.Int64ValueFromUintOrPrev(item.RequestTimeout, m.RequestTimeout)
	m.SecurityPolicyId = fwHelper.StringValueOrPrev(item.TlsSecurityPolicyId, m.SecurityPolicyId)
	m.Tags = flattenAlbListenerTags(item.Tags, m.Tags)
	m.XForwardedForConfig = flattenAlbListenerXForwardedForConfig(item.XForwardedForConfig, m.XForwardedForConfig)
}

// findAlbListener picks this resource's own item from a describe list.
func findAlbListener(listeners []*albv20251030.ListenerOutput, listenerId string) *albv20251030.ListenerOutput {
	for _, l := range listeners {
		if l.ListenerId != nil && *l.ListenerId == listenerId {
			return l
		}
	}

	return nil
}

// albListenerAttributesChanged reports whether any attribute covered by
// ModifyListenerAttributes differs between the plan and the state.
func albListenerAttributesChanged(plan *AlbListenerResourceModel, state *AlbListenerResourceModel) bool {
	return !plan.ListenerName.Equal(state.ListenerName) ||
		!plan.CaCertificateIds.Equal(state.CaCertificateIds) ||
		!plan.CaEnabled.Equal(state.CaEnabled) ||
		!plan.CertificateIds.Equal(state.CertificateIds) ||
		!plan.GzipEnabled.Equal(state.GzipEnabled) ||
		!plan.Http2Enabled.Equal(state.Http2Enabled) ||
		!plan.IdleTimeout.Equal(state.IdleTimeout) ||
		!plan.RequestTimeout.Equal(state.RequestTimeout) ||
		!plan.SecurityPolicyId.Equal(state.SecurityPolicyId) ||
		!reflect.DeepEqual(plan.DefaultActions, state.DefaultActions) ||
		!reflect.DeepEqual(plan.XForwardedForConfig, state.XForwardedForConfig)
}

// expandAlbListenerDefaultActions converts the model list into SDK
// DefaultAction items.
func expandAlbListenerDefaultActions(models []AlbListenerDefaultActionModel) []*albv20251030.DefaultAction {
	if models == nil {
		return nil
	}

	actions := make([]*albv20251030.DefaultAction, 0, len(models))
	for _, m := range models {
		action := &albv20251030.DefaultAction{
			Type: m.Type.ValueStringPointer(),
		}
		if m.TargetGroupConfig != nil {
			action.TargetGroupConfig = expandAlbListenerTargetGroupConfig(m.TargetGroupConfig)
		}
		actions = append(actions, action)
	}

	return actions
}

// expandAlbListenerTargetGroupConfig converts the model into the SDK
// TargetGroupConfig object.
func expandAlbListenerTargetGroupConfig(m *AlbListenerTargetGroupConfigModel) *albv20251030.TargetGroupConfig {
	if m == nil {
		return nil
	}

	config := &albv20251030.TargetGroupConfig{}
	if m.TargetGroups != nil {
		config.TargetGroups = make([]*albv20251030.TargetGroupTuple, 0, len(m.TargetGroups))
		for _, tg := range m.TargetGroups {
			tuple := &albv20251030.TargetGroupTuple{
				TargetGroupId: tg.TargetGroupId.ValueStringPointer(),
			}
			if !tg.Weight.IsNull() && !tg.Weight.IsUnknown() {
				tuple.Weight = fwHelper.Int64PointerFromValue(tg.Weight)
			}
			config.TargetGroups = append(config.TargetGroups, tuple)
		}
	}
	if m.TargetGroupStickySession != nil {
		sticky := &albv20251030.TargetGroupStickySession{
			Enabled: m.TargetGroupStickySession.Enabled.ValueBoolPointer(),
		}
		if !m.TargetGroupStickySession.Timeout.IsNull() && !m.TargetGroupStickySession.Timeout.IsUnknown() {
			sticky.Timeout = fwHelper.Int64PointerFromValue(m.TargetGroupStickySession.Timeout)
		}
		config.TargetGroupStickySession = sticky
	}

	return config
}

// expandAlbListenerTags converts the model list into SDK TagInfo items.
func expandAlbListenerTags(models []AlbListenerTagModel) []*albv20251030.TagInfo {
	if models == nil {
		return nil
	}

	tags := make([]*albv20251030.TagInfo, 0, len(models))
	for _, m := range models {
		tags = append(tags, &albv20251030.TagInfo{
			TagKey:   m.TagKey.ValueStringPointer(),
			TagValue: m.TagValue.ValueStringPointer(),
		})
	}

	return tags
}

// expandAlbListenerXForwardedForConfig converts the model into the SDK
// XForwardedForConfig object.
func expandAlbListenerXForwardedForConfig(m *AlbListenerXForwardedForConfigModel) *albv20251030.XForwardedForConfig {
	if m == nil {
		return nil
	}

	return &albv20251030.XForwardedForConfig{
		XForwardedForAlbIdEnabled:         fwHelper.BoolPointerFromValue(m.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: fwHelper.BoolPointerFromValue(m.XForwardedForClientSrcPortEnabled),
		XForwardedForHostEnabled:          fwHelper.BoolPointerFromValue(m.XForwardedForHostEnabled),
		XForwardedForMode:                 fwHelper.StringPointerFromValue(m.XForwardedForMode),
		XForwardedForPortEnabled:          fwHelper.BoolPointerFromValue(m.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         fwHelper.BoolPointerFromValue(m.XForwardedForProtoEnabled),
		XTencentClientIDNEnabled:          fwHelper.BoolPointerFromValue(m.XTencentClientIDNEnabled),
		XTencentClientSDNEnabled:          fwHelper.BoolPointerFromValue(m.XTencentClientSDNEnabled),
		XTencentClientSerialEnabled:       fwHelper.BoolPointerFromValue(m.XTencentClientSerialEnabled),
		XTencentClientVerifyEnabled:       fwHelper.BoolPointerFromValue(m.XTencentClientVerifyEnabled),
	}
}

// flattenAlbListenerTags rebuilds the tag model list from the describe
// response; a cloud-omitted list keeps the previous value.
func flattenAlbListenerTags(cloud []*albv20251030.TagInfo, prev []AlbListenerTagModel) []AlbListenerTagModel {
	if cloud == nil {
		return prev
	}

	models := make([]AlbListenerTagModel, 0, len(cloud))
	for _, t := range cloud {
		if t == nil {
			continue
		}
		models = append(models, AlbListenerTagModel{
			TagKey:   fwHelper.StringPointerValueOrNull(t.TagKey),
			TagValue: fwHelper.StringPointerValueOrNull(t.TagValue),
		})
	}

	return models
}

// flattenAlbListenerXForwardedForConfig rebuilds the model object from the
// describe response; a cloud-omitted object keeps the previous value.
func flattenAlbListenerXForwardedForConfig(cloud *albv20251030.XForwardedForConfig, prev *AlbListenerXForwardedForConfigModel) *AlbListenerXForwardedForConfigModel {
	if cloud == nil {
		return prev
	}

	var prevModel AlbListenerXForwardedForConfigModel
	if prev != nil {
		prevModel = *prev
	}

	return &AlbListenerXForwardedForConfigModel{
		XForwardedForAlbIdEnabled:         fwHelper.BoolValueOrPrev(cloud.XForwardedForAlbIdEnabled, prevModel.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: fwHelper.BoolValueOrPrev(cloud.XForwardedForClientSrcPortEnabled, prevModel.XForwardedForClientSrcPortEnabled),
		XForwardedForHostEnabled:          fwHelper.BoolValueOrPrev(cloud.XForwardedForHostEnabled, prevModel.XForwardedForHostEnabled),
		XForwardedForMode:                 fwHelper.StringValueOrPrev(cloud.XForwardedForMode, prevModel.XForwardedForMode),
		XForwardedForPortEnabled:          fwHelper.BoolValueOrPrev(cloud.XForwardedForPortEnabled, prevModel.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         fwHelper.BoolValueOrPrev(cloud.XForwardedForProtoEnabled, prevModel.XForwardedForProtoEnabled),
		XTencentClientIDNEnabled:          fwHelper.BoolValueOrPrev(cloud.XTencentClientIDNEnabled, prevModel.XTencentClientIDNEnabled),
		XTencentClientSDNEnabled:          fwHelper.BoolValueOrPrev(cloud.XTencentClientSDNEnabled, prevModel.XTencentClientSDNEnabled),
		XTencentClientSerialEnabled:       fwHelper.BoolValueOrPrev(cloud.XTencentClientSerialEnabled, prevModel.XTencentClientSerialEnabled),
		XTencentClientVerifyEnabled:       fwHelper.BoolValueOrPrev(cloud.XTencentClientVerifyEnabled, prevModel.XTencentClientVerifyEnabled),
	}
}
