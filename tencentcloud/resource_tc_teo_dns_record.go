/*
Provides a resource to create a teo dnsRecord

Example Usage

```hcl
resource "tencentcloud_teo_dns_record" "dnsRecord"
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

teo dnsRecord can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_dns_record.dnsRecord dnsRecord_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoDnsRecord() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoDnsRecordRead,
		Create: resourceTencentCloudTeoDnsRecordCreate,
		Update: resourceTencentCloudTeoDnsRecordUpdate,
		Delete: resourceTencentCloudTeoDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"record_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS Record Type.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS Record Name.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS Record Content.",
			},

			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Proxy mode. Valid values: dns_only, cdn_only, and secure_cdn.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: ".",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Priority.",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time.",
			},

			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the DNS record is locked.",
			},

			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site Name.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resolution status.",
			},

			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME address.",
			},

			"domain_status": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: ".",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_record.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateDnsRecordRequest()
		response *teo.CreateDnsRecordResponse
	)

	if v, ok := d.GetOk("record_type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ttl"); ok {
		request.Ttl = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateDnsRecord(request)
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
		log.Printf("[CRITAL]%s create teo dnsRecord failed, reason:%+v", logId, err)
		return err
	}

	dnsRecordId := *response.Response.Id

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::teo:%s:uin/:zone/%s", region, dnsRecordId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	d.SetId(dnsRecordId)
	return resourceTencentCloudTeoDnsRecordRead(d, meta)
}

func resourceTencentCloudTeoDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dnsRecord.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	dnsRecordId := d.Id()

	dnsRecord, err := service.DescribeTeoDnsRecord(ctx, dnsRecordId)

	if err != nil {
		return err
	}

	if dnsRecord == nil {
		d.SetId("")
		return fmt.Errorf("resource `dnsRecord` %s does not exist", dnsRecordId)
	}

	if dnsRecord.Type != nil {
		_ = d.Set("record_type", dnsRecord.Type)
	}

	if dnsRecord.Name != nil {
		_ = d.Set("name", dnsRecord.Name)
	}

	if dnsRecord.Content != nil {
		_ = d.Set("content", dnsRecord.Content)
	}

	if dnsRecord.Mode != nil {
		_ = d.Set("mode", dnsRecord.Mode)
	}

	if dnsRecord.Ttl != nil {
		_ = d.Set("ttl", dnsRecord.Ttl)
	}

	if dnsRecord.Priority != nil {
		_ = d.Set("priority", dnsRecord.Priority)
	}

	if dnsRecord.CreatedOn != nil {
		_ = d.Set("created_on", dnsRecord.CreatedOn)
	}

	if dnsRecord.ModifiedOn != nil {
		_ = d.Set("modified_on", dnsRecord.ModifiedOn)
	}

	if dnsRecord.Locked != nil {
		_ = d.Set("locked", dnsRecord.Locked)
	}

	if dnsRecord.ZoneId != nil {
		_ = d.Set("zone_id", dnsRecord.ZoneId)
	}

	if dnsRecord.ZoneName != nil {
		_ = d.Set("zone_name", dnsRecord.ZoneName)
	}

	if dnsRecord.Status != nil {
		_ = d.Set("status", dnsRecord.Status)
	}

	if dnsRecord.Cname != nil {
		_ = d.Set("cname", dnsRecord.Cname)
	}

	if dnsRecord.DomainStatus != nil {
		_ = d.Set("domain_status", dnsRecord.DomainStatus)
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

func resourceTencentCloudTeoDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_record.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := teo.NewModifyDnsRecordRequest()

	request.Id = helper.String(d.Id())

	if d.HasChange("record_type") {
		if v, ok := d.GetOk("record_type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}
	}

	if d.HasChange("mode") {
		if v, ok := d.GetOk("mode"); ok {
			request.Mode = helper.String(v.(string))
		}
	}

	if d.HasChange("ttl") {
		if v, ok := d.GetOk("ttl"); ok {
			request.Ttl = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("priority") {
		if v, ok := d.GetOk("priority"); ok {
			request.Priority = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDnsRecord(request)
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

	return resourceTencentCloudTeoDnsRecordRead(d, meta)
}

func resourceTencentCloudTeoDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_record.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	dnsRecordId := d.Id()

	if err := service.DeleteTeoDnsRecordById(ctx, dnsRecordId); err != nil {
		return err
	}

	return nil
}
