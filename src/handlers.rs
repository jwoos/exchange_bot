use std::convert::Infallible;

use crate::processor::Processor;
use crate::slack::{event, event_wrapper};

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

struct EventsError<'a> {
    message: &'a str,
}

#[derive(serde::Serialize)]
struct UrlVerificationResponse<'a> {
    challenge: &'a str,
}

impl<'a> UrlVerificationResponse<'a> {
    fn new(challenge: &'a str) -> UrlVerificationResponse {
        UrlVerificationResponse { challenge }
    }
}

pub async fn events(
    event: serde_json::Value,
    client: &reqwest::Client,
) -> Result<impl warp::Reply, Infallible> {
    let event_type_opt = event
        .get("type")
        .and_then(|event_type: &serde_json::Value| event_type.as_str());

    if let Some(event_type) = event_type_opt {
        match event_type {
            "url_verification" => {
                if let Ok(url_verification) =
                    serde_json::from_value::<event::url_verification::UrlVerification>(event)
                {
                    return Ok(warp::reply::with_status(
                        warp::reply::json(&UrlVerificationResponse::new(
                            url_verification.get_challenge(),
                        )),
                        http::status::StatusCode::OK,
                    ));
                } else {
                    return Ok(warp::reply::with_status(
                        warp::reply::json(&serde_json::json!({"message": "Invalid event"})),
                        http::status::StatusCode::BAD_REQUEST,
                    ));
                }
            }
            "event_callback" => {
                if let Ok(event_wrapper) =
                    serde_json::from_value::<event_wrapper::EventWrapper>(event)
                {
                    let proessor = Processor::new(event_wrapper);
                    return proessor.process().await;
                } else {
                    return Ok(warp::reply::with_status(
                        warp::reply::json(&serde_json::json!({"message": "Invalid event"})),
                        http::status::StatusCode::BAD_REQUEST,
                    ));
                }
            }
            _ => {
                return Ok(warp::reply::with_status(
                    warp::reply::json(&serde_json::json!({"message": "Invalid message"})),
                    http::status::StatusCode::BAD_REQUEST,
                ))
            }
        }
    }

    Ok(warp::reply::with_status(
        warp::reply::json(&serde_json::json!({"message": "Invalid message"})),
        http::status::StatusCode::BAD_REQUEST,
    ))
}
