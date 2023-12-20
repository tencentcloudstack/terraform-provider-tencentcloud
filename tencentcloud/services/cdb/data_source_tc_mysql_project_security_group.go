package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlProjectSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlProjectSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "project id.",
			},

			"groups": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Security group details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "project id.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time, time format: yyyy-mm-dd hh:mm:sss.",
						},
						"inbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "inbound rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy, ACCEPT or DROPs.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source IP or IP range, such as 192.168.0.0/16.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "port.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network protocol, support UDP, TCP, etc.",
									},
									"dir": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The direction defined by the rule, the inbound rule is INPUT.",
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
							Description: "outbound rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy, ACCEPT or DROP.",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Destination IP or IP segment, such as 172.16.0.0/12.",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "port or port range.",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network protocol, support UDP, TCP, etc.",
									},
									"dir": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The direction defined by the rule, the inbound rule is OUTPUT.",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule description.",
									},
								},
							},
						},
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

func dataSourceTencentCloudMysqlProjectSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_project_security_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var groups []*cdb.SecurityGroup
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlProjectSecurityGroupByFilter(ctx, paramMap)
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

					if inbound.Dir != nil {
						inboundMap["dir"] = inbound.Dir
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

					if outbound.Dir != nil {
						outboundMap["dir"] = outbound.Dir
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
