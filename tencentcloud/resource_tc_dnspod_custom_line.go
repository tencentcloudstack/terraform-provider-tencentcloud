/*
Provides a resource to create a dnspod custom_line

Example Usage

```hcl
resource "tencentcloud_dnspod_custom_line" "custom_line" {
  domain = "dnspod.com"
  name = "testline8"
  area = "6.6.6.1-6.6.6.3"
  domain_id = 1009
}
```

Import

dnspod custom_line can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_custom_line.custom_line custom_line_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudDnspodCustom_line() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodCustom_lineCreate,
		Read:   resourceTencentCloudDnspodCustom_lineRead,
		Update: resourceTencentCloudDnspodCustom_lineUpdate,
		Delete: resourceTencentCloudDnspodCustom_lineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
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

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},
		},
	}
}

func resourceTencentCloudDnspodCustom_lineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_custom_line.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dnspod.NewCreateDomainCustomLineRequest()
		response = dnspod.NewCreateDomainCustomLineResponse()
		domain   string
		name     string
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

	if v, ok := d.GetOkExists("domain_id"); ok {
		request.DomainId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnspodClient().CreateDomainCustomLine(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dnspod custom_line failed, reason:%+v", logId, err)
		return err
	}

	domain = *response.Response.Domain
	d.SetId(strings.Join([]string{domain, name}, FILED_SP))

	return resourceTencentCloudDnspodCustom_lineRead(d, meta)
}

func resourceTencentCloudDnspodCustom_lineRead(d *schema.ResourceData, meta interface{}) error {
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

	custom_line, err := service.DescribeDnspodCustom_lineById(ctx, domain, name)
	if err != nil {
		return err
	}

	if custom_line == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DnspodCustom_line` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if custom_line.Domain != nil {
		_ = d.Set("domain", custom_line.Domain)
	}

	if custom_line.Name != nil {
		_ = d.Set("name", custom_line.Name)
	}

	if custom_line.Area != nil {
		_ = d.Set("area", custom_line.Area)
	}

	if custom_line.DomainId != nil {
		_ = d.Set("domain_id", custom_line.DomainId)
	}

	return nil
}

func resourceTencentCloudDnspodCustom_lineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_custom_line.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dnspod.NewModifyDomainCustomLineRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domain := idSplit[0]
	name := idSplit[1]

	request.Domain = &domain
	request.Name = &name

	immutableArgs := []string{"domain", "name", "area", "domain_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("domain") {
		if v, ok := d.GetOk("domain"); ok {
			request.Domain = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("area") {
		if v, ok := d.GetOk("area"); ok {
			request.Area = helper.String(v.(string))
		}
	}

	if d.HasChange("domain_id") {
		if v, ok := d.GetOkExists("domain_id"); ok {
			request.DomainId = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnspodClient().ModifyDomainCustomLine(request)
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

	return resourceTencentCloudDnspodCustom_lineRead(d, meta)
}

func resourceTencentCloudDnspodCustom_lineDelete(d *schema.ResourceData, meta interface{}) error {
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

	if err := service.DeleteDnspodCustom_lineById(ctx, domain, name); err != nil {
		return err
	}

	return nil
}
