package rum

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRumTawInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumTawInstanceRead,
		Schema: map[string]*schema.Schema{
			"charge_statuses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "Billing status.",
			},

			"charge_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "Billing type.",
			},

			"area_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "Region ID.",
			},

			"instance_statuses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "Instance status (`1`: creating; `2`: running; `3`: exceptional; `4`: restarting; `5`: stopping; `6`: stopped; `7`: terminating; `8`: terminated).",
			},

			"instance_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Instance ID.",
			},

			"instance_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status (`1` = creating, `2` = running, `3` = exception, `4` = restarting, `5` = stopping, `6` = stopped, `7` = deleted).",
						},
						"area_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Area ID.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag List.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag Value.",
									},
								},
							},
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"cluster_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"instance_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance Desc.",
						},
						"charge_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing status (`1` = in use, `2` = expired, `3` = destroyed, `4` = assigning, `5` = failed).",
						},
						"charge_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing type (`1` = free version, `2` = prepaid, `3` = postpaid).",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"data_retention_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data retention time (days).",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
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

func dataSourceTencentCloudRumTawInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_rum_taw_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("charge_statuses"); ok {
		chargeStatusSet := []*int64{}
		charge_statusSet := v.(*schema.Set).List()
		for i := range charge_statusSet {
			charge_status := charge_statusSet[i].(int)
			chargeStatusSet = append(chargeStatusSet, helper.IntInt64(charge_status))
		}
		paramMap["charge_statuses"] = chargeStatusSet
	}

	if v, ok := d.GetOk("charge_types"); ok {
		chargeTypesSet := []*int64{}
		charge_typesSet := v.(*schema.Set).List()
		for i := range charge_typesSet {
			charge_type := charge_typesSet[i].(int)
			chargeTypesSet = append(chargeTypesSet, helper.IntInt64(charge_type))
		}
		paramMap["charge_types"] = chargeTypesSet
	}

	if v, ok := d.GetOk("area_ids"); ok {
		areaIdsSet := []*int64{}
		area_idsSet := v.(*schema.Set).List()
		for i := range area_idsSet {
			area_id := area_idsSet[i].(int)
			areaIdsSet = append(areaIdsSet, helper.IntInt64(area_id))
		}
		paramMap["area_ids"] = areaIdsSet
	}

	if v, ok := d.GetOk("instance_statuses"); ok {
		instanceStatusSet := []*int64{}
		instance_statusSet := v.(*schema.Set).List()
		for i := range instance_statusSet {
			instance_status := instance_statusSet[i].(int)
			instanceStatusSet = append(instanceStatusSet, helper.IntInt64(instance_status))
		}
		paramMap["instance_statuses"] = instanceStatusSet
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := []*string{}
		instance_idsSet := v.(*schema.Set).List()
		for i := range instance_idsSet {
			instance_id := instance_idsSet[i].(string)
			instanceIdsSet = append(instanceIdsSet, &instance_id)
		}
		paramMap["instance_ids"] = instanceIdsSet
	}

	rumService := RumService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceSet []*rum.RumInstanceInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := rumService.DescribeRumTawInstanceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceSet = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Rum instanceSet failed, reason:%+v", logId, err)
		return err
	}

	instanceSetList := []interface{}{}
	ids := make([]string, 0, len(instanceSet))
	if instanceSet != nil {
		for _, instanceSet := range instanceSet {
			instanceSetMap := map[string]interface{}{}
			if instanceSet.InstanceStatus != nil {
				instanceSetMap["instance_status"] = instanceSet.InstanceStatus
			}
			if instanceSet.AreaId != nil {
				instanceSetMap["area_id"] = instanceSet.AreaId
			}
			if instanceSet.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range instanceSet.Tags {
					tagsMap := map[string]interface{}{}
					if tags.Key != nil {
						tagsMap["key"] = tags.Key
					}
					if tags.Value != nil {
						tagsMap["value"] = tags.Value
					}

					tagsList = append(tagsList, tagsMap)
				}
				instanceSetMap["tags"] = tagsList
			}
			if instanceSet.InstanceId != nil {
				instanceSetMap["instance_id"] = instanceSet.InstanceId
			}
			if instanceSet.ClusterId != nil {
				instanceSetMap["cluster_id"] = instanceSet.ClusterId
			}
			if instanceSet.InstanceDesc != nil {
				instanceSetMap["instance_desc"] = instanceSet.InstanceDesc
			}
			if instanceSet.ChargeStatus != nil {
				instanceSetMap["charge_status"] = instanceSet.ChargeStatus
			}
			if instanceSet.ChargeType != nil {
				instanceSetMap["charge_type"] = instanceSet.ChargeType
			}
			if instanceSet.UpdatedAt != nil {
				instanceSetMap["updated_at"] = instanceSet.UpdatedAt
			}
			if instanceSet.DataRetentionDays != nil {
				instanceSetMap["data_retention_days"] = instanceSet.DataRetentionDays
			}
			if instanceSet.InstanceName != nil {
				instanceSetMap["instance_name"] = instanceSet.InstanceName
			}
			if instanceSet.CreatedAt != nil {
				instanceSetMap["created_at"] = instanceSet.CreatedAt
			}

			instanceSetList = append(instanceSetList, instanceSetMap)
			ids = append(ids, *instanceSet.InstanceId)
		}
		_ = d.Set("instance_set", instanceSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceSetList); e != nil {
			return e
		}
	}

	return nil
}
