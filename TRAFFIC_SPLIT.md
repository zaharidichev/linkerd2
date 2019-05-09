# Traffic Split Demo

1. Install traffic split CRD: `kubectl apply -f trafficsplit-crd.yaml`
2. (use Linkerd images with tag `dev-f23dcd96-alex` or build images off of this branch)
  a. To build off of this branch, generate a github access token and then edit Dockerfile-go-deps and replace
     USER with your github user name and TOKEN with your access token
3. Install Linkerd `bin/linkerd install | kubectl apply -f -`
4. Install mysql backend: `kubectl apply -f mysql-backend.yml`
5. Install modified booksapp `bin/linkerd inject mysql-app.yml | kubectl apply -f -`
6. Navigate to books deploy detail in dashboard.  Note books is sending to authors-v1 and SR sucks.
7. Create 50-50 traffic split: `kubectl apply -f mysplit.yaml`
8. Observe on books deploy detail that traffic is now split between authors-v1 and authors-v2.  Note the SR improvement.
9. Go full v2: `kubectl edit ts/authors-split` and remove the authors-v1 backend.
10. Observe effects in the dashboard: all green!
