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

    <v-dialog
      v-model="showUpdateWorkerDialog"
      width="500"
    >
      <v-card>
        <v-card-title class="grey">
          Update Worker
        </v-card-title>

        <v-card-text>
          <v-form
            ref="updateWorkerForm"
            v-model="updateWorkerFormValid"
            lazy-validation
          >
            <v-text-field
              type="text"
              v-model="updateWorkerURL"
              :rules="updateWorkerFormURLRules"
              label="Update URL"
              required
            ></v-text-field>

          </v-form>
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            @click="updateWorker(name)"
          >
            Update
          </v-btn>
          <v-btn
            @click="showUpdateWorkerDialog = false"
          >
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-icon
      @click="showUpdateWorkerDialog = true"
    >
      mdi-update
    </v-icon>
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
    showUpdateWorkerDialog: false,
    showStopWorkerDialog: false,
    updateWorkerFormValid: false,
    updateWorkerURL: '',
    updateWorkerFormURLRules: [
      v => !!v || 'Update URL is required'
    ],
  }),
  methods: {
    updateWorker(worker) {
      if (!this.$refs.updateWorkerForm.validate()) {
        return
      }
      let formData = new FormData()
      formData.append("url", this.updateWorkerURL)
      axios.post(this.$apiHost+'/api/v1/workers/' + worker + '/update', formData)
        .then( res => {
          this.$root.$emit('showInfo', 'Worker ' + worker + ' will update after current jobs.')
          this.showUpdateWorkerDialog = false
          this.updateWorkerURL = '';
      })
      .catch(err => {
        console.log(err)
        var msg
        if (err.response != null) {
          msg = err.response.data
        } else {
          msg = err
        }
	      this.$root.$emit('showError', 'Error updating worker: ' + msg)
        this.showUpdateWorkerDialog = false
        this.updateWorkerURL = '';
      })
    },
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
