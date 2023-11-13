/*
Use this data source to query detailed information of cdb project_security_group

Example Usage

```hcl
data "tencentcloud_cdb_project_security_group" "project_security_group" {
  project_id =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbProjectSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbProjectSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project id.",
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
							Description: "Project id.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time, time format: yyyy-mm-dd hh:mm:sss.",
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
										Description: "Port.",
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
							Description: "Outbound rules.",
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
										Description: "Port or port range.",
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

func dataSourceTencentCloudCdbProjectSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_project_security_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbProjectSecurityGroupByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		totalCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(totalCount))
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

				securityGroupMap["inbound"] = []interface{}{inboundList}
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

				securityGroupMap["outbound"] = []interface{}{outboundList}
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

			ids = append(ids, *securityGroup.ProjectId)
			tmpList = append(tmpList, securityGroupMap)
		}

		_ = d.Set("groups", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
