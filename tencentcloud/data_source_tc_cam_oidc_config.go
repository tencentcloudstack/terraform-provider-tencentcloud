/*
Use this data source to query detailed information of cam oidc_config

Example Usage

```hcl
data "tencentcloud_cam_oidc_config" "oidc_config" {
  name = "cls-kzilgv5m"
}

output "identity_key" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_key
}

output "identity_url" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_url
}

```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamOidcConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamOidcConfigRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name.",
			},

			"provider_type": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "IdP type. 11: Role IdP.",
			},

			"identity_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "IdP URL.",
			},

			"identity_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Public key for signature.",
			},

			"client_id": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Client ID.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Status. 0: Not set; 2: Disabled; 11: Enabled.",
			},

			"description": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Description.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCamOidcConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_oidc_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	var name string
	result := make(map[string]interface{})

	request := cam.NewDescribeOIDCConfigRequest()

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		request.Name = helper.String(v.(string))
	}

	var response *cam.DescribeOIDCConfigResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DescribeOIDCConfig(request)
		if e != nil {
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role SSO failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.ProviderType != nil {
		_ = d.Set("provider_type", response.Response.ProviderType)
		result["provider_type"] = response.Response.ProviderType
	}

	if response.Response.IdentityUrl != nil {
		_ = d.Set("identity_url", response.Response.IdentityUrl)
		result["identity_url"] = response.Response.IdentityUrl
	}

	if response.Response.IdentityKey != nil {
		_ = d.Set("identity_key", response.Response.IdentityKey)
		result["identity_key"] = response.Response.IdentityKey
	}

	if response.Response.ClientId != nil {
		_ = d.Set("client_id", response.Response.ClientId)
		result["client_id"] = response.Response.ClientId
	}

	if response.Response.Status != nil {
		_ = d.Set("status", response.Response.Status)
		result["status"] = response.Response.Status
	}

	if response.Response.Description != nil {
		_ = d.Set("description", response.Response.Description)
		result["description"] = response.Response.Description
	}

	d.SetId(name)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
