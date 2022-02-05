import Vue from 'vue'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify'
import axios from 'axios'
import md5 from 'md5'
import 'material-design-icons-iconfont/dist/material-design-icons.css'

Vue.config.productionTip = false
Vue.prototype.$appTitle = 'Fynca'
Vue.prototype.$apiHost = process.env.VUE_APP_API_HOST ? process.env.VUE_APP_API_HOST : ''
Vue.prototype.$filters = Vue.options.filters

// filters
Vue.filter('round', function(value, decimals) {
  if(!value) {
    value = 0;
  }

  if(!decimals) {
    decimals = 0;
  }

  value = Math.round(value);
  return value;
});

Vue.filter('gravatar', function(v, size) {
  if(v === undefined || v === '') {
    return ''
  }
  if (size === undefined ) {
    size = '128'
  }
  const hash = md5(v.trim().toLowerCase())
  return `https://www.gravatar.com/avatar/${hash}` + '?s=' + size;
});

Vue.filter('toTitleCase', function(v) {
  return v.replace(
    /\w\S*/g,
    function(txt) {
      return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
    }
  );
});

Vue.filter('formatRenderEngineName', function(v) {
  return v.replace('BLENDER_', '');
});

// format bytes
Vue.filter('formatBytes', function (num) {
  if (typeof num !== 'number' || isNaN(num)) {
    return 0;
  }

  var exponent;
  var unit;
  var neg = num < 0;
  var units = ['B', 'kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  if (neg) {
    num = -num;
  }

  if (num < 1) {
    return (neg ? '-' : '') + num + ' B';
  }

  exponent = Math.min(Math.floor(Math.log(num) / Math.log(1000)), units.length - 1);
  num = (num / Math.pow(1000, exponent)).toFixed(2) * 1;
  unit = units[exponent];

  return (neg ? '-' : '') + num + ' ' + unit;
});

// auth tokens
axios.defaults.headers.common['X-Session-Token'] = localStorage.token !== undefined ? localStorage.token : "";
axios.interceptors.request.use(function (config) {
  return config;
}, function (error) {
  return Promise.reject(error);
});
axios.interceptors.response.use(function (response) {
  return response;
}, function (err) {
  var r = err.response;
  if (r.status === 401 || r.status === 403) {
    router.push({name: "login"});
    return
  }
  if (r.status === 500) {
    if (r.data.includes('invalid or missing token')) {
      router.push({name: "login"});
      return
    }
  }
  return Promise.reject(err);
});

// global helper funcs
Vue.prototype.$getAppConfig = function() {
  return axios.get(Vue.prototype.$apiHost + '/versionz')
    .then(resp => {
      return resp.data
    })
    .catch(err => {
      var msg
      if (err.response != null) {
        msg = err.response.data
      } else {
        msg = err
      }
      this.loading = false
    });
}

Vue.prototype.$getUserConfig = function() {
  return axios.get(Vue.prototype.$apiHost + '/api/v1/accounts/profile')
    .then(resp => {
      return resp.data
    })
    .catch(err => {
      var msg
      if (err.response != null) {
        msg = err.response.data
      } else {
        msg = err
      }
      this.loading = false
    });
}

// app
new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
