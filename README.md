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

## `out`: Nothing
