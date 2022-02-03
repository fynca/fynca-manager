<template>
  <router-view></router-view>
</template>

<script>
  import axios from "axios"

  export default {
    computed: {},
    created () {
      this.$vuetify.theme.dark = true
    }
  }
</script>
