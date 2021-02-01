pub mod help;

use std::convert::Infallible;

use super::slack::event_wrapper;

pub struct Processor {
    event: event_wrapper::EventWrapper,
}

// TODO make Processor a trait and have a processor factory delegate to various processors
impl Processor {
    pub fn new(event: event_wrapper::EventWrapper) -> Processor {
        Processor { event }
    }

    pub async fn process(&self) -> Result<warp::reply::WithStatus<warp::reply::Json>, Infallible> {
        println!("yay");
        return Ok(warp::reply::with_status(
            warp::reply::json(&serde_json::json!({"message": "okay"})),
            http::status::StatusCode::OK,
        ));
    }
}
