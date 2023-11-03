/*
Provides a resource to create a dnspod domain_alias

Example Usage

```hcl
resource "tencentcloud_dnspod_domain_alias" "domain_alias" {
  domain_alias = "dnspod.com"
  domain = "dnspod.cn"
  domain_id = 123
}
```

Import

dnspod domain_alias can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_domain_alias.domain_alias domain_alias_id
```
*/
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

func resourceTencentCloudDnspodDomainAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodDomainAliasCreate,
		Read:   resourceTencentCloudDnspodDomainAliasRead,
		Update: resourceTencentCloudDnspodDomainAliasUpdate,
		Delete: resourceTencentCloudDnspodDomainAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_alias": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain alias.",
			},

			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain.",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},
		},
	}
}

func resourceTencentCloudDnspodDomainAliasCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_alias.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = dnspod.NewCreateDomainAliasRequest()
		response      = dnspod.NewCreateDomainAliasResponse()
		domainAlias   string
		domainAliasId string
	)
	if v, ok := d.GetOk("domain_alias"); ok {
		request.DomainAlias = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		request.DomainId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnspodClient().CreateDomainAlias(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dnspod domain_alias failed, reason:%+v", logId, err)
		return err
	}

	domainAliasId = *response.Response.DomainAliasId
	// d.SetId(helper.String(domainAlias))
	d.SetId(strings.Join([]string{domain, helper.UInt64ToStr(groupId)}, FILED_SP))

	return resourceTencentCloudDnspodDomainAliasRead(d, meta)
}

func resourceTencentCloudDnspodDomainAliasRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_alias.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainAlias := idSplit[0]
	domainAliasId := idSplit[1]

	domain_alias, err := service.DescribeDnspodDomainAliasById(ctx, domainAlias, domainAliasId)
	if err != nil {
		return err
	}

	if domain_alias == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DnspodDomain_alias` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domain_alias.DomainAlias != nil {
		_ = d.Set("domain_alias", domain_alias.DomainAlias)
	}

	if domain_alias.Domain != nil {
		_ = d.Set("domain", domain_alias.Domain)
	}

	if domain_alias.DomainId != nil {
		_ = d.Set("domain_id", domain_alias.DomainId)
	}

	return nil
}

func resourceTencentCloudDnspodDomainAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_alias.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"domain_alias", "domain", "domain_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudDnspodDomainAliasRead(d, meta)
}

func resourceTencentCloudDnspodDomainAliasDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_domain_alias.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainAlias := idSplit[0]
	domainAliasId := idSplit[1]

	if err := service.DeleteDnspodDomainAliasById(ctx, domainAlias, domainAliasId); err != nil {
		return err
	}

	return nil
}
