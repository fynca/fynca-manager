import Vue from 'vue'
import Router from 'vue-router'
import Jobs from './views/Jobs.vue'
import Job from './views/Job.vue'
import Workers from './views/Workers.vue'
import Login from './views/Login.vue'
import Profile from './views/Profile.vue'
import About from './views/About.vue'

Vue.use(Router)

export default new Router({
  routes: [
    { path: '/',
      redirect: '/jobs',
      props: true,
      name: 'main'
    },
    {
      path: '/jobs',
      name: 'jobs',
      props: true,
      component: Jobs
    },
    {
      path: '/jobs/:id',
      name: 'job',
      props: true,
      component: Job
    },
    {
      path: '/workers',
      name: 'workers',
      props: true,
      component: Workers
    },
    {
      path: '/login',
      name: 'login',
      props: true,
      component: Login
    },
    {
      path: '/profile',
      name: 'profile',
      props: true,
      component: Profile
    },
    {
      path: '/about',
      name: 'about',
      props: true,
      component: About
    }
  ]
})
