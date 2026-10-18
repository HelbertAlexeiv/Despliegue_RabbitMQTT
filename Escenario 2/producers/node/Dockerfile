FROM node:22-alpine
WORKDIR /app
RUN npm init -y && npm install amqplib
COPY producer.js .
CMD ["node", "producer.js"]