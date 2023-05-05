/*
Use this data source to query detailed information of mariadb securityGroups

Example Usage

```hcl
data "tencentcloud_mariadb_security_groups" "securityGroups" {
  instance_id = "tdsql-4pzs5b67"
  product = "mariadb"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "product name, fixed to mariadb.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "SecurityGroup list.",
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
							Description: "Creation time, time format: `yyyy-mm-dd hh:mm:ss`.",
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
							Description: "Security Group Notes.",
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
										Description: "Policy, ACCEPT or DROP.",
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
										Description: "Network protocols, support `UDP`, `TCP`, etc.",
									},
								},
							},
						},
						"outbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Outbound Rules.",
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
										Description: "Network protocols, support `UDP`, `TCP`, etc.",
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

func dataSourceTencentCloudMariadbSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_security_groups.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var instanceId string
	var product string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
		paramMap["product"] = helper.String(v.(string))
	}

	mariadbService := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var groups []*mariadb.SecurityGroup
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := mariadbService.DescribeMariadbSecurityGroupsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		groups = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Mariadb groups failed, reason:%+v", logId, err)
		return err
	}

	groupList := []interface{}{}
	if groups != nil {
		for _, group := range groups {
			groupMap := map[string]interface{}{}
			if group.ProjectId != nil {
				groupMap["project_id"] = group.ProjectId
			}
			if group.CreateTime != nil {
				groupMap["create_time"] = group.CreateTime
			}
			if group.SecurityGroupId != nil {
				groupMap["security_group_id"] = group.SecurityGroupId
			}
			if group.SecurityGroupName != nil {
				groupMap["security_group_name"] = group.SecurityGroupName
			}
			if group.SecurityGroupRemark != nil {
				groupMap["security_group_remark"] = group.SecurityGroupRemark
			}
			if group.Inbound != nil {
				inboundList := []interface{}{}
				for _, inbound := range group.Inbound {
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

					inboundList = append(inboundList, inboundMap)
				}
				groupMap["inbound"] = inboundList
			}
			if group.Outbound != nil {
				outboundList := []interface{}{}
				for _, outbound := range group.Outbound {
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

					outboundList = append(outboundList, outboundMap)
				}
				groupMap["outbound"] = outboundList
			}

			groupList = append(groupList, groupMap)
		}
		err = d.Set("list", groupList)
		if err != nil {
			log.Printf("[CRITAL]%s provider set instances list fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	d.SetId(instanceId + FILED_SP + product)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), groupList); e != nil {
			return e
		}
	}

	return nil
}
