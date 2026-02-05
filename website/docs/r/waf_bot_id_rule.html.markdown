---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_bot_id_rule"
sidebar_current: "docs-tencentcloud-resource-waf_bot_id_rule"
description: |-
  Provides a resource to create a WAF bot id rule
---

# tencentcloud_waf_bot_id_rule

Provides a resource to create a WAF bot id rule

~> **NOTE:** When using the current `tencentcloud_waf_bot_id_rule` resource configuration, if you need to customize configuration field `data`, it is recommended to first import all bot rules into Terraform for management using the `terraform import` command.

## Example Usage

### Configure using only protect_level and global_switch(Global Configuration)

```hcl
resource "tencentcloud_waf_bot_id_rule" "example" {
  domain        = "demo.com"
  scene_id      = "3000000001"
  protect_level = "normal"
  global_switch = 5
}
```

### Configure data details(Custom configuration)

```hcl
resource "tencentcloud_waf_bot_id_rule" "example" {
  domain        = "demo.com"
  scene_id      = "3000000001"
  global_switch = 0
  data {
    action  = "monitor"
    bot_id  = "Abot"
    rule_id = "3300002262"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "AppScan"
    rule_id = "3300002263"
    status  = false
  }
  data {
    action  = "monitor"
    bot_id  = "Astra"
    rule_id = "3300002264"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "BBScan"
    rule_id = "3300002265"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Bugscan"
    rule_id = "3300002266"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "C-SpamMasal"
    rule_id = "3300002267"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Censys"
    rule_id = "3300002268"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Corsy"
    rule_id = "3300002269"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "DDoS-Ripper"
    rule_id = "3300002270"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "DSXS"
    rule_id = "3300002271"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "DirBuster"
    rule_id = "3300002272"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "DotnetSpider"
    rule_id = "3300002273"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "EasySpider"
    rule_id = "3300002274"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "FinalRecon"
    rule_id = "3300002275"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "GoBot2"
    rule_id = "3300002276"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "HTTrack"
    rule_id = "3300002277"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Hawk"
    rule_id = "3300002278"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Impulse_DDoS"
    rule_id = "3300002279"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Impulse_SMS"
    rule_id = "3300002280"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "InfinityCrawler"
    rule_id = "3300002281"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "JSKY"
    rule_id = "3300002282"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Jorgee"
    rule_id = "3300002283"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "MechanicalSoup"
    rule_id = "3300002285"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Mysqloit"
    rule_id = "3300002286"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Netsparker"
    rule_id = "3300002287"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Nettacker"
    rule_id = "3300002288"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Nuclei"
    rule_id = "3300002289"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Paros"
    rule_id = "3300002290"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Photon"
    rule_id = "3300002291"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Pker"
    rule_id = "3300002292"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Raven-Storm"
    rule_id = "3300002293"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Recon-ng"
    rule_id = "3300002294"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "SMS_bomber"
    rule_id = "3300002295"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "SMS_bomber_version2"
    rule_id = "3300002296"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "SpamSms_alodok"
    rule_id = "3300002297"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "SpamSms_matahari"
    rule_id = "3300002298"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "SpamSms_olx"
    rule_id = "3300002299"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "SpamSms_payu"
    rule_id = "3300002300"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "SpamSms_socil"
    rule_id = "3300002301"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "SqlPower"
    rule_id = "3300002302"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Striker"
    rule_id = "3300002303"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Strong-Web-Crawler"
    rule_id = "3300002304"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Symfony"
    rule_id = "3300002305"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "TBomb"
    rule_id = "3300002306"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "TBomb_flipkart"
    rule_id = "3300002307"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "TBomb_makaan"
    rule_id = "3300002308"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "TBomb_referer"
    rule_id = "3300002309"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "Volnx"
    rule_id = "3300002310"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "WebCollector"
    rule_id = "3300002311"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "YetAnotherSMSBomber"
    rule_id = "3300002312"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "ZmEu"
    rule_id = "3300002313"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "abot"
    rule_id = "3300002314"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "arachni"
    rule_id = "3300002315"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "aspider"
    rule_id = "3300002316"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "autoscraper"
    rule_id = "3300002317"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "bad-user-agents"
    rule_id = "3300002318"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "bomb3r_flipkart"
    rule_id = "3300002320"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "burpsuite"
    rule_id = "3300002321"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "cocrawler"
    rule_id = "3300002322"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "cola"
    rule_id = "3300002323"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "colly"
    rule_id = "3300002324"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "commix"
    rule_id = "3300002325"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "crawlee"
    rule_id = "3300002326"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "crawley"
    rule_id = "3300002327"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "dark-fantasy-hack-tool"
    rule_id = "3300002328"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "device_detector_Crawler"
    rule_id = "3300002329"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "device_detector_Feed_Fetcher"
    rule_id = "3300002330"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "device_detector_Search_bot"
    rule_id = "3300002331"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "device_detector_Security_Checker"
    rule_id = "3300002332"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "dirhunt"
    rule_id = "3300002334"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "dirsearch"
    rule_id = "3300002335"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "emptyua"
    rule_id = "3300002336"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "feedparser"
    rule_id = "3300002337"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "feroxbuster"
    rule_id = "3300002338"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "ferret"
    rule_id = "3300002339"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "ffuf"
    rule_id = "3300002340"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "fluxay"
    rule_id = "3300002341"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "fofa"
    rule_id = "3300002342"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "gain"
    rule_id = "3300002343"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "gecco"
    rule_id = "3300002344"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "geziyor"
    rule_id = "3300002345"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "gobuster"
    rule_id = "3300002346"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "gocrawl"
    rule_id = "3300002347"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "hakrawler"
    rule_id = "3300002348"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "httpx"
    rule_id = "3300002349"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "ice_sms_boomber"
    rule_id = "3300002350"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "identYwaf"
    rule_id = "3300002351"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "jaeles"
    rule_id = "3300002352"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "katana"
    rule_id = "3300002353"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "lux"
    rule_id = "3300002354"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "masscan"
    rule_id = "3300002355"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "multi_platform"
    rule_id = "3300002356"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "news-please"
    rule_id = "3300002357"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "newspaper"
    rule_id = "3300002358"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "ni_bomber"
    rule_id = "3300002359"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "nikto"
    rule_id = "3300002360"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "nmap"
    rule_id = "3300002361"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "nuclei"
    rule_id = "3300002362"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "numspy_bomber"
    rule_id = "3300002363"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "openvas"
    rule_id = "3300002364"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "osv-scanner"
    rule_id = "3300002365"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "overload"
    rule_id = "3300002366"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "paloalto"
    rule_id = "3300002367"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "pulsarr"
    rule_id = "3300002368"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "pyspider"
    rule_id = "3300002369"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "scrapy"
    rule_id = "3300002370"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "scrapy-redis"
    rule_id = "3300002371"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "scylla"
    rule_id = "3300002372"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "shodan"
    rule_id = "3300002373"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "skipfish"
    rule_id = "3300002374"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "skrape.it"
    rule_id = "3300002375"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "sms-repeater"
    rule_id = "3300002376"
    status  = true
  }
  data {
    action = "monitor"
    bot_id = "smsBomb"

    rule_id = "3300002377"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "sms_boom_platform"
    rule_id = "3300002378"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "smsbomb"
    rule_id = "3300002379"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "smsboom"
    rule_id = "3300002380"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "spider-flow"
    rule_id = "3300002381"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "sqlmap"
    rule_id = "3300002382"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "sukhoi"
    rule_id = "3300002383"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "supercrawler"
    rule_id = "3300002384"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "suspicious_browser_chrome"
    rule_id = "3300002385"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "suspicious_browser_firefox"
    rule_id = "3300002386"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "suspicious_browser_ie"
    rule_id = "3300002387"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "suspicious_os_mobile"
    rule_id = "3300002388"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "suspicious_os_windows"
    rule_id = "3300002389"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "tsunami-security-scanner"
    rule_id = "3300002390"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "ufonet"
    rule_id = "3300002391"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "uil2pn"
    rule_id = "3300002392"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "w3af"
    rule_id = "3300002393"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "w9scan"
    rule_id = "3300002394"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "webBee"
    rule_id = "3300002395"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "webmagic"
    rule_id = "3300002396"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "webster"
    rule_id = "3300002397"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "wfuzz"
    rule_id = "3300002398"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "whatwaf"
    rule_id = "3300002399"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "whatweb"
    rule_id = "3300002400"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "wombat"
    rule_id = "3300002401"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "wpscan"
    rule_id = "3300002402"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "xxl-crawler"
    rule_id = "3300002403"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "you-get"
    rule_id = "3300002404"
    status  = true
  }
  data {
    action  = "monitor"
    bot_id  = "zgrab"
    rule_id = "3300002405"
    status  = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `scene_id` - (Required, String, ForceNew) Scene ID.
* `data` - (Optional, List) Configuration information, supports batch processing.
* `global_switch` - (Optional, Int) 0-global settings do not take effect 1-global switch configuration field takes effect 2-global action configuration field takes effect 3-both global switch and action fields take effect 4-only modify global redirect path 5-only modify global protection level.
* `protect_level` - (Optional, String) Protection level: normal-normal; strict-strict.

The `data` object supports the following:

* `action` - (Required, String) Action configuration.
* `rule_id` - (Required, String) Rule ID.
* `status` - (Required, Bool) Rule switch.
* `bot_id` - (Optional, String) Rule name.
* `redirect` - (Optional, String) Redirect path.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

WAF bot id rule can be imported using the domain#sceneId, e.g.

```
terraform import tencentcloud_waf_bot_id_rule.example demo.com#3000000001
```

