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
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CBS_STORAGE_USAGE),
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"percent": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypt": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
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
