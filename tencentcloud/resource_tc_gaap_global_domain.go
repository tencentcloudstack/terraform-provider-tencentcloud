package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudGaapGlobalDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapGlobalDomainCreate,
		Read:   resourceTencentCloudGaapGlobalDomainRead,
		Update: resourceTencentCloudGaapGlobalDomainUpdate,
		Delete: resourceTencentCloudGaapGlobalDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Domain Name Project ID.",
			},

			"default_value": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain name default entry.",
			},

			"alias": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "alias.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tags.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{GLOBAL_DOMAIN_STATUS_OPEN, GLOBAL_DOMAIN_STATUS_CLOSE}),
				Description:  "Global domain statue. Available values: open and close, default is open.",
			},
		},
	}
}

func resourceTencentCloudGaapGlobalDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_global_domain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request   = gaap.NewCreateGlobalDomainRequest()
		response  = gaap.NewCreateGlobalDomainResponse()
		domainId  string
		projectId int
	)
	if v, ok := d.GetOkExists("project_id"); ok {
		projectId = v.(int)
		request.ProjectId = helper.IntInt64(projectId)
	}

	if v, ok := d.GetOk("default_value"); ok {
		request.DefaultValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alias"); ok {
		request.Alias = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagPair := gaap.TagPair{}
			if v, ok := dMap["tag_key"]; ok {
				tagPair.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				tagPair.TagValue = helper.String(v.(string))
			}
			request.TagSet = append(request.TagSet, &tagPair)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseGaapClient().CreateGlobalDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create gaap globalDomain failed, reason:%+v", logId, err)
		return err
	}

	domainId = *response.Response.DomainId

	d.SetId(strconv.Itoa(projectId) + FILED_SP + domainId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("gaap", "domain", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapGlobalDomainRead(d, meta)
}

func resourceTencentCloudGaapGlobalDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_global_domain.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	projectId, err := strconv.Atoi(idSplit[0])
	if err != nil {
		return err
	}
	domainId := idSplit[1]

	globalDomain, err := service.DescribeGaapGlobalDomainById(ctx, domainId, projectId)
	if err != nil {
		return err
	}

	if globalDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `GaapGlobalDomain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if globalDomain.ProjectId != nil {
		_ = d.Set("project_id", globalDomain.ProjectId)
	}

	if globalDomain.DefaultValue != nil {
		_ = d.Set("default_value", globalDomain.DefaultValue)
	}

	if globalDomain.Alias != nil {
		_ = d.Set("alias", globalDomain.Alias)
	}

	if globalDomain.Status != nil {
		statusInt := int(*globalDomain.Status)
		if statusInt == 0 {
			_ = d.Set("status", GLOBAL_DOMAIN_STATUS_OPEN)
		} else if statusInt == 1 {
			_ = d.Set("status", GLOBAL_DOMAIN_STATUS_CLOSE)
		}
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "gaap", "domain", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudGaapGlobalDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_global_domain.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := gaap.NewModifyGlobalDomainAttributeRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	projectId, err := strconv.Atoi(idSplit[0])
	if err != nil {
		return err
	}
	domainId := idSplit[1]

	request.DomainId = &domainId
	request.ProjectId = helper.IntUint64(projectId)
	immutableArgs := []string{"project_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	isAttributeChange := false
	if d.HasChange("default_value") {
		if v, ok := d.GetOk("default_value"); ok {
			request.DefaultValue = helper.String(v.(string))
			isAttributeChange = true
		}
	}

	if d.HasChange("alias") {
		if v, ok := d.GetOk("alias"); ok {
			request.Alias = helper.String(v.(string))
			isAttributeChange = true
		}
	}
	if isAttributeChange {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseGaapClient().ModifyGlobalDomainAttribute(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update gaap globalDomain failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("gaap", "domain", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}
			if status == GLOBAL_DOMAIN_STATUS_OPEN {
				if err := service.EnableGlobalDomain(ctx, domainId); err != nil {
					return err
				}

				conf := BuildStateChangeConf([]string{}, []string{"0"}, 1*readRetryTimeout, time.Second, service.DomainInstanceStateRefreshFunc(domainId, projectId, []string{}))
				if _, e := conf.WaitForState(); e != nil {
					return e
				}
			}
			if status == GLOBAL_DOMAIN_STATUS_CLOSE {
				if err := service.DisableGlobalDomain(ctx, domainId); err != nil {
					return err
				}

				conf := BuildStateChangeConf([]string{}, []string{"1"}, 1*readRetryTimeout, time.Second, service.DomainInstanceStateRefreshFunc(domainId, projectId, []string{}))
				if _, e := conf.WaitForState(); e != nil {
					return e
				}
			}
		}
	}

	return resourceTencentCloudGaapGlobalDomainRead(d, meta)
}

func resourceTencentCloudGaapGlobalDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_global_domain.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId, err := strconv.Atoi(idSplit[0])
	if err != nil {
		return err
	}
	domainId := idSplit[1]

	globalDomain, err := service.DescribeGaapGlobalDomainById(ctx, domainId, projectId)
	if err != nil {
		return err
	}

	if globalDomain != nil && globalDomain.Status != nil && int(*globalDomain.Status) == 0 {
		if err := service.DisableGlobalDomain(ctx, domainId); err != nil {
			return err
		}
		conf := BuildStateChangeConf([]string{}, []string{"1"}, 1*readRetryTimeout, time.Second, service.DomainInstanceStateRefreshFunc(domainId, projectId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if err := service.DeleteGaapGlobalDomainById(ctx, domainId); err != nil {
		return err
	}

	return nil
}
