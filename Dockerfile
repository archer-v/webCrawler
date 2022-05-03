# stage 2: scratch
FROM alpine:latest as scratch
RUN apk --no-cache add ca-certificates libc6-compat 
WORKDIR /opt/webCrawler
ADD ./distr/webCrawler.tar.gz .
USER 1000
CMD ["./webCrawler"]
