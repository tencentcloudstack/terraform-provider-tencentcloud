/*
Use this data source to query detailed information of lighthouse disk

Example Usage

```hcl
data "tencentcloud_lighthouse_disks" "disks" {
  disk_ids = ["lhdisk-xxxxxx"]
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseInstanceDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseInstanceDisksRead,
		Schema: map[string]*schema.Schema{
			"disk_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of disk ids.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields to be filtered. Valid names: `disk-id`: Filters by disk id; `instance-id`: Filter by instance id; `disk-name`: Filter by disk name; `zone`: Filter by zone; `disk-usage`: Filter by disk usage(Values: `SYSTEM_DISK` or `DATA_DISK`); `disk-state`: Filter by disk state.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Value of the field.",
						},
					},
				},
			},

			"disk_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cloud disk information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk id.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance id.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"disk_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk name.",
						},
						"disk_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk usage.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk type.",
						},
						"disk_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk charge type.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Renew flag.",
						},
						"disk_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk state. Valid values:`PENDING`, `UNATTACHED`, `ATTACHING`, `ATTACHED`, `DETACHING`, `SHUTDOWN`, `CREATED_FAILED`, `TERMINATING`, `DELETING`, `FREEZING`.",
						},
						"attached": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Disk attach state.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to release with the instance.",
						},
						"latest_operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest operation.",
						},
						"latest_operation_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest operation state.",
						},
						"latest_operation_request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest operation request id.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created time. Expressed according to the ISO8601 standard, and using UTC time. The format is `YYYY-MM-DDThh:mm:ssZ`.",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expired time. Expressed according to the ISO8601 standard, and using UTC time. The format is `YYYY-MM-DDThh:mm:ssZ`.",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Isolated time. Expressed according to the ISO8601 standard, and using UTC time. The format is `YYYY-MM-DDThh:mm:ssZ`.",
						},
						"disk_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of existing backup points of cloud disk.",
						},
						"disk_backup_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of backup points quota for cloud disk.",
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

func dataSourceTencentCloudLighthouseInstanceDisksRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_instance_disks.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	diskIds := make([]string, 0)
	for _, diskId := range d.Get("disk_ids").(*schema.Set).List() {
		diskIds = append(diskIds, diskId.(string))
	}
	filters := make([]*lighthouse.Filter, 0)
	if v, ok := d.GetOk("filters"); ok {
		filterSet := v.([]interface{})

		for _, item := range filterSet {
			filter := lighthouse.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			filters = append(filters, &filter)
		}
	}
	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	disks, err := service.DescribeLighthouseDisk(ctx, diskIds, filters)
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	diskList := make([]map[string]interface{}, 0)
	for _, disk := range disks {
		diskMap := make(map[string]interface{})
		if disk.DiskId != nil {
			diskMap["disk_id"] = disk.DiskId
			ids = append(ids, *disk.DiskId)
		}

		if disk.InstanceId != nil {
			diskMap["instance_id"] = disk.InstanceId
		}

		if disk.Zone != nil {
			diskMap["zone"] = disk.Zone
		}

		if disk.DiskName != nil {
			diskMap["disk_name"] = disk.DiskName
		}

		if disk.DiskUsage != nil {
			diskMap["disk_usage"] = disk.DiskUsage
		}

		if disk.DiskType != nil {
			diskMap["disk_type"] = disk.DiskType
		}

		if disk.DiskChargeType != nil {
			diskMap["disk_charge_type"] = disk.DiskChargeType
		}

		if disk.DiskSize != nil {
			diskMap["disk_size"] = disk.DiskSize
		}

		if disk.RenewFlag != nil {
			diskMap["renew_flag"] = disk.RenewFlag
		}

		if disk.DiskState != nil {
			diskMap["disk_state"] = disk.DiskState
		}

		if disk.Attached != nil {
			diskMap["attached"] = disk.Attached
		}

		if disk.DeleteWithInstance != nil {
			diskMap["delete_with_instance"] = disk.DeleteWithInstance
		}

		if disk.LatestOperation != nil {
			diskMap["latest_operation"] = disk.LatestOperation
		}

		if disk.LatestOperationState != nil {
			diskMap["latest_operation_state"] = disk.LatestOperationState
		}

		if disk.LatestOperationRequestId != nil {
			diskMap["latest_operation_request_id"] = disk.LatestOperationRequestId
		}

		if disk.CreatedTime != nil {
			diskMap["created_time"] = disk.CreatedTime
		}

		if disk.ExpiredTime != nil {
			diskMap["expired_time"] = disk.ExpiredTime
		}

		if disk.IsolatedTime != nil {
			diskMap["isolated_time"] = disk.IsolatedTime
		}

		if disk.DiskBackupCount != nil {
			diskMap["disk_backup_count"] = disk.DiskBackupCount
		}

		if disk.DiskBackupQuota != nil {
			diskMap["disk_backup_quota"] = disk.DiskBackupQuota
		}

		diskList = append(diskList, diskMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("disk_list", diskList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), diskList); e != nil {
			return e
		}
	}
	return nil
}
