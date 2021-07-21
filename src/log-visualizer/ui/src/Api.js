import axios from "axios";

const baseURL = "http://localhost:5000/api/";

export async function getMsgs() {
    var response = await axios.get(baseURL + "packets");
    return response.data;
}

export async function getMsgData(hash) {
    var response = await axios.get(baseURL + "packet-data/" + hash);
    return response.data;
}

