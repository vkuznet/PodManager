# PodManager
PodManager service to manage k8s pods based on alert information provided by AlertManager.
The flow of the PodManager is following:
- it fetches all active alerts from AlertManager
- it matched its rules with alerts
- if matched alert is found it will perform an action (provided by the rule) on
  a given pod within its namespace, e.g. delete and restart the pod

The rules are defined as following:
- `name` is used to match alert name
- `namespace` defines k8s namespace
- `pod` is used to identify pod value from alert attributes
- `env` defines k8s environment
- `action` defines which action to apply for a given pod

Here is an example of configuration:
```
{
    "alert_manager": "http://alert-manager.url",
    "interval": 10,
    "rules": [
        {"name": "service is down", "namespace": "xxx", "pod": "apod", "action": "restart", "env": "k8s-prod"},
        {"name": "number of workflows is high", "namespace": "xyz", "pod": "apod", "action": "print"}
    ],
    "verbose": 1
}
```
The `interval` defines periodicity of the service checks with given
AlertManager. The rules define the alert name, namespace, env, pod attribute in
alert, and appropriate action. The first rule will watch for alert with
`service is down` name within env, if found, it will use `apod` attribute of alert
to fetch the pod name to use, and it will apply `restart` action on that pod
within given namespace.
