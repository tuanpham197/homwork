#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Vui lòng cung cấp đường dẫn đến tệp hoặc thư mục để mở trong Visual Studio Code."
    exit 1
fi

if ! command -v code &> /dev/null; then
    echo "Visual Studio Code không được cài đặt."
    exit 1
fi

code "$@"

echo "Done open vs"

# start docker porject
PROJECT_PATH="home\tuan\source\scrapmarket"

# Change directory to your project path
cd "$PROJECT_PATH" || { echo "Error: Could not change directory to $PROJECT_PATH"; exit 1; }

# Start Docker Compose
docker-compose up -d


