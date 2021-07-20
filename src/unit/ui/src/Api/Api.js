import axios from "axios";

export const baseURL = "http://localhost:8000/";
const apiURL = baseURL + "api/";

export async function getMsgs(ip) {
    var response = await axios.get(apiURL + "msgs/"+ ip);
    return response.data;
}

export async function postMsg(ip, msg) {
    axios.post(apiURL + "msg/"+ ip, msg);
}

export async function postAudioMsg(ip, audio) {
    const formData = new FormData();
    formData.append("audio", audio);

    var response = await axios.post(apiURL + "audio-msg/"+ ip, formData, {headers: {
        'Content-Type': 'multipart/form-data'
      }});
    return response.data;
}

export async function getAudioMsgs(ip) {
    var response = await axios.get(apiURL + "audio-msgs/"+ ip);
    return response.data;
}
