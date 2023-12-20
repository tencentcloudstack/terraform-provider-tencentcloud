package cbs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCbsSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCbsSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the snapshot to be queried.",
			},
			"snapshot_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the snapshot to be queried.",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the the CBS which this snapshot created from.",
			},
			"storage_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CBS_STORAGE_USAGE),
				Description:  "Types of CBS which this snapshot created from, and available values include SYSTEM_DISK and DATA_DISK.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the project within the snapshot.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the CBS instance locates at.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"snapshot_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of snapshot. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the snapshot.",
						},
						"snapshot_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the snapshot.",
						},
						"storage_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the the CBS which this snapshot created from.",
						},
						"storage_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of CBS which this snapshot created from.",
						},
						"storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume of storage which this snapshot created from.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone that the CBS instance locates at.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project within the snapshot.",
						},
						"percent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Snapshot creation progress percentage.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of snapshot.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the snapshot is encrypted.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cbs_snapshots.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]string)
	if v, ok := d.GetOk("snapshot_id"); ok {
		params["snapshot-id"] = v.(string)
	}
	if v, ok := d.GetOk("snapshot_name"); ok {
		params["snapshot-name"] = v.(string)
	}
	if v, ok := d.GetOk("storage_id"); ok {
		params["disk-id"] = v.(string)
	}
	if v, ok := d.GetOk("storage_usage"); ok {
		params["disk-usage"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		params["project-id"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		params["zone"] = v.(string)
	}

	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		snapshots, e := cbsService.DescribeSnapshotsByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		ids := make([]string, 0, len(snapshots))
		snapshotList := make([]map[string]interface{}, 0, len(snapshots))
		for _, snapshot := range snapshots {
			mapping := map[string]interface{}{
				"snapshot_id":       *snapshot.SnapshotId,
				"snapshot_name":     *snapshot.SnapshotName,
				"storage_id":        *snapshot.DiskId,
				"storage_usage":     *snapshot.DiskUsage,
				"storage_size":      *snapshot.DiskSize,
				"availability_zone": *snapshot.Placement.Zone,
				"project_id":        *snapshot.Placement.ProjectId,
				"percent":           *snapshot.Percent,
				"create_time":       *snapshot.CreateTime,
				"encrypt":           *snapshot.Encrypt,
			}
			snapshotList = append(snapshotList, mapping)
			ids = append(ids, *snapshot.SnapshotId)
		}

		d.SetId(helper.DataResourceIdsHash(ids))
		if e = d.Set("snapshot_list", snapshotList); e != nil {
			log.Printf("[CRITAL]%s provider set snapshot list fail, reason:%s\n ", logId, e.Error())
			return resource.NonRetryableError(e)
		}

		output, ok := d.GetOk("result_output_file")
		if ok && output.(string) != "" {
			if e := tccommon.WriteToFile(output.(string), snapshotList); e != nil {
				return resource.NonRetryableError(e)
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs snapshots failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
