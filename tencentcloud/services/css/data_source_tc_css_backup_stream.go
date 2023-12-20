package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssBackupStream() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssBackupStreamRead,
		Schema: map[string]*schema.Schema{
			"stream_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Stream id.",
			},

			"stream_info_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Backup stream group info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stream name.",
						},
						"backup_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Backup stream info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Push domain.",
									},
									"app_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Push path.",
									},
									"publish_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "UTC time, eg, 2018-06-29T19:00:00Z.",
									},
									"upstream_sequence": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Push stream sequence.",
									},
									"source_from": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source from.",
									},
									"master_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Master stream flag.",
									},
								},
							},
						},
						"optimal_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Optimal switch, 1-enable, 0-disable.",
						},
						"host_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group name.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCssBackupStreamRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_backup_stream.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("stream_name"); ok {
		paramMap["StreamName"] = helper.String(v.(string))
	}

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var streamInfoList []*css.BackupStreamGroupInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssBackupStreamByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		streamInfoList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(streamInfoList))
	tmpList := make([]map[string]interface{}, 0, len(streamInfoList))

	if streamInfoList != nil {
		for _, backupStreamGroupInfo := range streamInfoList {
			backupStreamGroupInfoMap := map[string]interface{}{}

			if backupStreamGroupInfo.StreamName != nil {
				backupStreamGroupInfoMap["stream_name"] = backupStreamGroupInfo.StreamName
			}

			if backupStreamGroupInfo.BackupList != nil {
				backupListList := []interface{}{}
				for _, backupList := range backupStreamGroupInfo.BackupList {
					backupListMap := map[string]interface{}{}

					if backupList.DomainName != nil {
						backupListMap["domain_name"] = backupList.DomainName
					}

					if backupList.AppName != nil {
						backupListMap["app_name"] = backupList.AppName
					}

					if backupList.PublishTime != nil {
						backupListMap["publish_time"] = backupList.PublishTime
					}

					if backupList.UpstreamSequence != nil {
						backupListMap["upstream_sequence"] = backupList.UpstreamSequence
					}

					if backupList.SourceFrom != nil {
						backupListMap["source_from"] = backupList.SourceFrom
					}

					if backupList.MasterFlag != nil {
						backupListMap["master_flag"] = backupList.MasterFlag
					}

					backupListList = append(backupListList, backupListMap)
				}

				backupStreamGroupInfoMap["backup_list"] = backupListList
			}

			if backupStreamGroupInfo.OptimalEnable != nil {
				backupStreamGroupInfoMap["optimal_enable"] = backupStreamGroupInfo.OptimalEnable
			}

			if backupStreamGroupInfo.HostGroupName != nil {
				backupStreamGroupInfoMap["host_group_name"] = backupStreamGroupInfo.HostGroupName
			}

			ids = append(ids, *backupStreamGroupInfo.StreamName)
			tmpList = append(tmpList, backupStreamGroupInfoMap)
		}

		_ = d.Set("stream_info_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
