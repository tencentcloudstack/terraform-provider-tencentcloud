package teo

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ action.ActionWithConfigure = &TeoLogAnalysisDownloadTask{}

func NewTeoLogAnalysisDownloadTask() action.Action {
	return &TeoLogAnalysisDownloadTask{}
}

type TeoLogAnalysisDownloadTask struct {
	fw.ActionWithConfigure
}

type TeoLogAnalysisDownloadTaskModel struct {
	ZoneId    types.String `tfsdk:"zone_id"`
	Area      types.String `tfsdk:"area"`
	StartTime types.String `tfsdk:"start_time"`
	EndTime   types.String `tfsdk:"end_time"`
	LogType   types.String `tfsdk:"log_type"`
	Condition types.String `tfsdk:"condition"`
	Format    types.String `tfsdk:"format"`
	Sort      types.String `tfsdk:"sort"`
}

var (
	teoLogAnalysisDownloadTaskAreas    = []string{"mainland", "overseas"}
	teoLogAnalysisDownloadTaskLogTypes = []string{"l7-access-logs", "web-attack"}
	teoLogAnalysisDownloadTaskFormats  = []string{"csv"}
	teoLogAnalysisDownloadTaskSorts    = []string{"asc", "desc"}
)

func (a *TeoLogAnalysisDownloadTask) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_log_analysis_download_task"
}

func (a *TeoLogAnalysisDownloadTask) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to create a TEO (TencentCloud EdgeOne) log analysis download task. " +
			"The action creates the download task and waits for the task to be packed, then reports the download URL. " +
			"Download tasks are retained for 3 days after creation and cannot be deleted or modified.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Site ID.",
			},
			"area": schema.StringAttribute{
				Required:    true,
				Description: "Data region. Valid values: `mainland`, `overseas`. For a site serving global zones, query `mainland` and `overseas` separately.",
			},
			"start_time": schema.StringAttribute{
				Required:    true,
				Description: "Start time of the logs to download, for example `2020-04-29T00:00:00Z`. The maximum span between start time and end time is 31 days.",
			},
			"end_time": schema.StringAttribute{
				Required:    true,
				Description: "End time of the logs to download, for example `2020-04-30T00:00:00Z`.",
			},
			"log_type": schema.StringAttribute{
				Optional:    true,
				Description: "Log type. Valid values: `l7-access-logs` (L7 access logs), `web-attack` (managed rules logs). Default: `l7-access-logs`.",
			},
			"condition": schema.StringAttribute{
				Optional:    true,
				Description: "Log match condition, maximum length 12KB.",
			},
			"format": schema.StringAttribute{
				Optional:    true,
				Description: "File format. Valid values: `csv`. Default: `csv`.",
			},
			"sort": schema.StringAttribute{
				Optional:    true,
				Description: "Time sort of the raw logs. Valid values: `asc`, `desc`. Default: `desc`.",
			},
		},
	}
}

func (a *TeoLogAnalysisDownloadTask) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TeoLogAnalysisDownloadTaskModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, ok := teoLogAnalysisDownloadTaskRequiredString("zone_id", data.ZoneId, resp)
	if !ok {
		return
	}

	area, ok := teoLogAnalysisDownloadTaskRequiredString("area", data.Area, resp)
	if !ok {
		return
	}
	if !slices.Contains(teoLogAnalysisDownloadTaskAreas, area) {
		resp.Diagnostics.AddError(
			"Invalid area",
			fmt.Sprintf("area must be one of %v, got %s.", teoLogAnalysisDownloadTaskAreas, area),
		)
		return
	}

	startTime, ok := teoLogAnalysisDownloadTaskRequiredString("start_time", data.StartTime, resp)
	if !ok {
		return
	}

	endTime, ok := teoLogAnalysisDownloadTaskRequiredString("end_time", data.EndTime, resp)
	if !ok {
		return
	}

	request := teo.NewCreateLogAnalysisDownloadTaskRequest()
	request.ZoneId = &zoneId
	request.Area = &area
	request.StartTime = &startTime
	request.EndTime = &endTime
	request.LogType = teoLogAnalysisDownloadTaskOptionalString("log_type", data.LogType, teoLogAnalysisDownloadTaskLogTypes, resp)
	request.Condition = teoLogAnalysisDownloadTaskOptionalString("condition", data.Condition, nil, resp)
	request.Format = teoLogAnalysisDownloadTaskOptionalString("format", data.Format, teoLogAnalysisDownloadTaskFormats, resp)
	request.Sort = teoLogAnalysisDownloadTaskOptionalString("sort", data.Sort, teoLogAnalysisDownloadTaskSorts, resp)
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

	service := NewTeoService(a.Client())
	teoLogAnalysisDownloadTaskSendProgress(resp, "Creating TEO log analysis download task for zone %s (area: %s, %s ~ %s)", zoneId, area, startTime, endTime)

	taskId, err := service.CreateTeoLogAnalysisDownloadTask(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("creating TEO log analysis download task for zone (%s)", zoneId),
			err.Error(),
		)
		return
	}

	teoLogAnalysisDownloadTaskSendProgress(resp, "Created TEO log analysis download task: %s. Waiting for the download file to be packed.", taskId)

	task, err := service.WaitTeoLogAnalysisDownloadTaskById(ctx, tccommon.ReadRetryTimeout, zoneId, area, taskId)
	if err != nil {
		teoLogAnalysisDownloadTaskSendProgress(resp,
			"TEO log analysis download task %s was created, but its final status could not be confirmed: %s. "+
				"Query it with the DescribeLogAnalysisDownloadTasks API to get the download URL.",
			taskId, err.Error())
		return
	}

	if task.Status != nil && *task.Status == "failed" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("TEO log analysis download task (%s) failed", taskId),
			"The download task was created but failed to pack the logs. Please check the log match condition and the log type, then invoke the action again.",
		)
		return
	}

	if task.Url != nil && strings.TrimSpace(*task.Url) != "" {
		teoLogAnalysisDownloadTaskSendProgress(resp, "TEO log analysis download task %s is completed, download URL: %s", taskId, *task.Url)
	} else {
		teoLogAnalysisDownloadTaskSendProgress(resp,
			"TEO log analysis download task %s is completed, but no download URL was returned. "+
				"Query it with the DescribeLogAnalysisDownloadTasks API to get the download URL.",
			taskId)
	}
}

func teoLogAnalysisDownloadTaskSendProgress(resp *action.InvokeResponse, format string, args ...interface{}) {
	if resp.SendProgress == nil {
		return
	}
	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf(format, args...)})
}

func teoLogAnalysisDownloadTaskRequiredString(name string, value types.String, resp *action.InvokeResponse) (string, bool) {
	if value.IsNull() || value.IsUnknown() || strings.TrimSpace(value.ValueString()) == "" {
		resp.Diagnostics.AddError(
			"Missing "+name,
			name+" is required and cannot be empty.",
		)
		return "", false
	}
	return value.ValueString(), true
}

func teoLogAnalysisDownloadTaskOptionalString(name string, value types.String, allowed []string, resp *action.InvokeResponse) *string {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return nil
	}

	result := value.ValueString()
	if len(allowed) > 0 && !slices.Contains(allowed, result) {
		resp.Diagnostics.AddError(
			"Invalid "+name,
			fmt.Sprintf("%s must be one of %v, got %s.", name, allowed, result),
		)
		return nil
	}

	return &result
}
