import axios from "axios";

export const baseURL = "http://localhost:8000/";
const apiURL = baseURL + "api/";


export function getMsgs(ip) {
    axios.get(apiURL + "msgs"+ ip).then((response) => {
        return response.data;
    });
}

export function postMsg(ip, msg) {
    axios.post(apiURL + "msgs"+ ip, msg).then((response) => {
        return response;
    });
}

export function postAudioMsg(ip, audio) {
    axios.post(apiURL + "audio-msgs"+ ip, audio).then((response) => {
        return response;
    });
}

export function getAudioMsgs(ip) {
    axios.get(apiURL + "audio-msgs/"+ ip).then((response) => {
        return response.data;
    });
}
