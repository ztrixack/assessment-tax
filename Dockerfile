############################
# STEP 1 prepare the source
############################
FROM golang:1.21-alpine AS builder

# Set the environment variables for the go command:
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64

# Create a non-root user and group
RUN adduser -D -g '' appuser

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Copy go.mod, go.sum and dependencies if exists
COPY go.mod go.sum vendor? ./

# Install dependencies if the vendor folder is not present
RUN if [ ! -d "vendor" ]; then go mod tidy; fi

# Import the code from the context.
COPY . .

# Build the Go application
RUN go build -o /bin/app .

############################
# STEP 2 the running container
############################
FROM scratch AS runner
LABEL maintainer="Tanawat Hongthai <ztrixack.th@gmail.com>"

# Import the timezone, user, and group files
COPY --from=builder /etc/passwd /etc/group /etc/

# Copy the compiled Go application and set ownership to appuser
COPY --from=builder --chown=appuser:appuser /bin/app /bin/app

# Switch to the non-root user
USER appuser

# Start the application
CMD ["/bin/app"]