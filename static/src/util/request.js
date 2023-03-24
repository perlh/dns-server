import axios from 'axios'
import store from "@/util/store";

const Axios = axios.create({
    timeout: 3000
})

Axios.interceptors.request.use(req => {
    let token = store.getToken();
    if(token){
        req.headers["Authorization"] = token
    }
    return req
})

Axios.interceptors.response.use(res => res)

export default Axios