FROM python:3.14-slim
WORKDIR /app
RUN pip install --no-cache-dir pika
COPY producer.py .
CMD ["python", "-u", "producer.py"]
