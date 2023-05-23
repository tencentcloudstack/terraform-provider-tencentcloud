/*
Use this data source to query detailed information of sqlserver project_security_groups

Example Usage

```hcl
data "tencentcloud_sqlserver_project_security_groups" "project_security_groups" {
  project_id = 0
}
```
*/
package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverProjectSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverProjectSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID, which can be viewed through the console project management.",
			},
			"security_group_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Security group details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "project ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time, time format: yyyy-mm-dd hh:mm:ss.",
						},
						"inbound_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "inbound rules.",
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
										Description: "The direction defined by the rules, OUTPUT-outgoing rules INPUT-inbound rules.",
									},
								},
							},
						},
						"outbound_set": {
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
										Description: "The direction defined by the rules, OUTPUT-outgoing rules INPUT-inbound rules.",
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
							Description: "security group name.",
						},
						"security_group_remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security Group Remarks.",
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

func dataSourceTencentCloudSqlserverProjectSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_project_security_groups.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		projectId string
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
		projectId = strconv.Itoa(v.(int))
	}

	var securityGroupSet []*sqlserver.SecurityGroup

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverProjectSecurityGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		securityGroupSet = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(securityGroupSet))

	if securityGroupSet != nil {
		for _, securityGroup := range securityGroupSet {
			securityGroupMap := map[string]interface{}{}

			if securityGroup.ProjectId != nil {
				securityGroupMap["project_id"] = securityGroup.ProjectId
			}

			if securityGroup.CreateTime != nil {
				securityGroupMap["create_time"] = securityGroup.CreateTime
			}

			if securityGroup.InboundSet != nil {
				inboundSetList := []interface{}{}
				for _, inboundSet := range securityGroup.InboundSet {
					inboundSetMap := map[string]interface{}{}

					if inboundSet.Action != nil {
						inboundSetMap["action"] = inboundSet.Action
					}

					if inboundSet.CidrIp != nil {
						inboundSetMap["cidr_ip"] = inboundSet.CidrIp
					}

					if inboundSet.PortRange != nil {
						inboundSetMap["port_range"] = inboundSet.PortRange
					}

					if inboundSet.IpProtocol != nil {
						inboundSetMap["ip_protocol"] = inboundSet.IpProtocol
					}

					if inboundSet.Dir != nil {
						inboundSetMap["dir"] = inboundSet.Dir
					}

					inboundSetList = append(inboundSetList, inboundSetMap)
				}

				securityGroupMap["inbound_set"] = inboundSetList
			}

			if securityGroup.OutboundSet != nil {
				outboundSetList := []interface{}{}
				for _, outboundSet := range securityGroup.OutboundSet {
					outboundSetMap := map[string]interface{}{}

					if outboundSet.Action != nil {
						outboundSetMap["action"] = outboundSet.Action
					}

					if outboundSet.CidrIp != nil {
						outboundSetMap["cidr_ip"] = outboundSet.CidrIp
					}

					if outboundSet.PortRange != nil {
						outboundSetMap["port_range"] = outboundSet.PortRange
					}

					if outboundSet.IpProtocol != nil {
						outboundSetMap["ip_protocol"] = outboundSet.IpProtocol
					}

					if outboundSet.Dir != nil {
						outboundSetMap["dir"] = outboundSet.Dir
					}

					outboundSetList = append(outboundSetList, outboundSetMap)
				}

				securityGroupMap["outbound_set"] = outboundSetList
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

			tmpList = append(tmpList, securityGroupMap)
		}

		_ = d.Set("security_group_set", tmpList)
	}

	d.SetId(projectId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
