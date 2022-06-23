/*
Provides a resource to create a CDN domain.

Example Usage

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "xxxx.com"
  service_type   = "web"
  area           = "mainland"
  full_url_cache = false

  origin {
    origin_type          = "ip"
    origin_list          = ["127.0.0.1"]
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"

    force_redirect {
      switch               = "on"
      redirect_type        = "http"
      redirect_status_code = 302
    }
  }

  tags = {
    hello = "world"
  }
}
```

Example Usage of cdn uses cache and request headers

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "xxxx.com"
  service_type   = "web"
  area           = "mainland"
  full_url_cache = false
  range_origin_switch = "off"

  rule_cache{
  	cache_time = 10000
  	no_cache_switch="on"
  	re_validate="on"
  }

  request_header{
  	switch = "on"

  	header_rules {
  		header_mode = "add"
  		header_name = "tf-header-name"
  		header_value = "tf-header-value"
  		rule_type = "all"
  		rule_paths = ["*"]
  	}
  }

  origin {
    origin_type          = "ip"
    origin_list          = ["127.0.0.1"]
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"

    force_redirect {
      switch               = "on"
      redirect_type        = "http"
      redirect_status_code = 302
    }
  }

  tags = {
    hello = "world"
  }
}
```

Example Usage of COS bucket url as origin

```hcl
resource "tencentcloud_cos_bucket" "bucket" {
  # Bucket format should be [custom name]-[appid].
  bucket = "demo-bucket-1251234567"
  acl    = "private"
}

# Create cdn domain
resource "tencentcloud_cdn_domain" "cdn" {
  domain         = "abc.com"
  service_type   = "web"
  area           = "mainland"
  full_url_cache = false

  origin {
    origin_type          = "cos"
    origin_list          = [tencentcloud_cos_bucket.bucket.cos_bucket_url]
    server_name          = tencentcloud_cos_bucket.bucket.cos_bucket_url
    origin_pull_protocol = "follow"
    cos_private_access   = "on"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"
  }
}
```

Import

CDN domain can be imported using the id, e.g.

```
$ terraform import tencentcloud_cdn_domain.foo xxxx.com
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudCdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdnDomainCreate,
		Read:   resourceTencentCloudCdnDomainRead,
		Update: resourceTencentCloudCdnDomainUpdate,
		Delete: resourceTencentCloudCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("authentication", []interface{}{map[string]interface{}{
					"switch": "off",
				}})
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the acceleration domain.",
			},
			"service_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CDN_SERVICE_TYPE),
				Description:  "Acceleration domain name service type. `web`: static acceleration, `download`: download acceleration, `media`: streaming media VOD acceleration.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The project CDN belongs to, default to 0.",
			},
			"area": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CDN_AREA),
				Description:  "Domain name acceleration region. `mainland`: acceleration inside mainland China, `overseas`: acceleration outside mainland China, `global`: global acceleration. Overseas acceleration service must be enabled to use overseas acceleration and global acceleration.",
			},
			"full_url_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable full-path cache. Default value is `true`.",
			},
			"origin": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Origin server configuration. It's a list and consist of at most one item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue(CDN_ORIGIN_TYPE),
							Description:  "Master origin server type. The following types are supported: `domain`: domain name type, `cos`: COS origin, `ip`: IP list used as origin server, `ipv6`: origin server list is a single IPv6 address, `ip_ipv6`: origin server list is multiple IPv4 addresses and an IPv6 address.",
						},
						"origin_list": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Master origin server list. Valid values can be ip or domain name. When modifying the origin server, you need to enter the corresponding `origin_type`.",
						},
						"server_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Host header used when accessing the master origin server. If left empty, the acceleration domain name will be used by default.",
						},
						"cos_private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     CDN_SWITCH_OFF,
							Description: "When OriginType is COS, you can specify if access to private buckets is allowed. Valid values are `on` and `off`. and default value is `off`.",
						},
						"origin_pull_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_ORIGIN_PULL_PROTOCOL_HTTP,
							ValidateFunc: validateAllowedStringValue(CDN_ORIGIN_PULL_PROTOCOL),
							Description:  "Origin-pull protocol configuration. `http`: forced HTTP origin-pull, `follow`: protocol follow origin-pull, `https`: forced HTTPS origin-pull. This only supports origin server port 443 for origin-pull.",
						},
						"backup_origin_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateAllowedStringValue(CDN_BACKUP_ORIGIN_TYPE),
							Description:  "Backup origin server type, which supports the following types: `domain`: domain name type, `ip`: IP list used as origin server.",
						},
						"backup_origin_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Backup origin server list. Valid values can be ip or domain name. When modifying the backup origin server, you need to enter the corresponding `backup_origin_type`.",
						},
						"backup_server_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Host header used when accessing the backup origin server. If left empty, the ServerName of master origin server will be used by default.",
						},
					},
				},
			},
			"https_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "HTTPS acceleration configuration. It's a list and consist of at most one item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"https_switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "HTTPS configuration switch. Valid values are `on` and `off`.",
						},
						"http2_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "HTTP2 configuration switch. Valid values are `on` and `off`. and default value is `off`.",
						},
						"ocsp_stapling_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "OCSP configuration switch. Valid values are `on` and `off`. and default value is `off`.",
						},
						"spdy_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Spdy configuration switch. Valid values are `on` and `off`. and default value is `off`. This parameter is for white-list customer.",
						},
						"verify_client": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Client certificate authentication feature. Valid values are `on` and `off`. and default value is `off`.",
						},
						"server_certificate_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Server certificate configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server certificate ID.",
									},
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Server certificate name.",
									},
									"certificate_content": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server certificate information. This is required when uploading an external certificate, which should contain the complete certificate chain.",
									},
									"private_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server key information. This is required when uploading an external certificate.",
									},
									"message": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Certificate remarks.",
									},
									"deploy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy time of server certificate.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire time of server certificate.",
									},
								},
							},
						},
						"client_certificate_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Client certificate configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client certificate name.",
									},
									"certificate_content": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Client Certificate PEM format, requires Base64 encoding.",
									},
									"deploy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy time of client certificate.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire time of client certificate.",
									},
								},
							},
						},
						"force_redirect": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Configuration of forced HTTP or HTTPS redirects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      CDN_SWITCH_OFF,
										ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
										Description:  "Forced redirect configuration switch. Valid values are `on` and `off`. Default value is `off`.",
									},
									"redirect_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      CDN_ORIGIN_PULL_PROTOCOL_HTTP,
										ValidateFunc: validateAllowedStringValue(CDN_FORCE_REDIRECT_TYPE),
										Description: "Forced redirect type. Valid values are `http` and `https`. `http` means a forced redirect from HTTPS to HTTP, `https` means a forced redirect from HTTP to HTTPS. " +
											"When `switch` setting `off`, this property does not need to be set or set to `http`. Default value is `http`.",
									},
									"redirect_status_code": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      302,
										ValidateFunc: validateAllowedIntValue([]int{301, 302}),
										Description: "Forced redirect status code. Valid values are `301` and `302`. " +
											"When `switch` setting `off`, this property does not need to be set or set to `302`. Default value is `302`.",
									},
								},
							},
						},
					},
				},
			},
			"range_origin_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_ON,
				ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
				Description:  "Sharding back to source configuration switch. Valid values are `on` and `off`. Default value is `on`.",
			},
			"ipv6_access_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_OFF,
				ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
				Description:  "ipv6 access configuration switch. Only available when area set to `mainland`. Valid values are `on` and `off`. Default value is `off`.",
			},
			"follow_redirect_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_OFF,
				ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
				Description:  "301/302 redirect following switch, available values: `on`, `off` (default).",
			},
			"authentication": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specify timestamp hotlink protection configuration, NOTE: only one type can choose for the sub elements.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Authentication switching, available values: `on`, `off`.",
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
						},
						"type_a": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp hotlink protection mode A configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.",
									},
									"sign_param": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Signature parameter name. Only upper and lower-case letters, digits, and underscores (_) are allowed. It cannot start with a digit. Length limit: 1-100 characters.",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Signature expiration time in second. The maximum value is 630720000.",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.",
									},
								},
							},
						},
						"type_b": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp hotlink protection mode B configuration. NOTE: according to upgrading of TencentCloud Platform, TypeB is unavailable for now.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Signature expiration time in second. The maximum value is 630720000.",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.",
									},
								},
							},
						},
						"type_c": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp hotlink protection mode C configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Signature expiration time in second. The maximum value is 630720000.",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.",
									},
									"time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Timestamp formation, available values: `dec`, `hex`.",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.",
									},
								},
							},
						},
						"type_d": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp hotlink protection mode D configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Signature expiration time in second. The maximum value is 630720000.",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.",
									},
									"time_param": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Timestamp parameter name. Only upper and lower-case letters, digits, and underscores (_) are allowed. It cannot start with a digit. Length limit: 1-100 characters.",
									},
									"time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Timestamp formation, available values: `dec`, `hex`.",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.",
									},
								},
							},
						},
					},
				},
			},
			"rule_cache": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Advanced path cache configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Matching content under the corresponding type of CacheType: `all`: fill *, `file`: fill in the suffix name, such as jpg, txt, " +
								"`directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html, `index`: fill /, `default`: Fill `no max-age`.",
						},
						"rule_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_RULE_TYPE_DEFAULT,
							ValidateFunc: validateAllowedStringValue(CDN_RULE_TYPE),
							Description: "Rule type. The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, " +
								"`directory`: the specified path takes effect, `path`: specify the absolute path to take effect, `index`: home page, `default`: effective when the source site has no max-age.",
						},
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Cache configuration switch. Valid values are `on` and `off`.",
						},
						"cache_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Cache expiration time setting, the unit is second, the maximum can be set to 365 days.",
						},
						"compare_max_age": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description: "Advanced cache expiration configuration. When it is turned on, it will compare the max-age value returned by the origin site with the cache expiration time set in CacheRules, " +
								"and take the minimum value to cache at the node. Valid values are `on` and `off`. Default value is `off`.",
						},
						"ignore_cache_control": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description: "Force caching. After opening, the no-store and no-cache resources returned by the origin site will also be cached in accordance with the CacheRules " +
								"rules. Valid values are `on` and `off`. Default value is `off`.",
						},
						"ignore_set_cookie": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Ignore the Set-Cookie header of the origin site. Valid values are `on` and `off`. Default value is `off`. This parameter is for white-list customer.",
						},
						"no_cache_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Cache configuration switch. Valid values are `on` and `off`.",
						},
						"re_validate": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Always check back to origin. Valid values are `on` and `off`. Default value is `off`.",
						},
						"follow_origin_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Follow the source station configuration switch. Valid values are `on` and `off`.",
						},
					},
				},
			},
			"request_header": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Request header configuration. It's a list and consist of at most one item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: validateAllowedStringValue(CDN_SWITCH),
							Description:  "Custom request header configuration switch. Valid values are `on` and `off`. and default value is `off`.",
						},
						"header_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom request header configuration rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Http header setting method. The following types are supported: `add`: add a head, if a head already exists, there will be a duplicate head, `del`: delete the head.",
									},
									"header_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateStringLengthInRange(1, 100),
										Description:  "Http header name.",
									},
									"header_value": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateStringLengthInRange(1, 1000),
										Description:  "Http header value, optional when Mode is `del`, Required when Mode is `add`/`set`.",
									},
									"rule_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateAllowedStringValue(CDN_HEADER_RULE),
										Description: "Rule type. The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, " +
											"`directory`: the specified path takes effect, `path`: specify the absolute path to take effect.",
									},
									"rule_paths": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Description: "Matching content under the corresponding type of CacheType: `all`: fill *, `file`: fill in the suffix name, such as jpg, txt, " +
											"`directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html.",
									},
								},
							},
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of cdn domain.",
			},

			// computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration service status.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME address of domain name.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of domain name.",
			},
		},
	}
}

func resourceTencentCloudCdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_domain.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cdnService := CdnService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := cdn.NewAddCdnDomainRequest()
	domain := d.Get("domain").(string)
	request.Domain = &domain
	request.ServiceType = helper.String(d.Get("service_type").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}
	fullUrlCache := d.Get("full_url_cache").(bool)
	request.CacheKey = &cdn.CacheKey{}
	if fullUrlCache {
		request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_ON)
	} else {
		request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_OFF)
	}
	// Range Origin Pull
	request.RangeOriginPull = &cdn.RangeOriginPull{}
	request.RangeOriginPull.Switch = helper.String(d.Get("range_origin_switch").(string))

	if v, ok := d.GetOk("ipv6_access_switch"); ok {
		request.Ipv6Access = &cdn.Ipv6Access{
			Switch: helper.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("follow_redirect_switch"); ok {
		request.FollowRedirect = &cdn.FollowRedirect{
			Switch: helper.String(v.(string)),
		}
	}

	if v, ok := helper.InterfacesHeadMap(d, "authentication"); ok {
		switchOn := v["switch"].(string)
		request.Authentication = &cdn.Authentication{
			Switch: &switchOn,
		}

		if v, ok := v["type_a"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				signParam      = item["sign_param"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeA = &cdn.AuthenticationTypeA{
				SecretKey:      &secretKey,
				SignParam:      &signParam,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeA.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_b"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeB = &cdn.AuthenticationTypeB{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeB.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_c"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeC = &cdn.AuthenticationTypeC{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if timeFormat, ok := item["time_format"].(string); ok {
				request.Authentication.TypeC.TimeFormat = &timeFormat
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeC.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_d"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
				timeParam      = item["time_param"].(string)
			)

			request.Authentication.TypeD = &cdn.AuthenticationTypeD{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
				TimeParam:      &timeParam,
			}

			if timeFormat, ok := item["time_format"].(string); ok {
				request.Authentication.TypeD.TimeFormat = &timeFormat
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeD.BackupSecretKey = &backupSecretKey
			}
		}
	}

	// rule_cache
	if v, ok := d.GetOk("rule_cache"); ok {
		ruleCache := v.([]interface{})
		var ruleCaches []*cdn.RuleCache
		for _, v := range ruleCache {
			re := &cdn.RuleCache{}
			ruleCacheMap := v.(map[string]interface{})
			rulePaths := ruleCacheMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			ruleType := ruleCacheMap["rule_type"].(string)
			if ruleType == CDN_RULE_TYPE_DEFAULT {
				rulePathList = append(rulePathList, helper.String(CDN_RULE_PATH))
			} else {
				for _, value := range rulePaths {
					rulePathList = append(rulePathList, helper.String(value.(string)))
				}
			}
			switchFlag := ruleCacheMap["switch"].(string)
			cacheTime := ruleCacheMap["cache_time"].(int)
			compareMaxAge := ruleCacheMap["compare_max_age"].(string)
			ignoreCacheControl := ruleCacheMap["ignore_cache_control"].(string)
			ignoreSetCookie := ruleCacheMap["ignore_set_cookie"].(string)
			noCacheSwitch := ruleCacheMap["no_cache_switch"].(string)
			reValidate := ruleCacheMap["re_validate"].(string)
			followOriginSwitch := ruleCacheMap["follow_origin_switch"].(string)
			ruleCacheConfig := &cdn.RuleCacheConfig{}
			cache := &cdn.CacheConfigCache{}
			noCache := &cdn.CacheConfigNoCache{}
			followOrigin := &cdn.CacheConfigFollowOrigin{}
			ruleCacheConfig.Cache = cache
			ruleCacheConfig.NoCache = noCache
			ruleCacheConfig.FollowOrigin = followOrigin
			re.CacheConfig = ruleCacheConfig
			re.RulePaths = rulePathList
			re.RuleType = &ruleType
			re.CacheConfig.Cache.Switch = &switchFlag
			re.CacheConfig.Cache.CacheTime = helper.IntInt64(cacheTime)
			re.CacheConfig.Cache.CompareMaxAge = &compareMaxAge
			re.CacheConfig.Cache.IgnoreCacheControl = &ignoreCacheControl
			re.CacheConfig.Cache.IgnoreSetCookie = &ignoreSetCookie
			re.CacheConfig.NoCache.Switch = &noCacheSwitch
			re.CacheConfig.NoCache.Revalidate = &reValidate
			re.CacheConfig.FollowOrigin.Switch = &followOriginSwitch
			ruleCaches = append(ruleCaches, re)
		}
		request.Cache = &cdn.Cache{}
		request.Cache.RuleCache = ruleCaches
	}

	if v, ok := d.GetOk("request_header"); ok {
		requestHeaders := v.([]interface{})
		requestHeader := requestHeaders[0].(map[string]interface{})
		headerRule := requestHeader["header_rules"].([]interface{})
		var headerRules []*cdn.HttpHeaderPathRule
		for _, value := range headerRule {
			hr := &cdn.HttpHeaderPathRule{}
			headerRuleMap := value.(map[string]interface{})
			headerMode := headerRuleMap["header_mode"].(string)
			headerName := headerRuleMap["header_name"].(string)
			headerValue := headerRuleMap["header_value"].(string)
			ruleType := headerRuleMap["rule_type"].(string)
			rulePaths := headerRuleMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			for _, value := range rulePaths {
				rulePathList = append(rulePathList, helper.String(value.(string)))
			}
			hr.HeaderMode = &headerMode
			hr.HeaderName = &headerName
			hr.HeaderValue = &headerValue
			hr.RuleType = &ruleType
			hr.RulePaths = rulePathList
			headerRules = append(headerRules, hr)
		}
		request.RequestHeader = &cdn.RequestHeader{}
		request.RequestHeader.Switch = helper.String(requestHeader["switch"].(string))
		request.RequestHeader.HeaderRules = headerRules
	}

	// origin
	origins := d.Get("origin").([]interface{})
	if len(origins) < 1 {
		return fmt.Errorf("origin is required")
	}
	origin := origins[0].(map[string]interface{})
	request.Origin = &cdn.Origin{}
	request.Origin.OriginType = helper.String(origin["origin_type"].(string))
	originList := origin["origin_list"].([]interface{})
	request.Origin.Origins = make([]*string, 0, len(originList))
	for _, item := range originList {
		request.Origin.Origins = append(request.Origin.Origins, helper.String(item.(string)))
	}
	if v := origin["server_name"]; v.(string) != "" {
		request.Origin.ServerName = helper.String(v.(string))
	}
	if v := origin["cos_private_access"]; v.(string) != "" {
		request.Origin.CosPrivateAccess = helper.String(v.(string))
	}
	if v := origin["origin_pull_protocol"]; v.(string) != "" {
		request.Origin.OriginPullProtocol = helper.String(v.(string))
	}
	if v := origin["backup_origin_type"]; v.(string) != "" {
		request.Origin.BackupOriginType = helper.String(v.(string))
	}
	if v := origin["backup_server_name"]; v.(string) != "" {
		request.Origin.BackupServerName = helper.String(v.(string))
	}
	if v := origin["backup_origin_list"]; len(v.([]interface{})) > 0 {
		backupOriginList := v.([]interface{})
		request.Origin.BackupOrigins = make([]*string, 0, len(backupOriginList))
		for _, item := range backupOriginList {
			request.Origin.BackupOrigins = append(request.Origin.BackupOrigins, helper.String(item.(string)))
		}
	}

	// https config
	if v, ok := d.GetOk("https_config"); ok {
		httpsConfigs := v.([]interface{})
		if len(httpsConfigs) > 0 {
			config := httpsConfigs[0].(map[string]interface{})
			request.Https = &cdn.Https{}
			request.Https.Switch = helper.String(config["https_switch"].(string))
			if v := config["http2_switch"]; v.(string) != "" {
				request.Https.Http2 = helper.String(v.(string))
			}
			request.Https.OcspStapling = helper.String(config["ocsp_stapling_switch"].(string))
			request.Https.Spdy = helper.String(config["spdy_switch"].(string))
			request.Https.VerifyClient = helper.String(config["verify_client"].(string))
			if v := config["server_certificate_config"]; len(v.([]interface{})) > 0 {
				serverCerts := v.([]interface{})
				if len(serverCerts) > 0 && serverCerts[0] != nil {
					serverCert := serverCerts[0].(map[string]interface{})
					request.Https.CertInfo = &cdn.ServerCert{}
					if v := serverCert["certificate_id"]; v.(string) != "" {
						request.Https.CertInfo.CertId = helper.String(v.(string))
					}
					if v := serverCert["certificate_content"]; v.(string) != "" {
						request.Https.CertInfo.Certificate = helper.String(v.(string))
					}
					if v := serverCert["private_key"]; v.(string) != "" {
						request.Https.CertInfo.PrivateKey = helper.String(v.(string))
					}
					if v := serverCert["message"]; v.(string) != "" {
						request.Https.CertInfo.Message = helper.String(v.(string))
					}
				}
			}
			if v := config["client_certificate_config"]; len(v.([]interface{})) > 0 {
				clientCerts := v.([]interface{})
				if len(clientCerts) > 0 && clientCerts[0] != nil {
					clientCert := clientCerts[0].(map[string]interface{})
					request.Https.ClientCertInfo = &cdn.ClientCert{}
					if v := clientCert["certificate_content"]; v.(string) != "" {
						request.Https.ClientCertInfo.Certificate = helper.String(v.(string))
					}
				}
			}
			if v, ok := config["force_redirect"]; ok {
				forceRedirect := v.([]interface{})
				if len(forceRedirect) > 0 {
					var redirect cdn.ForceRedirect
					redirectMap := forceRedirect[0].(map[string]interface{})
					if sw := redirectMap["switch"]; sw.(string) != "" {
						redirect.Switch = helper.String(sw.(string))
					}
					if rt := redirectMap["redirect_type"]; rt.(string) != "" {
						redirect.RedirectType = helper.String(rt.(string))
					}
					if rsc := redirectMap["redirect_status_code"]; rsc.(int) != 0 {
						redirect.RedirectStatusCode = helper.Int64(int64(rsc.(int)))
					}
					request.ForceRedirect = &redirect
				}
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := meta.(*TencentCloudClient).apiV3Conn.UseCdnClient().AddCdnDomain(request)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROE || sdkErr.Code == CDN_HOST_EXISTS {
					return resource.NonRetryableError(err)
				}
			}
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(domain)

	time.Sleep(1 * time.Second)
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if err != nil {
			return retryError(err, InternalError)
		}
		if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
			return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
		}
		return nil
	})
	if err != nil {
		return err
	}

	// tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		client := meta.(*TencentCloudClient).apiV3Conn
		tagService := TagService{client: client}
		region := client.Region
		resourceName := BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		err := tagService.ModifyTags(ctx, resourceName, tags, nil)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudCdnDomainRead(d, meta)
}

func resourceTencentCloudCdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	region := client.Region
	cdnService := CdnService{client: client}
	tagService := TagService{client: client}

	domain := d.Id()
	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if domainConfig == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("service_type", domainConfig.ServiceType)
	_ = d.Set("project_id", domainConfig.ProjectId)
	_ = d.Set("area", domainConfig.Area)
	_ = d.Set("status", domainConfig.Status)
	_ = d.Set("create_time", domainConfig.CreateTime)
	_ = d.Set("cname", domainConfig.Cname)
	_ = d.Set("range_origin_switch", domainConfig.RangeOriginPull.Switch)

	if domainConfig.Ipv6Access != nil {
		_ = d.Set("ipv6_access_switch", domainConfig.Ipv6Access.Switch)
	}
	if domainConfig.FollowRedirect != nil {
		_ = d.Set("follow_redirect_switch", domainConfig.FollowRedirect.Switch)
	}
	if *domainConfig.CacheKey.FullUrlCache == CDN_SWITCH_OFF {
		_ = d.Set("full_url_cache", false)
	} else {
		_ = d.Set("full_url_cache", true)
	}

	origins := make([]map[string]interface{}, 0, 1)
	origin := make(map[string]interface{}, 8)
	origin["origin_type"] = domainConfig.Origin.OriginType
	origin["origin_list"] = domainConfig.Origin.Origins
	origin["server_name"] = domainConfig.Origin.ServerName
	origin["cos_private_access"] = domainConfig.Origin.CosPrivateAccess
	origin["origin_pull_protocol"] = domainConfig.Origin.OriginPullProtocol
	origin["backup_origin_type"] = domainConfig.Origin.BackupOriginType
	origin["backup_origin_list"] = domainConfig.Origin.BackupOrigins
	origin["backup_server_name"] = domainConfig.Origin.BackupServerName
	origins = append(origins, origin)
	_ = d.Set("origin", origins)

	if len(domainConfig.Cache.RuleCache) > 0 {
		ruleCaches := make([]map[string]interface{}, len(domainConfig.Cache.RuleCache))
		for index, value := range domainConfig.Cache.RuleCache {
			ruleCache := make(map[string]interface{})
			ruleCache["rule_paths"] = value.RulePaths
			ruleCache["rule_type"] = value.RuleType
			ruleCache["switch"] = value.CacheConfig.Cache.Switch
			ruleCache["cache_time"] = value.CacheConfig.Cache.CacheTime
			ruleCache["compare_max_age"] = value.CacheConfig.Cache.CompareMaxAge
			ruleCache["ignore_cache_control"] = value.CacheConfig.Cache.IgnoreCacheControl
			ruleCache["ignore_set_cookie"] = value.CacheConfig.Cache.IgnoreSetCookie
			ruleCache["no_cache_switch"] = value.CacheConfig.NoCache.Switch
			ruleCache["re_validate"] = value.CacheConfig.NoCache.Revalidate
			ruleCache["follow_origin_switch"] = value.CacheConfig.FollowOrigin.Switch
			ruleCaches[index] = ruleCache
		}
		_ = d.Set("rule_cache", ruleCaches)
	}

	requestHeaders := make([]map[string]interface{}, 1)
	requestHeader := make(map[string]interface{})
	if domainConfig.RequestHeader != nil {
		requestHeader["switch"] = domainConfig.RequestHeader.Switch
		if len(domainConfig.RequestHeader.HeaderRules) > 0 {
			headerRules := make([]map[string]interface{}, len(domainConfig.RequestHeader.HeaderRules))
			headerRuleList := domainConfig.RequestHeader.HeaderRules
			for index, value := range headerRuleList {
				headerRule := make(map[string]interface{})
				headerRule["header_mode"] = value.HeaderMode
				headerRule["header_name"] = value.HeaderName
				headerRule["header_value"] = value.HeaderValue
				headerRule["rule_type"] = value.RuleType
				headerRule["rule_paths"] = value.RulePaths
				headerRules[index] = headerRule
			}
			requestHeader["header_rules"] = headerRules
		}
		requestHeaders[0] = requestHeader
		_ = d.Set("request_header", requestHeaders)
	}

	httpsConfigs := make([]map[string]interface{}, 0, 1)
	httpsConfig := make(map[string]interface{}, 8)
	httpsConfig["https_switch"] = domainConfig.Https.Switch
	httpsConfig["http2_switch"] = domainConfig.Https.Http2
	httpsConfig["ocsp_stapling_switch"] = domainConfig.Https.OcspStapling
	httpsConfig["spdy_switch"] = domainConfig.Https.Spdy
	httpsConfig["verify_client"] = domainConfig.Https.VerifyClient

	oldHttpsConfigs := make([]interface{}, 0)
	if _, ok := d.GetOk("https_config"); ok {
		oldHttpsConfigs = d.Get("https_config").([]interface{})
	}
	oldHttpsConfig := make(map[string]interface{})
	if len(oldHttpsConfigs) > 0 {
		oldHttpsConfig = oldHttpsConfigs[0].(map[string]interface{})
	}
	oldServerConfigs := make([]interface{}, 0)
	if _, ok := oldHttpsConfig["server_certificate_config"]; ok {
		oldServerConfigs = oldHttpsConfig["server_certificate_config"].([]interface{})
	}
	oldServerConfig := make(map[string]interface{})
	if len(oldServerConfigs) > 0 {
		oldServerConfig = oldServerConfigs[0].(map[string]interface{})
	}
	oldClientConfigs := make([]interface{}, 0)
	if _, ok := oldHttpsConfig["client_certificate_config"]; ok {
		oldClientConfigs = oldHttpsConfig["client_certificate_config"].([]interface{})
	}
	oldClientConfig := make(map[string]interface{})
	if len(oldClientConfigs) > 0 {
		oldClientConfig = oldClientConfigs[0].(map[string]interface{})
	}

	if domainConfig.Https.CertInfo != nil && domainConfig.Https.CertInfo.CertName != nil {
		serverCertConfigs := make([]map[string]interface{}, 0, 1)
		serverCertConfig := make(map[string]interface{}, 5)
		serverCertConfig["certificate_id"] = domainConfig.Https.CertInfo.CertId
		serverCertConfig["certificate_name"] = domainConfig.Https.CertInfo.CertName
		serverCertConfig["certificate_content"] = oldServerConfig["certificate_content"]
		serverCertConfig["private_key"] = oldServerConfig["private_key"]
		serverCertConfig["message"] = domainConfig.Https.CertInfo.Message
		serverCertConfig["deploy_time"] = domainConfig.Https.CertInfo.DeployTime
		serverCertConfig["expire_time"] = domainConfig.Https.CertInfo.ExpireTime
		serverCertConfigs = append(serverCertConfigs, serverCertConfig)
		httpsConfig["server_certificate_config"] = serverCertConfigs
	}
	if domainConfig.Https.ClientCertInfo != nil && domainConfig.Https.ClientCertInfo.CertName != nil {
		clientCertConfigs := make([]map[string]interface{}, 0, 1)
		clientCertConfig := make(map[string]interface{}, 2)
		clientCertConfig["certificate_content"] = oldClientConfig["certificate_content"]
		clientCertConfig["certificate_name"] = domainConfig.Https.ClientCertInfo.CertName
		clientCertConfig["deploy_time"] = domainConfig.Https.ClientCertInfo.DeployTime
		clientCertConfig["expire_time"] = domainConfig.Https.ClientCertInfo.ExpireTime
		clientCertConfigs = append(clientCertConfigs, clientCertConfig)
		httpsConfig["client_certificate_config"] = clientCertConfigs
	}
	if domainConfig.ForceRedirect != nil {
		httpsConfig["force_redirect"] = []map[string]interface{}{
			{
				"switch":               domainConfig.ForceRedirect.Switch,
				"redirect_type":        domainConfig.ForceRedirect.RedirectType,
				"redirect_status_code": domainConfig.ForceRedirect.RedirectStatusCode,
			},
		}
	}
	httpsConfigs = append(httpsConfigs, httpsConfig)
	_ = d.Set("https_config", httpsConfigs)

	authRaw := d.Get("authentication").([]interface{})
	if authentication := domainConfig.Authentication; authentication != nil && len(authRaw) > 0 {
		auth := make(map[string]interface{})
		auth["switch"] = authentication.Switch
		if authType := authentication.TypeA; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"sign_param":        authType.SignParam,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_a"] = []interface{}{dMap}
		}
		if authType := authentication.TypeB; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_b"] = []interface{}{dMap}
		}
		if authType := authentication.TypeC; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"time_format":       authType.TimeFormat,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_c"] = []interface{}{dMap}
		}
		if authType := authentication.TypeD; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"time_param":        authType.TimeParam,
				"time_format":       authType.TimeFormat,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_d"] = []interface{}{dMap}
		}
		_ = d.Set("authentication", []interface{}{auth})
	}

	tags, errRet := tagService.DescribeResourceTags(ctx, CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
	if errRet != nil {
		return errRet
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_domain.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	cdnService := CdnService{client: client}

	d.Partial(true)
	updateAttrs := make([]string, 0)

	domain := d.Id()
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = &domain

	if d.HasChange("service_type") {
		request.ServiceType = helper.String(d.Get("service_type").(string))
		updateAttrs = append(updateAttrs, "service_type")
	}
	if d.HasChange("project_id") {
		request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
		updateAttrs = append(updateAttrs, "project_id")
	}
	if d.HasChange("area") {
		request.Area = helper.String(d.Get("area").(string))
		updateAttrs = append(updateAttrs, "area")
	}
	if d.HasChange("full_url_cache") {
		fullUrlCache := d.Get("full_url_cache").(bool)
		request.CacheKey = &cdn.CacheKey{}
		if fullUrlCache {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_ON)
		} else {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_OFF)
		}
		updateAttrs = append(updateAttrs, "full_url_cache")
	}
	if d.HasChange("range_origin_switch") {
		request.RangeOriginPull = &cdn.RangeOriginPull{}
		request.RangeOriginPull.Switch = helper.String(d.Get("range_origin_switch").(string))
		updateAttrs = append(updateAttrs, "range_origin_switch")
	}
	if d.HasChange("ipv6_access_switch") {
		request.Ipv6Access = &cdn.Ipv6Access{}
		request.Ipv6Access.Switch = helper.String(d.Get("ipv6_access_switch").(string))
		updateAttrs = append(updateAttrs, "ipv6_access_switch")
	}
	if d.HasChange("follow_redirect_switch") {
		request.FollowRedirect = &cdn.FollowRedirect{}
		request.FollowRedirect.Switch = helper.String(d.Get("follow_redirect_switch").(string))
		updateAttrs = append(updateAttrs, "follow_redirect_switch")
	}
	if d.HasChange("origin") {
		updateAttrs = append(updateAttrs, "origin")
		origins := d.Get("origin").([]interface{})
		if len(origins) < 1 {
			return fmt.Errorf("origin is required")
		}
		origin := origins[0].(map[string]interface{})
		request.Origin = &cdn.Origin{}
		request.Origin.OriginType = helper.String(origin["origin_type"].(string))
		originList := origin["origin_list"].([]interface{})
		request.Origin.Origins = make([]*string, 0, len(originList))
		for _, item := range originList {
			request.Origin.Origins = append(request.Origin.Origins, helper.String(item.(string)))
		}
		if v := origin["server_name"]; v.(string) != "" {
			request.Origin.ServerName = helper.String(v.(string))
		}
		if v := origin["cos_private_access"]; v.(string) != "" {
			request.Origin.CosPrivateAccess = helper.String(v.(string))
		}
		if v := origin["origin_pull_protocol"]; v.(string) != "" {
			request.Origin.OriginPullProtocol = helper.String(v.(string))
		}
		if v := origin["backup_origin_type"]; v.(string) != "" {
			request.Origin.BackupOriginType = helper.String(v.(string))
		}
		if v := origin["backup_server_name"]; v.(string) != "" {
			request.Origin.BackupServerName = helper.String(v.(string))
		}
		if v := origin["backup_origin_list"]; len(v.([]interface{})) > 0 {
			backupOriginList := v.([]interface{})
			request.Origin.BackupOrigins = make([]*string, 0, len(backupOriginList))
			for _, item := range backupOriginList {
				request.Origin.BackupOrigins = append(request.Origin.BackupOrigins, helper.String(item.(string)))
			}
		}
	}
	if d.HasChange("request_header") {
		updateAttrs = append(updateAttrs, "request_header")
		requestHeaders := d.Get("request_header").([]interface{})
		requestHeader := requestHeaders[0].(map[string]interface{})
		headerRule := requestHeader["header_rules"].([]interface{})
		var headerRules []*cdn.HttpHeaderPathRule
		for _, value := range headerRule {
			hr := &cdn.HttpHeaderPathRule{}
			headerRuleMap := value.(map[string]interface{})
			headerMode := headerRuleMap["header_mode"].(string)
			headerName := headerRuleMap["header_name"].(string)
			headerValue := headerRuleMap["header_value"].(string)
			ruleType := headerRuleMap["rule_type"].(string)
			rulePaths := headerRuleMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			for _, value := range rulePaths {
				rulePathList = append(rulePathList, helper.String(value.(string)))
			}
			hr.HeaderMode = &headerMode
			hr.HeaderName = &headerName
			hr.HeaderValue = &headerValue
			hr.RuleType = &ruleType
			hr.RulePaths = rulePathList
			headerRules = append(headerRules, hr)
		}
		request.RequestHeader = &cdn.RequestHeader{}
		request.RequestHeader.Switch = helper.String(requestHeader["switch"].(string))
		request.RequestHeader.HeaderRules = headerRules
	}
	if d.HasChange("rule_cache") {
		updateAttrs = append(updateAttrs, "rule_cache")
		ruleCache := d.Get("rule_cache").([]interface{})
		var ruleCaches []*cdn.RuleCache
		for _, v := range ruleCache {
			re := &cdn.RuleCache{}
			ruleCacheMap := v.(map[string]interface{})
			rulePaths := ruleCacheMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			ruleType := ruleCacheMap["rule_type"].(string)
			if ruleType == CDN_RULE_TYPE_DEFAULT {
				rulePathList = append(rulePathList, helper.String(CDN_RULE_PATH))
			} else {
				for _, value := range rulePaths {
					rulePathList = append(rulePathList, helper.String(value.(string)))
				}
			}
			switchFlag := ruleCacheMap["switch"].(string)
			cacheTime := ruleCacheMap["cache_time"].(int)
			compareMaxAge := ruleCacheMap["compare_max_age"].(string)
			ignoreCacheControl := ruleCacheMap["ignore_cache_control"].(string)
			ignoreSetCookie := ruleCacheMap["ignore_set_cookie"].(string)
			noCacheSwitch := ruleCacheMap["no_cache_switch"].(string)
			reValidate := ruleCacheMap["re_validate"].(string)
			followOriginSwitch := ruleCacheMap["follow_origin_switch"].(string)
			ruleCacheConfig := &cdn.RuleCacheConfig{}
			cache := &cdn.CacheConfigCache{}
			noCache := &cdn.CacheConfigNoCache{}
			followOrigin := &cdn.CacheConfigFollowOrigin{}
			ruleCacheConfig.Cache = cache
			ruleCacheConfig.NoCache = noCache
			ruleCacheConfig.FollowOrigin = followOrigin
			re.CacheConfig = ruleCacheConfig
			re.RulePaths = rulePathList
			re.RuleType = &ruleType
			re.CacheConfig.Cache.Switch = &switchFlag
			re.CacheConfig.Cache.CacheTime = helper.IntInt64(cacheTime)
			re.CacheConfig.Cache.CompareMaxAge = &compareMaxAge
			re.CacheConfig.Cache.IgnoreCacheControl = &ignoreCacheControl
			re.CacheConfig.Cache.IgnoreSetCookie = &ignoreSetCookie
			re.CacheConfig.NoCache.Switch = &noCacheSwitch
			re.CacheConfig.NoCache.Revalidate = &reValidate
			re.CacheConfig.FollowOrigin.Switch = &followOriginSwitch
			ruleCaches = append(ruleCaches, re)
		}
		request.Cache = &cdn.Cache{}
		request.Cache.RuleCache = ruleCaches
	}
	if d.HasChange("https_config") {
		updateAttrs = append(updateAttrs, "https_config")
		httpsConfigs := d.Get("https_config").([]interface{})
		if len(httpsConfigs) > 0 {
			config := httpsConfigs[0].(map[string]interface{})
			request.Https = &cdn.Https{}
			request.Https.Switch = helper.String(config["https_switch"].(string))
			if v := config["http2_switch"]; v.(string) != "" {
				request.Https.Http2 = helper.String(v.(string))
			}
			request.Https.OcspStapling = helper.String(config["ocsp_stapling_switch"].(string))
			request.Https.Spdy = helper.String(config["spdy_switch"].(string))
			request.Https.VerifyClient = helper.String(config["verify_client"].(string))
			if v := config["server_certificate_config"]; len(v.([]interface{})) > 0 {
				serverCerts := v.([]interface{})
				if len(serverCerts) > 0 && serverCerts[0] != nil {
					serverCert := serverCerts[0].(map[string]interface{})
					request.Https.CertInfo = &cdn.ServerCert{}
					if v := serverCert["certificate_id"]; v.(string) != "" {
						request.Https.CertInfo.CertId = helper.String(v.(string))
					}
					if v := serverCert["certificate_content"]; v.(string) != "" {
						request.Https.CertInfo.Certificate = helper.String(v.(string))
					}
					if v := serverCert["private_key"]; v.(string) != "" {
						request.Https.CertInfo.PrivateKey = helper.String(v.(string))
					}
					if v := serverCert["message"]; v.(string) != "" {
						request.Https.CertInfo.Message = helper.String(v.(string))
					}
				}
			}
			if v := config["client_certificate_config"]; len(v.([]interface{})) > 0 {
				clientCerts := v.([]interface{})
				if len(clientCerts) > 0 && clientCerts[0] != nil {
					clientCert := clientCerts[0].(map[string]interface{})
					request.Https.ClientCertInfo = &cdn.ClientCert{}
					if v := clientCert["certificate_content"]; v.(string) != "" {
						request.Https.ClientCertInfo.Certificate = helper.String(v.(string))
					}
				}
			}
			if v, ok := config["force_redirect"]; ok {
				forceRedirect := v.([]interface{})
				if len(forceRedirect) > 0 {
					var redirect cdn.ForceRedirect
					redirectMap := forceRedirect[0].(map[string]interface{})
					if sw := redirectMap["switch"]; sw.(string) != "" {
						redirect.Switch = helper.String(sw.(string))
					}
					if rt := redirectMap["redirect_type"]; rt.(string) != "" {
						redirect.RedirectType = helper.String(rt.(string))
					}
					if rsc := redirectMap["redirect_status_code"]; rsc.(int) != 0 {
						redirect.RedirectStatusCode = helper.Int64(int64(rsc.(int)))
					}
					request.ForceRedirect = &redirect
				}
			}
		}
	}

	if d.HasChange("authentication") {
		updateAttrs = append(updateAttrs, "authentication")
		if v, ok := helper.InterfacesHeadMap(d, "authentication"); ok {
			switchOn := v["switch"].(string)
			request.Authentication = &cdn.Authentication{
				Switch: &switchOn,
			}

			if v, ok := v["type_a"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					signParam      = item["sign_param"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeA = &cdn.AuthenticationTypeA{
					SecretKey:      &secretKey,
					SignParam:      &signParam,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeA.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_b"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeB = &cdn.AuthenticationTypeB{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeB.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_c"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeC = &cdn.AuthenticationTypeC{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if timeFormat, ok := item["time_format"].(string); ok {
					request.Authentication.TypeC.TimeFormat = &timeFormat
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeC.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_d"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
					timeParam      = item["time_param"].(string)
				)

				request.Authentication.TypeD = &cdn.AuthenticationTypeD{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
					TimeParam:      &timeParam,
				}

				if timeFormat, ok := item["time_format"].(string); ok {
					request.Authentication.TypeD.TimeFormat = &timeFormat
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeD.BackupSecretKey = &backupSecretKey
				}
			}
		}
	}

	if len(updateAttrs) > 0 {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err := meta.(*TencentCloudClient).apiV3Conn.UseCdnClient().UpdateDomainConfig(request)
			if err != nil {
				if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROE {
						return resource.NonRetryableError(err)
					}
				}
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		for _, attr := range updateAttrs {
			d.SetPartial(attr)
		}

		err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
			domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
			if err != nil {
				return retryError(err, InternalError)
			}
			if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := TagService{client: client}
		region := client.Region
		resourceName := BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudCdnDomainRead(d, meta)
}

func resourceTencentCloudCdnDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_domain.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	cdnService := CdnService{client: client}

	domain := d.Id()
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: client}
		region := client.Region
		resourceName := BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		deleteTags := make([]string, 0, len(tags))
		for key := range tags {
			deleteTags = append(deleteTags, key)
		}
		err := tagService.ModifyTags(ctx, resourceName, nil, deleteTags)
		if err != nil {
			return err
		}
	}

	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if domainConfig == nil {
		return nil
	}

	if *domainConfig.Status == CDN_DOMAIN_STATUS_ONLINE {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet = cdnService.StopDomain(ctx, domain)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
			domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
			if err != nil {
				return retryError(err, InternalError)
			}
			if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet = cdnService.DeleteDomain(ctx, domain)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
