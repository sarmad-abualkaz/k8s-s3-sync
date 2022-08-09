# k8s-s3-sync
A project to sync kubernetes secrets to an s3 bucket (and vise-versa)

## How to user
This project can be deployed to a Kubenretes cluster via Helm using the following

```
helm install <Release-Name> --set args[0]="--s3-bucket-name=<bucket-name>"
```

## How it works 
This can run in two modes -  `sync-secret-to-s3` or  `sync-s3-to-secret` by switching the flag `--sync-s3=true`. 

Setting this flag to true will perform `sync-secret-to-s3` and false will do the reverse.

*Note this flag is set to false by default.*

Flags to note:


| flag | purpose | default |
| --- | --- | --- | 
|`--aws-profile` | aws profile to use in `~/<user>/.aws` folder. If set to empty it will perform proper aws-cred cascade. Set to empty to make use of AssumeWebIdentity through a service account. |`tools`| 
|`--aws-region` | aws region to target. |`us-east-1` |
|`--cert-object-name` | object name where cert is stored in s3 bucket. | `"cert.text"`|
| `--cert-key-object-name` |  object name where cert key is stored in s3 bucket. |`"key.text"`|
| `--kube-config` | where the process is running, i.e. how kubeconfig will be setup. `"in-cluster"` and `local` are the only other acceptable options. | `in-cluster` |
| `--namespace` | namespace where the secret is stored. | `"local"` |
| `--s3-bucket-name` | Name of s3 bucket where objects are synced. | |
| `--secret-name` | kubernetes secret name to sync from/to. | |
| `--sleeping` | sleep time between syncs in seconds. | 20 |
| `--syncS3` | should sync s3 from secret. This syncs a kubernetes secret from s3 by dafault. | `false` |

