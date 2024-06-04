package cvm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEipsRead,

		Schema: map[string]*schema.Schema{
			"eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the EIP to be queried.",
			},
			"eip_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the EIP to be queried.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The elastic ip address.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of EIP.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"eip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of EIP. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the EIP.",
						},
						"eip_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the EIP.",
						},
						"eip_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the EIP.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The EIP current status.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The elastic ip address.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance id to bind with the EIP.",
						},
						"eni_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The eni id to bind with the EIP.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the EIP.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the EIP.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudEipsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eips.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	tagService := svctag.NewTagService(client)
	region := client.Region

	filter := make(map[string][]string)
	if v, ok := d.GetOk("eip_id"); ok {
		filter["address-id"] = []string{v.(string)}
	}
	if v, ok := d.GetOk("eip_name"); ok {
		filter["address-name"] = []string{v.(string)}
	}
	if v, ok := d.GetOk("public_ip"); ok {
		filter["public-ip"] = []string{v.(string)}
	}

	tags := helper.GetTags(d, "tags")

	var eips []*vpc.Address
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eips, errRet = vpcService.DescribeEipByFilter(ctx, filter)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	eipList := make([]map[string]interface{}, 0, len(eips))
	ids := make([]string, 0, len(eips))

EIP_LOOP:
	for _, eip := range eips {
		respTags, err := tagService.DescribeResourceTags(ctx, svcvpc.VPC_SERVICE_TYPE, svcvpc.EIP_RESOURCE_TYPE, region, *eip.AddressId)
		if err != nil {
			log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
			return err
		}

		for k, v := range tags {
			if respTags[k] != v {
				continue EIP_LOOP
			}
		}

		mapping := map[string]interface{}{
			"eip_id":      eip.AddressId,
			"eip_name":    eip.AddressName,
			"eip_type":    eip.AddressType,
			"status":      eip.AddressStatus,
			"public_ip":   eip.AddressIp,
			"instance_id": eip.InstanceId,
			"eni_id":      eip.NetworkInterfaceId,
			"create_time": eip.CreatedTime,
			"tags":        respTags,
		}

		eipList = append(eipList, mapping)
		ids = append(ids, *eip.AddressId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("eip_list", eipList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set eip list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), eipList); err != nil {
			return err
		}
	}
	return nil
}
