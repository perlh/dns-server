import request from '@/util/request'

export default {
    dnsSearch(params){
        return request.get("/api/dns/search",{
            params: params
        })
    }
}