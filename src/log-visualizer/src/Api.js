import axios from "axios";

const baseURL = "http://localhost:3170/api/";

export async function getMsgs() {
    var response = await axios.get(baseURL + "msgs");
    return response.data;
     
    /*[{hash:..}, {hash:..}]*/
}

export async function getMsgData(hash) {
    var response = await axios.get(baseURL + "msg-data" + hash);
    return response.data;
    /*
    {
    locations: { 
        10.0.0.1: {lng: , lat:} 
        10.0.0.2: {lng: , lat:} 
        }, 
    path: [10.0.0.1, 10.0.0.2, ...] 
    } */
    
}

