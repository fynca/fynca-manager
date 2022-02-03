<template>
  <v-row class="pb-12">
    <v-overlay :value="loading" class="text-center">
      <v-progress-circular
        :size="100"
        :width="10"
        color="blue"
        indeterminate
      ></v-progress-circular>
      <div class="text-center text-h4 font-weight-light grey--text">Loading</div>
    </v-overlay>

    <v-col v-if="showEmptyJobsMessage">
      <v-card align="center" justify="center" elevation="0">
      <v-img
        contain
        max-width="75%"
        src="/img/fynca-default.png"
      >
      </v-img>
        <v-card-text class="text-h4 font-weight-light">
          Start a render job by clicking the "Render" button in the top right.
        </v-card-text>
      </v-card>
    </v-col>

    <JobCard v-for="job in jobs" :key="job.id" v-bind:job="job" />

  </v-row>
</template>

<script>
import axios from "axios";
import JobCard from './JobCard.vue';

export default {
  components: {
    JobCard,
  },
  props: {},
  data: () => ({
    jobs: [],
    loading: true,
    snackbar: false,
    text: '',
    cacheKey: '',
    timer: '',
    color: 'error',
    showEmptyJobsMessage: false,
  }),
  methods: {
    loadJobs() {
      axios.get(this.$apiHost + '/api/v1/jobs')
        .then(resp => {
          this.jobs = []
          var x
          if (resp.data.jobs.length === 0) {
            this.showEmptyJobsMessage = true
            this.loading = false
            return
          }
          this.showEmptyJobsMessage = false
          this.loading = false

          var jobs = resp.data.jobs
          for (x=0; x<jobs.length; x++) {
            this.jobs.push(jobs[x])
          }
        })
        .catch(err => {
          var msg
          if (err.response != null) {
            msg = err.response.data
          } else {
            msg = err
          }
          this.loading = false
          console.log(err)
        })
    },
    enableLoadJobs() {
      this.loadJobs()
      this.timer = setInterval(this.loadJobs, 10000)
    },
    cancelLoadJobs() {
      clearInterval(this.timer)
    }
  },
  mounted () {
    this.$root.$on('cancelLoadJobs', () => {
      this.cancelLoadJobs()
    })
    this.$root.$on('enableLoadJobs', () => {
      this.enableLoadJobs()
    })
    this.$root.$on('loadJobs', () => {
      this.loadJobs()
    })
  },
  created () {
    this.loadJobs()
    this.enableLoadJobs()
  },
  beforeDestroy() {
    this.cancelLoadJobs()
  },
}
</script>
