package dev.bench.config;

import org.springframework.context.annotation.Configuration;
import org.yaml.snakeyaml.Yaml;
import jakarta.annotation.PostConstruct;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Map;

@Configuration
public class YamlConfig {

    private static final String[] PATHS = {
        "config.yaml",
        "../../config.yaml",
        "/config/config.yaml"
    };

    public static String DB_URL;
    public static int DB_MAX_CONNECTIONS = 400;

    @PostConstruct
    public void load() throws Exception {
        Yaml yaml = new Yaml();
        for (String p : PATHS) {
            var path = Paths.get(p);
            if (Files.exists(path)) {
                try (var input = Files.newInputStream(path)) {
                    Map<String, Object> cfg = yaml.load(input);
                    DB_URL = (String) cfg.get("db");
                    DB_MAX_CONNECTIONS = ((Number) cfg.getOrDefault("db_max_connections", 400)).intValue();
                    return;
                }
            }
        }
        throw new IllegalStateException("config.yaml not found");
    }
}