import axios from "axios";

export async function getName(port) {
    var response = await axios.get("http://localhost:" + port + "/api/name");
    return response.data;
}

export async function postMsg(msg, port) {
    axios.post("http://localhost:" + port + "/api/code-msg", msg, {headers: {
        'Content-Type': 'application/json'
    }});
}

export async function postAudioMsg(audio, port) {
    const formData = new FormData();
    formData.append("audio", audio);
    axios.post("http://localhost:" + port + "/api/audio-msg", formData, {headers: {
        'Content-Type': 'multipart/form-data'
      }});
}
