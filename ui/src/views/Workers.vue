<template>
  <v-app>
    <v-main>
      <v-container>
        <NavBar v-bind:appConfig="appConfig" v-bind:userConfig="userConfig" />
        <Workers v-bind:appConfig="appConfig" v-bind:userConfig="userConfig" />
        <Footer v-bind:appConfig="appConfig" v-bind:userConfig="userConfig" />
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import NavBar from '../components/NavBar';
import Workers from '../components/Workers';
import Footer from '../components/Footer';

export default {
  data: () => ({
    appConfig: {},
    userConfig: {
      account: {},
    },
  }),
  components: {
    NavBar,
    Workers,
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
