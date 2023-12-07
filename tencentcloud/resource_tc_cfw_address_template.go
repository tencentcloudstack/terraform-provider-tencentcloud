package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCfwAddressTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwAddressTemplateCreate,
		Read:   resourceTencentCloudCfwAddressTemplateRead,
		Update: resourceTencentCloudCfwAddressTemplateUpdate,
		Delete: resourceTencentCloudCfwAddressTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name.",
			},
			"detail": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template Detail.",
			},
			"ip_string": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Type is 1, ip template eg: 1.1.1.1,2.2.2.2; Type is 5, domain name template eg: www.qq.com, www.tencent.com.",
			},
			"type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue(ADDRESS_TEMPLATE_TYPE),
				Description:  "1: ip template; 5: domain name templates.",
			},
		},
	}
}

func resourceTencentCloudCfwAddressTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_address_template.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = cfw.NewCreateAddressTemplateRequest()
		response = cfw.NewCreateAddressTemplateResponse()
		uuid     string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("detail"); ok {
		request.Detail = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_string"); ok {
		request.IpString = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().CreateAddressTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw addressTemplate failed, reason:%+v", logId, err)
		return err
	}

	uuid = *response.Response.Uuid
	d.SetId(uuid)

	return resourceTencentCloudCfwAddressTemplateRead(d, meta)
}

func resourceTencentCloudCfwAddressTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_address_template.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		uuid    = d.Id()
	)

	addressTemplate, err := service.DescribeCfwAddressTemplateById(ctx, uuid)
	if err != nil {
		return err
	}

	if addressTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwAddressTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if addressTemplate.Name != nil {
		_ = d.Set("name", addressTemplate.Name)
	}

	if addressTemplate.Detail != nil {
		_ = d.Set("detail", addressTemplate.Detail)
	}

	if addressTemplate.IpString != nil {
		_ = d.Set("ip_string", addressTemplate.IpString)
	}

	if addressTemplate.Type != nil {
		_ = d.Set("type", addressTemplate.Type)
	}

	return nil
}

func resourceTencentCloudCfwAddressTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_address_template.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = cfw.NewModifyAddressTemplateRequest()
		uuid    = d.Id()
	)

	request.Uuid = &uuid

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("detail"); ok {
		request.Detail = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_string"); ok {
		request.IpString = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyAddressTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw addressTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwAddressTemplateRead(d, meta)
}

func resourceTencentCloudCfwAddressTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_address_template.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		uuid    = d.Id()
	)

	if err := service.DeleteCfwAddressTemplateById(ctx, uuid); err != nil {
		return err
	}

	return nil
}
