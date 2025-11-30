package dev.bench.web;

import dev.bench.config.YamlConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.io.buffer.DataBufferFactory;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.r2dbc.core.DatabaseClient;
import reactor.core.publisher.Mono;

import java.util.Random;

@RestController
public class TxController {

    private final DatabaseClient client;
    private final Random rnd = new Random();
    private static final int SCALE = 10;

    public TxController(@Autowired DatabaseClient client) {
        this.client = client;
    }

    @GetMapping(value = "/", produces = MediaType.TEXT_PLAIN_VALUE)
    public Mono<String> tx() {
        int aid = rnd.nextInt(SCALE * 100_000) + 1;
        int tid = rnd.nextInt(SCALE * 10) + 1;
        int bid = rnd.nextInt(SCALE) + 1;
        int delta = rnd.nextInt(10_000) - 5_000;

        return client.sql("SELECT pgbench_tx($1, $2, $3, $4)")
                .bind(0, aid)
                .bind(1, tid)
                .bind(2, bid)
                .bind(3, delta)
                .fetch()
                .first()
                .then(Mono.fromSupplier(() -> "OK"));
    }
}