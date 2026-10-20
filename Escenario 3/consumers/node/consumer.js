const amqp = require('amqplib');

const EXCHANGE = 'anime_exchange';
const COLA = 'cola_node';
const CLAVES = ['anime.#'];

const sleep = ms => new Promise(r => setTimeout(r, ms));

async function conectar() {
    const { RABBIT_USER: u, RABBIT_PASS: p, RABBIT_HOST: h } = process.env;
    while (true) {
        try {
        return await amqp.connect(`amqp://${u}:${p}@${h}`);
        } catch (e) {
        console.log('Esperando RabbitMQ...');
        await sleep(3000);
        }
    }
}

(async () => {
    const conn = await conectar();
    const ch = await conn.createChannel();
    await ch.assertExchange(EXCHANGE, 'topic', { durable: true });
    await ch.assertQueue(COLA, { durable: true });
    for (const k of CLAVES) await ch.bindQueue(COLA, EXCHANGE, k);
    console.log('Node escuchando: ' + CLAVES);
    ch.consume(COLA, m => {
        const d = JSON.parse(m.content.toString());
        console.log(`[NODE] (${m.fields.routingKey}) ${d.titulo} - ${d.evento} (de ${d.productor} a las ${d.hora})`);
    }, { noAck: true });
})();