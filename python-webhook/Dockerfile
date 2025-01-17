# Use the official Python image as the base image
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy the requirements file into the container
COPY requirements.txt .

# Install any dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code
COPY python-webhook.py .

# Expose port 443 for the webhook server
EXPOSE 443

# Run the webhook server
CMD ["python", "python-webhook.py"]
