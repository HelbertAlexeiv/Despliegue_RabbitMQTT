FROM node:22-alpine
WORKDIR /app
RUN npm init -y && npm install amqplib
COPY consumer.js .
CMD ["node", "consumer.js"]