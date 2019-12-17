# gce-node-auth
A sidecar app which use gce node service credentials to provide Git Auth Callback URL.

# Build Image

```
export PROJECT=MY_GCP_PROJECT
gcloud builds submit --tag gcr.io/${PROJECT}/gce-node-auth  .
```

# Test

## Create test repo on GCP

```
export PROJECT=MY_GCP_PROJECT
export REPO_NAME=MY_REPO_NAME
gcloud config set project ${PROJECT}
gcloud source repos create ${REPO_NAME}
# clone and some files into the repo
```

## Apply yaml


```
eval "cat git-sync-with-gce-node-auth.yaml | sed 's/\$PROJECT/$PROJECT/g' | sed 's/\$REPO_NAME/$REPO_NAME/g'" | kubectl apply -f -
```

## Check git-sync succeed

```
kubectl get pod/git-sync-with-gce-node-auth
kubectl logs pod/git-sync-with-gce-node-auth git-sync
kubectl exec pod/git-sync-with-gce-node-auth ls /tmp/git/git-data
```
