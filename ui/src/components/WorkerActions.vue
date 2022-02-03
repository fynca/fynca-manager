<template>
  <v-row class="pl-2">
    <v-dialog
      v-model="showStopWorkerDialog"
      width="500"
    >
      <v-card>
        <v-card-title class="grey">
          Stop Worker
        </v-card-title>

        <v-card-text>
          <div class="text-center">
          Are you sure you want to stop the worker {{name}}?
          </div>
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="red"
            @click="stopWorker(name)"
          >
            Stop
          </v-btn>
          <v-btn
            @click="showStopWorkerDialog = false"
          >
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-icon
      @click="showStopWorkerDialog = true"
    >
      mdi-stop
    </v-icon>
  </v-row>
</template>

<script>
import axios from "axios";

export default {
  components: {},
  props: ['name'],
  data: () => ({
    showStopWorkerDialog: false,
  }),
  methods: {
    stopWorker(worker) {
      axios.post(this.$apiHost + '/api/v1/workers/'+worker+'/stop')
        .then(resp => {
          this.$root.$emit('showInfo', 'Worker ' + worker + ' will stop after current jobs.')
          this.showStopWorkerDialog = false
        })
        .catch(err => {
          var msg
          if (err.response != null) {
            msg = err.response.data
          } else {
            msg = err
          }
          this.showStopWorkerDialog = false
          this.$root.$emit('showError', 'Error stopping worker: ' + msg)
        })
    }
  },
  created () {
  },
  updated() {
  },
}
</script>
