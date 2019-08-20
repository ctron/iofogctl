## iofogctl deploy connector

Deploy a Connector

### Synopsis

Deploy a Connector.

```
iofogctl deploy connector [flags]
```

### Examples

```
iofogctl deploy connector -f connector.yaml
```

### Options

```
  -f, --file string   YAML file containing resource definitions for Connector
  -h, --help          help for connector
```

### Options inherited from parent commands

```
      --config string      CLI configuration file (default is ~/.iofog/config.yaml)
  -n, --namespace string   Namespace to execute respective command within (default "default")
  -v, --verbose            Toggle for displaying verbose output of API client
```

### SEE ALSO

* [iofogctl deploy](iofogctl_deploy.md)	 - Deploy ioFog platform or components on existing infrastructure

###### Auto generated by spf13/cobra on 19-Aug-2019