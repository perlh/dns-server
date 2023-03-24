import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import Apis from '@/apis'
import Store from '@/util/store'
import Tool from '@/util/tool'
import router from "@/router";
createApp(App).use(ElementPlus).use(router).use(Store).use(Apis).use(Tool).mount('#app')
