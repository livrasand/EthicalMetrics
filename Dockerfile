# Usa imagen de Go con SQLite + CGO habilitado
FROM golang:1.21-bullseye

# Instala SQLCipher
RUN apt-get update && apt-get install -y libsqlcipher-dev pkg-config sqlite3 && \
    rm -rf /var/lib/apt/lists/*

# Establece el directorio de trabajo
WORKDIR /app

# Copia todos los archivos
COPY . .

# Habilita sqlcipher al compilar
ENV CGO_ENABLED=1
ENV GOOS=linux

# Descarga dependencias
RUN go mod tidy

# Compila la app con soporte para SQLCipher
RUN go build -tags "sqlcipher" -o server ./cmd/server

# Expone el puerto
EXPOSE 8080

# Ejecuta el binario
CMD ["./server"]
