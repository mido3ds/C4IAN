import axios from "axios";

export const baseURL = "http://localhost:8000/";
const apiURL = baseURL + "api/";

export async function postMsg(msg) {
    axios.post(apiURL + "code-msg/", msg);
}

export async function postAudioMsg(audio) {
    const formData = new FormData();
    formData.append("audio", audio);

    var response = await axios.post(apiURL + "audio-msg/", formData, {headers: {
        'Content-Type': 'multipart/form-data'
      }});
    return response.data;
}
