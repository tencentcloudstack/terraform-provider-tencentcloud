package cdwpg

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdwpgInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdwpgInstancesRead,
		Schema: map[string]*schema.Schema{
			"search_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search instance id.",
			},

			"search_instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search instance name.",
			},

			"search_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Search tags.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"instances_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instances list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "id.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance id.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance name.",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Version.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Region id.",
						},
						"region_desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region description.",
						},
						"zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Zone.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Zone id.",
						},
						"zone_desc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Zone description.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vpc id.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Subnet id.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Create time, such as 2022-09-05 20:00:01.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Expire time, such as 2022-09-05 20:00:01.",
						},
						"access_info": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access information.",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Pay mode.",
						},
						"renew_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Renew flag.",
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

func dataSourceTencentCloudCdwpgInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdwpg_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_instance_id"); ok {
		paramMap["SearchInstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_instance_name"); ok {
		paramMap["SearchInstanceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_tags"); ok {
		searchTagsList := []*string{}
		searchTagsSet := v.(*schema.Set).List()
		for i := range searchTagsSet {
			searchTags := searchTagsSet[i].(string)
			searchTagsList = append(searchTagsList, helper.String(searchTags))
		}
		paramMap["SearchTags"] = searchTagsList
	}

	var instances []*cdwpg.InstanceSimpleInfoNew
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdwpgInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instances))
	instancesList := make([]map[string]interface{}, 0, len(instances))
	for _, instance := range instances {
		instancesListMap := map[string]interface{}{}

		if instance.ID != nil {
			instancesListMap["id"] = instance.ID
		}

		if instance.InstanceId != nil {
			ids = append(ids, *instance.InstanceId)
			instancesListMap["instance_id"] = instance.InstanceId
		}

		if instance.InstanceName != nil {
			instancesListMap["instance_name"] = instance.InstanceName
		}

		if instance.Version != nil {
			instancesListMap["version"] = instance.Version
		}

		if instance.Region != nil {
			instancesListMap["region"] = instance.Region
		}

		if instance.RegionId != nil {
			instancesListMap["region_id"] = instance.RegionId
		}

		if instance.RegionDesc != nil {
			instancesListMap["region_desc"] = instance.RegionDesc
		}

		if instance.Zone != nil {
			instancesListMap["zone"] = instance.Zone
		}

		if instance.ZoneId != nil {
			instancesListMap["zone_id"] = instance.ZoneId
		}

		if instance.ZoneDesc != nil {
			instancesListMap["zone_desc"] = instance.ZoneDesc
		}

		if instance.VpcId != nil {
			instancesListMap["vpc_id"] = instance.VpcId
		}

		if instance.SubnetId != nil {
			instancesListMap["subnet_id"] = instance.SubnetId
		}

		if instance.CreateTime != nil {
			instancesListMap["create_time"] = instance.CreateTime
		}

		if instance.ExpireTime != nil {
			instancesListMap["expire_time"] = instance.ExpireTime
		}

		if instance.AccessInfo != nil {
			instancesListMap["access_info"] = instance.AccessInfo
		}

		if instance.PayMode != nil {
			instancesListMap["pay_mode"] = instance.PayMode
		}

		if instance.RenewFlag != nil {
			instancesListMap["renew_flag"] = instance.RenewFlag
		}

		instancesList = append(instancesList, instancesListMap)
	}

	_ = d.Set("instances_list", instancesList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
