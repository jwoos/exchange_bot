pub mod help;

use std::collections::HashMap;
use std::convert::Infallible;

use super::slack;
use crate::user::User;

pub struct Processor {
    event_wrapper: slack::event_wrapper::EventWrapper,
}

// TODO make Processor a trait and have a processor factory delegate to various processors
impl Processor {
    pub fn new(event_wrapper: slack::event_wrapper::EventWrapper) -> Processor {
        Processor { event_wrapper }
    }

    pub async fn process(
        &self,
        client: &reqwest::Client,
        user_table: &HashMap<String, User>,
    ) -> Result<warp::reply::WithStatus<warp::reply::Json>, Infallible> {
        match self.event_wrapper.get_event().get_type() {
            slack::event::EventType::AppMention => {
                if let Some(event) = self
                    .event_wrapper
                    .get_event()
                    .as_any()
                    .downcast_ref::<slack::event::app_mention::AppMention>()
                {
                    let user_id = event.get_user();
                    return Ok(warp::reply::with_status(
                        warp::reply::json(&serde_json::json!({"message": "invalid event"})),
                        http::status::StatusCode::BAD_REQUEST,
                    ));
                } else {
                    return Ok(warp::reply::with_status(
                        warp::reply::json(&serde_json::json!({"message": "invalid event"})),
                        http::status::StatusCode::BAD_REQUEST,
                    ));
                }
            }
            _ => {
                return Ok(warp::reply::with_status(
                    warp::reply::json(&serde_json::json!({"message": "invalid event"})),
                    http::status::StatusCode::BAD_REQUEST,
                ));
            }
        }
    }
}
