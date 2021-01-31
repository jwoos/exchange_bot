use warp::Filter;

use crate::handlers;
use crate::slack::{event, event_wrapper};

pub fn compose_api(
) -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    status().or(events()).or(events_url_verification())
}

fn status() -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("slack" / "status")
        .and(warp::get())
        .and_then(handlers::status)
}

fn post_json_filter(
) -> impl warp::Filter<Extract = (serde_json::Value,), Error = warp::Rejection> + Clone {
    warp::post().and(warp::body::json())
}

fn events_url_verification_filter(
) -> impl warp::Filter<Extract = (event::url_verification::UrlVerification,), Error = warp::Rejection>
       + Clone {
    post_json_filter().and_then(|body: serde_json::value::Value| async move {
        let opt = body
            .get("type")
            .and_then(|event_type: &serde_json::Value| event_type.as_str())
            .and_then(|event_type: &str| {
                if event_type == "url_verification" {
                    return Some(());
                } else {
                    return None;
                }
            });

        if opt.is_none() {
            return Err(warp::reject::reject());
        } else {
            return serde_json::from_value(body).or(Err(warp::reject::reject()));
        }
    })
}

fn events_filter(
) -> impl warp::Filter<Extract = (event_wrapper::EventWrapper,), Error = warp::Rejection> + Clone {
    post_json_filter().and_then(|body: serde_json::value::Value| async move {
        let opt = body
            .get("type")
            .and_then(|event_type: &serde_json::Value| event_type.as_str())
            .and_then(|event_type: &str| {
                if event_type == "event_callback" {
                    return Some(());
                } else {
                    return None;
                }
            });

        if opt.is_none() {
            return Err(warp::reject::reject());
        } else {
            return serde_json::from_value(body).or(Err(warp::reject::reject()));
        }
    })
}

fn events_url_verification(
) -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("slack" / "events")
        .and(events_url_verification_filter())
        .and_then(handlers::events_url_verification)
}

// TODO fix
fn events() -> impl warp::Filter<Extract = impl warp::Reply, Error = warp::Rejection> + Clone {
    warp::path!("slack" / "events").and(post_json_filter().and_then(handlers::echo))
}
