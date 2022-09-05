/*
Provides a resource to create a teo dnsSec

Example Usage

```hcl
resource "tencentcloud_teo_dns_sec" "dns_sec" {
  zone_id = tencentcloud_teo_zone.zone.id
  status  = "disabled"
}

```
Import

teo dns_sec can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_dns_sec.dns_sec zoneId
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

func resourceTencentCloudTeoDnsSec() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoDnsSecRead,
		Create: resourceTencentCloudTeoDnsSecCreate,
		Update: resourceTencentCloudTeoDnsSecUpdate,
		Delete: resourceTencentCloudTeoDnsSecDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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
				Required:    true,
				Description: "DNSSEC status. Valid values: `enabled`, `disabled`.",
			},

			"dnssec": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "DNSSEC infos.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flags": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Flag.",
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Encryption algorithm.",
						},
						"key_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Encryption type.",
						},
						"digest_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Digest type.",
						},
						"digest_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Digest algorithm.",
						},
						"digest": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Digest message.",
						},
						"d_s": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DS record value.",
						},
						"key_tag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Key tag.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public key.",
						},
					},
				},
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification date.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsSecCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_sec.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = teo.NewModifyDnssecRequest()
		zoneId  string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.Id = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDnssec(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo dnsSec failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(zoneId)
	return resourceTencentCloudTeoDnsSecRead(d, meta)
}

func resourceTencentCloudTeoDnsSecRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_sec.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	zoneId := d.Id()

	dnsSec, err := service.DescribeTeoDnsSec(ctx, zoneId)

	if err != nil {
		return err
	}

	if dnsSec == nil {
		d.SetId("")
		return fmt.Errorf("resource `dnsSec` %s does not exist", zoneId)
	}

	if dnsSec.Name != nil {
		_ = d.Set("zone_name", dnsSec.Name)
	}

	if dnsSec.Status != nil {
		_ = d.Set("status", dnsSec.Status)
	}

	if dnsSec.Dnssec != nil {
		dnssecMap := map[string]interface{}{}
		if dnsSec.Dnssec.Flags != nil {
			dnssecMap["flags"] = dnsSec.Dnssec.Flags
		}
		if dnsSec.Dnssec.Algorithm != nil {
			dnssecMap["algorithm"] = dnsSec.Dnssec.Algorithm
		}
		if dnsSec.Dnssec.KeyType != nil {
			dnssecMap["key_type"] = dnsSec.Dnssec.KeyType
		}
		if dnsSec.Dnssec.DigestType != nil {
			dnssecMap["digest_type"] = dnsSec.Dnssec.DigestType
		}
		if dnsSec.Dnssec.DigestAlgorithm != nil {
			dnssecMap["digest_algorithm"] = dnsSec.Dnssec.DigestAlgorithm
		}
		if dnsSec.Dnssec.Digest != nil {
			dnssecMap["digest"] = dnsSec.Dnssec.Digest
		}
		if dnsSec.Dnssec.DS != nil {
			dnssecMap["d_s"] = dnsSec.Dnssec.DS
		}
		if dnsSec.Dnssec.KeyTag != nil {
			dnssecMap["key_tag"] = dnsSec.Dnssec.KeyTag
		}
		if dnsSec.Dnssec.PublicKey != nil {
			dnssecMap["public_key"] = dnsSec.Dnssec.PublicKey
		}

		_ = d.Set("dnssec", []interface{}{dnssecMap})
	}

	if dnsSec.ModifiedOn != nil {
		_ = d.Set("modified_on", dnsSec.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoDnsSecUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_sec.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := teo.NewModifyDnssecRequest()

	zoneId := d.Id()
	request.Id = &zoneId

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDnssec(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo dnsSec failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoDnsSecRead(d, meta)
}

func resourceTencentCloudTeoDnsSecDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_dns_sec.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
