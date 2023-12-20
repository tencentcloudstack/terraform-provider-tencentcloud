package apigateway

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func ResourceTencentCloudAPIGatewayServiceRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayServiceReleaseCreate,
		Read:   resourceTencentCloudAPIGatewayServiceReleaseRead,
		Delete: resourceTencentCloudAPIGatewayServiceReleaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of API gateway service.",
			},
			"environment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(API_GATEWAY_SERVICE_ENVS),
				Description:  "API gateway service environment name to be released. Valid values: `test`, `prepub`, `release`.",
			},
			"release_desc": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "This release description of the API gateway service.",
			},
			"release_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The release version.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayServiceReleaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_service_release.create")()

	var (
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService    = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId            = d.Get("service_id").(string)
		environmentName      = d.Get("environment_name").(string)
		releaseDesc          = d.Get("release_desc").(string)
		err                  error
		has, existEnv        bool
		releaseResponse      *apigateway.ReleaseServiceResponse
		serviceResponse      apigateway.DescribeServiceResponse
		checkServiceResponse apigateway.DescribeServiceResponse
	)

	//check API gateway serviceid and service contains api
	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		checkServiceResponse, has, err = apiGatewayService.DescribeService(ctx, serviceId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("API gateway service %s not found on server", serviceId))
	}); err != nil {
		return err
	}
	if *checkServiceResponse.Response.ApiTotalCount < 1 {
		return fmt.Errorf("there is no API under the current service, please create an API before publishing")
	}

	// release API gateway service
	if releaseResponse, err = apiGatewayService.ReleaseService(ctx, serviceId, environmentName, releaseDesc); err != nil {
		return err
	}

	//wait service release ok
	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		serviceResponse, has, err = apiGatewayService.DescribeService(ctx, serviceId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("API gateway service %s not found on server", serviceId))
	}); err != nil {
		return err
	}

	for _, v := range serviceResponse.Response.AvailableEnvironments {
		if *v == environmentName {
			existEnv = true
			break
		}
	}

	if !existEnv {
		return fmt.Errorf("API gateway service not release success")
	}

	d.SetId(strings.Join([]string{serviceId, environmentName, *releaseResponse.Response.Result.ReleaseVersion}, tccommon.FILED_SP))

	return resourceTencentCloudAPIGatewayServiceReleaseRead(d, meta)
}

func resourceTencentCloudAPIGatewayServiceReleaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_service_release.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                = d.Id()
		info              []*apigateway.ServiceReleaseHistoryInfo
		err               error
	)

	ids := strings.Split(id, tccommon.FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", id)
	}
	var (
		serviceId  = ids[0]
		envName    = ids[1]
		envVersion = ids[2]
	)

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, _, err = apiGatewayService.DescribeServiceEnvironmentReleaseHistory(ctx, serviceId, envName)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	if len(info) == 0 {
		d.SetId("")
		return nil
	}

	for _, v := range info {
		if *v.VersionName == envVersion {
			_ = d.Set("service_id", serviceId)
			_ = d.Set("environment_name", envName)
			_ = d.Set("release_desc", v.VersionDesc)
			_ = d.Set("release_version", v.VersionName)
		}
	}

	return nil
}

func resourceTencentCloudAPIGatewayServiceReleaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_service_release.delete")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                = d.Id()
		err               error
	)

	ids := strings.Split(id, tccommon.FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", id)
	}
	var (
		serviceId = ids[0]
		envName   = ids[1]
	)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if err = apiGatewayService.UnReleaseService(ctx, serviceId, envName); err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
