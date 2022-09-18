# 注意，这里的构建上下文，是在git源代码的根目录
FROM golang:latest AS build
WORKDIR /src

COPY . .

WORKDIR /src/${git_name}/${project_path}
RUN go build -o /app/${git_name} .

FROM alpine:latest AS base
WORKDIR /app
COPY --from=build /app .

#设置时区
RUN cp /usr/share/zoneinfo/GMT /etc/localtime
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai    /etc/localtime

EXPOSE ${entry_port}
EXPOSE 443

ENTRYPOINT ["./${git_name}"]