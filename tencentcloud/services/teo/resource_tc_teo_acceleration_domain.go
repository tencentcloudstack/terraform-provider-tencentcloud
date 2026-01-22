package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoAccelerationDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoAccelerationDomainCreate,
		Read:   resourceTencentCloudTeoAccelerationDomainRead,
		Update: resourceTencentCloudTeoAccelerationDomainUpdate,
		Delete: resourceTencentCloudTeoAccelerationDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the site related with the accelerated domain name.",
			},

			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Accelerated domain name.",
			},

			"origin_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Details of the origin.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin server type, with values: IP_DOMAIN: IPv4, IPv6, or domain name type origin server; COS: Tencent Cloud COS origin server; AWS_S3: AWS S3 origin server; ORIGIN_GROUP: origin server group type origin server; VOD: Video on Demand; SPACE: origin server uninstallation. Currently only available to the allowlist; LB: load balancing. Currently only available to the allowlist.",
						},
						"origin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin server address, which varies according to the value of OriginType: When OriginType = IP_DOMAIN, fill in an IPv4 address, an IPv6 address, or a domain name; When OriginType = COS, fill in the access domain name of the COS bucket; When OriginType = AWS_S3, fill in the access domain name of the S3 bucket; When OriginType = ORIGIN_GROUP, fill in the origin server group ID; When OriginType = VOD, fill in the VOD application ID; When OriginType = LB, fill in the Cloud Load Balancer instance ID. This feature is currently only available to the allowlist; When OriginType = SPACE, fill in the origin server uninstallation space ID. This feature is currently only available to the allowlist.",
						},
						"backup_origin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the secondary origin group. This parameter is valid only when OriginType is ORIGIN_GROUP. This field indicates the old version capability, which cannot be configured or modified on the control panel after being called. Please submit a ticket if required.",
						},
						"private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether access to the private Cloud Object Storage origin server is allowed. This parameter is valid only when OriginType is COS or AWS_S3. Valid values: on: Enable private authentication; off: Disable private authentication. If it is not specified, the default value is off.",
						},
						"private_parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Private authentication parameter. This parameter is valid only when `private_access` is on.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameter name. Valid values: `AccessKeyId`: Access Key ID; `SecretAccessKey`: Secret Access Key; `SignatureVersion`: authentication version, v2 or v4; `Region`: bucket region.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameter value.",
									},
								},
							},
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Custom origin server HOST header. this parameter is valid only when OriginType=IP_DOMAIN.If the OriginType is another type of origin, this parameter does not need to be passed in, otherwise an error will be reported. If OriginType is COS or AWS_S3, the HOST header for origin-pull will remain consistent with the origin server domain name. If OriginType is ORIGIN_GROUP, the HOST header follows the ORIGIN site GROUP configuration. if not configured, it defaults to the acceleration domain name. If OriginType is VOD or SPACE, no configuration is required for this header, and the domain name takes effect based on the corresponding origin.",
						},
						"vod_origin_scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The scope of cloud on-demand back-to-source. This parameter is effective when OriginType = VOD. The possible values are: all: all files in the cloud on-demand application corresponding to the current origin station. The default value is all; bucket: files in a specified bucket under the cloud on-demand application corresponding to the current origin station. The bucket is specified by the parameter VodBucketId.",
						},
						"vod_bucket_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "VOD bucket ID. This parameter is required when OriginType = VOD and VodOriginScope = bucket. Data source: the storage ID of the bucket in the Cloud VOD Professional Edition application.",
						},
					},
				},
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"online", "offline"}),
				Description:  "Accelerated domain name status, the values are: `online`: enabled; `offline`: disabled. Default is `online`.",
			},

			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Origin return protocol, possible values are: `FOLLOW`: protocol follow; `HTTP`: HTTP protocol back to source; `HTTPS`: HTTPS protocol back to source. If not filled in, the default is: `FOLLOW`.",
			},

			"http_origin_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "HTTP back-to-origin port, the value is 1-65535, effective when OriginProtocol=FOLLOW/HTTP, if not filled in, the default value is 80.",
			},

			"https_origin_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "HTTPS back-to-origin port. The value range is 1-65535. It takes effect when OriginProtocol=FOLLOW/HTTPS. If it is not filled in, the default value is 443.",
			},

			"ipv6_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv6 status, the value is: `follow`: follow the site IPv6 configuration; `on`: on; `off`: off. If not filled in, the default is: `follow`.",
			},

			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME address.",
			},
		},
	}
}

func resourceTencentCloudTeoAccelerationDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_acceleration_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = teo.NewCreateAccelerationDomainRequest()
		response   = teo.NewCreateAccelerationDomainResponse()
		zoneId     string
		domainName string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("domain_name"); ok {
		request.DomainName = helper.String(v.(string))
		domainName = v.(string)
	}

	if originInfoMap, ok := helper.InterfacesHeadMap(d, "origin_info"); ok {
		originInfo := teo.OriginInfo{}
		var originType string
		if v, ok := originInfoMap["origin_type"]; ok {
			originInfo.OriginType = helper.String(v.(string))
			originType = v.(string)
		}

		if v, ok := originInfoMap["origin"]; ok {
			originInfo.Origin = helper.String(v.(string))
		}

		if v, ok := originInfoMap["backup_origin"]; ok {
			originInfo.BackupOrigin = helper.String(v.(string))
		}

		if v, ok := originInfoMap["private_access"]; ok {
			originInfo.PrivateAccess = helper.String(v.(string))
		}

		if v, ok := originInfoMap["private_parameters"]; ok {
			for _, item := range v.([]interface{}) {
				privateParametersMap := item.(map[string]interface{})
				privateParameter := teo.PrivateParameter{}
				if v, ok := privateParametersMap["name"]; ok {
					privateParameter.Name = helper.String(v.(string))
				}

				if v, ok := privateParametersMap["value"]; ok {
					privateParameter.Value = helper.String(v.(string))
				}

				originInfo.PrivateParameters = append(originInfo.PrivateParameters, &privateParameter)
			}
		}

		if v, ok := originInfoMap["host_header"].(string); ok && v != "" {
			if originType == "IP_DOMAIN" {
				originInfo.HostHeader = helper.String(v)
			} else {
				return fmt.Errorf("Only `origin_type` is `IP_DOMAIN` can set `host_header`.")
			}
		}

		if v, ok := originInfoMap["vod_origin_scope"].(string); ok && v != "" {
			originInfo.VodOriginScope = helper.String(v)
		}

		if v, ok := originInfoMap["vod_bucket_id"].(string); ok && v != "" {
			originInfo.VodBucketId = helper.String(v)
		}

		request.OriginInfo = &originInfo
	}

	if v, ok := d.GetOk("origin_protocol"); ok {
		request.OriginProtocol = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("http_origin_port"); ok {
		request.HttpOriginPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("https_origin_port"); ok {
		request.HttpsOriginPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("ipv6_status"); ok {
		request.IPv6Status = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateAccelerationDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo acceleration domain failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo acceleration domain failed, reason:%+v", logId, err)
		return err
	}

	// wait
	if err := resourceTencentCloudTeoAccelerationDomainCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{zoneId, domainName}, tccommon.FILED_SP))

	// offline
	if v, ok := d.GetOk("status"); ok {
		if v.(string) == "offline" {
			request := teo.NewModifyAccelerationDomainStatusesRequest()
			request.ZoneId = helper.String(zoneId)
			request.DomainNames = []*string{helper.String(domainName)}
			if v, ok := d.GetOk("status"); ok {
				request.Status = helper.String(v.(string))
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAccelerationDomainStatusesWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update teo acceleration domain status failed, reason:%+v", logId, err)
				return err
			}

			// wait
			if err = resourceTencentCloudTeoAccelerationDomainUpdateOnExit(ctx); err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudTeoAccelerationDomainRead(d, meta)
}

func resourceTencentCloudTeoAccelerationDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_acceleration_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	domainName := idSplit[1]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("domain_name", domainName)

	respData, err := service.DescribeTeoAccelerationDomainById(ctx, zoneId, domainName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_acceleration_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.DomainName != nil {
		_ = d.Set("domain_name", respData.DomainName)
	}

	if respData.OriginDetail != nil {
		originDetailMap := map[string]interface{}{}
		if respData.OriginDetail.OriginType != nil {
			originDetailMap["origin_type"] = respData.OriginDetail.OriginType
		}

		if respData.OriginDetail.Origin != nil {
			originDetailMap["origin"] = respData.OriginDetail.Origin
		}

		if respData.OriginDetail.BackupOrigin != nil {
			originDetailMap["backup_origin"] = respData.OriginDetail.BackupOrigin
		}

		if respData.OriginDetail.PrivateAccess != nil {
			originDetailMap["private_access"] = respData.OriginDetail.PrivateAccess
		}

		if respData.OriginDetail.PrivateParameters != nil {
			privateParametersList := make([]map[string]interface{}, 0, len(respData.OriginDetail.PrivateParameters))
			for _, privateParameters := range respData.OriginDetail.PrivateParameters {
				privateParametersMap := map[string]interface{}{}
				if privateParameters.Name != nil {
					privateParametersMap["name"] = privateParameters.Name
				}

				if privateParameters.Value != nil {
					privateParametersMap["value"] = privateParameters.Value
				}

				privateParametersList = append(privateParametersList, privateParametersMap)
			}

			originDetailMap["private_parameters"] = privateParametersList
		}

		if respData.OriginDetail.HostHeader != nil {
			originDetailMap["host_header"] = respData.OriginDetail.HostHeader
		}

		if respData.OriginDetail.VodOriginScope != nil {
			originDetailMap["vod_origin_scope"] = respData.OriginDetail.VodOriginScope
		}

		if respData.OriginDetail.VodBucketId != nil {
			originDetailMap["vod_bucket_id"] = respData.OriginDetail.VodBucketId
		}

		_ = d.Set("origin_info", []interface{}{originDetailMap})
	}

	if respData.DomainStatus != nil {
		_ = d.Set("status", respData.DomainStatus)
	}

	if respData.OriginProtocol != nil {
		_ = d.Set("origin_protocol", respData.OriginProtocol)
	}

	if respData.HttpOriginPort != nil {
		_ = d.Set("http_origin_port", respData.HttpOriginPort)
	}

	if respData.HttpsOriginPort != nil {
		_ = d.Set("https_origin_port", respData.HttpsOriginPort)
	}

	if respData.IPv6Status != nil {
		_ = d.Set("ipv6_status", respData.IPv6Status)
	}

	if respData.Cname != nil {
		_ = d.Set("cname", respData.Cname)
	}

	return nil
}

func resourceTencentCloudTeoAccelerationDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_acceleration_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	immutableArgs := []string{"https_origin_port"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	domainName := idSplit[1]

	needChange := false
	mutableArgs := []string{"origin_info", "origin_protocol", "http_origin_port", "https_origin_port", "ipv6_status"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teo.NewModifyAccelerationDomainRequest()
		request.ZoneId = helper.String(zoneId)
		request.DomainName = helper.String(domainName)
		if originInfoMap, ok := helper.InterfacesHeadMap(d, "origin_info"); ok {
			originInfo := teo.OriginInfo{}
			var originType string
			if v, ok := originInfoMap["origin_type"]; ok {
				originInfo.OriginType = helper.String(v.(string))
				originType = v.(string)
			}

			if v, ok := originInfoMap["origin"]; ok {
				originInfo.Origin = helper.String(v.(string))
			}

			if v, ok := originInfoMap["backup_origin"]; ok {
				originInfo.BackupOrigin = helper.String(v.(string))
			}

			if v, ok := originInfoMap["private_access"]; ok {
				originInfo.PrivateAccess = helper.String(v.(string))
			}

			if v, ok := originInfoMap["private_parameters"]; ok {
				for _, item := range v.([]interface{}) {
					privateParametersMap := item.(map[string]interface{})
					privateParameter := teo.PrivateParameter{}
					if v, ok := privateParametersMap["name"]; ok {
						privateParameter.Name = helper.String(v.(string))
					}

					if v, ok := privateParametersMap["value"]; ok {
						privateParameter.Value = helper.String(v.(string))
					}

					originInfo.PrivateParameters = append(originInfo.PrivateParameters, &privateParameter)
				}
			}

			if v, ok := originInfoMap["host_header"].(string); ok && v != "" {
				if originType == "IP_DOMAIN" {
					originInfo.HostHeader = helper.String(v)
				}
			}

			if v, ok := originInfoMap["vod_origin_scope"].(string); ok && v != "" {
				originInfo.VodOriginScope = helper.String(v)
			}

			if v, ok := originInfoMap["vod_bucket_id"].(string); ok && v != "" {
				originInfo.VodBucketId = helper.String(v)
			}

			request.OriginInfo = &originInfo
		}

		if v, ok := d.GetOk("origin_protocol"); ok {
			request.OriginProtocol = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("http_origin_port"); ok {
			request.HttpOriginPort = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("https_origin_port"); ok {
			request.HttpsOriginPort = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("ipv6_status"); ok {
			request.IPv6Status = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAccelerationDomainWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo acceleration domain failed, reason:%+v", logId, err)
			return err
		}

		// wait
		if err := resourceTencentCloudTeoAccelerationDomainUpdateOnExit(ctx); err != nil {
			return err
		}
	}

	if d.HasChange("status") {
		request := teo.NewModifyAccelerationDomainStatusesRequest()
		request.ZoneId = helper.String(zoneId)
		request.DomainNames = []*string{helper.String(domainName)}
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAccelerationDomainStatusesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo acceleration domain status failed, reason:%+v", logId, err)
			return err
		}

		// wait
		if err := resourceTencentCloudTeoAccelerationDomainUpdateOnExit(ctx); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoAccelerationDomainRead(d, meta)
}

func resourceTencentCloudTeoAccelerationDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_acceleration_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teo.NewModifyAccelerationDomainStatusesRequest()
		response = teo.NewModifyAccelerationDomainStatusesResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	domainName := idSplit[1]

	// check offline first
	if v, ok := d.GetOk("status"); ok {
		if v.(string) == "online" {
			request.ZoneId = helper.String(zoneId)
			request.DomainNames = []*string{helper.String(domainName)}
			request.Status = helper.String("offline")
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAccelerationDomainStatusesWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify teo acceleration domain status failed, Response is nil."))
				}

				response = result
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify teo acceleration domain status failed, reason:%+v", logId, err)
				return err
			}

			// wait
			if err := resourceTencentCloudTeoAccelerationDomainDeletePostHandleResponse0(ctx, response); err != nil {
				return err
			}
		}
	}

	// delete
	delRequest := teo.NewDeleteAccelerationDomainsRequest()
	delRequest.ZoneId = helper.String(zoneId)
	delRequest.DomainNames = []*string{helper.String(domainName)}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteAccelerationDomainsWithContext(ctx, delRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, delRequest.GetAction(), delRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo acceleration domain failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo acceleration domain failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
