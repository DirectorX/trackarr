import Vue from 'vue'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify';
import axios from 'axios'
import VueNativeSock from 'vue-native-websocket'

Vue.config.productionTip = false
Vue.prototype.$axios = axios


Vue.use(require('vue-moment'));


Vue.use(VueNativeSock, process.env.VUE_APP_WEBSOCKET, {
  reconnection: true, // (Boolean) whether to reconnect automatically (false)
})

new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
