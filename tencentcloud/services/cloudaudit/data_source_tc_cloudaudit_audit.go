package cloudaudit

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
)

var _ datasource.DataSourceWithConfigure = &CloudauditAuditDataSource{}

// NewCloudauditAuditDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewCloudauditAuditDataSource() datasource.DataSource {
	return &CloudauditAuditDataSource{}
}

// CloudauditAuditDataSource implements datasource.DataSource for
// tencentcloud_cloudaudit_audit.
type CloudauditAuditDataSource struct {
	fw.DataSourceWithConfigure
}

// CloudauditAuditDataSourceModel maps the schema attributes.
type CloudauditAuditDataSourceModel struct {
	AuditName          types.String `tfsdk:"audit_name"`
	IsEnableCmqNotify  types.Int64  `tfsdk:"is_enable_cmq_notify"`
	ReadWriteAttribute types.Int64  `tfsdk:"read_write_attribute"`
	KeyId              types.String `tfsdk:"key_id"`
	AuditStatus        types.Int64  `tfsdk:"audit_status"`
	CosRegion          types.String `tfsdk:"cos_region"`
	CmqQueueName       types.String `tfsdk:"cmq_queue_name"`
	KmsAlias           types.String `tfsdk:"kms_alias"`
	KmsRegion          types.String `tfsdk:"kms_region"`
	IsEnableKmsEncry   types.Int64  `tfsdk:"is_enable_kms_encry"`
	CosBucketName      types.String `tfsdk:"cos_bucket_name"`
	CmqRegion          types.String `tfsdk:"cmq_region"`
	LogFilePrefix      types.String `tfsdk:"log_file_prefix"`
}

func (d *CloudauditAuditDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_cloudaudit_audit"
}

func (d *CloudauditAuditDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query detailed information of a CloudAudit audit track set.",
		Attributes: map[string]schema.Attribute{
			// ---- Query input ----
			"audit_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the audit track set.",
			},
			// ---- Computed outputs (read from cloud) ----
			"is_enable_cmq_notify": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether CMQ message notification is enabled. `1`: yes, `0`: no.",
			},
			"read_write_attribute": schema.Int64Attribute{
				Computed:    true,
				Description: "Management event read/write attribute. `1`: read only, `2`: write only, `3`: all.",
			},
			"key_id": schema.StringAttribute{
				Computed:    true,
				Description: "Global unique identifier of the CMK.",
			},
			"audit_status": schema.Int64Attribute{
				Computed:    true,
				Description: "Status of the audit track set. `1`: enabled, `0`: disabled.",
			},
			"cos_region": schema.StringAttribute{
				Computed:    true,
				Description: "Region of the COS bucket.",
			},
			"cmq_queue_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the CMQ queue.",
			},
			"kms_alias": schema.StringAttribute{
				Computed:    true,
				Description: "Alias of the CMK.",
			},
			"kms_region": schema.StringAttribute{
				Computed:    true,
				Description: "Region of KMS.",
			},
			"is_enable_kms_encry": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether KMS encryption is enabled. `1`: yes, `0`: no. When KMS encryption is enabled, data is encrypted before being delivered to COS.",
			},
			"cos_bucket_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the COS bucket.",
			},
			"cmq_region": schema.StringAttribute{
				Computed:    true,
				Description: "Region of the CMQ queue.",
			},
			"log_file_prefix": schema.StringAttribute{
				Computed:    true,
				Description: "Log file prefix.",
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *CloudauditAuditDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CloudauditAuditDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	auditName := config.AuditName.ValueString()
	if auditName == "" {
		resp.Diagnostics.AddError(
			"Missing audit_name",
			"audit_name is required and cannot be empty.",
		)
		return
	}

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	logId := tccommon.GetLogId(ctx)
	request := audit.NewDescribeAuditRequest()
	request.AuditName = helper.String(auditName)

	client := d.Client().UseAuditClient()
	response, err := client.DescribeAuditWithContext(ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFound.AuditNotExist") {
			resp.Diagnostics.AddError(
				"CloudAudit audit not found",
				fmt.Sprintf("Audit track set %q not found: %v", auditName, err),
			)
			return
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		resp.Diagnostics.AddError(
			"Error reading CloudAudit audit",
			fmt.Sprintf("DescribeAudit failed: %v", err),
		)
		return
	}

	// nil response is an API failure, not an empty result.
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error reading CloudAudit audit",
			"DescribeAudit returned an empty response.",
		)
		return
	}

	// Single-object query: a missing audit name signals an inconsistent
	// result and is surfaced as an error rather than a silent empty state.
	if response.Response.AuditName == nil || *response.Response.AuditName == "" {
		resp.Diagnostics.AddError(
			"CloudAudit audit not found",
			fmt.Sprintf("Audit track set %q returned an empty result.", auditName),
		)
		return
	}

	config.IsEnableCmqNotify = types.Int64PointerValue(response.Response.IsEnableCmqNotify)
	config.ReadWriteAttribute = types.Int64PointerValue(response.Response.ReadWriteAttribute)
	config.KeyId = types.StringPointerValue(response.Response.KeyId)
	config.AuditStatus = types.Int64PointerValue(response.Response.AuditStatus)
	config.CosRegion = types.StringPointerValue(response.Response.CosRegion)
	config.CmqQueueName = types.StringPointerValue(response.Response.CmqQueueName)
	config.KmsAlias = types.StringPointerValue(response.Response.KmsAlias)
	config.KmsRegion = types.StringPointerValue(response.Response.KmsRegion)
	config.IsEnableKmsEncry = types.Int64PointerValue(response.Response.IsEnableKmsEncry)
	config.CosBucketName = types.StringPointerValue(response.Response.CosBucketName)
	config.CmqRegion = types.StringPointerValue(response.Response.CmqRegion)
	config.LogFilePrefix = types.StringPointerValue(response.Response.LogFilePrefix)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
