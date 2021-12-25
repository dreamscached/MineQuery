FROM alpine:3.15.0

ARG JAVA_PACKAGE=openjdk17-jre-headless
ARG SERVER_URL

RUN apk add $JAVA_PACKAGE
RUN mkdir /server && wget -O /server/server.jar $SERVER_URL

COPY docker_entrypoint.sh /
RUN adduser -D server && chown -R server:server /server

USER server:server
EXPOSE 25565
ENV EULA=false JAVA_OPTS="-Xmx256M -Xms128M"
STOPSIGNAL SIGINT
ENTRYPOINT [ "/docker_entrypoint.sh" ]