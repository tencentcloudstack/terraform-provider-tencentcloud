package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDnspodCustomLine() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodCustomLineCreate,
		Read:   resourceTencentCloudDnspodCustomLineRead,
		Update: resourceTencentCloudDnspodCustomLineUpdate,
		Delete: resourceTencentCloudDnspodCustomLineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The Name of custom line.",
			},

			"area": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The IP segment of custom line, split with `-`.",
			},
		},
	}
}

func resourceTencentCloudDnspodCustomLineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_custom_line.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = dnspod.NewCreateDomainCustomLineRequest()
		domain  string
		name    string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().CreateDomainCustomLine(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dnspod custom_line failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{domain, name}, FILED_SP))

	return resourceTencentCloudDnspodCustomLineRead(d, meta)
}

func resourceTencentCloudDnspodCustomLineRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_custom_line.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domain := idSplit[0]
	name := idSplit[1]

	customLineInfo, err := service.DescribeDnspodCustomLineById(ctx, domain, name)
	if err != nil {
		return err
	}

	if customLineInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DnspodCustom_line` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("name", name)

	if customLineInfo.Area != nil {
		_ = d.Set("area", customLineInfo.Area)
	}

	return nil
}

func resourceTencentCloudDnspodCustomLineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_custom_line.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dnspod.NewModifyDomainCustomLineRequest()

	var (
		newName string
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domain := idSplit[0]
	name := idSplit[1]

	request.Domain = &domain
	request.Name = &name

	if v, ok := d.GetOk("name"); ok {
		newName = v.(string)
		request.PreName = helper.String(name)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().ModifyDomainCustomLine(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dnspod custom_line failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("name") {
		d.SetId(strings.Join([]string{domain, newName}, FILED_SP))
	}

	return resourceTencentCloudDnspodCustomLineRead(d, meta)
}

func resourceTencentCloudDnspodCustomLineDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_custom_line.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domain := idSplit[0]
	name := idSplit[1]

	if err := service.DeleteDnspodCustomLineById(ctx, domain, name); err != nil {
		return err
	}

	return nil
}
