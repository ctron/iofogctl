## iofogctl delete

Delete an existing ioFog resource

### Synopsis

Delete an existing ioFog resource.

Deleting Agents or Controllers will result in the respective deployments being torn down.

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
      --config string      CLI configuration file (default is ~/.iofog/config.yaml)
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -q, --quiet              Toggle for displaying verbose output
  -v, --verbose            Toggle for displaying verbose output of API client
```

### SEE ALSO

* [iofogctl](iofogctl.md)	 - 
* [iofogctl delete agent](iofogctl_delete_agent.md)	 - Delete an Agent
* [iofogctl delete all](iofogctl_delete_all.md)	 - Delete all resources within a namespace
* [iofogctl delete application](iofogctl_delete_application.md)	 - Delete an application
* [iofogctl delete controller](iofogctl_delete_controller.md)	 - Delete a Controller
* [iofogctl delete microservice](iofogctl_delete_microservice.md)	 - Delete a Microservice
* [iofogctl delete namespace](iofogctl_delete_namespace.md)	 - Delete a Namespace

###### Auto generated by spf13/cobra on 9-Aug-2019