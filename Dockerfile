FROM golang:1.23-alpine as build

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add tzdata \
	&& apk --no-cache add ca-certificates \
    && update-ca-certificates

#设置 GO 环境变量
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

WORKDIR /jjsdapi

COPY . .

RUN go build -o jjsdapi ./cmd/jjsdapi/jjsdapi.go

FROM scratch as final

# 设置时区为上海
COPY --from=build /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ=Asia/Shanghai

COPY --from=build /jjsdapi/jjsdapi /
COPY --from=build /jjsdapi/etc/ /etc/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 13100
ENTRYPOINT [ "/jjsdapi", "-env" ]
CMD ["prod"]
