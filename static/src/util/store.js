export default {
    install(app) {
        app.config.globalProperties.$store = this
    },
    get(key) {
        let item = localStorage.getItem(key);
        if (item) {
            return JSON.parse(item)
        }
        return ""
    },
    set(key, val) {
        localStorage.setItem(key, JSON.stringify(val))
    },
    del(key) {
        localStorage.removeItem(key)
    },
    getToken() {
        return this.get("token")
    },
    clearToken() {
        this.del("token")
    },
    setToken(val) {
        return this.set("token", val)
    }
}