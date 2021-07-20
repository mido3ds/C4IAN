import axios from "axios";

export const baseURL = "http://localhost:8000/";
const apiURL = baseURL + "api/";

export async function postMsg(msg) {
    await axios.post(apiURL + "code-msg/", msg, {headers: {
        'Content-Type': 'application/x-msgpack'
    }});
}

export async function postAudioMsg(audio) {
    // const formData = new FormData();
    // formData.append("audio", audio);

    await axios.post(apiURL + "audio-msg/", audio, {headers: {
        'Content-Type': 'application/x-msgpack'
      }});
}
