# Dockerfile for User Service (Customer Service Building)
FROM python:3.13-slim

# Set working directory inside container
WORKDIR /app

# Copy requirements first (for better caching)
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY main.py .

# Expose port 8001 (Customer Service Building entrance)
EXPOSE 8001

# Command to start the service
CMD ["python", "main.py"]