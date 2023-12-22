package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcSnapshotFiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSnapshotFilesRead,
		Schema: map[string]*schema.Schema{
			"business_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Business type, currently supports security group:securitygroup.",
			},

			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"start_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start date in the format %Y-%m-%d %H:%M:%S.",
			},

			"end_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End date in the format %Y-%m-%d %H:%M:%S.",
			},

			"snapshot_file_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "snap shot file set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Snapshot Policy Id.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"snapshot_file_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "snap shot file id.",
						},
						"backup_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "backup time.",
						},
						"operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Uin of operator.",
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

func dataSourceTencentCloudVpcSnapshotFilesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_snapshot_files.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("business_type"); ok {
		paramMap["BusinessType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_date"); ok {
		paramMap["StartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var snapshotFileSet []*vpc.SnapshotFileInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcSnapshotFilesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		snapshotFileSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(snapshotFileSet))
	tmpList := make([]map[string]interface{}, 0, len(snapshotFileSet))

	if snapshotFileSet != nil {
		for _, snapshotFileInfo := range snapshotFileSet {
			snapshotFileInfoMap := map[string]interface{}{}

			if snapshotFileInfo.SnapshotPolicyId != nil {
				snapshotFileInfoMap["snapshot_policy_id"] = snapshotFileInfo.SnapshotPolicyId
			}

			if snapshotFileInfo.InstanceId != nil {
				snapshotFileInfoMap["instance_id"] = snapshotFileInfo.InstanceId
			}

			if snapshotFileInfo.SnapshotFileId != nil {
				snapshotFileInfoMap["snapshot_file_id"] = snapshotFileInfo.SnapshotFileId
			}

			if snapshotFileInfo.BackupTime != nil {
				snapshotFileInfoMap["backup_time"] = snapshotFileInfo.BackupTime
			}

			if snapshotFileInfo.Operator != nil {
				snapshotFileInfoMap["operator"] = snapshotFileInfo.Operator
			}

			ids = append(ids, *snapshotFileInfo.SnapshotFileId)
			tmpList = append(tmpList, snapshotFileInfoMap)
		}

		_ = d.Set("snapshot_file_set", tmpList)
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
