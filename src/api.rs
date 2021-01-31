use warp::Filter;

use crate::handlers;

pub fn compose_api(
) -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    status().or(events())
}

fn status() -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("slack" / "status")
        .and(warp::get())
        .and_then(handlers::status)
}

fn events() -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("slack" / "events")
        .and(warp::post())
        .and(warp::body::json())
        .and_then(handlers::events)
}
