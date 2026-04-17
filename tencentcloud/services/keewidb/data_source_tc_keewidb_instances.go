package keewidb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	keewidbv20220308 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/keewidb/v20220308"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKeewidbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKeewidbInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance ID, e.g. `kee-6ubhg****`.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by instance name.",
			},

			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fuzzy search keyword. Supports instance ID or instance name.",
			},

			"uniq_vpc_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by VPC ID (string format, e.g. vpc-xxx).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"uniq_subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by subnet ID (string format, e.g. subnet-xxx).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"project_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by project IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by instance status. 0: pending init; 1: in process; 2: running; -2: isolated; -3: to be deleted.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"billing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by billing mode. postpaid: pay-as-you-go; prepaid: prepaid.",
			},

			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field. Valid values: projectId, createtime, instancename, type, curDeadline.",
			},

			"order_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sort direction. 1: descending (default); 0: ascending.",
			},

			"type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter by instance type. 13: standard; 14: cluster.",
			},

			"auto_renew": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by renewal mode. 0: manual renewal; 1: auto-renewal; 2: no renewal on expiry.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"vpc_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by VPC ID (numeric format).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by subnet ID (numeric format).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"search_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search keywords. Supports instance ID, instance name, and private network IP.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tag_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by tag keys.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by tag key and value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			// Computed
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of KeeWiDB instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region ID.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Availability zone ID.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID (string format).",
						},
						"uniq_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID (string format).",
						},
						"wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance VIP.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance port.",
						},
						"createtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance creation time.",
						},
						"size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total persistent memory capacity (MB).",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance type. 13: standard; 14: cluster.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto-renewal flag. 1: enabled; 0: disabled.",
						},
						"deadline_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance expiry time.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage engine.",
						},
						"product_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product type. standalone or cluster.",
						},
						"billing_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing mode. 0: pay-as-you-go; 1: prepaid.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project name.",
						},
						"no_auth": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the instance is password-free.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk capacity (GB).",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region, e.g. ap-guangzhou.",
						},
						"machine_memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory capacity (GB).",
						},
						"compression": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data compression switch. ON or OFF.",
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

func dataSourceTencentCloudKeewidbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_keewidb_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = KeewidbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		paramMap["InstanceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_ids"); ok {
		items := v.([]interface{})
		strs := make([]*string, 0, len(items))
		for _, item := range items {
			strs = append(strs, helper.String(item.(string)))
		}
		paramMap["UniqVpcIds"] = strs
	}

	if v, ok := d.GetOk("uniq_subnet_ids"); ok {
		items := v.([]interface{})
		strs := make([]*string, 0, len(items))
		for _, item := range items {
			strs = append(strs, helper.String(item.(string)))
		}
		paramMap["UniqSubnetIds"] = strs
	}

	if v, ok := d.GetOk("project_ids"); ok {
		items := v.([]interface{})
		ids := make([]*int64, 0, len(items))
		for _, item := range items {
			ids = append(ids, helper.IntInt64(item.(int)))
		}
		paramMap["ProjectIds"] = ids
	}

	if v, ok := d.GetOk("status"); ok {
		items := v.([]interface{})
		statuses := make([]*int64, 0, len(items))
		for _, item := range items {
			statuses = append(statuses, helper.IntInt64(item.(int)))
		}
		paramMap["Status"] = statuses
	}

	if v, ok := d.GetOk("billing_mode"); ok {
		paramMap["BillingMode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_type"); ok {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("auto_renew"); ok {
		items := v.([]interface{})
		vals := make([]*int64, 0, len(items))
		for _, item := range items {
			vals = append(vals, helper.IntInt64(item.(int)))
		}
		paramMap["AutoRenew"] = vals
	}

	if v, ok := d.GetOk("vpc_ids"); ok {
		items := v.([]interface{})
		strs := make([]*string, 0, len(items))
		for _, item := range items {
			strs = append(strs, helper.String(item.(string)))
		}
		paramMap["VpcIds"] = strs
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		items := v.([]interface{})
		strs := make([]*string, 0, len(items))
		for _, item := range items {
			strs = append(strs, helper.String(item.(string)))
		}
		paramMap["SubnetIds"] = strs
	}

	if v, ok := d.GetOk("search_keys"); ok {
		items := v.([]interface{})
		strs := make([]*string, 0, len(items))
		for _, item := range items {
			strs = append(strs, helper.String(item.(string)))
		}
		paramMap["SearchKeys"] = strs
	}

	if v, ok := d.GetOk("tag_keys"); ok {
		items := v.([]interface{})
		strs := make([]*string, 0, len(items))
		for _, item := range items {
			strs = append(strs, helper.String(item.(string)))
		}
		paramMap["TagKeys"] = strs
	}

	if v, ok := d.GetOk("tag_list"); ok {
		tagListRaw := v.([]interface{})
		tagList := make([]*keewidbv20220308.InstanceTagInfo, 0, len(tagListRaw))
		for _, item := range tagListRaw {
			m := item.(map[string]interface{})
			tag := &keewidbv20220308.InstanceTagInfo{}
			if k, ok := m["tag_key"].(string); ok && k != "" {
				tag.TagKey = helper.String(k)
			}
			if val, ok := m["tag_value"].(string); ok && val != "" {
				tag.TagValue = helper.String(val)
			}
			tagList = append(tagList, tag)
		}
		paramMap["TagList"] = tagList
	}

	var respData []*keewidbv20220308.InstanceInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKeewidbInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if reqErr != nil {
		return reqErr
	}

	instanceList := make([]map[string]interface{}, 0, len(respData))
	for _, inst := range respData {
		m := map[string]interface{}{}
		if inst.InstanceId != nil {
			m["instance_id"] = inst.InstanceId
		}
		if inst.InstanceName != nil {
			m["instance_name"] = inst.InstanceName
		}
		if inst.Status != nil {
			m["status"] = int(*inst.Status)
		}
		if inst.RegionId != nil {
			m["region_id"] = int(*inst.RegionId)
		}
		if inst.ZoneId != nil {
			m["zone_id"] = int(*inst.ZoneId)
		}
		if inst.UniqVpcId != nil {
			m["uniq_vpc_id"] = inst.UniqVpcId
		}
		if inst.UniqSubnetId != nil {
			m["uniq_subnet_id"] = inst.UniqSubnetId
		}
		if inst.WanIp != nil {
			m["wan_ip"] = inst.WanIp
		}
		if inst.Port != nil {
			m["port"] = int(*inst.Port)
		}
		if inst.Createtime != nil {
			m["createtime"] = inst.Createtime
		}
		if inst.Size != nil {
			m["size"] = *inst.Size
		}
		if inst.Type != nil {
			m["type"] = int(*inst.Type)
		}
		if inst.AutoRenewFlag != nil {
			m["auto_renew_flag"] = int(*inst.AutoRenewFlag)
		}
		if inst.DeadlineTime != nil {
			m["deadline_time"] = inst.DeadlineTime
		}
		if inst.Engine != nil {
			m["engine"] = inst.Engine
		}
		if inst.ProductType != nil {
			m["product_type"] = inst.ProductType
		}
		if inst.BillingMode != nil {
			m["billing_mode"] = int(*inst.BillingMode)
		}
		if inst.ProjectId != nil {
			m["project_id"] = int(*inst.ProjectId)
		}
		if inst.ProjectName != nil {
			m["project_name"] = inst.ProjectName
		}
		if inst.NoAuth != nil {
			m["no_auth"] = *inst.NoAuth
		}
		if inst.DiskSize != nil {
			m["disk_size"] = int(*inst.DiskSize)
		}
		if inst.Region != nil {
			m["region"] = inst.Region
		}
		if inst.MachineMemory != nil {
			m["machine_memory"] = int(*inst.MachineMemory)
		}
		if inst.Compression != nil {
			m["compression"] = inst.Compression
		}
		instanceList = append(instanceList, m)
	}

	_ = d.Set("instance_list", instanceList)
	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil
}
