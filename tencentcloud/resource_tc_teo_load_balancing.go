/*
Provides a resource to create a teo loadBalancing

Example Usage

```hcl
resource "tencentcloud_teo_load_balancing" "loadBalancing" {
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

teo loadBalancing can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_load_balancing.loadBalancing loadBalancing_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoLoadBalancing() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoLoadBalancingRead,
		Create: resourceTencentCloudTeoLoadBalancingCreate,
		Update: resourceTencentCloudTeoLoadBalancingUpdate,
		Delete: resourceTencentCloudTeoLoadBalancingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancing_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CLB instance ID.",
			},

			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subdomain name. You can use @ to represent the root domain.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Proxy mode. Valid values: dns_only: Only DNS, proxied: Enable proxy.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Indicates DNS TTL time when Type=dns_only.",
			},

			"origin_id": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "ID of the origin group used.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},

			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Schedules domain names, Note: This field may return null, indicating that no valid value can be obtained.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTeoLoadBalancingCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_load_balancing.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateLoadBalancingRequest()
		response *teo.CreateLoadBalancingResponse
		zoneId   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("origin_id"); ok {
		originIdSet := v.(*schema.Set).List()
		for i := range originIdSet {
			originId := originIdSet[i].(string)
			request.OriginId = append(request.OriginId, &originId)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateLoadBalancing(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo loadBalancing failed, reason:%+v", logId, err)
		return err
	}

	loadBalancingId := *response.Response.LoadBalancingId

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::teo:%s:uin/:zone/%s", region, loadBalancingId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	d.SetId(zoneId + "#" + loadBalancingId)
	return resourceTencentCloudTeoLoadBalancingRead(d, meta)
}

func resourceTencentCloudTeoLoadBalancingRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_loadBalancing.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	loadBalancingId := idSplit[1]

	loadBalancing, err := service.DescribeTeoLoadBalancing(ctx, zoneId, loadBalancingId)

	if err != nil {
		return err
	}

	if loadBalancing == nil {
		d.SetId("")
		return fmt.Errorf("resource `loadBalancing` %s does not exist", loadBalancingId)
	}

	if loadBalancing.LoadBalancingId != nil {
		_ = d.Set("load_balancing_id", loadBalancing.LoadBalancingId)
	}

	if loadBalancing.ZoneId != nil {
		_ = d.Set("zone_id", loadBalancing.ZoneId)
	}

	if loadBalancing.Host != nil {
		_ = d.Set("host", loadBalancing.Host)
	}

	if loadBalancing.Type != nil {
		_ = d.Set("type", loadBalancing.Type)
	}

	if loadBalancing.TTL != nil {
		_ = d.Set("ttl", loadBalancing.TTL)
	}

	if loadBalancing.OriginId != nil {
		_ = d.Set("origin_id", loadBalancing.OriginId)
	}

	if loadBalancing.UpdateTime != nil {
		_ = d.Set("update_time", loadBalancing.UpdateTime)
	}

	if loadBalancing.Cname != nil {
		_ = d.Set("cname", loadBalancing.Cname)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "teo", "zone", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTeoLoadBalancingUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_load_balancing.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := teo.NewModifyLoadBalancingRequest()

	request.ZoneId = helper.String(d.Id())

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if d.HasChange("host") {
		return fmt.Errorf("`host` do not support change now.")
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			request.TTL = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("origin_id") {
		if v, ok := d.GetOk("origin_id"); ok {
			originIdSet := v.(*schema.Set).List()
			for i := range originIdSet {
				originId := originIdSet[i].(string)
				request.OriginId = append(request.OriginId, &originId)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyLoadBalancing(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("teo", "zone", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoLoadBalancingRead(d, meta)
}

func resourceTencentCloudTeoLoadBalancingDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_load_balancing.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	loadBalancingId := idSplit[1]

	if err := service.DeleteTeoLoadBalancingById(ctx, zoneId, loadBalancingId); err != nil {
		return err
	}

	return nil
}
