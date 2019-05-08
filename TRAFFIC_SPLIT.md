# Traffic Split Demo

1. Install traffic split CRD: `kubectl apply -f trafficsplit-crd.yaml`
2. (use Linkerd images with tag `dev-f23dcd96-alex` or build images off of this branch)
3. Install Linkerd `bin/linkerd install | kubectl apply -f -`
4. Install modified booksapp from this branch: `bin/linkerd inject bookapp.yml | kubectl apply -f -`
5. Observe that authors is sending to books v1: `bin/linkerd stat deploy --from deploy/authors` or view authors deployment detail in dashboard
6. Create 50-50 traffic split: `kubectl apply -f mysplit.yaml`
7. Observe that authors is sending to books v1 and v2: `bin/linkerd stat deploy --from deploy/authors` (or observe in dashboard)
8. Delete traffic split `kubectl delete -f mysplit.yaml`
9. Observe that authors is sending to books v1: `bin/linkerd stat deploy --from deploy/authors` (or observe in dashboard)
