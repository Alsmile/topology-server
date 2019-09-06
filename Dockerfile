FROM alpine:latest as certs

# Install the CA certificates
RUN apk --update add ca-certificates


FROM scratch AS prod

MAINTAINER Alsmile "alsmile123@qq.com"
ENV REFRESHED_AT 2019-08-02

# 从certs阶段拷贝CA证书
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# 拷贝主程序
COPY . .

WORKDIR server
#RUN chmod +x ./topology
EXPOSE 8210
ENTRYPOINT ["./topology"]

# docker build -t registry.local/topology:0.1 .
# docker run --name topology -d -p 8210:8210 [-v /etc/le5leTopology.yaml:/etc/le5leTopology.yaml] <image name:tag>
