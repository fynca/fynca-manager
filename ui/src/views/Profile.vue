<template>
  <v-app>
    <v-main>
      <v-container fluid>
        <NavBar v-bind:appConfig="appConfig" v-bind:userConfig="userConfig" />
        <Profile v-bind:userConfig="userConfig" />
        <Footer v-bind:appConfig="appConfig"  v-bind:userConfig="userConfig" />
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import NavBar from '../components/NavBar';
import Profile from '../components/Profile';
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
    Profile,
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
