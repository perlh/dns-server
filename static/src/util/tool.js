export default {
    install(app){
        app.config.globalProperties.$tool = this
    },
    splitText(str,len=30){
        if(str.length > len){
            return str.substr(0,len) + "..."
        }
        return str
    }
}