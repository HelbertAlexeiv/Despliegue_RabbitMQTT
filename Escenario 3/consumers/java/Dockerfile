FROM eclipse-temurin:21-jdk
WORKDIR /app
ADD https://repo1.maven.org/maven2/com/rabbitmq/amqp-client/5.21.0/amqp-client-5.21.0.jar lib/
ADD https://repo1.maven.org/maven2/org/slf4j/slf4j-api/1.7.36/slf4j-api-1.7.36.jar lib/
ADD https://repo1.maven.org/maven2/org/slf4j/slf4j-simple/1.7.36/slf4j-simple-1.7.36.jar lib/
COPY Consumer.java .
RUN javac -cp "lib/*" Consumer.java
CMD ["java", "-cp", "lib/*:.", "Consumer"]