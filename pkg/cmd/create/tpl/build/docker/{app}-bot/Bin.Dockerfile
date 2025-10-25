FROM alpine:3.22
LABEL maintainer="<brooksyang@outlook.com>"

WORKDIR /opt/{[.AppName]}

# Tools
RUN apk add curl

# Timezone
# RUN apk --no-cache add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
#      echo "Asia/Shanghai" > /etc/timezone \

COPY _output/platforms/linux/amd64/{[.AppName]}-bot bin/

EXPOSE 8080

ENTRYPOINT ["/opt/{[.AppName]}/bin/{[.AppName]}-bot"]
