package oceanus

import (
	"context"
	"fmt"
	"log"
	"strconv"

	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	regionv20220627 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// OceanusRegionProduct is the product name used to scope the region
// DescribeRegions / DescribeZones queries to the Oceanus numbering space.
const OceanusRegionProduct = "oceanus"

// DescribeOceanusRegions queries the region management service for the
// regions where Oceanus is available.
func (me *OceanusService) DescribeOceanusRegions(ctx context.Context) (regions []*regionv20220627.RegionInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := regionv20220627.NewDescribeRegionsRequest()
	request.Product = helper.String(OceanusRegionProduct)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *regionv20220627.DescribeRegionsResponse
	err := sdkresource.Retry(tccommon.ReadRetryTimeout, func() *sdkresource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseRegionClient().DescribeRegionsWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeRegions returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	regions = response.Response.RegionSet
	return
}

// DescribeOceanusZones queries the region management service for the
// availability zones where Oceanus is available.
func (me *OceanusService) DescribeOceanusZones(ctx context.Context) (zones []*regionv20220627.ZoneInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := regionv20220627.NewDescribeZonesRequest()
	request.Product = helper.String(OceanusRegionProduct)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *regionv20220627.DescribeZonesResponse
	err := sdkresource.Retry(tccommon.ReadRetryTimeout, func() *sdkresource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseRegionClient().DescribeZonesWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeZones returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	zones = response.Response.ZoneSet
	return
}

// DescribeOceanusRegionId translates a region name (e.g. `ap-guangzhou`) into
// the numeric region id the Oceanus cluster APIs expect (e.g. `1`).
func (me *OceanusService) DescribeOceanusRegionId(ctx context.Context, region string) (regionId int64, errRet error) {
	regions, err := me.DescribeOceanusRegions(ctx)
	if err != nil {
		errRet = err
		return
	}
	for _, item := range regions {
		if item == nil || item.Region == nil || *item.Region != region {
			continue
		}
		if item.RegionIdMC == nil {
			break
		}
		regionId, err = strconv.ParseInt(*item.RegionIdMC, 10, 64)
		if err != nil {
			errRet = fmt.Errorf("region %s reports a non-numeric id %q", region, *item.RegionIdMC)
			return
		}
		return
	}
	errRet = fmt.Errorf("region %s is not available for oceanus", region)
	return
}

// DescribeOceanusZoneId translates a zone name (e.g. `ap-guangzhou-7`) into
// the numeric zone id the Oceanus cluster APIs expect (e.g. `100007`).
func (me *OceanusService) DescribeOceanusZoneId(ctx context.Context, zone string) (zoneId int64, errRet error) {
	zones, err := me.DescribeOceanusZones(ctx)
	if err != nil {
		errRet = err
		return
	}
	for _, item := range zones {
		if item == nil || item.Zone == nil || *item.Zone != zone {
			continue
		}
		if item.ZoneId == nil {
			break
		}
		zoneId, err = strconv.ParseInt(*item.ZoneId, 10, 64)
		if err != nil {
			errRet = fmt.Errorf("zone %s reports a non-numeric id %q", zone, *item.ZoneId)
			return
		}
		return
	}
	errRet = fmt.Errorf("zone %s is not available for oceanus", zone)
	return
}
