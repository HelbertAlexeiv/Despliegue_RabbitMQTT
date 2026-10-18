# Escenario 2: dos nodos

Este escenario conserva la lógica de publicación y suscripción del escenario 1,
pero separa el despliegue en dos máquinas:

```text
Tu máquina local                            Instancia OCI
┌───────────────────────────────┐          ┌────────────────────────┐
│ 4 productores + 4 consumidores│ ─ AMQP ─▶│ RabbitMQ                │
│ docker-compose.apps.yml       │ Internet │ docker-compose.broker.yml│
└───────────────────────────────┘  TCP 5672└────────────────────────┘
```

Los productores y consumidores publican y consumen en `anime_exchange`, un
exchange `topic`. Las routing keys, colas y bindings son las mismas del
escenario 1. La diferencia es que `RABBIT_HOST` ya no es `rabbitmq`: debe ser
la IP o el DNS del nodo OCI.

## 1. Configurar el nodo broker en OCI

En la instancia de OCI:

1. Instala Docker Engine y Docker Compose.
2. Copia `docker-compose.broker.yml` y crea `.env` a partir de `.env.example`.
3. Define una contraseña fuerte en `RABBIT_PASS`.
4. Abre en la VCN y en el firewall del sistema únicamente el puerto `5672`
   desde la IP pública de tu conexión a Internet.
5. No expongas `15672` a Internet. El compose lo publica en `127.0.0.1`;
   para acceder al panel usa un túnel SSH:

   ```powershell
   ssh -L 15672:127.0.0.1:15672 opc@<IP_PUBLICA_OCI>
   ```

   Luego abre <http://localhost:15672>.

Ejemplo de `.env` en el broker:

```dotenv
RABBIT_USER=admin
RABBIT_PASS=una-contrasena-larga-y-unica
```

Inicia RabbitMQ:

```bash
docker compose --env-file .env -f docker-compose.broker.yml config
docker compose --env-file .env -f docker-compose.broker.yml up -d
docker compose --env-file .env -f docker-compose.broker.yml ps
```

Comprueba desde tu máquina local que AMQP sea accesible. En PowerShell:

```powershell
Test-NetConnection <IP_PUBLICA_O_DNS_DEL_BROKER> -Port 5672
```

El resultado debe mostrar `TcpTestSucceeded : True`.

## 2. Configurar las aplicaciones en tu máquina local

En tu máquina local, abre PowerShell en esta carpeta y crea `.env`:

```powershell
cd "C:\ruta\al\proyecto\Escenario 2"
cp .env.example .env
```

Edita `.env` con la dirección del broker:

```dotenv
RABBIT_HOST=<IP_PUBLICA_O_DNS_DEL_BROKER>
RABBIT_USER=admin
RABBIT_PASS=una-contrasena-larga-y-unica
```

`RABBIT_HOST` debe ser resoluble y enrutable desde los contenedores Docker
locales. No uses `localhost` ni `rabbitmq`, porque `localhost` apuntaría al
contenedor de la aplicación y `rabbitmq` no existe en este Compose.

Valida y levanta las aplicaciones:

```powershell
docker compose --env-file .env -f docker-compose.apps.yml config
docker compose --env-file .env -f docker-compose.apps.yml up --build -d
docker compose --env-file .env -f docker-compose.apps.yml ps
```

Los procesos reintentan la conexión cada tres segundos mientras RabbitMQ no
esté disponible. Para observar el flujo:

```powershell
docker compose --env-file .env -f docker-compose.apps.yml logs -f
```

Para revisar solo consumidores:

```powershell
docker compose --env-file .env -f docker-compose.apps.yml logs -f consumer-python consumer-node consumer-go consumer-java
```

## 3. Configuración importante de OCI

- Usa una regla de entrada TCP `5672` limitada a tu IP pública actual con
  máscara `/32`; evita `0.0.0.0/0`.
- Configura la misma regla en la lista de seguridad de la VCN y en `iptables`
  o `firewalld` del sistema operativo.
- Como la aplicación corre fuera de OCI, `RABBIT_HOST` debe ser la IP pública
  reservada o un DNS que apunte a la instancia OCI. No uses la IP privada de
  OCI salvo que tu equipo esté conectado a esa VCN mediante VPN o una red
  equivalente.
- Si tu proveedor cambia tu IP pública, actualiza la regla de OCI y `.env`.
  Para evitarlo, usa una IP fija, una VPN o un DNS dinámico controlado.
- Mantén `15672` cerrado públicamente; usa túnel SSH o una VPN para la
  administración.
- Conserva el volumen `rabbitmq_data`: contiene el estado, usuarios, colas y
  mensajes persistentes.
- No incluyas `.env` en el repositorio. Usa `.env.example` como plantilla y
  rota la contraseña si se expone.
- Verifica que la hora de ambas máquinas esté sincronizada para interpretar
  correctamente el campo `hora` de los mensajes.

## 4. Alternativa si no quieres exponer el puerto 5672

Puedes crear un túnel SSH desde PowerShell:

```powershell
ssh -N -L 5673:127.0.0.1:5672 opc@<IP_PUBLICA_OCI>
```

En ese caso, cambia temporalmente `.env` en tu máquina local a:

```dotenv
RABBIT_HOST=host.docker.internal
```

y modifica los servicios o el código para usar el puerto `5673`. La
configuración actual usa el puerto estándar `5672`, por lo que la opción
recomendada para esta entrega es permitir `5672` únicamente desde tu IP
 pública `/32`.

## 5. Detener el despliegue

En tu máquina local:

```powershell
docker compose --env-file .env -f docker-compose.apps.yml down
```

En el nodo broker:

```bash
docker compose --env-file .env -f docker-compose.broker.yml down
```

No uses `down -v` en el broker salvo que quieras borrar deliberadamente el
estado persistente de RabbitMQ.
