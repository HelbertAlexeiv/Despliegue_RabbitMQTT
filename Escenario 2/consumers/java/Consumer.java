import com.rabbitmq.client.*;

public class Consumer {
    public static void main(String[] args) throws Exception {
        String exchange = "anime_exchange";
        String cola = "cola_java";
        String[] claves = {"anime.isekai", "anime.romance"};

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
        ch.exchangeDeclare(exchange, "topic", true);
        ch.queueDeclare(cola, true, false, false, null);
        for (String k : claves) {
            ch.queueBind(cola, exchange, k);
        }

        System.out.println("Java escuchando: " + String.join(", ", claves));
        DeliverCallback cb = (tag, d) ->
            System.out.println("[JAVA] (" + d.getEnvelope().getRoutingKey() + ") " + new String(d.getBody(), "UTF-8"));
        ch.basicConsume(cola, true, cb, tag -> {});
    }
}