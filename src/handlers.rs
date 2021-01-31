use crate::slack::{event, event_wrapper};
use std::convert::Infallible;

pub async fn echo(data: serde_json::Value) -> Result<impl warp::Reply, Infallible> {
    Ok(warp::reply::json(&data))
}

#[derive(serde::Serialize)]
struct Status<'a> {
    status: &'a str,
}

pub async fn status() -> Result<impl warp::Reply, Infallible> {
    Ok(warp::reply::json(&Status { status: "okay" }))
}

pub async fn events_url_verification(
    event: event::url_verification::UrlVerification,
) -> Result<impl warp::Reply, Infallible> {
    Ok(warp::reply::json(
        &serde_json::json!({"challenge": event.get_challenge()}),
    ))
}

pub async fn events(event: event_wrapper::EventWrapper) -> Result<impl warp::Reply, Infallible> {
    Ok(warp::reply::json(
        &serde_json::json!({"status": "received"}),
    ))
}
