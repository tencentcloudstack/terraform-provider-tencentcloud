package dbdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomNodeSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomNodeSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DB Custom node ID.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Security group list bound to the DB Custom node.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group ID.",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group name.",
						},
						"security_group_remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group remark.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group creation time.",
						},
						"inbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Inbound rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule action, ACCEPT or DROP.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source/Destination IP or CIDR.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Port range.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol type, e.g., tcp, udp, icmp, ALL.",
									},
									"service_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol port template ID.",
									},
									"address_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP address template ID.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Security group ID.",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule description.",
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
										Description: "Rule action, ACCEPT or DROP.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source/Destination IP or CIDR.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Port range.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol type, e.g., tcp, udp, icmp, ALL.",
									},
									"service_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol port template ID.",
									},
									"address_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP address template ID.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Security group ID.",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule description.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDbdcDbCustomNodeSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_node_security_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	nodeId := ""
	if v, ok := d.GetOk("node_id"); ok {
		nodeId = v.(string)
	}

	var respData []*dbdcv20201029.SecurityGroup
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDBCustomNodeSecurityGroupsById(ctx, nodeId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	groupsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, sg := range respData {
			sgMap := map[string]interface{}{}
			if sg.SecurityGroupId != nil {
				sgMap["security_group_id"] = sg.SecurityGroupId
			}

			if sg.SecurityGroupName != nil {
				sgMap["security_group_name"] = sg.SecurityGroupName
			}

			if sg.SecurityGroupRemark != nil {
				sgMap["security_group_remark"] = sg.SecurityGroupRemark
			}

			if sg.ProjectId != nil {
				sgMap["project_id"] = sg.ProjectId
			}

			if sg.CreateTime != nil {
				sgMap["create_time"] = sg.CreateTime
			}

			if sg.Inbound != nil {
				inboundList := make([]map[string]interface{}, 0, len(sg.Inbound))
				for _, rule := range sg.Inbound {
					ruleMap := map[string]interface{}{}
					if rule.Action != nil {
						ruleMap["action"] = rule.Action
					}
					if rule.CidrIp != nil {
						ruleMap["cidr_ip"] = rule.CidrIp
					}
					if rule.PortRange != nil {
						ruleMap["port_range"] = rule.PortRange
					}
					if rule.IpProtocol != nil {
						ruleMap["ip_protocol"] = rule.IpProtocol
					}
					if rule.ServiceModule != nil {
						ruleMap["service_module"] = rule.ServiceModule
					}
					if rule.AddressModule != nil {
						ruleMap["address_module"] = rule.AddressModule
					}
					if rule.Id != nil {
						ruleMap["id"] = rule.Id
					}
					if rule.Desc != nil {
						ruleMap["desc"] = rule.Desc
					}
					inboundList = append(inboundList, ruleMap)
				}
				sgMap["inbound"] = inboundList
			}

			if sg.Outbound != nil {
				outboundList := make([]map[string]interface{}, 0, len(sg.Outbound))
				for _, rule := range sg.Outbound {
					ruleMap := map[string]interface{}{}
					if rule.Action != nil {
						ruleMap["action"] = rule.Action
					}
					if rule.CidrIp != nil {
						ruleMap["cidr_ip"] = rule.CidrIp
					}
					if rule.PortRange != nil {
						ruleMap["port_range"] = rule.PortRange
					}
					if rule.IpProtocol != nil {
						ruleMap["ip_protocol"] = rule.IpProtocol
					}
					if rule.ServiceModule != nil {
						ruleMap["service_module"] = rule.ServiceModule
					}
					if rule.AddressModule != nil {
						ruleMap["address_module"] = rule.AddressModule
					}
					if rule.Id != nil {
						ruleMap["id"] = rule.Id
					}
					if rule.Desc != nil {
						ruleMap["desc"] = rule.Desc
					}
					outboundList = append(outboundList, ruleMap)
				}
				sgMap["outbound"] = outboundList
			}

			groupsList = append(groupsList, sgMap)
		}

		_ = d.Set("groups", groupsList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
