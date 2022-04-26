FROM alpine

RUN apk update && apk add ca-certificates libc6-compat && rm -rf /var/cache/apk/*

ENV PORT 3033
EXPOSE $PORT

RUN mkdir /app
RUN mkdir /app/log

COPY ./klepa-shop-backend /app/klepa-shop-backend
COPY ./.env /app/.env

WORKDIR /app

RUN chmod +x /app/klepa-shop-backend
