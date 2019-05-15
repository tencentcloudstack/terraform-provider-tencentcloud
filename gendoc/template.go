package main

const (
	docTPL = `---
layout: "{{.cloud_mark}}"
page_title: "{{.cloud_title}}: {{.name}}"
sidebar_current: "docs-{{.cloud_mark}}-{{.dtype}}-{{.resource}}"
description: |-
  {{.description_short}}
---

# {{.name}}

{{.description}}

## Example Usage

{{.example}}

## Argument Reference

The following arguments are supported:

{{.arguments}}
{{if ne .attributes ""}}
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

{{.attributes}}
{{end}}
{{if ne .import ""}}
## Import

{{.import}}
{{end}}
`
	idxTPL = `
<% wrap_layout :inner do %>
    <% content_for :sidebar do %>
        <div class="docs-sidebar hidden-print affix-top" role="complementary">
            <ul class="nav docs-sidenav">
                <li<%= sidebar_current("docs-home") %>>
                    <a href="/docs/providers/index.html">All Providers</a>
                </li>
                <li<%= sidebar_current("docs-{{.cloud_mark}}-index") %>>
                    <a href="/docs/providers/{{.cloud_mark}}/index.html">{{.cloud_title}} Provider</a>
                </li>
                {{range $k, $v := .datasource}}
                <li<%= sidebar_current("docs-{{$.cloud_mark}}-{{$v.ResType}}{{if ne $v.NameShort ""}}-{{$v.NameShort}}{{end}}") %>>
                    <a href="#">{{$v.Name}}</a>
                    <ul class="nav nav-visible">
                        {{range $kk, $vv := $v.Resources}}
                        <li<%= sidebar_current("docs-{{$.cloud_mark}}-{{$v.ResType}}-{{index $vv 1}}") %>>
                            <a href="/docs/providers/{{$.cloud_mark}}/{{$v.ResTypeShort}}/{{index $vv 1}}.html">{{index $vv 0}}</a>
                        </li>{{end}}
                    </ul>
                </li>
                {{end}}
            </ul>
        </div>
    <% end %>
    <%= yield %>
<% end %>
`
)
