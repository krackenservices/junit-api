# JUnit API Service

This is a Go application designed to serve JUnit test result files through an API. The application reads JUnit XML files from a location where they are stored as flat files and exposes them via an HTTP API, allowing them to be consumed by other applications or services.

## Features

- Supports reading JUnit files from both:
    - **Amazon S3**
    - **Local Filesystem**
- Exposes an HTTP API to fetch JUnit test results in JSON format.
- Written entirely in **Golang** for performance and simplicity.

## Limitations

- Currently only single location will be used and all files must be in top level. 

## Supported File Systems

### 1. **Amazon S3**
The application can list and retrieve JUnit files stored in an Amazon S3 bucket.

### 2. **Local Filesystem**
It can also read JUnit files stored on the local file system, providing flexibility for local and cloud-based workflows.

## Prerequisites

- **Go 1.18 or higher**
- **AWS SDK for Go v2** (if using S3)
- AWS credentials for S3 access, or local access permissions for file operations.

## Installation

To get started, clone this repository and install the necessary dependencies:

```bash
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/s3


