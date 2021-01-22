use warp::Filter;

use crate::handlers;

pub fn compose_api(
) -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    status()
    //status().or(events())
}

fn status() -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("status")
        .and(warp::get())
        .and_then(handlers::status)
}

/*
 *fn events() -> impl warp::Filter + Clone {
 *
 *}
 */
