/*
Use this data source to query detailed information of CBS snapshots.

Example Usage

```hcl
data "tencentcloud_cbs_snapshots" "snapshots" {
    snapshot_id        = "snap-f3io7adt"
    result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudCbsSnapshots() *schema.Resource {
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
				ValidateFunc: validateAllowedStringValue(CBS_STORAGE_USAGE),
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
	defer logElapsed("data_source.tencentcloud_cbs_snapshots.read")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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
	if v, ok := d.GetOk("project_id"); ok {
		params["project-id"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		params["zone"] = v.(string)
	}

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	snapshots, err := cbsService.DescribeSnapshotsByFilter(ctx, params)
	if err != nil {
		return err
	}

	snapshotList := make([]map[string]interface{}, 0, len(snapshots))
	ids := make([]string, len(snapshots))
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

	d.SetId(dataResourceIdsHash(ids))
	if err = d.Set("snapshot_list", snapshotList); err != nil {
		log.Printf("[CRITAL]%s provider set snapshot list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), snapshotList); err != nil {
			return err
		}
	}
	return nil
}
