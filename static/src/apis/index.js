import auth from '@/apis/auth'
import dns from "@/apis/dns";
export default {
    install(app){
        app.config.globalProperties.$api = this
    },
    ...auth,
    ...dns
}