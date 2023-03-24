import request from '@/util/request'

export default {
    checkLogin: () => {
        return request.get("/api/auth/checkLogin")
    },
    login: (username, password) => {
        return request.postForm("/api/auth/login",{
            username,password
        })
    },
    logout: () => {
        return request.get("/api/auth/logout")
    }
}