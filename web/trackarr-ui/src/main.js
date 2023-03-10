/*eslint no-undef: "off"*/
// Overwrite Vue publicPath with a runtime dynamic var
let CORE_APP_BASE_URL = '';

if (process.env.VUE_APP_BASE_URL) {
    __webpack_public_path__ = process.env.VUE_APP_BASE_URL.endsWith("/") ? process.env.VUE_APP_BASE_URL : process.env.VUE_APP_BASE_URL + "/";
    CORE_APP_BASE_URL = process.env.VUE_APP_BASE_URL;
} else {
    __webpack_public_path__ = window.baseurl.endsWith("/") ? window.baseurl : window.baseurl + "/";
    CORE_APP_BASE_URL = window.baseurl;
}

import Vue from 'vue'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify';
import axios from 'axios'
import VueNativeSock from 'vue-native-websocket'
import VueToastr from "vue-toastr";

var path = require('path');

/* generate dynamic variables */

// - api
let CORE_API_URL = '';
let CORE_API_KEY = '';
let CORE_APP_VERSION = '';

if (process.env.VUE_APP_API_URL && process.env.VUE_APP_API_KEY) {
    CORE_API_URL = process.env.VUE_APP_API_URL;
    CORE_API_KEY = process.env.VUE_APP_API_KEY;
    CORE_APP_VERSION = 'Local Build';
} else {
    CORE_API_URL = new URL(path.join(window.baseurl, '/api'), window.location.href).href;
    CORE_API_KEY = window.apikey;
    CORE_APP_VERSION = window.version;
}

// - version

// - websocket
let CORE_WEBSOCKET_URL = '';
if (process.env.VUE_APP_WEBSOCKET) {
    CORE_WEBSOCKET_URL = process.env.VUE_APP_WEBSOCKET;
} else {
    let socketUrl = new URL(path.join(window.baseurl, '/api/ws'), window.location.href);
    socketUrl.protocol = socketUrl.protocol.replace('http', 'ws');
    socketUrl.searchParams.set('apikey', CORE_API_KEY);
    CORE_WEBSOCKET_URL = socketUrl.href;
}


/* eslint-disable no-console */

// log dynamically generated variables
console.log('Using WEBSOCKET_URL =', CORE_WEBSOCKET_URL);
console.log('Using API_URL =', CORE_API_URL);
console.log('Using API_KEY =', CORE_API_KEY);
console.log('Using APP_BASE_URL =', CORE_APP_BASE_URL);
console.log('using APP_VERSION =', CORE_APP_VERSION);

/* eslint-enable no-console */


/* Vue init */

Vue.prototype.CORE_API_URL = CORE_API_URL;
Vue.prototype.CORE_API_KEY = CORE_API_KEY;
Vue.prototype.CORE_WEBSOCKET_URL = CORE_WEBSOCKET_URL;
Vue.prototype.CORE_APP_VERSION = CORE_APP_VERSION;

Vue.prototype.$axios = axios.create({
    baseURL: CORE_API_URL
});
Vue.config.productionTip = false;

Vue.use(require('vue-moment'));
Vue.use(VueToastr, {
    defaultTimeout: 5000,
    defaultPosition: 'toast-bottom-right',
    defaultType: 'info',
    defaultStyle: {'font-family':'roboto','opacity':'100%'}
});

Vue.use(VueNativeSock, CORE_WEBSOCKET_URL, {
    format: 'json',
    reconnection: true, // (Boolean) whether to reconnect automatically (false)
});


Vue.filter('capitalize',function (value) {
        if (!value) return '';
        value = value.toString();
        return value.charAt(0).toUpperCase() + value.slice(1)
    })

new Vue({
    router,
    vuetify,
    render: h => h(App),
    mounted: function () {
        // global websocket message handler
        this.$options.sockets.onmessage = (message) => {
            // parse message
            let event = JSON.parse(message.data);
            
            // alert events
            if (event.type === 'alert') {
                switch (event.data.level) {
                    case 'info':
                        this.$toastr.i(event.data.msg, event.data.title);
                        break;
                    case 'success':
                        this.$toastr.s(event.data.msg, event.data.title);
                        break;
                    case 'warn':
                        this.$toastr.w(event.data.msg, event.data.title);
                        break;
                    case 'error':
                        this.$toastr.e(event.data.msg, event.data.title);
                        break;
                    default:
                        break;
                }
            }
        }
    }
}).$mount('#app');
