package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbProjectSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbProjectSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},
			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search Keywords.",
			},
			"groups": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Security Group Details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time, time format: yyyy mm dd hh: mm: ss.",
						},
						"inbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inbound Rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CidrIp.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "PortRange.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ip protocol.",
									},
									"service_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service Module.",
									},
									"address_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AddressModule.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "id.",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description.",
									},
								},
							},
						},
						"outbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Outbound rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cidr Ip.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Port range.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ip protocol.",
									},
									"service_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service module.",
									},
									"address_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Address module.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "id.",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description.",
									},
								},
							},
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security Group ID.",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security Group Name.",
						},
						"security_group_remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security Group Notes.",
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

func dataSourceTencentCloudCynosdbProjectSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_project_security_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		groups  []*cynosdb.SecurityGroup
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbProjectSecurityGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		groups = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groups))
	tmpList := make([]map[string]interface{}, 0, len(groups))

	if groups != nil {
		for _, securityGroup := range groups {
			securityGroupMap := map[string]interface{}{}

			if securityGroup.ProjectId != nil {
				securityGroupMap["project_id"] = securityGroup.ProjectId
			}

			if securityGroup.CreateTime != nil {
				securityGroupMap["create_time"] = securityGroup.CreateTime
			}

			if securityGroup.Inbound != nil {
				inboundList := []interface{}{}
				for _, inbound := range securityGroup.Inbound {
					inboundMap := map[string]interface{}{}

					if inbound.Action != nil {
						inboundMap["action"] = inbound.Action
					}

					if inbound.CidrIp != nil {
						inboundMap["cidr_ip"] = inbound.CidrIp
					}

					if inbound.PortRange != nil {
						inboundMap["port_range"] = inbound.PortRange
					}

					if inbound.IpProtocol != nil {
						inboundMap["ip_protocol"] = inbound.IpProtocol
					}

					if inbound.ServiceModule != nil {
						inboundMap["service_module"] = inbound.ServiceModule
					}

					if inbound.AddressModule != nil {
						inboundMap["address_module"] = inbound.AddressModule
					}

					if inbound.Id != nil {
						inboundMap["id"] = inbound.Id
					}

					if inbound.Desc != nil {
						inboundMap["desc"] = inbound.Desc
					}

					inboundList = append(inboundList, inboundMap)
				}

				securityGroupMap["inbound"] = inboundList
			}

			if securityGroup.Outbound != nil {
				outboundList := []interface{}{}
				for _, outbound := range securityGroup.Outbound {
					outboundMap := map[string]interface{}{}

					if outbound.Action != nil {
						outboundMap["action"] = outbound.Action
					}

					if outbound.CidrIp != nil {
						outboundMap["cidr_ip"] = outbound.CidrIp
					}

					if outbound.PortRange != nil {
						outboundMap["port_range"] = outbound.PortRange
					}

					if outbound.IpProtocol != nil {
						outboundMap["ip_protocol"] = outbound.IpProtocol
					}

					if outbound.ServiceModule != nil {
						outboundMap["service_module"] = outbound.ServiceModule
					}

					if outbound.AddressModule != nil {
						outboundMap["address_module"] = outbound.AddressModule
					}

					if outbound.Id != nil {
						outboundMap["id"] = outbound.Id
					}

					if outbound.Desc != nil {
						outboundMap["desc"] = outbound.Desc
					}

					outboundList = append(outboundList, outboundMap)
				}

				securityGroupMap["outbound"] = outboundList
			}

			if securityGroup.SecurityGroupId != nil {
				securityGroupMap["security_group_id"] = securityGroup.SecurityGroupId
			}

			if securityGroup.SecurityGroupName != nil {
				securityGroupMap["security_group_name"] = securityGroup.SecurityGroupName
			}

			if securityGroup.SecurityGroupRemark != nil {
				securityGroupMap["security_group_remark"] = securityGroup.SecurityGroupRemark
			}

			ids = append(ids, *securityGroup.SecurityGroupId)
			tmpList = append(tmpList, securityGroupMap)
		}

		_ = d.Set("groups", tmpList)
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
