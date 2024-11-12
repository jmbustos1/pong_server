# Usa una imagen base de Go que incluye herramientas esenciales para el desarrollo
FROM golang:1.22-bullseye

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
# RUN go build -o /docker-app # compila la aplicacion

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
# CMD ["/docker-app"]

# ejecutar codigo fuente
CMD ["go", "run", "."] 

