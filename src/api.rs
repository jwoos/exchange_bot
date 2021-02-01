use std::convert::Infallible;

use crate::handlers;

use warp::Filter;

lazy_static::lazy_static! {
    static ref REQUEST_CLIENT: reqwest::Client = {
        reqwest::Client::new()
    };
}

pub fn compose_api(
) -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    status().or(events())
}

fn status() -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("slack" / "status")
        .and(warp::get())
        .and_then(handlers::status)
}

fn with_request_client() -> impl Filter<Extract = (&'static reqwest::Client,), Error = Infallible> + Clone {
    warp::any().map(|| &(*REQUEST_CLIENT))
}

fn events() -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("slack" / "events")
        .and(warp::post())
        .and(warp::body::json())
        .and(with_request_client())
        .and_then(handlers::events)
}
