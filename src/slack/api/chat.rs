/* this module represents the chat api
 */

pub struct PostMessageRequest {
    token: String,
    channel: String,
}

pub struct PostMessageResponse {
    ok: bool,
    channel: String,
}

pub async fn post_message() {

}
