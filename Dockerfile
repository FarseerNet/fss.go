# 注意，这里的构建上下文，是在git源代码的根目录
FROM golang:alpine AS build
# 设置github代理
ENV GOPROXY https://goproxy.cn,direct
# 进入到项目目录中
WORKDIR /src/${git_name}
# 复制go.mod文件
COPY /${git_name}/go.mod .
# 下载依赖（支持docker缓存）
RUN go mod download

# 设置src目录，并将源代码复制到此
WORKDIR /src
COPY . .

# 设置项目的main路径
WORKDIR /src/${git_name}

# 删除go.work文件
RUN rm -rf go.work

# 更新go.sum
RUN go mod tidy

# 编译
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/${git_name} -ldflags="-w -s" .

FROM alpine:latest AS base
WORKDIR /app
COPY --from=build /app .
# 复制配置
COPY --from=build /src/${git_name}/farseer.yaml .

#设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai    /etc/localtime

EXPOSE ${entry_port}
EXPOSE 443

ENTRYPOINT ["./${git_name}"]