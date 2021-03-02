mod api;
mod asset;
mod handlers;
mod iex;
mod portfolio;
mod processor;
mod slack;
mod user;

use warp::Filter;

#[tokio::main]
async fn main() {
    pretty_env_logger::init();

    let api = api::compose_api().with(warp::log("exchange"));

    warp::serve(api).run(([127, 0, 0, 1], 8000)).await;
}
