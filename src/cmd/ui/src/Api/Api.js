import axios from "axios";

const baseURL = "http://localhost:3170/api/";

export async function postMsg(ip, msg) {
    axios.post(baseURL + "msg/"+ ip, msg);
}

export async function postAudioMsg(ip, audio) {
    const formData = new FormData();
    formData.append("audio", audio);

    var response = await axios.post(baseURL + "audio-msg/"+ ip, formData, {headers: {
        'Content-Type': 'multipart/form-data'
      }});
    return response.data;
}

export async function getAllMsgs() {
    var response = await axios.get(baseURL + "received-msgs");
    return response.data;
}

export async function getNames() {
    var response = await axios.get(baseURL + "units-names");
    return response.data;
}

export async function getUnits() {
    var response = await axios.get(baseURL + "units");
    return response.data;
}

export async function getGroups() {
    var response = await axios.get(baseURL + "groups");
    return response.data;
}

export async function getMembers() {
    var response = await axios.get(baseURL + "memberships");
    return response.data;
}

export async function getAudioMsgs(ip) {
    var response = await axios.get(baseURL + "audio-msgs/"+ ip);
    return response.data;
}

export async function getMsgs(ip) {
    var response = await axios.get(baseURL + "msgs/"+ ip);
    return response.data;
}

export async function getVideos(ip) {
    var response = await axios.get(baseURL + "videos/"+ ip);
    return response.data;
}

export async function getSensorsData(ip) {
    var response = await axios.get(baseURL + "sensors-data/"+ ip);
    return response.data;
}
