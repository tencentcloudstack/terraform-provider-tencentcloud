Provides a resource to create a teo config group version

Example Usage

```hcl
resource "tencentcloud_teo_config_group_version" "teo_config_group_version" {
  content = <<EOT
{
  "FormatVersion": "1.0",
  "ZoneConfig": {
    "SmartRouting": {
      "Switch": "off"
    },
    "Cache": {
      "NoCache": {
        "Switch": "off"
      },
      "FollowOrigin": {
        "Switch": "on",
        "DefaultCache": "on",
        "DefaultCacheStrategy": "on",
        "DefaultCacheTime": 0
      },
      "CustomTime": {
        "Switch": "off",
        "CacheTime": 2592000
      }
    },
    "MaxAge": {
      "FollowOrigin": "on",
      "CacheTime": 600
    },
    "CacheKey": {
      "FullURLCache": "on",
      "QueryString": {
        "Switch": "off",
        "Action": "includeCustom"
      },
      "IgnoreCase": "off"
    },
    "CachePrefresh": {
      "Switch": "off",
      "CacheTimePercent": 90
    },
    "OfflineCache": {
      "Switch": "on"
    },
    "Compression": {
      "Switch": "off",
      "Algorithms": [
        "brotli",
        "gzip"
      ]
    },
    "ForceRedirectHTTPS": {
      "Switch": "off",
      "RedirectStatusCode": 302
    },
    "HSTS": {
      "Switch": "off",
      "Timeout": 0,
      "IncludeSubDomains": "off",
      "Preload": "off"
    },
    "TLSConfig": {
      "Version": [
        "TLSv1",
        "TLSv1.1",
        "TLSv1.2",
        "TLSv1.3"
      ],
      "CipherSuite": "loose-v2023"
    },
    "OCSPStapling": {
      "Switch": "off"
    },
    "HTTP2": {
      "Switch": "off"
    },
    "QUIC": {
      "Switch": "off"
    },
    "UpstreamHTTP2": {
      "Switch": "off"
    },
    "IPv6": {
      "Switch": "off"
    },
    "WebSocket": {
      "Switch": "off",
      "Timeout": 30
    },
    "PostMaxSize": {
      "Switch": "on",
      "MaxSize": 838860800
    },
    "ClientIPHeader": {
      "Switch": "off"
    },
    "ClientIPCountry": {
      "Switch": "off"
    },
    "gRPC": {
      "Switch": "off"
    },
    "NetworkErrorLogging": {
      "Switch": "off"
    },
    "AccelerateMainland": {
      "Switch": "off"
    },
    "StandardDebug": {
      "Switch": "off",
      "AllowClientIPList": [
        "1.14.231.0/24",
        "1.194.255.0/24"
      ],
      "Expires": "2025-09-01T12:45:37Z"
    }
  },
  "Rules": [
    {
      "RuleName": "网站加速-cdn.defaultsetting.cn",
      "Branches": [
        {
          "Condition": "$${http.request.host} in ['cdn.defaultsetting.cn']",
          "Actions": [
            {
              "Name": "Cache",
              "CacheParameters": {
                "CustomTime": {
                  "Switch": "on",
                  "IgnoreCacheControl": "off",
                  "CacheTime": 2592000
                }
              }
            },
            {
              "Name": "CacheKey",
              "CacheKeyParameters": {
                "FullURLCache": "on",
                "QueryString": {
                  "Switch": "off"
                },
                "IgnoreCase": "off"
              }
            }
          ],
          "SubRules": [
            {
              "Branches": [
                {
                  "Condition": "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']",
                  "Actions": [
                    {
                      "Name": "Cache",
                      "CacheParameters": {
                        "NoCache": {
                          "Switch": "on"
                        }
                      }
                    }
                  ]
                }
              ]
            },
            {
              "Branches": [
                {
                  "Condition": "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']",
                  "Actions": [
                    {
                      "Name": "MaxAge",
                      "MaxAgeParameters": {
                        "FollowOrigin": "off",
                        "CacheTime": 3600
                      }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "RuleName": "网站加速-pages.migraine.com.cn",
      "Branches": [
        {
          "Condition": "$${http.request.host} in ['pages.migraine.com.cn']",
          "Actions": [
            {
              "Name": "Cache",
              "CacheParameters": {
                "CustomTime": {
                  "Switch": "on",
                  "IgnoreCacheControl": "off",
                  "CacheTime": 0
                }
              }
            },
            {
              "Name": "CacheKey",
              "CacheKeyParameters": {
                "FullURLCache": "on",
                "QueryString": {
                  "Switch": "off"
                },
                "IgnoreCase": "off"
              }
            }
          ],
          "SubRules": [
            {
              "Branches": [
                {
                  "Condition": "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']",
                  "Actions": [
                    {
                      "Name": "Cache",
                      "CacheParameters": {
                        "NoCache": {
                          "Switch": "on"
                        }
                      }
                    }
                  ]
                }
              ]
            },
            {
              "Branches": [
                {
                  "Condition": "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']",
                  "Actions": [
                    {
                      "Name": "MaxAge",
                      "MaxAgeParameters": {
                        "FollowOrigin": "off",
                        "CacheTime": 3600
                      }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "RuleName": "跳过pages鉴权验证",
      "Branches": [
        {
          "Condition": "$${http.request.host} in ['pages.migraine.com.cn']",
          "Actions": [
            {
              "Name": "ModifyRequestHeader",
              "ModifyRequestHeaderParameters": {
                "HeaderActions": [
                  {
                    "Action": "add",
                    "Name": "X-SKIP-TOKEN",
                    "Value": "eop-1022"
                  }
                ]
              }
            }
          ]
        }
      ]
    },
    {
      "RuleName": "分地域回源",
      "Branches": [
        {
          "Condition": "$${http.request.host} in ['abc.migraine.com.cn']",
          "SubRules": [
            {
              "Branches": [
                {
                  "Condition": "$${http.request.ip.country} in ['Asia']",
                  "Actions": [
                    {
                      "Name": "ModifyOrigin",
                      "ModifyOriginParameters": {
                        "OriginType": "IPDomain",
                        "Origin": "1.2.3.4",
                        "OriginProtocol": "follow",
                        "HTTPOriginPort": 80,
                        "HTTPSOriginPort": 443
                      }
                    }
                  ]
                }
              ]
            },
            {
              "Branches": [
                {
                  "Condition": "$${http.request.ip.country} in ['Africa']",
                  "Actions": [
                    {
                      "Name": "ModifyOrigin",
                      "ModifyOriginParameters": {
                        "OriginType": "IPDomain",
                        "Origin": "3.4.5.6",
                        "OriginProtocol": "follow",
                        "HTTPOriginPort": 80,
                        "HTTPSOriginPort": 443
                      }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
EOT
  description = "test version"
  group_id    = "cg-3lchxitnb5pb"
  zone_id     = "zone-2xkazzl8yf6k"
}
```

Example of reading a configuration file

```hcl
resource "tencentcloud_teo_config_group_version" "teo_config_group_version" {
  content = file("l7_config.json")
  description = "test version"
  group_id    = "cg-3lchxitnb5pb"
  zone_id     = "zone-2xkazzl8yf6k"
}
```