/*
Use this data source to query query ENIs.

Example Usage

```hcl
data "tencentcloud_enis" "name" {
  name = "test eni"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func dataSourceTencentCloudEnis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEnisRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"vpc_id", "subnet_id", "instance_id", "security_group", "name", "description", "ipv4", "tags"},
				Description:   "ID of the ENIs to be queried. Conflict with `vpc_id`,`subnet_id`,`instance_id`,`security_group`,`name`,`ipv4` and `tags`.",
			},
			"vpc_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "ID of the vpc to be queried. Conflict with `ids`.",
			},
			"subnet_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "ID of the subnet within this vpc to be queried. Conflict with `ids`.",
			},
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "ID of the instance which bind the ENI. Conflict with `ids`.",
			},
			"security_group": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "A set of security group IDs which bind the ENI. Conflict with `ids`.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Name of the ENI to be queried. Conflict with `ids`.",
			},
			"description": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Description of the ENI. Conflict with `ids`.",
			},
			"ipv4": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Intranet IP of the ENI. Conflict with `ids`.",
			},
			"tags": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Tags of the ENI. Conflict with `ids`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"enis": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of ENIs. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the ENI.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the ENI.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the ENI.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the vpc.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet within this vpc.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "A set of security group IDs which bind the ENI.",
						},
						"primary": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the IP is primary.",
						},
						"mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC address.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "States of the ENI.",
						},
						"ipv4s": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A set of intranet IPv4s.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Intranet IP.",
									},
									"primary": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the IP is primary.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the IP.",
									},
								},
							},
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the instance which bind the ENI.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the ENI.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the ENI.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudEnisRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_enis.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		ids      []string
		vpcId    *string
		subnetId *string
		cvmId    *string
		sgId     *string
		name     *string
		desc     *string
		ipv4     *string
	)

	if raw, ok := d.GetOk("ids"); ok {
		ids = expandStringList(raw.(*schema.Set).List())
	}

	if raw, ok := d.GetOk("vpc_id"); ok {
		vpcId = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("subnet_id"); ok {
		subnetId = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("instance_id"); ok {
		cvmId = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("security_group"); ok {
		sgId = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("description"); ok {
		desc = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("ipv4"); ok {
		ipv4 = stringToPointer(raw.(string))
	}
	tags := getTags(d, "tags")

	var (
		respEnis []*vpc.NetworkInterface
		err      error
	)

	if len(ids) > 0 {
		respEnis, err = service.DescribeEniById(ctx, ids)
	} else {
		respEnis, err = service.DescribeEniByFilters(ctx, vpcId, subnetId, cvmId, sgId, name, desc, ipv4, tags)
	}

	if err != nil {
		return err
	}

	enis := make([]map[string]interface{}, 0, len(respEnis))
	eniIds := make([]string, 0, len(respEnis))

	for _, eni := range respEnis {
		if nilFields := CheckNil(eni, map[string]string{
			"NetworkInterfaceId":          "id",
			"NetworkInterfaceName":        "name",
			"NetworkInterfaceDescription": "description",
			"VpcId":                       "vpc id",
			"SubnetId":                    "subnet id",
			"MacAddress":                  "mac address",
			"State":                       "state",
			"CreatedTime":                 "create time",
			"Primary":                     "primary",
		}); len(nilFields) > 0 {
			return fmt.Errorf("eni %v are nil", nilFields)
		}

		ipv4s := make([]map[string]interface{}, 0, len(eni.PrivateIpAddressSet))
		for _, ipv4 := range eni.PrivateIpAddressSet {
			if nilFields := CheckNil(ipv4, map[string]string{
				"PrivateIpAddress": "ip",
				"Primary":          "primary",
				"Description":      "description",
			}); len(nilFields) > 0 {
				return fmt.Errorf("eni ipv4 %v are nil", nilFields)
			}

			ipv4s = append(ipv4s, map[string]interface{}{
				"ip":          *ipv4.PrivateIpAddress,
				"primary":     *ipv4.Primary,
				"description": *eni.NetworkInterfaceDescription,
			})
		}

		sgs := make([]string, 0, len(eni.GroupSet))
		for _, sg := range eni.GroupSet {
			sgs = append(sgs, *sg)
		}

		respTags := make(map[string]string, len(eni.TagSet))
		for _, tag := range eni.TagSet {
			if tag.Key == nil {
				return errors.New("eni tag key is nil")
			}
			if tag.Value == nil {
				return errors.New("eni tag value is nil")
			}

			respTags[*tag.Key] = *tag.Value
		}

		eniIds = append(eniIds, *eni.NetworkInterfaceId)

		m := map[string]interface{}{
			"id":              *eni.NetworkInterfaceId,
			"name":            *eni.NetworkInterfaceName,
			"description":     *eni.NetworkInterfaceDescription,
			"vpc_id":          *eni.VpcId,
			"subnet_id":       *eni.SubnetId,
			"primary":         *eni.Primary,
			"mac":             *eni.MacAddress,
			"state":           *eni.State,
			"create_time":     *eni.CreatedTime,
			"ipv4s":           ipv4s,
			"security_groups": sgs,
			"tags":            respTags,
		}

		if eni.Attachment != nil {
			if eni.Attachment.InstanceId == nil {
				return errors.New("eni attach instance id is nil")
			}
			m["instance_id"] = *eni.Attachment.InstanceId
		}

		enis = append(enis, m)
	}

	d.Set("enis", enis)
	d.SetId(dataResourceIdsHash(eniIds))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), enis); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]",
				logId, output.(string), err)
			return err
		}
	}

	return nil
}
