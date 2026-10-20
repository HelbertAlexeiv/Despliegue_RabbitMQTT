import com.rabbitmq.client.*;
import java.text.SimpleDateFormat;
import java.util.*;

public class Producer {
    static final String EXCHANGE = "anime_exchange";
    static final String GENERO = "romance";
    static final String[] TITULOS = {"Toradora!", "Kaguya-sama", "Horimiya", "Your Lie in April", "Clannad"};
    static final String[] EVENTOS = {"Nuevo episodio", "Nuevo trailer", "Confirmada nueva temporada", "Estreno de pelicula"};

    public static void main(String[] args) throws Exception {
        ConnectionFactory f = new ConnectionFactory();
        f.setHost(System.getenv("RABBIT_HOST"));
        f.setUsername(System.getenv("RABBIT_USER"));
        f.setPassword(System.getenv("RABBIT_PASS"));

        Connection conn = null;
        while (conn == null) {
            try {
                conn = f.newConnection();
            } catch (Exception e) {
                System.out.println("Esperando RabbitMQ...");
                Thread.sleep(3000);
            }
        }
        Channel ch = conn.createChannel();
        ch.exchangeDeclare(EXCHANGE, "topic", true);
        Random r = new Random();

        while (true) {
            String t = TITULOS[r.nextInt(TITULOS.length)];
            String e = EVENTOS[r.nextInt(EVENTOS.length)];
            String hora = new SimpleDateFormat("HH:mm:ss").format(new Date());
            String body = String.format(
                "{\"genero\":\"%s\",\"titulo\":\"%s\",\"evento\":\"%s\",\"productor\":\"java\",\"hora\":\"%s\"}",
                GENERO, t, e, hora);
            ch.basicPublish(EXCHANGE, "anime." + GENERO, null, body.getBytes("UTF-8"));
            System.out.println("[PUBLICADO] anime." + GENERO + " -> " + t + " | " + e);
            Thread.sleep(2000 + r.nextInt(3000));
        }
    }
}