# gce-node-auth
A sidecar app which use gce node service credentials to provide Git Auth Callback URL.

## Build Image

```bash
export PROJECT=MY_GCP_PROJECT
gcloud config set project ${PROJECT}
gcloud builds submit --tag gcr.io/${PROJECT}/git-askpass-gce-node .
```

## Test

### Create test repo on GCP Cloud Source Repo

```bash
export REPO_NAME=MY_REPO_NAME
gcloud source repos create ${REPO_NAME}
# clone and add some files into the repo
```

### Creating testing pod

```bash
eval "cat git-askpass-gce-node.yaml | sed 's/\${PROJECT}/$PROJECT/g' | sed 's/\${REPO_NAME}/$REPO_NAME/g' " | kubectl apply -f -
```

### Check git-askpass-gce-node running

```bash
kubectl port-forward pod/git-askpass-gce-node 9102 9102
curl "http://localhost:9102/git_askpass"
```

Output should be something like:

```bash
username=xxx@example.com
password=ya29.xxxxyyyyzzzz
```

### Check git-sync running

```bash
kubectl get pod/git-askpass-gce-node
kubectl logs pod/git-askpass-gce-node git-askpass-gce-node
kubectl logs pod/git-askpass-gce-node git-sync
kubectl exec pod/git-askpass-gce-node ls /tmp/git/git-data
```
