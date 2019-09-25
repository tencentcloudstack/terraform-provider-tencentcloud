/*
Use this data source to query eip instances.

Example Usage

```hcl
data "tencentcloud_eips" "foo" {
  eip_id = "eip-ry9h95hg"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func dataSourceTencentCloudEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEipsRead,

		Schema: map[string]*schema.Schema{
			"eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the eip to be queried.",
			},
			"eip_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the eip to be queried.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The elastic ip address.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"eip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of eip. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the eip.",
						},
						"eip_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the eip.",
						},
						"eip_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the eip.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The eip current status.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The elastic ip address.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance id to bind with the eip.",
						},
						"eni_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The eni id to bind with the eip.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the eip.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudEipsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eips.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	filter := make(map[string]string)
	if v, ok := d.GetOk("eip_id"); ok {
		filter["address-id"] = v.(string)
	}
	if v, ok := d.GetOk("eip_name"); ok {
		filter["eip-name"] = v.(string)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		filter["public-ip"] = v.(string)
	}

	var eips []*vpc.Address
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eips, errRet = vpcService.DescribeEipByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		return nil
	})
	if err != nil {
		return err
	}

	eipList := make([]map[string]interface{}, 0, len(eips))
	ids := make([]string, 0, len(eips))
	for _, eip := range eips {
		mapping := map[string]interface{}{
			"eip_id":      eip.AddressId,
			"eip_name":    eip.AddressName,
			"eip_type":    eip.AddressType,
			"status":      eip.AddressStatus,
			"public_ip":   eip.AddressIp,
			"instance_id": eip.InstanceId,
			"eni_id":      eip.NetworkInterfaceId,
			"create_time": eip.CreatedTime,
		}
		eipList = append(eipList, mapping)
		ids = append(ids, *eip.AddressId)
	}

	d.SetId(dataResourceIdsHash(ids))
	err = d.Set("eip_list", eipList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set eip list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), eipList); err != nil {
			return err
		}
	}
	return nil
}
