[DEPRECATED] Terraform provider for mackerel.io
==================================

**[DEPRECATED]** Use official [terraform-provider-mackerel](https://github.com/mackerelio-labs/terraform-provider-mackerel)


----


A [Terraform](https://www.terraform.io/) plugin that provides resources for [mackerel.io](https://mackerel.io/).

[![Build Status](https://travis-ci.org/kjmkznr/terraform-provider-mackerel.svg?branch=master)](https://travis-ci.org/kjmkznr/terraform-provider-mackerel)

Install
-------

* Download the latest release for your platform.
* Rename the executable to `terraform-provider-mackerel`

Provider Configuration
----------------------

### Example

```
provider "mackerel" {
  api_key = "xxx"
}
```

or

```
provider "mackerel" {}
```

### Reference

* `api_key` - (Optional) API Key of the Mackerel. If empty, provider try use `MACKEREL_API_KEY` envirenment variable.

Resources
---------

### `mackerel_dashboard`

Configure a dashboard.

#### Example

```
resource "mackerel_dashboard" "foobar" {
    title         = "terraform_for_mackerel_test_foobar"
    url_path      = "foo/bar"
    body_markdown = <<EOF
# Head1
## Head2

* List1
* List2
EOF
}
```

### `mackerel_host_monitor`

Configure a host monitor.

#### Example

```
resource "mackerel_host_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar"
    duration              = 10
    metric                = "cpu%"
    operator              = ">"
    warning               = 85.5
    critical              = 95.5
    notification_interval = 10
}
```

### `mackerel_service_monitor`

Configure a service monitor.

#### Example

```
resource "mackerel_service_monitor" "foobar" {
    name                      = "terraform_for_mackerel_test_foobar_upd"
    service                   = "Blog"
    duration                  = 10
    metric                    = "cpu%"
    operator                  = ">"
    warning                   = 85.5
    critical                  = 95.5
    missing_duration_warning  = 10
    missing_duration_critical = 100
    notification_interval     = 10
}
```

### `mackerel_external_monitor`

Configure a external url monitor.

#### Example

```
resource "mackerel_external_monitor" "foobar" {
    name                   = "terraform_for_mackerel_test_foobar"
    url                    = "https://terraform.io/"
    method                 = "GET"
    service                = "Web"
    notification_interval  = 10
    response_time_duration = 5
    response_time_warning  = 500
    response_time_critical = 1000
    contains_string        = "terraform"
    max_check_attempts     = 2

    certification_expiration_warning  = 30
    certification_expiration_critical = 10

    skip_certificate_verification = false

    request_body = "{\"request\": \"body\"}"
    headers = {
        "Content-Type" = "application/json",
        "API-Key" = "xxxxxx",
    }

    memo = "This monitor is for Terraform.io."
}
```

### `mackerel_expression_monitor`

Configure a expression monitor. (experimental)

#### Example

```
resource "mackerel_expression_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar"
    expression            = "avg(roleSlots(\"server:role\",\"loadavg5\"))"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}
```

### `mackerel_service`

Configure a service.

#### Example

```
resource "mackerel_service" "foobar" {
    name = "Web"
    memo = "This is Web service"
}
```

### `mackerel_channel`

Configure a channel.

#### Example

```
resource "mackerel_channel" "foobar" {
  name    = "test_slack"
  type    = "slack"
  events  = ["alert"]
  url = "https://hooks.slack.com/services/"
  mentions = {
    "ok": "status_ok",
    "critical": "critical_alert",
  }
  enabled_graph_image = true
}
```


Build
-----

```
$ make build
```

Testing
-------

```
$ MACKEREL_API_KEY=xxxx make testacc
```

TODO
----

* Support data source


Licence
-------

Mozilla Public License, version 2.0

Author
------

[KOJIMA Kazunori](https://github.com/kjmkznr)

