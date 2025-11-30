FROM eclipse-temurin:21-jdk-alpine
ARG GRADLE_VERSION=9.2.1
WORKDIR /app
RUN wget -q https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-all.zip && \
    unzip -q gradle-${GRADLE_VERSION}-all.zip && \
    rm gradle-${GRADLE_VERSION}-all.zip
COPY spring/build.gradle spring/settings.gradle ./
RUN ./gradle-${GRADLE_VERSION}/bin/gradle wrapper --gradle-version ${GRADLE_VERSION} --distribution-type all
RUN ./gradlew --no-daemon classes
COPY spring/src ./src
RUN ./gradlew --no-daemon build -x test
EXPOSE 8080
CMD ["java", "-XX:+UseNUMA", "-XX:+UseParallelGC", "-jar", "build/libs/spring-1.0.0.jar"]