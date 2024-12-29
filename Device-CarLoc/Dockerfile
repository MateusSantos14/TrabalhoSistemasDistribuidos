# Start with an official Python image
FROM python:3.12-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the working directory
COPY . /app

# Install any dependencies, including protobuf
RUN pip install --no-cache-dir -r requirements.txt

# Install protobuf compiler (if you need to generate .py from .proto files)
RUN apt-get update && apt-get install -y protobuf-compiler

# Expose the port the app will run on (if needed for your application)
EXPOSE 9997

# Command to run your application
CMD ["python", "main.py"]
