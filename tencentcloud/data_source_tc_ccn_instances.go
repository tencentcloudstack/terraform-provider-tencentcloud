/*
Use this data source to query detailed information of CCN instances.

Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

data "tencentcloud_ccn_instances" "id_instances" {
  ccn_id = tencentcloud_ccn.main.id
}

data "tencentcloud_ccn_instances" "name_instances" {
  name = tencentcloud_ccn.main.name
}
```
*/
package tencentcloud

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudCcnInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnInstancesRead,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CCN to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CCN to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of CCN.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ccn_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the CCN.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the CCN.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the CCN.",
						},
						"qos": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service quality of CCN, and the available value include 'PT', 'AU', 'AG'. The default is 'AU'.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "States of instance. The available value include 'ISOLATED'(arrears) and 'AVAILABLE'.",
						},
						"attachment_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information list of instance is attached.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of attached instance network, and available values include VPC, DIRECTCONNECT, BMVPC and VPNGW.",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The region that the instance locates at.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of instance is attached.",
									},
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "States of instance is attached, and available values include PENDING, ACTIVE, EXPIRED, REJECTED, DELETED, FAILED(asynchronous forced disassociation after 2 hours), ATTACHING, DETACHING and DETACHFAILED(asynchronous forced disassociation after 2 hours).",
									},
									"attached_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time of attaching.",
									},
									"cidr_block": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "A network address block of the instance that is attached.",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of resource.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCcnInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ccn_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId = ""
		name  = ""
	)

	if temp, ok := d.GetOk("ccn_id"); ok {
		if tempStr := temp.(string); tempStr != "" {
			ccnId = tempStr
		}
	}

	if temp, ok := d.GetOk("name"); ok {
		if tempStr := temp.(string); tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeCcns(ctx, ccnId, name)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["ccn_id"] = item.ccnId
		infoMap["name"] = item.name
		infoMap["description"] = item.description
		infoMap["qos"] = item.qos
		infoMap["state"] = strings.ToUpper(item.state)
		infoMap["create_time"] = item.createTime
		infoList = append(infoList, infoMap)

		instances, err := service.DescribeCcnAttachedInstances(ctx, item.ccnId)
		if err != nil {
			return err
		}
		attachmentList := make([]interface{}, 0, len(instances))

		for _, instance := range instances {

			instanceMap := map[string]interface{}{
				"instance_type":   instance.instanceType,
				"instance_region": instance.instanceRegion,
				"instance_id":     instance.instanceId,
				"state":           strings.ToUpper(instance.state),
				"attached_time":   instance.attachedTime,
				"cidr_block":      instance.cidrBlock,
			}
			attachmentList = append(attachmentList, instanceMap)

		}

		infoMap["attachment_list"] = attachmentList

	}
	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  ccn instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	m := md5.New()
	_, err = m.Write([]byte("ccn_instances" + ccnId + "_" + name))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
