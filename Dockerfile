FROM docker.cnb.cool/yuwen-gueen/docker-images-chrom/golang:1.24-alpine_amd64 AS builder
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN set -evx -o pipefail        \
    && apk update               \
    && apk add --no-cache git   \
    && rm -rf /var/cache/apk/*  \
    && go build -ldflags="-s -w" -o busuanzi main.go

FROM docker.cnb.cool/yuwen-gueen/docker-images-chrom/node:21-alpine_amd64 AS ts-builder
WORKDIR /app

# 接收构建参数，默认 busuanzi.xxdevops.cn
ARG BSZ_API_DOMAIN=https://busuanzi.xxdevops.cn

COPY ./dist .

# 创建 .env 文件
RUN echo "BSZ_API_DOMAIN=${BSZ_API_DOMAIN}" > .env

RUN set -evx -o pipefail        \
    && npm install -g pnpm      \
    && pnpm install             \
    && pnpm run build           \
    && rm -rf node_modules      \
    && rm -rf pnpm-lock.yaml    \
    && rm -rf tsconfig.json

FROM docker.cnb.cool/yuwen-gueen/docker-images-chrom/alpine:3.16_amd64
WORKDIR /app

COPY --from=builder /app/busuanzi /app
COPY --from=ts-builder /app /app/dist
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/entrypoint.sh /app

RUN chmod +x /app/entrypoint.sh

EXPOSE 8080
ENTRYPOINT  [ "sh", "entrypoint.sh" ]