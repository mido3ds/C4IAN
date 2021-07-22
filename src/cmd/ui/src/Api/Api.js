import axios from "axios";


export async function postMsg(ip, msg, port) {
    axios.post("http://localhost:" + port + "/api/msg/"+ ip, msg);
}

export async function postAudioMsg(ip, audio, port) {
    const formData = new FormData();
    formData.append("audio", audio);

    var response = await axios.post("http://localhost:" + port + "/api/audio-msg/"+ ip, formData, {headers: {
        'Content-Type': 'multipart/form-data'
      }});
    return response.data;
}

export async function getAllMsgs(port) {
    var response = await axios.get("http://localhost:" + port + "/api/received-msgs");
    return response.data;
}

export async function getNames(port) {
    var response = await axios.get("http://localhost:" + port + "/api/units-names");
    return response.data;
}

export async function getUnits(port) {
    var response = await axios.get("http://localhost:" + port + "/api/units");
    return response.data;
}

export async function getGroups(port) {
    var response = await axios.get("http://localhost:" + port + "/api/groups");
    return response.data;
}

export async function getMembers(port) {
    var response = await axios.get("http://localhost:" + port + "/api/memberships");
    return response.data;
}

export async function getAudioMsgs(ip, port) {
    var response = await axios.get("http://localhost:" + port + "/api/audio-msgs/"+ ip);
    return response.data;
}

export async function getMsgs(ip, port) {
    var response = await axios.get("http://localhost:" + port + "/api/msgs/"+ ip);
    return response.data;
}

export async function getVideos(ip, port) {
    var response = await axios.get("http://localhost:" + port + "/api/videos/"+ ip);
    return response.data;
}

export async function getSensorsData(ip, port) {
    var response = await axios.get("http://localhost:" + port + "/api/sensors-data/"+ ip);
    return response.data;
}
