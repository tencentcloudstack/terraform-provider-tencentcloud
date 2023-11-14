package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoZoneSettingResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZoneSetting,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_zone_setting.zone_setting", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_zone_setting.zone_setting",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoZoneSetting = `

resource "tencentcloud_teo_zone_setting" "zone_setting" {
  zone_id = &lt;nil&gt;
    cache {
		cache {
			switch = &lt;nil&gt;
			cache_time = &lt;nil&gt;
			ignore_cache_control = &lt;nil&gt;
		}
		no_cache {
			switch = &lt;nil&gt;
		}
		follow_origin {
			switch = &lt;nil&gt;
		}

  }
  cache_key {
		full_url_cache = &lt;nil&gt;
		ignore_case = &lt;nil&gt;
		query_string {
			switch = &lt;nil&gt;
			action = &lt;nil&gt;
			value = &lt;nil&gt;
		}

  }
  max_age {
		max_age_time = &lt;nil&gt;
		follow_origin = &lt;nil&gt;

  }
  offline_cache {
		switch = &lt;nil&gt;

  }
  quic {
		switch = &lt;nil&gt;

  }
  post_max_size {
		switch = &lt;nil&gt;
		max_size = &lt;nil&gt;

  }
  compression {
		switch = &lt;nil&gt;
		algorithms = &lt;nil&gt;

  }
  upstream_http2 {
		switch = &lt;nil&gt;

  }
  force_redirect {
		switch = &lt;nil&gt;
		redirect_status_code = &lt;nil&gt;

  }
  https {
		http2 = &lt;nil&gt;
		ocsp_stapling = &lt;nil&gt;
		tls_version = &lt;nil&gt;
		hsts {
			switch = &lt;nil&gt;
			max_age = &lt;nil&gt;
			include_sub_domains = &lt;nil&gt;
			preload = &lt;nil&gt;
		}
		cert_info {
			cert_id = &lt;nil&gt;
			status = &lt;nil&gt;
		}

  }
  origin {
		origins = &lt;nil&gt;
		backup_origins = &lt;nil&gt;
		origin_pull_protocol = &lt;nil&gt;
		cos_private_access = &lt;nil&gt;

  }
  smart_routing {
		switch = &lt;nil&gt;

  }
  web_socket {
		switch = &lt;nil&gt;
		timeout = &lt;nil&gt;

  }
  client_ip_header {
		switch = &lt;nil&gt;
		header_name = &lt;nil&gt;

  }
  cache_prefresh {
		switch = &lt;nil&gt;
		percent = &lt;nil&gt;

  }
  ipv6 {
		switch = &lt;nil&gt;

  }
}

`
