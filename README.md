# localcluster 

early attempt at migrating bash scripts into an executable. this was created for my own personal dev workflow. 

the purpose of the binary is to wrap k3d, kubectl, and helm to quickly create a kubernetes dev cluster with ingress, monitoring, logging, and a postgres database ready for development. 


### quickstart 
1. get dependencies installed. if you have asdf already installed just run the `bin/run deploy setup` command for additional dependencies. otherwise ensure the dependencies below are installed.
2.  create the cluster `bin/run cluster create`

once the cluster is created you should now be able to access:

- [traefik dashboard](https://traefik.localdev.me/dashboard/#/) ingress
- [grafana](https://grafana.localdev.me/login) loki is configured for log aggregation
  login with user `admin` pwd `admin`
- [prometheus](https://prometheus.localdev.me/) for metrics collection
- postgres
  ```shell
  # example connection:
  > psql "postgres://postgres:admin@localhost:32500/postgres"
  ```
- [localstack](http://localstack.localdev.me) fully functional local aws cloud stack
  see the [docs here](https://docs.localstack.cloud/overview/)
  ```shell
    # helpful to add an alias like so:
    
    > export LOCALSTACK_HOST=localstack.localdev.me
    > alias awslocal="AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test  aws --endpoint-url=http://${LOCALSTACK_HOST:-locals}"
    
    # now you can use just like the aws cli like so:
    > awslocal s3api list-buckets
    ```



### dependencies
- [asdf](https://github.com/asdf-vm/asdf) used for dependency management
- [helm](https://helm.sh/) for deployments
- [kubectl](https://kubernetes.io/docs/reference/kubectl/) for deployment customizations
- [k3d](https://k3d.io/v5.6.0/) for cluster management
- [mkcert](https://github.com/FiloSottile/mkcert) for ssl in the cluster
 