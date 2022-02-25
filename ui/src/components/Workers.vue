<template>
  <v-row>

    <v-col />

    <v-row
      >
      <v-col align="center">
        <v-card
          elevation="5"
          class="ma-2 pa-5"
          >
          <v-card-text class="text-h1">
            {{workersCPU}}
          </v-card-text>
          <v-card-subtitle class="text-subtitle-1">Cores</v-card-subtitle>
        </v-card>
      </v-col>
      <v-col align="center">
        <v-card
          elevation="5"
          class="ma-2 pa-5"
          >
          <v-card-text class="text-h1">
            {{workersMemory}}
          </v-card-text>
          <v-card-subtitle class="text-subtitle-1">Memory</v-card-subtitle>
        </v-card>
      </v-col>

      <v-col xs="12" s="12" md="12" lg="12" xl="12" >
        <v-data-table
          dark
          :headers="headers"
          :items="workers"
          item-key="name"
          class="ma-2 elevation-1"
          :search="search"
          :loading="loading"
          sort-by="name"
          no-data-text="No workers available"
          loading-text="Loading... Please wait"
        >
          <template v-slot:top>
            <v-text-field
              v-model="search"
              label="Search"
              class="mx-4"
            ></v-text-field>
          </template>

          <template v-slot:item.load1="{ item }">
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-progress-linear
                  :value="getLoadPercentage(item)"
                  :color="getUsageColor(getLoadPercentage(item))"
                  v-bind="attrs"
                  v-on="on"
                  >
                </v-progress-linear>
              </template>
              <span>{{getLoadPercentage(item).toFixed(2)}}%</span>
            </v-tooltip>
          </template>

          <template v-slot:item.memoryUsage="{ item }">
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-progress-linear
                  :value="getMemoryPercentage(item)"
                  :color="getUsageColor(getMemoryPercentage(item))"
                  v-bind="attrs"
                  v-on="on"
                  >
                </v-progress-linear>
              </template>
              <span>{{getMemoryPercentage(item).toFixed(2)}}%</span>
            </v-tooltip>
          </template>

          <template v-slot:item.gpus="{ item }">
            <v-tooltip bottom v-for="gpu in item.gpus" :key="gpu">
              <template v-slot:activator="{ on, attrs }">
                <v-icon
                  :color="getGPUColor(gpu)"
                  dark
                  v-bind="attrs"
                  v-on="on"
                >
                  mdi-expansion-card
                </v-icon>
              </template>
              <span>{{gpu}}</span>
            </v-tooltip>
          </template>

          <template v-slot:item.actions="{ item }">
            <WorkerActions v-bind:worker="item" />
          </template>
        </v-data-table>
      </v-col>
    </v-row>

    <v-col />

  </v-row>
</template>

<script>
import axios from "axios";
import WorkerActions from './WorkerActions.vue';

export default {
  components: {
    WorkerActions,
  },
  props: {},
  data: () => ({
    workers: [],
    workersCPU: 0,
    workersMemory: 0,
    loading: true,
    snackbar: false,
    search: '',
  }),
  methods: {
    loadWorkers() {
      axios.get(this.$apiHost + '/api/v1/workers')
        .then(resp => {
          // reset
          this.workers = []
          var x
          var workerMem = 0
          this.workersMemory = this.humanSize(workerMem)
          if (resp.data.workers == null) {
            this.loading = false
            return
          }
          this.loading = false

          var workers = resp.data.workers

          for (x=0; x<workers.length; x++) {
            this.workers.push(workers[x])
            this.workersCPU += workers[x].cpus
            workerMem += parseInt(workers[x].memoryTotal)
            this.workers[x].memory = this.humanSize(this.workers[x].memoryTotal)
          }
          this.workersMemory = this.humanSize(workerMem)
        })
        .catch(err => {
          var msg
          if (err.response != null) {
            msg = err.response.data
          } else {
            msg = err
          }
          this.loading = false
          this.$root.$emit('showError', 'Error loading workers: ' + msg)
        })
    },
    getLoadPercentage(v) {
      if (v.load1 === undefined) {
        return 0
      }
      return ((parseFloat(v.load1) / v.cpus) * 100)
    },
    getMemoryPercentage(v) {
      return 100 - ((parseFloat(v.memoryAvailable) / v.memoryTotal) * 100)
    },
    getUsageColor(v) {
      if (v < 70) {
        return 'green'
      }
      if (v > 70 && v < 90) {
        return 'yellow'
      }
      return 'red'
    },
    humanSize(bytes) {
      var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
      if (bytes == 0) return '0';
      var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
      return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i];
    },
    getGPUColor(v) {
      let gpu = v.toLowerCase()
      if (gpu.includes('nvidia')) {
        return "green"
      }
      if (gpu.includes('amd')) {
        return "red"
      }
      return "blue"
    },
  },
  computed: {
    headers () {
        return [
          {
            text: 'Name',
            align: 'start',
            sortable: true,
            value: 'name',
          },
          {
            text: 'Version',
            value: 'version',
          },
          { text: 'CPUs', value: 'cpus' },
          { text: 'CPU Usage', value: 'load1' },
          { text: 'Memory', value: 'memory' },
          { text: 'Memory Usage', value: 'memoryUsage' },
          { text: 'GPUs', value: 'gpus' },
          { text: 'Actions', value: 'actions', sortable: false },
        ]
      },
  },
  mounted () {
  },
  created () {
    this.loadWorkers()
  },
  beforeDestroy() {
  },
}
</script>
