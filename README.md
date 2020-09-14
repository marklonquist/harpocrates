# Harpocrates
> Harpocrates was the god of silence, secrets and confidentiality

Harpocrates is a small application that can be used to pull secrets from [HashiCorp Vault](https://www.vaultproject.io/).
It can output the secrets in different formats:
 * JSON, which is simple key-values.
 * `source` ready env file e.g.
 ```bash
 export KEY=value
 export FOO=bar
 ```
 * Raw key values.
 * Raw value in separate file.


Harpocrates is designed such it can be used as an init- or sidecar container in [Kubernetes](https://kubernetes.io/). 
In this scenario it uses the ServiceAccount token in `/var/run/secrets/kubernetes.io/serviceaccount/token` and exchanges this for a Vault token by posting it to `/auth/kubernetes/login`.

This requires that the [Kubernetes Auth Method](https://www.vaultproject.io/docs/auth/kubernetes) is enabled in Vault.

## Authentication
The easiest way to authenticate is to use your Vault token:
```bash
harpocrates --vault-token "sometoken"
```
This can also be specified as the environment var `VAULT_TOKEN`


## Usage
In harpocrates can specify which secrets to pull in 3 different ways.
### YAML file

```yaml
format: env
output: "/secrets"
prefix: PREFIX_
secrets:
  - secret/data/secret/dev
  - secret/data/foo:
      keys:
       - APIKEY
```

```bash
harpocrates -f /path/to/file.yaml
```

### Inline JSON
```bash
harpocrates '{"format":"env","output":"/secrets","prefix":"PREFIX_","secrets":["secret/data/secret/dev",{"secret/data/foo":{"keys":["APIKEY"]}}]}'
```

### CLI Parameters
```bash
harpocrates --format "env" --secret "/secret/data/somesecret" --prefix "PREFIX_" --output "/secrets"
```


## CLI and ENV Options

| Flag          | Values                                                                                                     |                       Default                       |
| ------------- | ---------------------------------------------------------------------------------------------------------- | :-------------------------------------------------: |
| vault_address | https://vaulturl                                                                                           |                          -                          |
| cluster_name  | string                                                                                                     |                          -                          |
| token_path    | /path/to/token, uses clustername and path to login and exchange a vault token which is used in vault_token | /var/run/secrets/kubernetes.io/serviceaccount/token |
| vault_token   | token as a string. If empty token_path will be queried                                                     |                          -                          |
| format        | env, json or secret                                                                                        |                         env                         |
| output        | /path/to/output                                                                                            |                  /tmp/secrets.env                   |
| prefix        | prefix keys, eg. K8S_                                                                                      |                          -                          |
| secret        | vault path /secretengine/data/some/secret                                                                  |                          -                          |



## Note
We have to set the following annotation, in order for the kubernetes autoscaler to be able to scale down again.
```
annotations:
    "cluster-autoscaler.kubernetes.io/safe-to-evict": "true"
```
https://issuetracker.google.com/issues/148295270


## Contribution










# old
When using a ServiceAccount in Kubernetes, the jwt token can be retrieved by reading the file `/var/run/secrets/kubernetes.io/serviceaccount/token`

And then it can be exchanged to a Vault token by posting it to `/auth/kubernetes/login`

Example of a secret file:
```yaml
format: json
output: path/to/dir/to/save/secret/to
secrets:
  - path/to/secret1
  - path/to/secret2:
    - key1:
        saveAsFile: true
    - key2
```
At the moment it takes a json file as input, you can convert your secret to json by doing:
`yq read secret.yml -j`

Orb should read kustomize yaml from Vault


## Deployment.yml
To use this, you can add harpocrates as an initContainers like so:
```yaml
initContainers:
  - name: secret-dumper
    image: harbor.bestsellerit.com/library/harpocrates:68
    args:
      - '{"format":"env","output":"/secrets","prefix":"alfeios_","secrets":["ES/data/alfeios/prod"]}'
    volumeMounts:
      - name: secrets
        mountPath: /secrets
    env:
      - name: VAULT_ADDRESS
        value: $VAULT_ADDR
      - name: CLUSTER_NAME
        value: es03-prod
volumes:
  - name: secrets
    emptyDir: {}
```

CircleCI steps:
```yaml
- secret-injector:
    app-name: alfeios
    file: deployment.yml
    secretFile: alfeios-secrets.yml
- secret-injector:
    app-name: alfeios-db
    file: deployment.yml
    secretFile: alfeios-db-secrets.yml
```


## TO-DO

* Support files ?


## NOTES
We have to set the following annotation, in order for the autoscaler to be able to scale down again.
```
annotations:
    "cluster-autoscaler.kubernetes.io/safe-to-evict": "true"
```
https://issuetracker.google.com/issues/148295270