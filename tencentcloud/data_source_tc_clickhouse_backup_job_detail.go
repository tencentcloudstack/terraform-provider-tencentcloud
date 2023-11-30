package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func dataSourceTencentCloudClickhouseBackupJobDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClickhouseBackupJobDetailRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"back_up_job_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Back up job id.",
			},

			"table_contents": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Back up tables.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database.",
						},
						"table": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table.",
						},
						"total_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total bytes.",
						},
						"v_cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual cluster.",
						},
						"ips": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ips.",
						},
						"zoo_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ZK path.",
						},
						"rip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ip address of cvm.",
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

func dataSourceTencentCloudClickhouseBackupJobDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clickhouse_backup_job_detail.read")()
	defer inconsistentCheck(d, meta)()

	var (
		request     = clickhouse.NewDescribeBackUpJobDetailRequest()
		instanceId  string
		backUpJobId int
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, _ := d.GetOk("back_up_job_id"); v != nil {
		backUpJobId = v.(int)
		request.BackUpJobId = helper.IntInt64(backUpJobId)
	}

	var tableContents []*clickhouse.BackupTableContent

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseCdwchClient().DescribeBackUpJobDetail(request)
		if e != nil {
			return retryError(e)
		}
		tableContents = response.Response.TableContents
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(tableContents))

	if tableContents != nil {
		for _, backupTableContent := range tableContents {
			backupTableContentMap := map[string]interface{}{}

			if backupTableContent.Database != nil {
				backupTableContentMap["database"] = backupTableContent.Database
			}

			if backupTableContent.Table != nil {
				backupTableContentMap["table"] = backupTableContent.Table
			}

			if backupTableContent.TotalBytes != nil {
				backupTableContentMap["total_bytes"] = backupTableContent.TotalBytes
			}

			if backupTableContent.VCluster != nil {
				backupTableContentMap["v_cluster"] = backupTableContent.VCluster
			}

			if backupTableContent.Ips != nil {
				backupTableContentMap["ips"] = backupTableContent.Ips
			}

			if backupTableContent.ZooPath != nil {
				backupTableContentMap["zoo_path"] = backupTableContent.ZooPath
			}

			if backupTableContent.Rip != nil {
				backupTableContentMap["rip"] = backupTableContent.Rip
			}

			tmpList = append(tmpList, backupTableContentMap)
		}

		_ = d.Set("table_contents", tmpList)
	}

	d.SetId(instanceId + helper.IntToStr(backUpJobId))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
