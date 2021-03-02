use std::collections::HashMap;
use std::convert::Infallible;

use crate::handlers;
use crate::user::User;

use warp::Filter;

lazy_static::lazy_static! {
    static ref REQUEST_CLIENT: reqwest::Client = {
        reqwest::Client::new()
    };

    static ref USER_TABLE: HashMap<String, User> = {
        HashMap::new()
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

fn with_request_client(
) -> impl Filter<Extract = (&'static reqwest::Client,), Error = Infallible> + Clone {
    warp::any().map(|| &(*REQUEST_CLIENT))
}

fn with_user_table(
) -> impl Filter<Extract = (&'static HashMap<String, User>,), Error = Infallible> + Clone {
    warp::any().map(|| &(*USER_TABLE))
}

fn events() -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("slack" / "events")
        .and(warp::post())
        .and(warp::body::json())
        .and(with_request_client())
        .and(with_user_table())
        .and_then(handlers::events)
}
