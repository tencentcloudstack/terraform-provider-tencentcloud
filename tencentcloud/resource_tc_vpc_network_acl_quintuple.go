/*
Provides a resource to create a vpc network_acl_quintuple

Example Usage

```hcl
resource "tencentcloud_vpc_network_acl_quintuple" "network_acl_quintuple" {
  network_acl_id = ""
  network_acl_quintuple_set {
		ingress {
			protocol = ""
			description = ""
			source_port = ""
			source_cidr = ""
			destination_port = ""
			destination_cidr = ""
			action = ""
			network_acl_quintuple_entry_id = ""
			priority =
			create_time = ""
			network_acl_direction = ""
		}
		egress {
			protocol = ""
			description = ""
			source_port = ""
			source_cidr = ""
			destination_port = ""
			destination_cidr = ""
			action = ""
			network_acl_quintuple_entry_id = ""
			priority =
			create_time = ""
			network_acl_direction = ""
		}

  }
}
```

Import

vpc network_acl_quintuple can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_network_acl_quintuple.network_acl_quintuple network_acl_quintuple_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcNetworkAclQuintuple() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcNetworkAclQuintupleCreate,
		Read:   resourceTencentCloudVpcNetworkAclQuintupleRead,
		Update: resourceTencentCloudVpcNetworkAclQuintupleUpdate,
		Delete: resourceTencentCloudVpcNetworkAclQuintupleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_acl_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Network ACL instance ID. For example:acl-12345678.",
			},

			"network_acl_quintuple_set": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Network quintuple ACL rule set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ingress": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Network ACL quintuple inbound rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Protocol, value: TCP,UDP, ICMP, ALL.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description.",
									},
									"source_port": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "source port (all, single port, range). When the protocol is ALL or ICMP, the port cannot be specified.",
									},
									"source_cidr": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "源CIDR。.",
									},
									"destination_port": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Destination port (all, single port, range). When Protocol is ALL or ICMP, Port cannot be specified.",
									},
									"destination_cidr": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Destination CIDR.",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action, ACCEPT or DROP.",
									},
									"network_acl_quintuple_entry_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Unique ID of a network ACL entry.",
									},
									"priority": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Priority, starting from 1.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Creation time, used as an output parameter of DescribeNetworkAclQuintupleEntries.",
									},
									"network_acl_direction": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Direction, INGRESS or EGRESS, is used as an output parameter of DescribeNetworkAclQuintupleEntries.",
									},
								},
							},
						},
						"egress": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Network ACL quintuple outbound rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Protocol, value: TCP,UDP, ICMP, ALL.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description.",
									},
									"source_port": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source port (all, single port, range). When Protocol is ALL or ICMP, Port cannot be specified.",
									},
									"source_cidr": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source CIDR.",
									},
									"destination_port": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Destination port (all, single port, range). When Protocol is ALL or ICMP, Port cannot be specified.",
									},
									"destination_cidr": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Destination CIDR.",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action, ACCEPT or DROP.",
									},
									"network_acl_quintuple_entry_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Unique ID of a network ACL entry.",
									},
									"priority": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Priority, starting from 1.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Creation time, used as an output parameter of DescribeNetworkAclQuintupleEntries.",
									},
									"network_acl_direction": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Direction, INGRESS or EGRESS, is used as an output parameter of DescribeNetworkAclQuintupleEntries.",
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

func resourceTencentCloudVpcNetworkAclQuintupleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_network_acl_quintuple.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = vpc.NewCreateNetworkAclQuintupleEntriesRequest()
		networkAclId string
	)
	if v, ok := d.GetOk("network_acl_id"); ok {
		networkAclId = v.(string)
		request.NetworkAclId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "network_acl_quintuple_set"); ok {
		networkAclQuintupleEntries := vpc.NetworkAclQuintupleEntries{}
		if v, ok := dMap["ingress"]; ok {
			for _, item := range v.([]interface{}) {
				ingressMap := item.(map[string]interface{})
				networkAclQuintupleEntry := vpc.NetworkAclQuintupleEntry{}
				if v, ok := ingressMap["protocol"]; ok {
					networkAclQuintupleEntry.Protocol = helper.String(v.(string))
				}
				if v, ok := ingressMap["description"]; ok {
					networkAclQuintupleEntry.Description = helper.String(v.(string))
				}
				if v, ok := ingressMap["source_port"]; ok {
					networkAclQuintupleEntry.SourcePort = helper.String(v.(string))
				}
				if v, ok := ingressMap["source_cidr"]; ok {
					networkAclQuintupleEntry.SourceCidr = helper.String(v.(string))
				}
				if v, ok := ingressMap["destination_port"]; ok {
					networkAclQuintupleEntry.DestinationPort = helper.String(v.(string))
				}
				if v, ok := ingressMap["destination_cidr"]; ok {
					networkAclQuintupleEntry.DestinationCidr = helper.String(v.(string))
				}
				if v, ok := ingressMap["action"]; ok {
					networkAclQuintupleEntry.Action = helper.String(v.(string))
				}
				if v, ok := ingressMap["network_acl_quintuple_entry_id"]; ok {
					networkAclQuintupleEntry.NetworkAclQuintupleEntryId = helper.String(v.(string))
				}
				if v, ok := ingressMap["priority"]; ok {
					networkAclQuintupleEntry.Priority = helper.IntInt64(v.(int))
				}
				if v, ok := ingressMap["create_time"]; ok {
					networkAclQuintupleEntry.CreateTime = helper.String(v.(string))
				}
				if v, ok := ingressMap["network_acl_direction"]; ok {
					networkAclQuintupleEntry.NetworkAclDirection = helper.String(v.(string))
				}
				networkAclQuintupleEntries.Ingress = append(networkAclQuintupleEntries.Ingress, &networkAclQuintupleEntry)
			}
		}
		if v, ok := dMap["egress"]; ok {
			for _, item := range v.([]interface{}) {
				egressMap := item.(map[string]interface{})
				networkAclQuintupleEntry := vpc.NetworkAclQuintupleEntry{}
				if v, ok := egressMap["protocol"]; ok {
					networkAclQuintupleEntry.Protocol = helper.String(v.(string))
				}
				if v, ok := egressMap["description"]; ok {
					networkAclQuintupleEntry.Description = helper.String(v.(string))
				}
				if v, ok := egressMap["source_port"]; ok {
					networkAclQuintupleEntry.SourcePort = helper.String(v.(string))
				}
				if v, ok := egressMap["source_cidr"]; ok {
					networkAclQuintupleEntry.SourceCidr = helper.String(v.(string))
				}
				if v, ok := egressMap["destination_port"]; ok {
					networkAclQuintupleEntry.DestinationPort = helper.String(v.(string))
				}
				if v, ok := egressMap["destination_cidr"]; ok {
					networkAclQuintupleEntry.DestinationCidr = helper.String(v.(string))
				}
				if v, ok := egressMap["action"]; ok {
					networkAclQuintupleEntry.Action = helper.String(v.(string))
				}
				if v, ok := egressMap["network_acl_quintuple_entry_id"]; ok {
					networkAclQuintupleEntry.NetworkAclQuintupleEntryId = helper.String(v.(string))
				}
				if v, ok := egressMap["priority"]; ok {
					networkAclQuintupleEntry.Priority = helper.IntInt64(v.(int))
				}
				if v, ok := egressMap["create_time"]; ok {
					networkAclQuintupleEntry.CreateTime = helper.String(v.(string))
				}
				if v, ok := egressMap["network_acl_direction"]; ok {
					networkAclQuintupleEntry.NetworkAclDirection = helper.String(v.(string))
				}
				networkAclQuintupleEntries.Egress = append(networkAclQuintupleEntries.Egress, &networkAclQuintupleEntry)
			}
		}
		request.NetworkAclQuintupleSet = &networkAclQuintupleEntries
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateNetworkAclQuintupleEntries(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc networkAclQuintuple failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(networkAclId)

	return resourceTencentCloudVpcNetworkAclQuintupleRead(d, meta)
}

func resourceTencentCloudVpcNetworkAclQuintupleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_network_acl_quintuple.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	networkAclId := d.Id()

	networkAclQuintuples, err := service.DescribeVpcNetworkAclQuintupleById(ctx, networkAclId)
	if err != nil {
		return err
	}

	if networkAclQuintuples == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcNetworkAclQuintuple` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("network_acl_id", networkAclId)

	networkAclQuintupleSetMap := map[string]interface{}{}
	ingressList := []interface{}{}
	egressList := []interface{}{}
	for _, quintuple := range networkAclQuintuples {
		if *quintuple.NetworkAclDirection == "INGRESS" {
			ingress := quintuple
			ingressMap := map[string]interface{}{}

			if ingress.Protocol != nil {
				ingressMap["protocol"] = ingress.Protocol
			}

			if ingress.Description != nil {
				ingressMap["description"] = ingress.Description
			}

			if ingress.SourcePort != nil {
				ingressMap["source_port"] = ingress.SourcePort
			}

			if ingress.SourceCidr != nil {
				ingressMap["source_cidr"] = ingress.SourceCidr
			}

			if ingress.DestinationPort != nil {
				ingressMap["destination_port"] = ingress.DestinationPort
			}

			if ingress.DestinationCidr != nil {
				ingressMap["destination_cidr"] = ingress.DestinationCidr
			}

			if ingress.Action != nil {
				ingressMap["action"] = ingress.Action
			}

			if ingress.NetworkAclQuintupleEntryId != nil {
				ingressMap["network_acl_quintuple_entry_id"] = ingress.NetworkAclQuintupleEntryId
			}

			if ingress.Priority != nil {
				ingressMap["priority"] = ingress.Priority
			}

			if ingress.CreateTime != nil {
				ingressMap["create_time"] = ingress.CreateTime
			}

			if ingress.NetworkAclDirection != nil {
				ingressMap["network_acl_direction"] = ingress.NetworkAclDirection
			}
			ingressList = append(ingressList, ingressMap)
		}

		if *quintuple.NetworkAclDirection == "EGRESS" {
			egress := quintuple
			egressMap := map[string]interface{}{}

			if egress.Protocol != nil {
				egressMap["protocol"] = egress.Protocol
			}

			if egress.Description != nil {
				egressMap["description"] = egress.Description
			}

			if egress.SourcePort != nil {
				egressMap["source_port"] = egress.SourcePort
			}

			if egress.SourceCidr != nil {
				egressMap["source_cidr"] = egress.SourceCidr
			}

			if egress.DestinationPort != nil {
				egressMap["destination_port"] = egress.DestinationPort
			}

			if egress.DestinationCidr != nil {
				egressMap["destination_cidr"] = egress.DestinationCidr
			}

			if egress.Action != nil {
				egressMap["action"] = egress.Action
			}

			if egress.NetworkAclQuintupleEntryId != nil {
				egressMap["network_acl_quintuple_entry_id"] = egress.NetworkAclQuintupleEntryId
			}

			if egress.Priority != nil {
				egressMap["priority"] = egress.Priority
			}

			if egress.CreateTime != nil {
				egressMap["create_time"] = egress.CreateTime
			}

			if egress.NetworkAclDirection != nil {
				egressMap["network_acl_direction"] = egress.NetworkAclDirection
			}
			egressList = append(egressList, egressMap)
		}
	}

	networkAclQuintupleSetMap["ingress"] = ingressList
	networkAclQuintupleSetMap["egress"] = egressList

	_ = d.Set("network_acl_quintuple_set", []interface{}{networkAclQuintupleSetMap})

	return nil
}

func resourceTencentCloudVpcNetworkAclQuintupleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_network_acl_quintuple.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewModifyNetworkAclQuintupleEntriesRequest()

	networkAclId := d.Id()

	request.NetworkAclId = &networkAclId

	if d.HasChange("network_acl_quintuple_set") {
		if dMap, ok := helper.InterfacesHeadMap(d, "network_acl_quintuple_set"); ok {
			networkAclQuintupleEntries := vpc.NetworkAclQuintupleEntries{}
			if v, ok := dMap["ingress"]; ok {
				for _, item := range v.([]interface{}) {
					ingressMap := item.(map[string]interface{})
					networkAclQuintupleEntry := vpc.NetworkAclQuintupleEntry{}
					if v, ok := ingressMap["protocol"]; ok {
						networkAclQuintupleEntry.Protocol = helper.String(v.(string))
					}
					if v, ok := ingressMap["description"]; ok {
						networkAclQuintupleEntry.Description = helper.String(v.(string))
					}
					if v, ok := ingressMap["source_port"]; ok {
						networkAclQuintupleEntry.SourcePort = helper.String(v.(string))
					}
					if v, ok := ingressMap["source_cidr"]; ok {
						networkAclQuintupleEntry.SourceCidr = helper.String(v.(string))
					}
					if v, ok := ingressMap["destination_port"]; ok {
						networkAclQuintupleEntry.DestinationPort = helper.String(v.(string))
					}
					if v, ok := ingressMap["destination_cidr"]; ok {
						networkAclQuintupleEntry.DestinationCidr = helper.String(v.(string))
					}
					if v, ok := ingressMap["action"]; ok {
						networkAclQuintupleEntry.Action = helper.String(v.(string))
					}
					if v, ok := ingressMap["network_acl_quintuple_entry_id"]; ok {
						networkAclQuintupleEntry.NetworkAclQuintupleEntryId = helper.String(v.(string))
					}
					if v, ok := ingressMap["priority"]; ok {
						networkAclQuintupleEntry.Priority = helper.IntInt64(v.(int))
					}
					if v, ok := ingressMap["create_time"]; ok {
						networkAclQuintupleEntry.CreateTime = helper.String(v.(string))
					}
					if v, ok := ingressMap["network_acl_direction"]; ok {
						networkAclQuintupleEntry.NetworkAclDirection = helper.String(v.(string))
					}
					networkAclQuintupleEntries.Ingress = append(networkAclQuintupleEntries.Ingress, &networkAclQuintupleEntry)
				}
			}
			if v, ok := dMap["egress"]; ok {
				for _, item := range v.([]interface{}) {
					egressMap := item.(map[string]interface{})
					networkAclQuintupleEntry := vpc.NetworkAclQuintupleEntry{}
					if v, ok := egressMap["protocol"]; ok {
						networkAclQuintupleEntry.Protocol = helper.String(v.(string))
					}
					if v, ok := egressMap["description"]; ok {
						networkAclQuintupleEntry.Description = helper.String(v.(string))
					}
					if v, ok := egressMap["source_port"]; ok {
						networkAclQuintupleEntry.SourcePort = helper.String(v.(string))
					}
					if v, ok := egressMap["source_cidr"]; ok {
						networkAclQuintupleEntry.SourceCidr = helper.String(v.(string))
					}
					if v, ok := egressMap["destination_port"]; ok {
						networkAclQuintupleEntry.DestinationPort = helper.String(v.(string))
					}
					if v, ok := egressMap["destination_cidr"]; ok {
						networkAclQuintupleEntry.DestinationCidr = helper.String(v.(string))
					}
					if v, ok := egressMap["action"]; ok {
						networkAclQuintupleEntry.Action = helper.String(v.(string))
					}
					if v, ok := egressMap["network_acl_quintuple_entry_id"]; ok {
						networkAclQuintupleEntry.NetworkAclQuintupleEntryId = helper.String(v.(string))
					}
					if v, ok := egressMap["priority"]; ok {
						networkAclQuintupleEntry.Priority = helper.IntInt64(v.(int))
					}
					if v, ok := egressMap["create_time"]; ok {
						networkAclQuintupleEntry.CreateTime = helper.String(v.(string))
					}
					if v, ok := egressMap["network_acl_direction"]; ok {
						networkAclQuintupleEntry.NetworkAclDirection = helper.String(v.(string))
					}
					networkAclQuintupleEntries.Egress = append(networkAclQuintupleEntries.Egress, &networkAclQuintupleEntry)
				}
			}
			request.NetworkAclQuintupleSet = &networkAclQuintupleEntries
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyNetworkAclQuintupleEntries(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc networkAclQuintuple failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcNetworkAclQuintupleRead(d, meta)
}

func resourceTencentCloudVpcNetworkAclQuintupleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_network_acl_quintuple.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	networkAclId := d.Id()

	if err := service.DeleteVpcNetworkAclQuintupleById(ctx, networkAclId); err != nil {
		return err
	}

	return nil
}
