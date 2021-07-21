import axios from "axios";

export const baseURL = "http://localhost:3270/";
const apiURL = baseURL + "api/";

export async function postMsg(msg) {
    axios.post(apiURL + "code-msg", msg, {headers: {
        'Content-Type': 'application/json'
    }});
}

export async function postAudioMsg(audio) {
    const formData = new FormData();
    formData.append("audio", audio);
    axios.post(apiURL + "audio-msg", formData, {headers: {
        'Content-Type': 'multipart/form-data'
      }});
}
