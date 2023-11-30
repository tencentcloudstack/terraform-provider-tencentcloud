package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoAccelerationDomain() *schema.Resource {
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
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the site related with the accelerated domain name.",
			},

			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Accelerated domain name.",
			},

			"origin_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Details of the origin.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The origin type. Values: `IP_DOMAIN`: IPv4/IPv6 address or domain name; `COS`: COS bucket address; `ORIGIN_GROUP`: Origin group; `AWS_S3`: AWS S3 bucket address; `SPACE`: EdgeOne Shield Space.",
						},
						"origin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The origin address. Enter the origin group ID if `OriginType=ORIGIN_GROUP`.",
						},
						"backup_origin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the secondary origin group (valid when `OriginType=ORIGIN_GROUP`). If it is not specified, it indicates that secondary origins are not used.",
						},
						"private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to authenticate access to the private object storage origin (valid when `OriginType=COS/AWS_S3`). Values: `on`: Enable private authentication; `off`: Disable private authentication. If this field is not specified, the default value `off` is used.",
						},
						"private_parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The private authentication parameters. This field is valid when `PrivateAccess=on`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameter name. Valid values: `AccessKeyId`: Access Key ID; `SecretAccessKey`: Secret Access Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameter value.",
									},
								},
							},
						},
					},
				},
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Accelerated domain name status, the values are: `online`: enabled; `offline`: disabled.",
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
	defer logElapsed("resource.tencentcloud_teo_acceleration_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = teo.NewCreateAccelerationDomainRequest()
		zoneId     string
		domainName string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "origin_info"); ok {
		originInfo := teo.OriginInfo{}
		if v, ok := dMap["origin_type"]; ok {
			originInfo.OriginType = helper.String(v.(string))
		}
		if v, ok := dMap["origin"]; ok {
			originInfo.Origin = helper.String(v.(string))
		}
		if v, ok := dMap["backup_origin"]; ok {
			originInfo.BackupOrigin = helper.String(v.(string))
		}
		if v, ok := dMap["private_access"]; ok {
			originInfo.PrivateAccess = helper.String(v.(string))
		}
		if v, ok := dMap["private_parameters"]; ok {
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
		request.OriginInfo = &originInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateAccelerationDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo accelerationDomain failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(zoneId + FILED_SP + domainName)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = service.CheckAccelerationDomainStatus(ctx, zoneId, domainName, "")
	if err != nil {
		return err
	}

	return resourceTencentCloudTeoAccelerationDomainRead(d, meta)
}

func resourceTencentCloudTeoAccelerationDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_acceleration_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	domainName := idSplit[1]

	accelerationDomain, err := service.DescribeTeoAccelerationDomainById(ctx, zoneId, domainName)
	if err != nil {
		return err
	}

	if accelerationDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoAccelerationDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if accelerationDomain.ZoneId != nil {
		_ = d.Set("zone_id", accelerationDomain.ZoneId)
	}

	if accelerationDomain.DomainName != nil {
		_ = d.Set("domain_name", accelerationDomain.DomainName)
	}

	if accelerationDomain.Cname != nil {
		_ = d.Set("cname", accelerationDomain.Cname)
	}

	if accelerationDomain.OriginDetail != nil {
		originInfoMap := map[string]interface{}{}
		originDetail := accelerationDomain.OriginDetail

		if originDetail.OriginType != nil {
			originInfoMap["origin_type"] = originDetail.OriginType
		}

		if originDetail.Origin != nil {
			originInfoMap["origin"] = originDetail.Origin
		}

		if originDetail.BackupOrigin != nil {
			originInfoMap["backup_origin"] = originDetail.BackupOrigin
		}

		if originDetail.PrivateAccess != nil {
			originInfoMap["private_access"] = originDetail.PrivateAccess
		}

		if originDetail.PrivateParameters != nil {
			privateParametersList := []interface{}{}
			for _, privateParameters := range originDetail.PrivateParameters {
				privateParametersMap := map[string]interface{}{}

				if privateParameters.Name != nil {
					privateParametersMap["name"] = privateParameters.Name
				}

				if privateParameters.Value != nil {
					privateParametersMap["value"] = privateParameters.Value
				}

				privateParametersList = append(privateParametersList, privateParametersMap)
			}

			originInfoMap["private_parameters"] = privateParametersList
		}

		_ = d.Set("origin_info", []interface{}{originInfoMap})
	}

	return nil
}

func resourceTencentCloudTeoAccelerationDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_acceleration_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	domainName := idSplit[1]

	if d.HasChange("origin_info") {
		request := teo.NewModifyAccelerationDomainRequest()
		request.ZoneId = &zoneId
		request.DomainName = &domainName

		if dMap, ok := helper.InterfacesHeadMap(d, "origin_info"); ok {
			originInfo := teo.OriginInfo{}
			if v, ok := dMap["origin_type"]; ok {
				originInfo.OriginType = helper.String(v.(string))
			}
			if v, ok := dMap["origin"]; ok {
				originInfo.Origin = helper.String(v.(string))
			}
			if v, ok := dMap["backup_origin"]; ok {
				originInfo.BackupOrigin = helper.String(v.(string))
			}
			if v, ok := dMap["private_access"]; ok {
				originInfo.PrivateAccess = helper.String(v.(string))
			}
			if v, ok := dMap["private_parameters"]; ok {
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
			request.OriginInfo = &originInfo
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyAccelerationDomain(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo accelerationDomain failed, reason:%+v", logId, err)
			return err
		}

		service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
		err = service.CheckAccelerationDomainStatus(ctx, zoneId, domainName, "")
		if err != nil {
			return err
		}
	}

	if d.HasChange("status") {
		request := teo.NewModifyAccelerationDomainStatusesRequest()
		request.ZoneId = &zoneId
		request.DomainNames = []*string{&domainName}

		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyAccelerationDomainStatuses(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo accelerationDomain failed, reason:%+v", logId, err)
			return err
		}

		service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
		err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := service.DescribeTeoAccelerationDomainById(ctx, zoneId, domainName)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			if *instance.DomainStatus == "online" || *instance.DomainStatus == "offline" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("AccelerationDomain status is %v, retry...", *instance.DomainStatus))
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoAccelerationDomainRead(d, meta)
}

func resourceTencentCloudTeoAccelerationDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_acceleration_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	domainName := idSplit[1]

	request := teo.NewModifyAccelerationDomainStatusesRequest()
	request.ZoneId = &zoneId
	request.DomainNames = []*string{&domainName}
	request.Status = helper.String("offline")

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyAccelerationDomainStatuses(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo accelerationDomain failed, reason:%+v", logId, err)
		return err
	}

	err = service.CheckAccelerationDomainStatus(ctx, zoneId, domainName, "delete")
	if err != nil {
		return err
	}

	if err = service.DeleteTeoAccelerationDomainById(ctx, zoneId, domainName); err != nil {
		return err
	}

	return nil
}
