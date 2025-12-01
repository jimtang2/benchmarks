use axum::{extract::State, routing::get, Router};
use futures_util::future::join_all;
use rand::{rngs::StdRng, SeedableRng};
use rand::Rng;
use serde::Deserialize;
use sqlx::{postgres::PgPoolOptions, PgPool};
use std::{net::SocketAddr, sync::Arc};
use axum::http::StatusCode;
use tokio::fs;

#[cfg(target_env = "musl")]
#[global_allocator]
static GLOBAL: mimalloc::MiMalloc = mimalloc::MiMalloc;

#[derive(Deserialize, Clone)]
struct Config {
    db: String,
    db_max_connections: u32,
}

struct AppState {
    pool: PgPool,
    rng: StdRng,
}

const SCALE: i64 = 10;

#[tokio::main(flavor = "multi_thread", worker_threads = 16)] 
async fn main() {
    let paths = ["./config.yaml", "../config.yaml", "/config/config.yaml"];
    let (config_path, config_text) = join_all(paths.iter().map(|&p| async move {
        match fs::read_to_string(p).await {
            Ok(content) => Some((p, content)),
            Err(_) => None,
        }
    }))
    .await
    .into_iter()
    .flatten()
    .next()
    .expect("Could not find config.yaml");

    println!("[CONFIG] Loaded config from: {config_path}");
    let config: Config = serde_yaml::from_str(&config_text).expect("Failed to parse config.yaml");

    let pool = PgPoolOptions::new()
        .max_connections(config.db_max_connections)
        .connect(&config.db)
        .await
        .expect("Failed to connect to database");

    sqlx::query("SELECT 1").fetch_one(&pool).await.unwrap();
    println!("[DB] Connected successfully");

    let state = Arc::new(AppState {
        pool,
        rng: StdRng::from_entropy(),
    });

    let app = Router::new()
        .route("/", get(handler))
        .route("/*_", get(handler))
        .with_state(state);

    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    println!("[HTTP] Listening on http://{addr}");

    axum::serve(tokio::net::TcpListener::bind(&addr).await.unwrap(), app)
        .await
        .unwrap();
}

#[axum::debug_handler]
async fn handler(State(state): State<Arc<AppState>>) -> StatusCode {
    let mut rng = rand::thread_rng();

    let aid   = rng.gen_range(1..=SCALE * 100_000);
    let tid   = rng.gen_range(1..=SCALE * 10);
    let bid   = rng.gen_range(1..=SCALE * 1);
    let delta = rng.gen_range(-5000..=4999i32);

    match sqlx::query("SELECT pgbench_tx($1,$2,$3,$4)")
        .bind(aid)
        .bind(tid)
        .bind(bid)
        .bind(delta)
        .persistent(true)
        .execute(&state.pool)
        .await
    {
        Ok(_) => StatusCode::OK,
        Err(e) => {
            eprintln!("tx failed: {e} | a={aid} t={tid} b={bid} Δ={delta}");
            StatusCode::INTERNAL_SERVER_ERROR
        }
    }
}