<template>
  <v-footer absolute bottom dark class="transparent">
    <v-col
      class="text-center"
      cols="12"
    >
    <div class="text-caption">&copy; Fynca Project {{ new Date().getFullYear() }}</div>
    <div class="text-caption">{{appConfig.name}} {{appConfig.version}}{{appConfig.build}} {{appConfig.commit}}</div>
    </v-col>
  </v-footer>
</template>

<script>
  import axios from "axios"

  export default {
    props: ['appConfig'],
    methods: {},
    created() {
    }
  }
</script>
