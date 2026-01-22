package privatedns

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPrivateDnsPrivateZoneList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPrivateDnsPrivateZoneListRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "values.",
						},
					},
				},
			},
			"private_zone_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Private Zone Set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PrivateZone ID.",
						},
						"owner_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Owner Uin.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"record_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record count.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark.",
						},
						"vpc_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Vpc list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Vpc Id.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private domain bound VPC status, not associated with vpc: SUSPEND, associated with VPC: ENABLED, associated with VPC failed: FAILED.",
						},
						"dns_forward_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name recursive resolution status: enabled: ENABLED, disabled, DISABLED.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "tag value.",
									},
								},
							},
						},
						"account_vpc_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "VPC list of bound associated accounts.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "uin.",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Vpc Id.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region.",
									},
								},
							},
						},
						"is_custom_tld": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Custom TLD.",
						},
						"cname_speedup_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CNAME acceleration status: enabled: ENABLED, off, DISABLED.",
						},
						"forward_rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding rule name.",
						},
						"forward_rule_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding rule type: from cloud to cloud, DOWN; From cloud to cloud, UP, currently only supports DOWN.",
						},
						"forward_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarded address.",
						},
						"end_point_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End point name.",
						},
						"deleted_vpc_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of deleted VPCs.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Vpc Id.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region.",
									},
								},
							},
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

func dataSourceTencentCloudPrivateDnsPrivateZoneListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_private_dns_private_zone_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = PrivateDnsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		privateZoneSet []*privatedns.PrivateZone
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*privatedns.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := privatedns.Filter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	privateZoneSet, err := service.DescribePrivatednsPrivateZoneListByFilter(ctx, paramMap)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(privateZoneSet))
	tmpList := make([]map[string]interface{}, 0, len(privateZoneSet))

	if privateZoneSet != nil {
		for _, privateZone := range privateZoneSet {
			privateZoneMap := map[string]interface{}{}

			if privateZone.ZoneId != nil {
				privateZoneMap["zone_id"] = privateZone.ZoneId
			}

			if privateZone.OwnerUin != nil {
				privateZoneMap["owner_uin"] = privateZone.OwnerUin
			}

			if privateZone.Domain != nil {
				privateZoneMap["domain"] = privateZone.Domain
			}

			if privateZone.CreatedOn != nil {
				privateZoneMap["created_on"] = privateZone.CreatedOn
			}

			if privateZone.UpdatedOn != nil {
				privateZoneMap["updated_on"] = privateZone.UpdatedOn
			}

			if privateZone.RecordCount != nil {
				privateZoneMap["record_count"] = privateZone.RecordCount
			}

			if privateZone.Remark != nil {
				privateZoneMap["remark"] = privateZone.Remark
			}

			if privateZone.VpcSet != nil {
				vpcSetList := []interface{}{}
				for _, vpcSet := range privateZone.VpcSet {
					vpcSetMap := map[string]interface{}{}

					if vpcSet.UniqVpcId != nil {
						vpcSetMap["uniq_vpc_id"] = vpcSet.UniqVpcId
					}

					if vpcSet.Region != nil {
						vpcSetMap["region"] = vpcSet.Region
					}

					vpcSetList = append(vpcSetList, vpcSetMap)
				}

				privateZoneMap["vpc_set"] = vpcSetList
			}

			if privateZone.Status != nil {
				privateZoneMap["status"] = privateZone.Status
			}

			if privateZone.DnsForwardStatus != nil {
				privateZoneMap["dns_forward_status"] = privateZone.DnsForwardStatus
			}

			if privateZone.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range privateZone.Tags {
					tagsMap := map[string]interface{}{}

					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				privateZoneMap["tags"] = tagsList
			}

			if privateZone.AccountVpcSet != nil {
				accountVpcSetList := []interface{}{}
				for _, accountVpcSet := range privateZone.AccountVpcSet {
					accountVpcSetMap := map[string]interface{}{}

					if accountVpcSet.Uin != nil {
						accountVpcSetMap["uin"] = accountVpcSet.Uin
					}

					if accountVpcSet.UniqVpcId != nil {
						accountVpcSetMap["uniq_vpc_id"] = accountVpcSet.UniqVpcId
					}

					if accountVpcSet.Region != nil {
						accountVpcSetMap["region"] = accountVpcSet.Region
					}

					accountVpcSetList = append(accountVpcSetList, accountVpcSetMap)
				}

				privateZoneMap["account_vpc_set"] = accountVpcSetList
			}

			if privateZone.IsCustomTld != nil {
				privateZoneMap["is_custom_tld"] = privateZone.IsCustomTld
			}

			if privateZone.CnameSpeedupStatus != nil {
				privateZoneMap["cname_speedup_status"] = privateZone.CnameSpeedupStatus
			}

			if privateZone.ForwardRuleName != nil {
				privateZoneMap["forward_rule_name"] = privateZone.ForwardRuleName
			}

			if privateZone.ForwardRuleType != nil {
				privateZoneMap["forward_rule_type"] = privateZone.ForwardRuleType
			}

			if privateZone.ForwardAddress != nil {
				privateZoneMap["forward_address"] = privateZone.ForwardAddress
			}

			if privateZone.EndPointName != nil {
				privateZoneMap["end_point_name"] = privateZone.EndPointName
			}

			if privateZone.DeletedVpcSet != nil {
				deletedVpcSetList := []interface{}{}
				for _, deletedVpcSet := range privateZone.DeletedVpcSet {
					deletedVpcSetMap := map[string]interface{}{}

					if deletedVpcSet.UniqVpcId != nil {
						deletedVpcSetMap["uniq_vpc_id"] = deletedVpcSet.UniqVpcId
					}

					if deletedVpcSet.Region != nil {
						deletedVpcSetMap["region"] = deletedVpcSet.Region
					}

					deletedVpcSetList = append(deletedVpcSetList, deletedVpcSetMap)
				}

				privateZoneMap["deleted_vpc_set"] = deletedVpcSetList
			}

			ids = append(ids, *privateZone.ZoneId)
			tmpList = append(tmpList, privateZoneMap)
		}

		_ = d.Set("private_zone_set", tmpList)
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
