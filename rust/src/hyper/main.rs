use deadpool_postgres::{Manager, Pool, Runtime};
use hyper::service::{make_service_fn, service_fn};
use hyper::{body::Body, Request, Response, server::Server, StatusCode};
use rand::RngCore;
use std::convert::Infallible;
use std::net::SocketAddr;
use tokio::fs;

#[cfg(target_env = "musl")]
#[global_allocator]
static GLOBAL: mimalloc::MiMalloc = mimalloc::MiMalloc;

#[derive(serde::Deserialize)]
struct Config {
    db: String,                  // PostgreSQL connection string
    db_max_connections: usize,   // pool size
    listen: String,              // e.g. "0.0.0.0:3000"
    scale: i64,                  // pgbench scale factor
}

async fn load_config() -> Config {
    let paths = [
        "./config.yaml",
        "../config.yaml",
        "/config/config.yaml",
    ];

    let mut config_text = None;
    for &path in &paths {
        if let Ok(text) = fs::read_to_string(path).await {
            config_text = Some(text);
            break;
        }
    }

    let text = config_text.expect("config.yaml not found");
    serde_yaml::from_str(&text).expect("failed to parse config.yaml")
}

async fn handle_request(
    _req: Request<Body>,
    pool: Pool,
    scale: i64,
) -> Result<Response<Body>, Infallible> {
    let mut client = match pool.get().await {
        Ok(c) => c,
        Err(_) => {
            return Ok(Response::builder()
                .status(StatusCode::SERVICE_UNAVAILABLE)
                .body(Body::from("failed to get db connection"))
                .unwrap());
        }
    };

    let mut rng = rand::thread_rng();

    let aid: i64 = rng.gen_range(1..=scale * 100_000);
    let tid: i64 = rng.gen_range(1..=scale * 10);
    let bid: i64 = rng.gen_range(1..=scale * 1);
    let delta: i64 = rng.gen_range(-5000..=4999);

    let result = client
        .execute("SELECT pgbench_tx($1,$2,$3,$4)", &[&aid, &tid, &bid, &delta])
        .await;

    match result {
        Ok(_) => Ok(Response::new(Body::from("ok"))),
        Err(e) => {
            eprintln!("pgbench_tx error: {e}");
            Ok(Response::builder()
                .status(StatusCode::INTERNAL_SERVER_ERROR)
                .body(Body::from("db error"))
                .unwrap())
        }
    }
}

#[tokio::main]
async fn main() {
    let config = load_config().await;

    let mgr = Manager::from_config(
        &config.db.parse().expect("invalid postgres connection string"),
        tokio_postgres::NoTls,
        deadpool_postgres::PoolConfig {
            max_size: config.db_max_connections,
            ..Default::default()
        },
    );

    let pool = Pool::builder(mgr)
        .runtime(Runtime::Tokio1)
        .build()
        .expect("failed to create deadpool pool");

    let addr: SocketAddr = config.listen.parse().expect("invalid listen address");

    let make_service = make_service_fn(move |_conn| {
        let pool = pool.clone();
        let scale = config.scale;
        async move {
            Ok::<_, Infallible>(service_fn(move |req| {
                handle_request(req, pool.clone(), scale)
            }))
        }
    });

    println!("Listening on http://{addr}");
    Server::bind(&addr)
        .serve(make_service)
        .await
        .unwrap_or_else(|e| eprintln!("server error: {e}"));
}