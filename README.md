# gce-node-auth
A sidecar app which use gce node service credentials to provide Git Auth Callback URL.

# Build Image

```
export PROJECT=MY_GCP_PROJECT
gcloud builds submit --tag gcr.io/${PROJECT}/git-askpass-gce-node .
```

# Test

## Create test repo on GCP

```bash
export REPO_NAME=MY_REPO_NAME
gcloud config set project ${PROJECT}
gcloud source repos create ${REPO_NAME}
# clone and some files into the repo
```

## Apply yaml


```bash
eval "cat git-askpass-gce-node.yaml | sed 's/\${PROJECT}/$PROJECT/g' | sed 's/\${REPO_NAME}/$REPO_NAME/g' " | kubectl apply -f -
```

## Check git-sync succeed

```bash
kubectl port-forward pod/git-askpass-gce-node 9102 9102
curl "http://localhost:9102/git_askpass"

kubectl get pod/git-askpass-gce-node
kubectl logs pod/git-askpass-gce-node git-sync
kubectl exec pod/git-askpass-gce-node ls /tmp/git/git-data
```
