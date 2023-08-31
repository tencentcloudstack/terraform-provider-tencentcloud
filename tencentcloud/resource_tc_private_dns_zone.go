/*
Provide a resource to create a Private Dns Zone.

Example Usage

Create a basic Private Dns Zone

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = tencentcloud_vpc.vpc.id
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
  }
}
```

Create a Private Dns Zone domain and bind associated accounts'VPC

```hcl
resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = tencentcloud_vpc.vpc.id
  }

  account_vpc_set {
    uin         = "123456789"
    uniq_vpc_id = "vpc-adsebmya"
    region      = "ap-guangzhou"
    vpc_name    = "vpc-name"
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Default:      DNS_FORWARD_STATUS_DISABLED,
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
			"cname_speedup_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CNAME_SPEEDUP_STATUS_ENABLED,
				ValidateFunc: validateAllowedStringValue(CNAME_SPEEDUP_STATUS),
				Description:  "CNAME acceleration: ENABLED, DISABLED, Default value is ENABLED.",
			},
		},
	}
}

func resourceTencentCloudDPrivateDnsZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_private_dns_zone.create")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		request = privatedns.NewCreatePrivateZoneRequest()
	)

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

	if v, ok := d.GetOk("cname_speedup_status"); ok {
		cnameSpeedupStatus := v.(string)
		request.CnameSpeedupStatus = helper.String(cnameSpeedupStatus)
	}

	result, err := meta.(*TencentCloudClient).apiV3Conn.UsePrivateDnsClient().CreatePrivateZone(request)

	if err != nil {
		log.Printf("[CRITAL]%s create PrivateDns failed, reason:%s\n", logId, err.Error())
		return err
	}

	response := result

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

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		request  = privatedns.NewDescribePrivateZoneRequest()
		response *privatedns.DescribePrivateZoneResponse
		id       = d.Id()
	)

	request.ZoneId = &id

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
	_ = d.Set("cname_speedup_status", info.CnameSpeedupStatus)

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

	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)
		id    = d.Id()
	)

	if d.HasChange("tag_set") {
		return fmt.Errorf("tag_set do not support change, please use tags instead.")
	}

	if d.HasChange("remark") || d.HasChange("dns_forward_status") || d.HasChange("cname_speedup_status") {
		request := privatedns.NewModifyPrivateZoneRequest()
		request.ZoneId = &id
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		if v, ok := d.GetOk("dns_forward_status"); ok {
			request.DnsForwardStatus = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cname_speedup_status"); ok {
			request.CnameSpeedupStatus = helper.String(v.(string))
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
		request.ZoneId = &id
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

	}

	return resourceTencentCloudDPrivateDnsZoneRead(d, meta)
}

func resourceTencentCloudDPrivateDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_private_dns_zone.delete")()

	var (
		logId   = getLogId(contextNil)
		request = privatedns.NewDeletePrivateZoneRequest()
		id      = d.Id()
	)

	request.ZoneId = &id

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
