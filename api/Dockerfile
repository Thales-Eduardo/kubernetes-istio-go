FROM golang
COPY . .

RUN apt-get update && apt-get install -y --no-install-recommends \
    postgresql-client \
    netcat-traditional && \
    rm -rf /var/lib/apt/lists/*

RUN go build -o server .
CMD ["./server"]