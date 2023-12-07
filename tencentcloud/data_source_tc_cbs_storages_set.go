package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCbsStoragesSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCbsStoragesSetRead,

		Schema: map[string]*schema.Schema{
			"storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CBS to be queried.",
			},
			"storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CBS to be queried.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the CBS instance locates at.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the project with which the CBS is associated.",
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CBS_STORAGE_TYPE),
				Description:  "Filter by cloud disk media type (`CLOUD_BASIC`: HDD cloud disk | `CLOUD_PREMIUM`: Premium Cloud Storage | `CLOUD_SSD`: SSD cloud disk).",
			},
			"storage_usage": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by cloud disk type (`SYSTEM_DISK`: system disk | `DATA_DISK`: data disk).",
			},
			"charge_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List filter by disk charge type (`POSTPAID_BY_HOUR` | `PREPAID`).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"portable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter by whether the disk is portable (Boolean `true` or `false`).",
			},
			"storage_state": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List filter by disk state (`UNATTACHED` | `ATTACHING` | `ATTACHED` | `DETACHING` | `EXPANDING` | `ROLLBACKING` | `TORECYCLE`).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"instance_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List filter by attached instance public or private IPs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"instance_name": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List filter by attached instance name.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tag_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List filter by tag keys.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tag_values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List filter by tag values.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"storage_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of storage. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CBS.",
						},
						"storage_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CBS.",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of storage medium.",
						},
						"storage_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of CBS.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The zone of CBS.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project.",
						},
						"storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume of CBS.",
						},
						"attached": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the CBS is mounted the CVM.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the CVM instance that be mounted by this CBS.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether CBS is encrypted.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of CBS.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of CBS.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this CBS.",
						},
						"prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The way that CBS instance will be renew automatically or not when it reach the end of the prepaid tenancy.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay type of the CBS instance.",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsStoragesSetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cbs_storages.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("storage_id"); ok {
		params["disk-id"] = v.(string)
	}
	if v, ok := d.GetOk("storage_name"); ok {
		params["disk-name"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		params["zone"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		params["project-id"] = fmt.Sprintf("%d", v.(int))
	}
	if v, ok := d.GetOk("storage_type"); ok {
		params["disk-type"] = v.(string)
	}
	if v, ok := d.GetOk("storage_usage"); ok {
		params["disk-usage"] = v.(string)
	}

	if v, ok := d.GetOk("charge_type"); ok {
		params["disk-charge-type"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("portable"); ok {
		if v.(bool) {
			params["portable"] = "TRUE"
		} else {
			params["portable"] = "FALSE"
		}
	}

	if v, ok := d.GetOk("storage_state"); ok {
		params["disk-state"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("instance_ips"); ok {
		params["instance-ip-address"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("instance_name"); ok {
		params["instance-name"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("tag_keys"); ok {
		params["tag-key"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("tag_values"); ok {
		params["tag-value"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	cbsService := CbsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	storages, e := cbsService.DescribeDisksInParallelByFilter(ctx, params)
	if e != nil {
		return e
	}
	ids := make([]string, 0, len(storages))
	storageList := make([]map[string]interface{}, 0, len(storages))
	for _, storage := range storages {
		mapping := map[string]interface{}{
			"storage_id":             storage.DiskId,
			"storage_name":           storage.DiskName,
			"storage_usage":          storage.DiskUsage,
			"storage_type":           storage.DiskType,
			"availability_zone":      storage.Placement.Zone,
			"project_id":             storage.Placement.ProjectId,
			"storage_size":           storage.DiskSize,
			"attached":               storage.Attached,
			"instance_id":            storage.InstanceId,
			"encrypt":                storage.Encrypt,
			"create_time":            storage.CreateTime,
			"status":                 storage.DiskState,
			"prepaid_renew_flag":     storage.RenewFlag,
			"charge_type":            storage.DiskChargeType,
			"throughput_performance": storage.ThroughputPerformance,
		}
		if storage.Tags != nil {
			tags := make(map[string]interface{}, len(storage.Tags))
			for _, t := range storage.Tags {
				tags[*t.Key] = *t.Value
			}
			mapping["tags"] = tags
		}
		storageList = append(storageList, mapping)
		ids = append(ids, *storage.DiskId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e = d.Set("storage_list", storageList); e != nil {
		log.Printf("[CRITAL]%s provider set storage list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), storageList); e != nil {
			return e
		}
	}

	return nil

}
