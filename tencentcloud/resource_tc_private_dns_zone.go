/*
Provide a resource to create a Private Dns Zone.

Example Usage

```hcl
resource "tencentcloud_private_dns_zone" "foo" {
  domain = "domain.com"
  tags {
    "created_by" : "terraform"
  }
  vpc_set {
    region = "ap-guangzhou"
    uniq_vpc_id = "vpc-xxxxx"
  }
  remark = "test"
  dns_forward_status = "DISABLED"
  account_vpc_set {
    uin = "454xxxxxxx"
    region = "ap-guangzhou"
    uniq_vpc_id = "vpc-xxxxx"
    vpc_name = "test-redis"
  }
}
```

Import

Private Dns Zone can be imported, e.g.

```
$ terraform import tencentcloud_private_dns_zone.foo zone_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
)

func resourceTencentCloudPrivateDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDPrivateDnsZoneCreate,
		Read:   resourceTencentCloudDPrivateDnsZoneRead,
		Update: resourceTencentCloudDPrivateDnsZoneUpdate,
		Delete: resourceTencentCloudDPrivateDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name, which must be in the format of standard TLD.",
			},
			"tag_set": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Description:   "Tags the private domain when it is created.",
				Deprecated:    "It has been deprecated from version 1.72.4. Use `tags` instead.",
				ConflictsWith: []string{"tags"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of Tag.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of Tag.",
						},
					},
				},
			},
			"tags": {
				Type:          schema.TypeMap,
				Optional:      true,
				Description:   "Tags of the private dns zone.",
				ConflictsWith: []string{"tag_set"},
			},
			"vpc_set": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Associates the private domain to a VPC when it is created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC REGION.",
						},
					},
				},
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remarks.",
			},
			"dns_forward_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(PRIVATE_DNS_FORWARD_STATUS),
				Description:  "Whether to enable subdomain recursive DNS. Valid values: ENABLED, DISABLED. Default value: DISABLED.",
			},
			"account_vpc_set": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of authorized accounts' VPCs to associate with the private domain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "UIN of the VPC account.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC NAME.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDPrivateDnsZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_private_dns_zone.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := privatedns.NewCreatePrivateZoneRequest()

	domain := d.Get("domain").(string)
	request.Domain = &domain

	if v, ok := d.GetOk("tag_set"); ok {
		tagSet := make([]*privatedns.TagInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			tagInfo := privatedns.TagInfo{
				TagKey:   helper.String(m["tag_key"].(string)),
				TagValue: helper.String(m["tag_value"].(string)),
			}
			tagSet = append(tagSet, &tagInfo)
		}
		request.TagSet = tagSet
	}

	if v, ok := d.GetOk("vpc_set"); ok {
		vpcSet := make([]*privatedns.VpcInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			vpcInfo := privatedns.VpcInfo{
				UniqVpcId: helper.String(m["uniq_vpc_id"].(string)),
				Region:    helper.String(m["region"].(string)),
			}
			vpcSet = append(vpcSet, &vpcInfo)
		}
		request.VpcSet = vpcSet
	}

	if v, ok := d.GetOk("remark"); ok {
		remark := v.(string)
		request.Remark = helper.String(remark)
	}

	if v, ok := d.GetOk("dns_forward_status"); ok {
		dnsForwardStatus := v.(string)
		request.DnsForwardStatus = helper.String(dnsForwardStatus)
	}

	if v, ok := d.GetOk("account_vpc_set"); ok {
		accountVpcSet := make([]*privatedns.AccountVpcInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			accountVpcInfo := privatedns.AccountVpcInfo{
				Uin:       helper.String(m["uin"].(string)),
				UniqVpcId: helper.String(m["uniq_vpc_id"].(string)),
				Region:    helper.String(m["region"].(string)),
				VpcName:   helper.String(m["vpc_name"].(string)),
			}
			accountVpcSet = append(accountVpcSet, &accountVpcInfo)
		}
		request.AccountVpcSet = accountVpcSet
	}

	result, err := meta.(*TencentCloudClient).apiV3Conn.UsePrivateDnsClient().CreatePrivateZone(request)

	if err != nil {
		log.Printf("[CRITAL]%s create PrivateDns failed, reason:%s\n", logId, err.Error())
		return err
	}

	var response *privatedns.CreatePrivateZoneResponse
	response = result

	id := *response.Response.ZoneId
	d.SetId(id)

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client: client}
	region := client.Region

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName("privatedns", "zone", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudDPrivateDnsZoneRead(d, meta)
}

func resourceTencentCloudDPrivateDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_private_dns_zone.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	request := privatedns.NewDescribePrivateZoneRequest()
	request.ZoneId = helper.String(id)

	var response *privatedns.DescribePrivateZoneResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePrivateDnsClient().DescribePrivateZone(request)
		if e != nil {
			return retryError(e)
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read DnsPod Domain failed, reason:%s\n", logId, err.Error())
		return err
	}

	info := response.Response.PrivateZone
	d.SetId(*info.ZoneId)

	_ = d.Set("domain", info.Domain)

	tagSets := make([]map[string]interface{}, 0, len(info.Tags))
	for _, item := range info.Tags {
		tagSets = append(tagSets, map[string]interface{}{
			"tag_key":   item.TagKey,
			"tag_value": item.TagValue,
		})
	}
	_ = d.Set("tag_set", tagSets)

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client: client}
	region := client.Region

	tags, err := tagService.DescribeResourceTags(ctx, "privatedns", "zone", region, id)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	vpcSet := make([]map[string]interface{}, 0, len(info.VpcSet))
	for _, item := range info.VpcSet {
		vpcSet = append(vpcSet, map[string]interface{}{
			"uniq_vpc_id": item.UniqVpcId,
			"region":      item.Region,
		})
	}
	_ = d.Set("vpc_set", vpcSet)
	_ = d.Set("remark", info.Remark)
	_ = d.Set("dns_forward_status", info.DnsForwardStatus)

	accountVpcSet := make([]map[string]interface{}, 0, len(info.AccountVpcSet))
	for _, item := range info.AccountVpcSet {
		accountVpcSet = append(accountVpcSet, map[string]interface{}{
			"uin":         item.Uin,
			"uniq_vpc_id": item.UniqVpcId,
			"region":      item.Region,
		})
	}
	_ = d.Set("account_vpc_set", accountVpcSet)
	return nil
}

func resourceTencentCloudDPrivateDnsZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_private_dns_zone.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	id := d.Id()

	if d.HasChange("remark") || d.HasChange("dns_forward_status") {
		request := privatedns.NewModifyPrivateZoneRequest()
		request.ZoneId = helper.String(id)
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}
		if v, ok := d.GetOk("dns_forward_status"); ok {
			request.DnsForwardStatus = helper.String(v.(string))
		}
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UsePrivateDnsClient().ModifyPrivateZone(request)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify privateDns zone info failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("vpc_set") || d.HasChange("account_vpc_set") {
		request := privatedns.NewModifyPrivateZoneVpcRequest()
		request.ZoneId = helper.String(id)
		if v, ok := d.GetOk("vpc_set"); ok {
			var vpcSets = make([]*privatedns.VpcInfo, 0)
			items := v.([]interface{})
			for _, item := range items {
				value := item.(map[string]interface{})
				vpcInfo := &privatedns.VpcInfo{
					UniqVpcId: helper.String(value["uniq_vpc_id"].(string)),
					Region:    helper.String(value["region"].(string)),
				}
				vpcSets = append(vpcSets, vpcInfo)
			}
			request.VpcSet = vpcSets
		}

		if v, ok := d.GetOk("account_vpc_set"); ok {
			var accVpcSets = make([]*privatedns.AccountVpcInfo, 0)
			items := v.([]interface{})
			for _, item := range items {
				value := item.(map[string]interface{})
				accVpcInfo := &privatedns.AccountVpcInfo{
					UniqVpcId: helper.String(value["uniq_vpc_id"].(string)),
					Region:    helper.String(value["region"].(string)),
					Uin:       helper.String(value["uin"].(string)),
					VpcName:   helper.String(value["vpc_name"].(string)),
				}
				accVpcSets = append(accVpcSets, accVpcInfo)
			}
			request.AccountVpcSet = accVpcSets
		}
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UsePrivateDnsClient().ModifyPrivateZoneVpc(request)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify privateDns zone vpc failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("tag_set") {
		return fmt.Errorf("tag_set do not support change, please use tags instead.")
	}

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client: client}
	region := client.Region

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("privatedns", "zone", region, id)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	return resourceTencentCloudDPrivateDnsZoneRead(d, meta)
}

func resourceTencentCloudDPrivateDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_private_dns_zone.delete")()

	logId := getLogId(contextNil)

	request := privatedns.NewDeletePrivateZoneRequest()
	request.ZoneId = helper.String(d.Id())

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UsePrivateDnsClient().DeletePrivateZone(request)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete privateDns zone failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
