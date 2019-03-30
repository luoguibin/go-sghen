sudo docker build -t go-sghen:1.1.0 .
sudo docker service update --image go-sghen:1.1.0 go-sghen-service

# docker images // 查看所有镜像
sudo docker rmi $(docker images -f "dangling=true" -q) // 删除构建镜像时暂时的临时的空名字镜像