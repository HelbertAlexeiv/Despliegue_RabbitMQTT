import pika  # type: ignore[reportMissingImports]
import json, time, os

EXCHANGE = "anime_exchange"
COLA = "cola_python"
CLAVES = ["anime.isekai", "anime.aventura"]

def conectar():
    cred = pika.PlainCredentials(os.getenv("RABBIT_USER"), os.getenv("RABBIT_PASS"))
    while True:
        try:
            return pika.BlockingConnection(pika.ConnectionParameters(os.getenv("RABBIT_HOST"), credentials=cred))
        except Exception:
            print("Esperando RabbitMQ...", flush=True)
            time.sleep(3)

conn = conectar()
ch = conn.channel()
ch.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)
ch.queue_declare(queue=COLA, durable=True)
for k in CLAVES:
    ch.queue_bind(queue=COLA, exchange=EXCHANGE, routing_key=k)

def callback(ch, method, props, body):
    m = json.loads(body)
    print("[PYTHON] (" + method.routing_key + ") " + m["titulo"] + " - " + m["evento"] + " (de " + m["productor"] + " a las " + m["hora"] + ")", flush=True)

ch.basic_consume(queue=COLA, on_message_callback=callback, auto_ack=True)
print("Python escuchando: " + str(CLAVES), flush=True)
ch.start_consuming()