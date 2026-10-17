import pika, json, time, random, os 
from datetime import datetime

EXCHANGE = "anime_exchange"
GENERO = "isekai"
TITULOS =["RE:Zero", "Sword Art Online", "Mushoku Tensei", "Overlord", "Thai Time I Got Reincarnated as a Slime"]
EVENTOS = ["Nuevo episodio", "Nuevo Trailer", "Confimadad nueva temporada", "Estreno de pelicula"]

def conectar():
    cred = pika.PlainCredentials(os.getenv("RABBIT_USER"), os.getenv("RABBIT_PASS"))
    while True:
        try:
            return pika.BlockingConnection(pika.ConnectionParameters(os.getenv("RABBIT_HOST"), credentials=cred))
        except Exception:
            print("Esperando RabbitMQ ...", flush=True)
            time.sleep(3)

conn = conectar()
ch = conn.channel()
ch.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)

while True:
    msg = {"genero": GENERO, "titulo": random.choice(TITULOS), "evento": random.choice(EVENTOS),
            "productor": "python","hora": datetime.now().strftime("%H:%M:%S")}
    ch.basic_publish(exchange=EXCHANGE, routing_key="anime." + GENERO, body=json.dumps(msg))
    print("[PUBLICADO] anime." + GENERO + " -> " + msg["titulo"] + " | " + msg["evento"], flush=True)
    time.sleep(random.randint(2, 5))