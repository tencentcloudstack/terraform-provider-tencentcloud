/*
Provides a resource to create a tse cngw_route

Example Usage

```hcl
resource "tencentcloud_tse_cngw_route" "cngw_route" {
  gateway_id = "gateway-xxxxxx"
  service_i_d = "451a9920-e67a-4519-af41-fccac0e72005"
  route_name = "routeA"
  methods =
  hosts =
  paths =
  protocols =
  preserve_host = true
  https_redirect_status_code = 302
  strip_path = true
  force_https =
  destination_ports =
  headers {
		key = "token"
		value = "xxxxxx"

  }
}
```

Import

tse cngw_route can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_route.cngw_route cngw_route_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTseCngwRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwRouteCreate,
		Read:   resourceTencentCloudTseCngwRouteRead,
		Update: resourceTencentCloudTseCngwRouteUpdate,
		Delete: resourceTencentCloudTseCngwRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"service_i_d": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the service which the route belongs to.",
			},

			"route_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The name of the route, unique in the instance.",
			},

			"methods": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Route methods. Reference value:- GET- POST- DELETE- PUT- OPTIONS- PATCH- HEAD- ANY- TRACE- COPY- MOVE- PROPFIND- PROPPATCH- MKCOL- LOCK- UNLOCK.",
			},

			"hosts": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Host list.",
			},

			"paths": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Path list.",
			},

			"protocols": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The protocol list of route.Reference value:- https- http.",
			},

			"preserve_host": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to keep the host when forwarding to the backend.",
			},

			"https_redirect_status_code": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Https redirection status code.",
			},

			"strip_path": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to strip path when forwarding to the backend.",
			},

			"force_https": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable forced HTTPS, no longer use.",
			},

			"destination_ports": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Destination port for Layer 4 matching.",
			},

			"headers": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The headers of route.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key of header.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value of header.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseCngwRouteCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayRouteRequest()
		response  = tse.NewCreateCloudNativeAPIGatewayRouteResponse()
		gatewayId string
		routeID   string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_i_d"); ok {
		request.ServiceID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_name"); ok {
		request.RouteName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("methods"); ok {
		methodsSet := v.(*schema.Set).List()
		for i := range methodsSet {
			methods := methodsSet[i].(string)
			request.Methods = append(request.Methods, &methods)
		}
	}

	if v, ok := d.GetOk("hosts"); ok {
		hostsSet := v.(*schema.Set).List()
		for i := range hostsSet {
			hosts := hostsSet[i].(string)
			request.Hosts = append(request.Hosts, &hosts)
		}
	}

	if v, ok := d.GetOk("paths"); ok {
		pathsSet := v.(*schema.Set).List()
		for i := range pathsSet {
			paths := pathsSet[i].(string)
			request.Paths = append(request.Paths, &paths)
		}
	}

	if v, ok := d.GetOk("protocols"); ok {
		protocolsSet := v.(*schema.Set).List()
		for i := range protocolsSet {
			protocols := protocolsSet[i].(string)
			request.Protocols = append(request.Protocols, &protocols)
		}
	}

	if v, ok := d.GetOkExists("preserve_host"); ok {
		request.PreserveHost = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("https_redirect_status_code"); ok {
		request.HttpsRedirectStatusCode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("strip_path"); ok {
		request.StripPath = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("force_https"); ok {
		request.ForceHttps = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("destination_ports"); ok {
		destinationPortsSet := v.(*schema.Set).List()
		for i := range destinationPortsSet {
			destinationPorts := destinationPortsSet[i].(int)
			request.DestinationPorts = append(request.DestinationPorts, helper.IntUint64(destinationPorts))
		}
	}

	if v, ok := d.GetOk("headers"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			kVMapping := tse.KVMapping{}
			if v, ok := dMap["key"]; ok {
				kVMapping.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				kVMapping.Value = helper.String(v.(string))
			}
			request.Headers = append(request.Headers, &kVMapping)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateCloudNativeAPIGatewayRoute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwRoute failed, reason:%+v", logId, err)
		return err
	}

	gatewayId = *response.Response.GatewayId
	d.SetId(strings.Join([]string{gatewayId, routeID}, FILED_SP))

	return resourceTencentCloudTseCngwRouteRead(d, meta)
}

func resourceTencentCloudTseCngwRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeID := idSplit[1]

	cngwRoute, err := service.DescribeTseCngwRouteById(ctx, gatewayId, routeID)
	if err != nil {
		return err
	}

	if cngwRoute == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwRoute` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cngwRoute.GatewayId != nil {
		_ = d.Set("gateway_id", cngwRoute.GatewayId)
	}

	if cngwRoute.ServiceID != nil {
		_ = d.Set("service_i_d", cngwRoute.ServiceID)
	}

	if cngwRoute.RouteName != nil {
		_ = d.Set("route_name", cngwRoute.RouteName)
	}

	if cngwRoute.Methods != nil {
		_ = d.Set("methods", cngwRoute.Methods)
	}

	if cngwRoute.Hosts != nil {
		_ = d.Set("hosts", cngwRoute.Hosts)
	}

	if cngwRoute.Paths != nil {
		_ = d.Set("paths", cngwRoute.Paths)
	}

	if cngwRoute.Protocols != nil {
		_ = d.Set("protocols", cngwRoute.Protocols)
	}

	if cngwRoute.PreserveHost != nil {
		_ = d.Set("preserve_host", cngwRoute.PreserveHost)
	}

	if cngwRoute.HttpsRedirectStatusCode != nil {
		_ = d.Set("https_redirect_status_code", cngwRoute.HttpsRedirectStatusCode)
	}

	if cngwRoute.StripPath != nil {
		_ = d.Set("strip_path", cngwRoute.StripPath)
	}

	if cngwRoute.ForceHttps != nil {
		_ = d.Set("force_https", cngwRoute.ForceHttps)
	}

	if cngwRoute.DestinationPorts != nil {
		_ = d.Set("destination_ports", cngwRoute.DestinationPorts)
	}

	if cngwRoute.Headers != nil {
		headersList := []interface{}{}
		for _, headers := range cngwRoute.Headers {
			headersMap := map[string]interface{}{}

			if cngwRoute.Headers.Key != nil {
				headersMap["key"] = cngwRoute.Headers.Key
			}

			if cngwRoute.Headers.Value != nil {
				headersMap["value"] = cngwRoute.Headers.Value
			}

			headersList = append(headersList, headersMap)
		}

		_ = d.Set("headers", headersList)

	}

	return nil
}

func resourceTencentCloudTseCngwRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tse.NewModifyCloudNativeAPIGatewayRouteRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeID := idSplit[1]

	request.GatewayId = &gatewayId
	request.RouteID = &routeID

	immutableArgs := []string{"gateway_id", "service_i_d", "route_name", "methods", "hosts", "paths", "protocols", "preserve_host", "https_redirect_status_code", "strip_path", "force_https", "destination_ports", "headers"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("gateway_id") {
		if v, ok := d.GetOk("gateway_id"); ok {
			request.GatewayId = helper.String(v.(string))
		}
	}

	if d.HasChange("service_i_d") {
		if v, ok := d.GetOk("service_i_d"); ok {
			request.ServiceID = helper.String(v.(string))
		}
	}

	if d.HasChange("route_name") {
		if v, ok := d.GetOk("route_name"); ok {
			request.RouteName = helper.String(v.(string))
		}
	}

	if d.HasChange("methods") {
		if v, ok := d.GetOk("methods"); ok {
			methodsSet := v.(*schema.Set).List()
			for i := range methodsSet {
				methods := methodsSet[i].(string)
				request.Methods = append(request.Methods, &methods)
			}
		}
	}

	if d.HasChange("hosts") {
		if v, ok := d.GetOk("hosts"); ok {
			hostsSet := v.(*schema.Set).List()
			for i := range hostsSet {
				hosts := hostsSet[i].(string)
				request.Hosts = append(request.Hosts, &hosts)
			}
		}
	}

	if d.HasChange("paths") {
		if v, ok := d.GetOk("paths"); ok {
			pathsSet := v.(*schema.Set).List()
			for i := range pathsSet {
				paths := pathsSet[i].(string)
				request.Paths = append(request.Paths, &paths)
			}
		}
	}

	if d.HasChange("protocols") {
		if v, ok := d.GetOk("protocols"); ok {
			protocolsSet := v.(*schema.Set).List()
			for i := range protocolsSet {
				protocols := protocolsSet[i].(string)
				request.Protocols = append(request.Protocols, &protocols)
			}
		}
	}

	if d.HasChange("preserve_host") {
		if v, ok := d.GetOkExists("preserve_host"); ok {
			request.PreserveHost = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("https_redirect_status_code") {
		if v, ok := d.GetOkExists("https_redirect_status_code"); ok {
			request.HttpsRedirectStatusCode = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("strip_path") {
		if v, ok := d.GetOkExists("strip_path"); ok {
			request.StripPath = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("force_https") {
		if v, ok := d.GetOkExists("force_https"); ok {
			request.ForceHttps = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("destination_ports") {
		if v, ok := d.GetOk("destination_ports"); ok {
			destinationPortsSet := v.(*schema.Set).List()
			for i := range destinationPortsSet {
				destinationPorts := destinationPortsSet[i].(int)
				request.DestinationPorts = append(request.DestinationPorts, helper.IntUint64(destinationPorts))
			}
		}
	}

	if d.HasChange("headers") {
		if v, ok := d.GetOk("headers"); ok {
			for _, item := range v.([]interface{}) {
				kVMapping := tse.KVMapping{}
				if v, ok := dMap["key"]; ok {
					kVMapping.Key = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					kVMapping.Value = helper.String(v.(string))
				}
				request.Headers = append(request.Headers, &kVMapping)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().ModifyCloudNativeAPIGatewayRoute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwRoute failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwRouteRead(d, meta)
}

func resourceTencentCloudTseCngwRouteDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_route.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeID := idSplit[1]

	if err := service.DeleteTseCngwRouteById(ctx, gatewayId, routeID); err != nil {
		return err
	}

	return nil
}
