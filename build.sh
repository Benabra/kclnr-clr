go mod tidy
go build -o kclnr
version=$(./kclnr version | awk '{print $3}')

git add .
git commit -m "Fix module import paths and update project structure"
git push origin main

git tag -a v$version -m "Release version v$version"
git push origin v$version