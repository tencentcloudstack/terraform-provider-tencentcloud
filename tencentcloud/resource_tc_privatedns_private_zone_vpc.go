/*
Provides a resource to create a privatedns private_zone_vpc

Example Usage

```hcl
resource "tencentcloud_privatedns_private_zone_vpc" "private_zone_vpc" {
  zone_id = "zone-xxxxxxx"
  vpc_set {
		uniq_vpc_id = "vpc-xadsafsdasd"
		region = "ap-guangzhou"

  }
  account_vpc_set {
		uniq_vpc_id = "vpc-xadsafsdasd"
		region = "ap-guangzhou"
		uin = "123456789"
		vpc_name = "testname"

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

privatedns private_zone_vpc can be imported using the id, e.g.

```
terraform import tencentcloud_privatedns_private_zone_vpc.private_zone_vpc private_zone_vpc_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPrivatednsPrivateZoneVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPrivatednsPrivateZoneVpcCreate,
		Read:   resourceTencentCloudPrivatednsPrivateZoneVpcRead,
		Delete: resourceTencentCloudPrivatednsPrivateZoneVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "PrivateZone ID.",
			},

			"vpc_set": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "New add vpc info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Uniq Vpc Id.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vpc region.",
						},
					},
				},
			},

			"account_vpc_set": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "New add account vpc info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Uniq Vpc Id.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vpc region.",
						},
						"uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vpc owner uin.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vpc name.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPrivatednsPrivateZoneVpcCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_privatedns_private_zone_vpc.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = privatedns.NewAddSpecifyPrivateZoneVpcRequest()
		response = privatedns.NewAddSpecifyPrivateZoneVpcResponse()
		zoneId   string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			vpcInfo := privatedns.VpcInfo{}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				vpcInfo.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["region"]; ok {
				vpcInfo.Region = helper.String(v.(string))
			}
			request.VpcSet = append(request.VpcSet, &vpcInfo)
		}
	}

	if v, ok := d.GetOk("account_vpc_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			accountVpcInfo := privatedns.AccountVpcInfo{}
			if v, ok := dMap["uniq_vpc_id"]; ok {
				accountVpcInfo.UniqVpcId = helper.String(v.(string))
			}
			if v, ok := dMap["region"]; ok {
				accountVpcInfo.Region = helper.String(v.(string))
			}
			if v, ok := dMap["uin"]; ok {
				accountVpcInfo.Uin = helper.String(v.(string))
			}
			if v, ok := dMap["vpc_name"]; ok {
				accountVpcInfo.VpcName = helper.String(v.(string))
			}
			request.AccountVpcSet = append(request.AccountVpcSet, &accountVpcInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePrivatednsClient().AddSpecifyPrivateZoneVpc(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create privatedns privateZoneVpc failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(zoneId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::privatedns:%s:uin/:zone/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudPrivatednsPrivateZoneVpcRead(d, meta)
}

func resourceTencentCloudPrivatednsPrivateZoneVpcRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_privatedns_private_zone_vpc.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PrivatednsService{client: meta.(*TencentCloudClient).apiV3Conn}

	privateZoneVpcId := d.Id()

	privateZoneVpc, err := service.DescribePrivatednsPrivateZoneVpcById(ctx, zoneId)
	if err != nil {
		return err
	}

	if privateZoneVpc == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PrivatednsPrivateZoneVpc` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if privateZoneVpc.ZoneId != nil {
		_ = d.Set("zone_id", privateZoneVpc.ZoneId)
	}

	if privateZoneVpc.VpcSet != nil {
		vpcSetList := []interface{}{}
		for _, vpcSet := range privateZoneVpc.VpcSet {
			vpcSetMap := map[string]interface{}{}

			if privateZoneVpc.VpcSet.UniqVpcId != nil {
				vpcSetMap["uniq_vpc_id"] = privateZoneVpc.VpcSet.UniqVpcId
			}

			if privateZoneVpc.VpcSet.Region != nil {
				vpcSetMap["region"] = privateZoneVpc.VpcSet.Region
			}

			vpcSetList = append(vpcSetList, vpcSetMap)
		}

		_ = d.Set("vpc_set", vpcSetList)

	}

	if privateZoneVpc.AccountVpcSet != nil {
		accountVpcSetList := []interface{}{}
		for _, accountVpcSet := range privateZoneVpc.AccountVpcSet {
			accountVpcSetMap := map[string]interface{}{}

			if privateZoneVpc.AccountVpcSet.UniqVpcId != nil {
				accountVpcSetMap["uniq_vpc_id"] = privateZoneVpc.AccountVpcSet.UniqVpcId
			}

			if privateZoneVpc.AccountVpcSet.Region != nil {
				accountVpcSetMap["region"] = privateZoneVpc.AccountVpcSet.Region
			}

			if privateZoneVpc.AccountVpcSet.Uin != nil {
				accountVpcSetMap["uin"] = privateZoneVpc.AccountVpcSet.Uin
			}

			if privateZoneVpc.AccountVpcSet.VpcName != nil {
				accountVpcSetMap["vpc_name"] = privateZoneVpc.AccountVpcSet.VpcName
			}

			accountVpcSetList = append(accountVpcSetList, accountVpcSetMap)
		}

		_ = d.Set("account_vpc_set", accountVpcSetList)

	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "privatedns", "zone", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudPrivatednsPrivateZoneVpcDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_privatedns_private_zone_vpc.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PrivatednsService{client: meta.(*TencentCloudClient).apiV3Conn}
	privateZoneVpcId := d.Id()

	if err := service.DeletePrivatednsPrivateZoneVpcById(ctx, zoneId); err != nil {
		return err
	}

	return nil
}
