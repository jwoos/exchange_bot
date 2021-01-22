use std::convert::Infallible;

#[derive(serde::Serialize)]
struct Status<'a> {
    status: &'a str,
}

pub async fn status() -> Result<impl warp::Reply, Infallible> {
    Ok(warp::reply::json(&Status { status: "okay" }))
}
