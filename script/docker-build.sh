image="$1"
tag="$2"
docker --version
docker build -t $image:$tag ./
docker tag $image:$tag $image:latest
docker push $image:$tag
docker push $image:latest
