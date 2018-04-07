# SonarQube Notification Resource

Gets Sonarqube results.

## Installing

Use this resource by adding the following to
the `resource_types` section of a pipeline config:

```yaml
---
resource_types:
- name: sonarqube-notifier
  type: docker-image
  source:
    repository: lgohr/sonarqube
    tag: latest
```

## Source configuration

Configure as follows:

```yaml
---
resources:
- name: my-sonarqube
  type: sonarqube-notifier
  source:
    target: https://my.sonar.server
    sonartoken: ((my-secret-token))
    component: my:component
    metrics: ncloc,complexity,violations,coverage
```

* `target`: *Required.* URL of your SonarQube instance e.g. `https://my-atlassian.com/sonar`.
* `sonartoken`: *Required.* [Security token](https://docs.sonarqube.org/display/SONAR/User+Token), which is used to connect to Sonarqube.
* `component`: *Required.* The component _key_ of your component. This is shown in the dashboard url as https://my-atlassian/sonar/dashboard?id=ComponentKey
* `metrics`: *Required.* The metrics you want to grab. See https://docs.sonarqube.org/display/SONAR/Metric+Definitions

## `in`: Get the latest result

Get the latest result; write it to the local working directory (e.g.
`/tmp/build/get`) with the filename result.json.

Example response (metrics: nloc,complexity,violations,coverage)
```json
{
  "component": {
    "id": "AWH_6osdce3G0HojaCW1",
    "key": "my:component",
    "name": "my-component",
    "qualifier": "TRK",
    "measures": [
      {
        "metric": "ncloc",
        "value": "824",
        "periods": [
          {
            "index": 1,
            "value": "299"
          }
        ]
      },
      {
        "metric": "complexity",
        "value": "90",
        "periods": [
          {
            "index": 1,
            "value": "27"
          }
        ]
      },
      {
        "metric": "violations",
        "value": "5",
        "periods": [
          {
            "index": 1,
            "value": "-6"
          }
        ]
      },
      {
        "metric": "coverage",
        "value": "91.4",
        "periods": [
          {
            "index": 1,
            "value": "40.7"
          }
        ]
      }
    ]
  },
  "periods": [
    {
      "index": 1,
      "mode": "previous_version",
      "date": "2018-03-07T16:58:31+0100"
    }
  ]
}
```

## `out`: Nothing
