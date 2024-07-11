import { createApp } from 'vue'
import App from './App.vue'
import './style.less'
import ElementPlus from 'element-plus' // 引入模块
import 'element-plus/dist/index.css' // 引入css

const app = createApp(App)
app.use(ElementPlus)
app.mount('#app')
