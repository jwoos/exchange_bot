use warp::Filter;

#[tokio::main]
async fn main() {
    pretty_env_logger::init();

    let log = warp::log("exchange");

    // GET /hello/warp => 200 OK with body "Hello, warp!"
    let hello = warp::path!("hello" / String)
        .map(|name| format!("Hello, {}!", name))
        .with(log);

    warp::serve(hello)
        .run(([127, 0, 0, 1], 3030))
        .await;
}
