/*
Provides a resource to create a tcr customized_domain

Example Usage

```hcl
resource "tencentcloud_tcr_customized_domain" "customized_domain" {
  registry_id = "tcr-xxx"
  domain_name = "xxx.test.com"
  certificate_id = "kWGTVuU3"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr customized_domain can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_customized_domain.customized_domain customized_domain_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudTcrCustomizedDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrCustomizedDomainCreate,
		Read:   resourceTencentCloudTcrCustomizedDomainRead,
		Update: resourceTencentCloudTcrCustomizedDomainUpdate,
		Delete: resourceTencentCloudTcrCustomizedDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Custom domain name.",
			},

			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate id.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTcrCustomizedDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_customized_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tcr.NewCreateInstanceCustomizedDomainRequest()
		response   = tcr.NewCreateInstanceCustomizedDomainResponse()
		registryId string
		domainName string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().CreateInstanceCustomizedDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcr CustomizedDomain failed, reason:%+v", logId, err)
		return err
	}

	registryId = *response.Response.RegistryId
	d.SetId(strings.Join([]string{registryId, domainName}, FILED_SP))

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 2*readRetryTimeout, time.Second, service.TcrCustomizedDomainStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tcr:%s:uin/:instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrCustomizedDomainRead(d, meta)
}

func resourceTencentCloudTcrCustomizedDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_customized_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	domainName := idSplit[1]

	CustomizedDomain, err := service.DescribeTcrCustomizedDomainById(ctx, registryId, domainName)
	if err != nil {
		return err
	}

	if CustomizedDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TcrCustomizedDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if CustomizedDomain.RegistryId != nil {
		_ = d.Set("registry_id", CustomizedDomain.RegistryId)
	}

	if CustomizedDomain.DomainName != nil {
		_ = d.Set("domain_name", CustomizedDomain.DomainName)
	}

	if CustomizedDomain.CertificateId != nil {
		_ = d.Set("certificate_id", CustomizedDomain.CertificateId)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tcr", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTcrCustomizedDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_customized_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"registry_id", "domain_name", "certificate_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tcr", "instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrCustomizedDomainRead(d, meta)
}

func resourceTencentCloudTcrCustomizedDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_customized_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	registryId := idSplit[0]
	domainName := idSplit[1]

	if err := service.DeleteTcrCustomizedDomainById(ctx, registryId, domainName); err != nil {
		return err
	}

	return nil
}
