tag="$1"
docker --version
docker build -t $tag ./
docker push $tag
