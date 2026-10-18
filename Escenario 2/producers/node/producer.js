const amqp = require('amqplib');

const EXCHANGE = 'anime_exchange';
const GENERO = 'aventura';
const TITULOS = ['One Piece', 'Hunter x Hunter', 'Made in Abyss', 'Fullmetal Alchemist', 'Frieren'];
const EVENTOS = ['Nuevo episodio', 'Nuevo trailer', 'Confirmada nueva temporada', 'Estreno de pelicula'];

const sleep = ms => new Promise(r => setTimeout(r, ms));
const rand = a => a[Math.floor(Math.random() * a.length)];

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
    while (true) {
        const msg = {
        genero: GENERO,
        titulo: rand(TITULOS),
        evento: rand(EVENTOS),
        productor: 'node',
        hora: new Date().toLocaleTimeString()
        };
        ch.publish(EXCHANGE, `anime.${GENERO}`, Buffer.from(JSON.stringify(msg)));
        console.log(`[PUBLICADO] anime.${GENERO} -> ${msg.titulo} | ${msg.evento}`);
        await sleep(2000 + Math.random() * 3000);
    }
})();