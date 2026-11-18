package makers

import (
	"context"
	"fmt"
	"time"

	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	makerssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ action.Action is intentionally omitted: the factory function
// NewMakersDeployment returns action.Action, which already performs the
// implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &MakersDeployment{}

// NewMakersDeployment is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
func NewMakersDeployment() action.Action {
	return &MakersDeployment{}
}

// defaultMakersDeploymentInvokeTimeout is the built-in maximum duration to wait
// for the asynchronous build to finish. It tracks the provider's global write
// timeout so an env override of the write timeout is honored. Practitioners
// can override it through the `timeouts` block, e.g. `timeouts { invoke = "30m" }`.
var defaultMakersDeploymentInvokeTimeout = tccommon.WriteRetryTimeout

// makersDeploymentPollInterval is the interval between two DescribeDeployments
// polls. It is a package-level var (not a const) so unit tests can shorten it.
var makersDeploymentPollInterval = 5 * time.Second

// MakersDeploymentPollIntervalForTest exposes the poll interval for tests to
// snapshot/restore. It is only meant to be used from the action's test file.
func MakersDeploymentPollIntervalForTest() time.Duration {
	return makersDeploymentPollInterval
}

// SetMakersDeploymentPollIntervalForTest overrides the poll interval for tests
// and always returns true. It is only meant to be used from the action's test
// file.
func SetMakersDeploymentPollIntervalForTest(d time.Duration) bool {
	makersDeploymentPollInterval = d
	return true
}

// MakersDeployment implements action.Action for
// tencentcloud_makers_deployment.
type MakersDeployment struct {
	fw.ActionWithConfigure
}

// MakersDeploymentModel maps the schema attributes.
//
// Action schemas cannot be Computed or Sensitive: the only marks available are
// Required / Optional / WriteOnly (WriteOnly needs framework v1.17.0+). Every
// attribute below maps 1:1 to a field of the CreateDeployment request; the
// provider forwards values as-is without rewriting or defaulting. `timeouts`
// is a provider-side block, not an API field.
type MakersDeploymentModel struct {
	ZoneId         types.String `tfsdk:"zone_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ViaMeta        types.String `tfsdk:"via_meta"`
	Env            types.String `tfsdk:"env"`
	Provider       types.String `tfsdk:"provider"`
	DistType       types.String `tfsdk:"dist_type"`
	TempBucketPath types.String `tfsdk:"temp_bucket_path"`
	RepoBranch     types.String `tfsdk:"repo_branch"`
	RepoCommitHash types.String `tfsdk:"repo_commit_hash"`
	RepoCommitMsg  types.String `tfsdk:"repo_commit_msg"`
	BuildFrom      types.String `tfsdk:"build_from"`
	TemplateId     types.String `tfsdk:"template_id"`
	CiSpec         types.String `tfsdk:"ci_spec"`

	// Timeouts is the value type defined by the official
	// terraform-plugin-framework-timeouts module. Only the `invoke` stage
	// exists because an action has a single lifecycle stage.
	Timeouts actiontimeouts.Value `tfsdk:"timeouts"`
}

func (a *MakersDeployment) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_makers_deployment"
}

func (a *MakersDeployment) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to trigger a Makers build via the CreateDeployment API. " +
			"This is a one-time invocation; it does not create, update or destroy any Terraform-managed resource, " +
			"and it exposes no attributes referenceable from other resources. " +
			"CreateDeployment is asynchronous: it returns a deployment id as soon as the build is accepted, " +
			"and the action polls DescribeDeployments until the build reaches a terminal state before returning; " +
			"a timeout diagnostic carries the last observed status. " +
			"Repeating the invocation triggers a brand-new build each time (the deployment id is not exposed and " +
			"not reused), and destroy does not roll back or cancel a build already accepted by the cloud. " +
			"Triggering this action through lifecycle.action_trigger runs it automatically during apply, " +
			"including on drift-driven after_update events.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Makers site ID.",
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "Project ID.",
			},
			"via_meta": schema.StringAttribute{
				Required:    true,
				Description: "Build trigger source. Valid values: `Github`, `ReDeploy`, `Api`, `Webhook`, `Upload`, `Coding`, `Gitee`, `Gitlab`.",
			},
			"env": schema.StringAttribute{
				Required:    true,
				Description: "Environment slug, e.g. `Production`, `Preview` or a custom environment.",
			},
			"provider": schema.StringAttribute{
				Required:    true,
				Description: "Code source. Use `Upload` when building from an uploaded artifact.",
			},
			"dist_type": schema.StringAttribute{
				Optional:    true,
				Description: "Artifact type. Valid values: `Zip`, `Folder`.",
			},
			"temp_bucket_path": schema.StringAttribute{
				Optional:    true,
				Description: "Object path of the uploaded artifact in the temp bucket, used when `provider` is `Upload`.",
			},
			"repo_branch": schema.StringAttribute{
				Optional:    true,
				Description: "Repository branch to build.",
			},
			"repo_commit_hash": schema.StringAttribute{
				Optional:    true,
				Description: "Commit hash to build.",
			},
			"repo_commit_msg": schema.StringAttribute{
				Optional:    true,
				Description: "Commit message.",
			},
			"build_from": schema.StringAttribute{
				Optional:    true,
				Description: "Build source.",
			},
			"template_id": schema.StringAttribute{
				Optional:    true,
				Description: "Template ID.",
			},
			"ci_spec": schema.StringAttribute{
				Optional:    true,
				Description: "Build specification. Reserved field: the server side currently still reads it from edgeone.json, so it has no effect yet.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": makersDeploymentTimeoutsBlock(ctx),
		},
	}
}

// makersDeploymentTimeoutsBlock builds the `timeouts` block from the official
// terraform-plugin-framework-timeouts module, whose single optional `invoke`
// attribute is validated as a time.Duration at plan time.
func makersDeploymentTimeoutsBlock(ctx context.Context) schema.Block {
	block := actiontimeouts.BlockWithOpts(ctx, actiontimeouts.Opts{
		InvokeDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `5m`. " +
			"It limits how long the action waits for the build to reach a terminal state. Default follows the provider write timeout.",
	})

	if nested, ok := block.(schema.SingleNestedBlock); ok {
		nested.Description = "The timeouts block allows you to specify the timeout for the invoke operation."
		return nested
	}

	return block
}

// Invoke is called to run the logic of the action.
func (a *MakersDeployment) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data MakersDeploymentModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	// An unset/null timeouts block falls back to the built-in default, and a
	// malformed duration is reported as an error diagnostic. The schema
	// validator normally rejects malformed values at plan time; this is the
	// runtime safety net.
	timeout, timeoutDiags := data.Timeouts.Invoke(ctx, defaultMakersDeploymentInvokeTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	progress := func(msg string) {
		if resp.SendProgress != nil {
			resp.SendProgress(action.InvokeProgressEvent{Message: msg})
		}
	}

	// Required attributes are forwarded to the API as-is: the provider performs
	// no enum check, no defaulting and no rewriting.
	request := makerssdk.NewCreateDeploymentRequest()
	request.ZoneId = data.ZoneId.ValueStringPointer()
	request.ProjectId = data.ProjectId.ValueStringPointer()
	request.ViaMeta = data.ViaMeta.ValueStringPointer()
	request.Env = data.Env.ValueStringPointer()
	request.Provider = data.Provider.ValueStringPointer()

	// Optional attributes are forwarded only when the practitioner set them,
	// so an omitted attribute stays absent from the cloud request instead of
	// being sent as an explicit empty string.
	if !data.DistType.IsNull() && !data.DistType.IsUnknown() {
		request.DistType = data.DistType.ValueStringPointer()
	}
	if !data.TempBucketPath.IsNull() && !data.TempBucketPath.IsUnknown() {
		request.TempBucketPath = data.TempBucketPath.ValueStringPointer()
	}
	if !data.RepoBranch.IsNull() && !data.RepoBranch.IsUnknown() {
		request.RepoBranch = data.RepoBranch.ValueStringPointer()
	}
	if !data.RepoCommitHash.IsNull() && !data.RepoCommitHash.IsUnknown() {
		request.RepoCommitHash = data.RepoCommitHash.ValueStringPointer()
	}
	if !data.RepoCommitMsg.IsNull() && !data.RepoCommitMsg.IsUnknown() {
		request.RepoCommitMsg = data.RepoCommitMsg.ValueStringPointer()
	}
	if !data.BuildFrom.IsNull() && !data.BuildFrom.IsUnknown() {
		request.BuildFrom = data.BuildFrom.ValueStringPointer()
	}
	if !data.TemplateId.IsNull() && !data.TemplateId.IsUnknown() {
		request.TemplateId = data.TemplateId.ValueStringPointer()
	}
	if !data.CiSpec.IsNull() && !data.CiSpec.IsUnknown() {
		request.CiSpec = data.CiSpec.ValueStringPointer()
	}

	zoneId := data.ZoneId.ValueString()
	projectId := data.ProjectId.ValueString()

	progress(fmt.Sprintf("calling CreateDeployment for Makers project (%s) of zone (%s)", projectId, zoneId))
	service := NewMakersService(a.Client())
	deploymentId, err := service.CreateDeployment(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("triggering Makers deployment for project (%s) of zone (%s)", projectId, zoneId),
			err.Error(),
		)
		return
	}
	if deploymentId == nil || *deploymentId == "" {
		// The API documents DeploymentId as nullable. Without an id there is no
		// way to confirm the build outcome, so this is surfaced as a failure
		// instead of a silent success.
		resp.Diagnostics.AddError(
			fmt.Sprintf("triggering Makers deployment for project (%s) of zone (%s)", projectId, zoneId),
			"CreateDeployment succeeded but returned no deployment id, so the build outcome cannot be confirmed.",
		)
		return
	}

	// CreateDeployment is asynchronous: acceptance only means the build was
	// queued. Poll DescribeDeployments to a terminal state before returning.
	progress(fmt.Sprintf("build %s accepted, waiting for it to finish", *deploymentId))
	last, err := service.WaitMakersDeploymentComplete(ctx, zoneId, projectId, *deploymentId, timeout, makersDeploymentPollInterval)
	if err != nil {
		snapshot := "no status observed"
		if last != nil {
			statusStr := "nil"
			if last.Status != nil {
				statusStr = *last.Status
			}
			codeStr := "nil"
			if last.Code != nil {
				codeStr = fmt.Sprintf("%d", *last.Code)
			}
			snapshot = fmt.Sprintf("status=%s, code=%s", statusStr, codeStr)
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("waiting for Makers deployment (%s) of project (%s) to complete", *deploymentId, projectId),
			fmt.Sprintf("%s; last observed: %s", err.Error(), snapshot),
		)
		return
	}

	progress(fmt.Sprintf("Makers deployment %s completed", *deploymentId))
	// No state to write: actions have no terraform state and no output attributes.
}
