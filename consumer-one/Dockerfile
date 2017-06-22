FROM openjdk:8-jre-alpine

COPY ./build/libs/app.jar /app/dist/app.jar
WORKDIR /app/dist

EXPOSE 8080

ENTRYPOINT exec java $JAVA_OPTS -jar app.jar
