package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbLogFiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbLogFilesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of `tdsql-ow728lmc`.",
			},
			"type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Requested log type. Valid values: 1 (binlog), 2 (cold backup), 3 (errlog), 4 (slowlog).",
			},
			"files": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information such as `uri`, `length`, and `mtime` (modification time).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mtime": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Last modified time of log.",
						},
						"length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File length.",
						},
						"uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Uniform resource identifier (URI) used during log download.",
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Filename.",
						},
					},
				},
			},
			"vpc_prefix": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "For an instance in a VPC, this prefix plus URI can be used as the download address.",
			},
			"normal_prefix": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "For an instance in a common network, this prefix plus URI can be used as the download address.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMariadbLogFilesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_log_files.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		logFiles   *mariadb.DescribeDBLogFilesResponseParams
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOkExists("type"); ok {
		paramMap["Type"] = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbLogFilesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		logFiles = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)
	if logFiles != nil {
		for _, logFileInfo := range logFiles.Files {
			logFileInfoMap := map[string]interface{}{}

			if logFileInfo.Mtime != nil {
				logFileInfoMap["mtime"] = logFileInfo.Mtime
			}

			if logFileInfo.Length != nil {
				logFileInfoMap["length"] = logFileInfo.Length
			}

			if logFileInfo.Uri != nil {
				logFileInfoMap["uri"] = logFileInfo.Uri
			}

			if logFileInfo.FileName != nil {
				logFileInfoMap["file_name"] = logFileInfo.FileName
			}

			tmpList = append(tmpList, logFileInfoMap)
		}

		_ = d.Set("files", tmpList)

		if logFiles.VpcPrefix != nil {
			_ = d.Set("vpc_prefix", logFiles.VpcPrefix)
		}

		if logFiles.NormalPrefix != nil {
			_ = d.Set("normal_prefix", logFiles.NormalPrefix)
		}
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
