# 🐇 Arquitectura Event-Driven con RabbitMQ

<div align="center">

![RabbitMQ](https://img.shields.io/badge/RabbitMQ-4.3-FF6600?logo=rabbitmq&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)
![Multi-language](https://img.shields.io/badge/Lenguajes-Python%20%7C%20Node.js%20%7C%20Go%20%7C%20Java-8A2BE2)

</div>

## 📌 Descripción general

Este proyecto implementa un sistema de mensajería basado en el patrón de publicación/suscripción con **RabbitMQ**. Se simula un flujo de **noticias de anime** donde varios productores envían eventos a un exchange tipo `topic` y varios consumidores los reciben según sus intereses.

En este caso, el escenario es **centralizado**: todos los servicios (productores, consumidores y broker) viven dentro de un único nodo de Docker Compose.

---

## 🧭 Escenario 1 – Centralizado

Sistema de mensajería publicar/suscribir sobre **RabbitMQ (AMQP)** con temática de **noticias de anime**. Productores y consumidores están escritos en distintos lenguajes y cada servicio corre en su propio contenedor Docker. En este escenario **todos los contenedores corren en un único nodo**.

---

## 🏗️ Arquitectura

```text
 Productores (4)                  Broker                     Consumidores (4)
┌───────────────────┐                                      ┌───────────────────┐
│ producer-python   │── anime.isekai ──┐                ┌──▶│ consumer-python   │
│ producer-node     │── anime.aventura ┤   RabbitMQ     │   │ consumer-node     │
│ producer-go       │── anime.accion ──┼─▶ anime_exchange ──┼──▶│ consumer-go       │
│ producer-java     │── anime.romance ─┘   (topic)      └──▶│ consumer-java     │
└───────────────────┘                                      └───────────────────┘
                    Un solo nodo · red por defecto de Docker Compose
```

- **Patrón:** `Publish/Subscribe` con un exchange tipo `topic` (`anime_exchange`)
- **Routing key:** el género del anime (`anime.<genero>`)
- **Colas:** cada consumidor tiene su propia cola durable y se enlaza a las claves que le interesan
- **Mensajes:** JSON con `genero`, `titulo`, `evento`, `productor` y `hora`

---

## 🚀 Servicios

### Productores

| Servicio | Lenguaje | Publica en |
|---|---|---|
| `producer-python` | Python | `anime.isekai` |
| `producer-node` | Node.js | `anime.aventura` |
| `producer-go` | Go | `anime.accion` |
| `producer-java` | Java | `anime.romance` |

### Consumidores

| Servicio | Lenguaje | Cola | Se suscribe a |
|---|---|---|---|
| `consumer-python` | Python | `cola_python` | `anime.isekai`, `anime.aventura` |
| `consumer-node` | Node.js | `cola_node` | `anime.#` (todos) |
| `consumer-go` | Go | `cola_go` | `anime.accion`, `anime.romance` |
| `consumer-java` | Java | `cola_java` | `anime.isekai`, `anime.romance` |

### Broker

| Servicio | Imagen | Puertos |
|---|---|---|
| `rabbitmq` | `rabbitmq:4.3-management` | `5672` (AMQP), `15672` (panel web) |

> **Nota:** el enunciado menciona `rabbitmq:3-management`; se utilizó la versión `4.3` de la imagen oficial, que es más reciente. El comportamiento de exchanges, colas y bindings es equivalente.

---

## 🗂️ Estructura del proyecto

```text
.
├── docker-compose.yml
├── producers/
│   ├── python/   (Dockerfile, producer.py)
│   ├── node/     (Dockerfile, producer.js)
│   ├── go/       (Dockerfile, go.mod, producer.go)
│   └── java/     (Dockerfile, Producer.java)
├── consumers/
│   ├── python/   (Dockerfile, consumer.py)
│   ├── node/     (Dockerfile, consumer.js)
│   ├── go/       (Dockerfile, go.mod, consumer.go)
│   └── java/     (Dockerfile, Consumer.java)
└── README.md
```

---

## ✅ Requisitos

- Docker y Docker Compose
- Docker Desktop abierto en Windows y con el motor en ejecución
- Conexión a internet para la primera construcción
- Puertos `5672` y `15672` libres

Verificación:

```bash
docker --version
docker compose version
```

---

## ▶️ Cómo ejecutar

Desde la carpeta que contiene `docker-compose.yml`:

```bash
# Validar el archivo
docker compose config

# Construir y levantar los 9 contenedores
docker compose up --build -d

# Ver el estado
docker compose ps
```

La primera vez puede tardar unos minutos porque se compilan Go y Java.

### Ver los mensajes en consola

```bash
# Todo junto
docker compose logs -f

# Solo productores
docker compose logs -f producer-python producer-node producer-go producer-java

# Solo consumidores
docker compose logs -f consumer-python consumer-node consumer-go consumer-java
```

Salida esperada:

```text
producer-python  | [PUBLICADO] anime.isekai -> Re:Zero | Nuevo episodio
consumer-python  | [PYTHON] (anime.isekai) Re:Zero - Nuevo episodio (de python a las 14:32:10)
consumer-node    | [NODE] (anime.aventura) One Piece - Nuevo trailer (de node a las 14:32:12)
```

### Detener el entorno

```bash
docker compose down          # detiene y elimina contenedores y red
docker compose down -v       # además borra los volúmenes
```

---

## 📊 Panel de administración de RabbitMQ

Accede a:

- URL: <http://localhost:15672>
- Usuario: `admin`
- Contraseña: `admin123`

Qué puedes observar:

- **Exchanges → `anime_exchange`:** tipo `topic`, durable. En *Bindings* aparecen las colas con sus claves.
- **Queues:** `cola_python`, `cola_node`, `cola_go`, `cola_java`, con sus tasas de mensajes.
- **Connections:** 8 conexiones (4 productores y 4 consumidores)

### Bindings esperados

| Cola | Routing keys |
|---|---|
| `cola_go` | `anime.accion`, `anime.romance` |
| `cola_java` | `anime.isekai`, `anime.romance` |
| `cola_node` | `anime.#` |
| `cola_python` | `anime.isekai`, `anime.aventura` |

---

## 🔁 Prueba de resiliencia

Como las colas son durables, si un consumidor se cae, los mensajes se acumulan hasta que vuelva a estar disponible.

```bash
docker compose stop consumer-go
# esperar ~20 s y ver en el panel cómo sube "Ready" en cola_go
docker compose start consumer-go
docker compose logs -f consumer-go
```

Al volver, el consumidor recibe de golpe los mensajes acumulados.

---

## 🧭 Escenario 2 – Distribuido en dos nodos

El escenario 2 conserva la lógica de publicación/suscripción del escenario 1,
pero separa el broker de las aplicaciones en dos máquinas virtuales conectadas
a la misma red privada:

- **Nodo de aplicaciones:** ejecuta los cuatro productores y los cuatro consumidores.
- **Nodo broker:** ejecuta únicamente RabbitMQ.

```text
 Nodo de aplicaciones                         Nodo broker
┌───────────────────────────────┐             ┌────────────────────────┐
│ producer-python               │             │ RabbitMQ               │
│ producer-node                 │── TCP 5672 ─▶│ anime_exchange          │
│ producer-go                   │             │ (topic)                │
│ producer-java                 │             └────────────────────────┘
│ consumer-python               │
│ consumer-node                 │
│ consumer-go                   │
│ consumer-java                 │
└───────────────────────────────┘
```

El nodo de aplicaciones se conecta al broker mediante la variable
`RABBIT_HOST`. Por ejemplo:

```text
Nodo aplicaciones: 192.168.56.10
Nodo broker:       192.168.56.20
RABBIT_HOST=192.168.56.20
```

### Arquitectura

- **Patrón:** `Publish/Subscribe` con el exchange `anime_exchange`.
- **Tipo de exchange:** `topic`.
- **Transporte:** AMQP sobre TCP.
- **Puerto del broker:** `5672`.
- **Colas:** `cola_python`, `cola_node`, `cola_go` y `cola_java`.
- **Red:** una red privada compartida entre las dos máquinas virtuales.

### Servicios

Los productores y consumidores son los mismos del escenario 1:

| Nodo | Servicios |
|---|---|
| Aplicaciones | `producer-python`, `producer-node`, `producer-go`, `producer-java` |
| Aplicaciones | `consumer-python`, `consumer-node`, `consumer-go`, `consumer-java` |
| Broker | `rabbitmq` |

Las routing keys y bindings no cambian:

| Productor | Routing key | Consumidores |
|---|---|---|
| Python | `anime.isekai` | Python, Node.js y Java |
| Node.js | `anime.aventura` | Python y Node.js |
| Go | `anime.accion` | Node.js y Go |
| Java | `anime.romance` | Node.js, Go y Java |

### Estructura del escenario

```text
Escenario 2/
├── docker-compose.apps.yml       # nodo de aplicaciones
├── docker-compose.broker.yml     # nodo RabbitMQ
├── .env.example
├── producers/
│   ├── python/
│   ├── node/
│   ├── go/
│   └── java/
└── consumers/
    ├── python/
    ├── node/
    ├── go/
    └── java/
```

### Requisitos

- Dos máquinas virtuales con Linux.
- Docker Engine y Docker Compose instalados en ambas máquinas.
- Una red privada compartida entre las máquinas virtuales.
- IPs fijas o reservadas para ambas máquinas.
- Acceso a Internet durante la primera construcción de las imágenes.

Se recomienda configurar un adaptador NAT para Internet

### Configurar y levantar el nodo broker

En la máquina virtual que ejecutará RabbitMQ:

```bash
cd "Escenario 2"
cp .env.example .env
```

Edita `.env` y configura las credenciales:

```dotenv
RABBIT_USER=admin
RABBIT_PASS=una-contrasena-fuerte
```

Valida y levanta RabbitMQ:

```bash
docker compose --env-file .env -f docker-compose.broker.yml config
docker compose --env-file .env -f docker-compose.broker.yml up -d
docker compose --env-file .env -f docker-compose.broker.yml ps
```

El archivo `docker-compose.broker.yml` crea el volumen persistente
`rabbitmq_data`. No uses `down -v` salvo que quieras borrar el estado del
broker.

### Configurar y levantar el nodo de aplicaciones

En la máquina virtual de aplicaciones:

```bash
cd "Escenario 2"
cp .env.example .env
```

Edita `.env` con la IP privada del nodo broker y las mismas credenciales:

```dotenv
RABBIT_HOST=192.168.56.20
RABBIT_USER=admin
RABBIT_PASS=una-contrasena-fuerte
```

No uses `localhost` ni `rabbitmq` como `RABBIT_HOST`, porque RabbitMQ está en
otra máquina virtual.

Valida y construye los ocho servicios:

```bash
docker compose --env-file .env -f docker-compose.apps.yml config
docker compose --env-file .env -f docker-compose.apps.yml up --build -d
docker compose --env-file .env -f docker-compose.apps.yml ps
```

### Ver los mensajes

Para observar todo el flujo:

```bash
docker compose --env-file .env -f docker-compose.apps.yml logs -f
```

Para observar únicamente los consumidores:

```bash
docker compose --env-file .env -f docker-compose.apps.yml logs -f \
  consumer-python consumer-node consumer-go consumer-java
```

Antes de levantar las aplicaciones puedes probar la comunicación desde el nodo
de aplicaciones:

```bash
ping 192.168.56.20
nc -vz 192.168.56.20 5672
```

### Panel de administración

El panel de RabbitMQ se publica en `127.0.0.1:15672` dentro del nodo broker.
Para acceder desde Windows, crea un túnel SSH:

```powershell
ssh -L 15672:127.0.0.1:15672 usuario@192.168.56.20
```

Después abre <http://localhost:15672> e inicia sesión con las credenciales de
`.env`.

En el panel puedes observar:

- `anime_exchange`, de tipo `topic`.
- Las colas `cola_python`, `cola_node`, `cola_go` y `cola_java`.
- Los bindings y sus routing keys.
- Conexiones y tasas de publicación/consumo.
- Mensajes `Ready` y `Unacked`.

### Prueba de resiliencia

Para demostrar que una cola continúa existiendo cuando su consumidor se
detiene:

```bash
docker compose --env-file .env -f docker-compose.apps.yml stop consumer-go
```

Espera algunos segundos, observa `cola_go` en el panel y luego inicia el
consumidor:

```bash
docker compose --env-file .env -f docker-compose.apps.yml start consumer-go
docker compose --env-file .env -f docker-compose.apps.yml logs -f consumer-go
```

También puedes detener los productores sin apagar RabbitMQ:

```bash
docker compose --env-file .env -f docker-compose.apps.yml stop \
  producer-python producer-node producer-go producer-java
```

Para reanudarlos:

```bash
docker compose --env-file .env -f docker-compose.apps.yml start \
  producer-python producer-node producer-go producer-java
```

### Variables de entorno

| Variable | Valor en el nodo de aplicaciones |
|---|---|
| `RABBIT_HOST` | IP privada o DNS del nodo broker |
| `RABBIT_USER` | Debe coincidir con el broker |
| `RABBIT_PASS` | Debe coincidir con el broker |

El nodo broker utiliza `RABBIT_USER` y `RABBIT_PASS` para crear el usuario
inicial. El archivo `.env` no debe subirse al repositorio.

### Detener el escenario

En el nodo de aplicaciones:

```bash
docker compose --env-file .env -f docker-compose.apps.yml down
```

En el nodo broker:

```bash
docker compose --env-file .env -f docker-compose.broker.yml down
```

---

## ⚙️ Variables de entorno

Todos los productores y consumidores leen su conexión desde estas variables definidas en el `docker-compose.yml`:

| Variable | Valor en el Escenario 1 |
|---|---|
| `RABBIT_HOST` | `rabbitmq` (nombre del servicio en la red de Compose) |
| `RABBIT_USER` | `admin` |
| `RABBIT_PASS` | `admin123` |

Si RabbitMQ aún no está listo, cada servicio reintenta la conexión cada 3 segundos y muestra el mensaje `Esperando RabbitMQ...`.

---

## 🛠️ Problemas comunes

| Síntoma | Causa y solución |
|---|---|
| `Esperando RabbitMQ...` que no termina | Servicios en redes distintas. Deja una sola red y reinicia con `docker compose down -v --remove-orphans`. |
| `can't open file '/app/producer.py'` | Nombre o carpeta del archivo incorrectos. Reconstruye con `docker compose build --no-cache`. |
| Errores rojos de `com.rabbitmq` en VS Code | Son normales: la librería se descarga dentro del contenedor, no en la máquina local. |
| Error `must be a mapping` en el compose | Indentación incorrecta en el YAML. Usa espacios, no tabulaciones. |
| Los comandos `docker` fallan en Windows | Docker Desktop no está abierto o el servicio no está ejecutándose. |

---

## 🔭 Escenarios siguientes

- **Escenario 2 – Distribuido en 2 nodos:** Nodo 1 con productores y consumidores, Nodo 2 con RabbitMQ (ver [`Escenario 2/README.md`](Escenario%202/README.md))
- **Escenario 3 – Distribuido en 3 nodos:** Nodos 1 y 2 con conjuntos de productores y consumidores, Nodo 3 con RabbitMQ (3 archivos `docker-compose.yml`)

En ambos casos, `RABBIT_HOST` pasará a ser la IP del nodo del broker y el puerto `5672` de ese nodo deberá estar accesible.

---

## 📌 Resumen

Este proyecto demuestra cómo integrar varios lenguajes y servicios con una arquitectura basada en eventos, usando **RabbitMQ** como backbone de comunicación. Además de ejemplificar el patrón `topic` y la distribución de mensajes por routing keys, ilustra una base sólida para construir sistemas escalables y desacoplados.
