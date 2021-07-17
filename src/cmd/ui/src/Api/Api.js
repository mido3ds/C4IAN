import axios from "axios";

const baseURL = "http://localhost:3170/api/";

export function postMsg(ip, msg) {
    axios.post(baseURL + "msg/"+ ip, msg).then((response) => {
        return response;
    });
}

export function postAudioMsg(ip, audio) {
    axios.post(baseURL + "audio-msg/"+ ip, audio).then((response) => {
        return response;
    });
}

export function getUnits() {
    axios.get(baseURL + "units/").then((response) => {
        return response.data;
    });
}

export function getGroups() {
    axios.get(baseURL + "groups/").then((response) => {
        return response.data;
    });
}

export function getMembers() {
    axios.get(baseURL + "groups/").then((response) => {
        return response.data;
    });
}

export function getAudioMsgs(ip) {
    axios.get(baseURL + "audio-msgs/"+ ip).then((response) => {
        return response.data;
    });
}

export function getMsgs(ip) {
    axios.get(baseURL + "msgs"+ ip).then((response) => {
        return response.data;
    });
}

export function getVideos(ip) {
    axios.get(baseURL + "videos"+ ip).then((response) => {
        return response.data;
    });
}

export function getSensorsData(ip) {
    axios.get(baseURL + "videos"+ ip).then((response) => {
        return response.data;
    });
}
