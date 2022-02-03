<template>
  <v-app>
    <v-main>
      <v-container fluid>
        <NavBar v-bind:appConfig="appConfig" v-bind:userConfig="userConfig" />
        <Job v-bind:id="id"  v-bind:userConfig="userConfig" />
        <Footer v-bind:appConfig="appConfig"  v-bind:userConfig="userConfig" />
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import NavBar from '../components/NavBar';
import Job from '../components/Job';
import Footer from '../components/Footer';

export default {
  props: ['id'],
  data: () => ({
    appConfig: {},
    userConfig: {
      account: {},
    },
  }),
  components: {
    NavBar,
    Job,
    Footer,
  },
  methods: {
    loadAppConfig: async function() {
      this.appConfig = await this.$getAppConfig()
      this.userConfig = await this.$getUserConfig()
    },
  },
  created () {
    this.loadAppConfig()
    this.$vuetify.theme.dark = true
  },
};
</script>
