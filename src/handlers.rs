use crate::slack::event;
use std::convert::Infallible;

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
    return Ok(warp::reply::json(
        &serde_json::json!({"challenge": event.get_challenge()}),
    ));
}

pub async fn events(event: impl event::Event) -> Result<impl warp::Reply, Infallible> {
    Ok(warp::reply::json(
        &serde_json::json!({"status": "Invalid event type!"}),
    ))
}
