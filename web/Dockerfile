FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

# Copy the built jar file
COPY target/osint-gateway-1.0.0.jar app.jar

# Create non-root user
RUN addgroup -S spring && adduser -S spring -G spring
USER spring

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["java", "-jar", "app.jar"]
