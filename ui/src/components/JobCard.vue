<template>
  <v-card
    :loading="loading"
    class="mx-2 mb-4"
    elevation="4"
    min-width="350"
    max-width="350"
  >
    <router-link v-if="job.status !== 'ERROR'" tag="span" class="link" :to="{name: 'job', params: { id: job.id }}">
        <v-img
          height="256"
          v-bind:src="jobImage"
        >
          <template v-slot:placeholder>
            <v-row
              class="fill-height ma-0"
              width="256"
              align="center"
              justify="center"
            >
              <v-progress-circular
                indeterminate
              ></v-progress-circular>
            </v-row>
          </template>
          <v-overlay
              :absolute="true"
              :value="rendering"
              opacity="0.35"
            >
            <v-progress-circular
              indeterminate
              color="grey lighten-5"
            ></v-progress-circular>
          </v-overlay>
        </v-img>
    </router-link>
    <v-img
      v-else
      height="256"
      src="/img/fynca-default.png"
    >
    </v-img>


    <v-card-title class="text-h4">
      <router-link v-if="job.status !== 'ERROR'" class="link" :to="{name: 'job', params: { id: job.id }}">
        {{job.request.name}}
      </router-link>
      <span v-else>{{job.request.name}}</span>
    </v-card-title>

    <v-card-text>
      <v-row
        align="center"
        class="mx-0"
      >
      </v-row>
      <v-list class="transparent ml-n4" dense>
        <v-list-item class="mb-n4">
          <v-list-item-title>Frames</v-list-item-title>
          <v-list-item-subtitle>
            <v-tooltip bottom v-if="frames > 1">
              <template v-slot:activator="{ on, attrs }">
                {{frames}}
                <v-icon
                  small
                  class="ml-2 mt-n1"
                  dark
                  v-bind="attrs"
                  v-on="on"
                  color="primary"
                >
                  mdi-animation-play
                </v-icon>
              </template>
              <span v-if="frames > 1">{{job.request.renderStartFrame}} - {{job.request.renderEndFrame}}</span>
            </v-tooltip>
            <span v-else>{{frames}}</span>
          </v-list-item-subtitle>
        </v-list-item>
        <v-list-item class="mb-n4">
          <v-list-item-title>Resolution</v-list-item-title>
          <v-list-item-subtitle>
            {{job.request.resolutionX}} x {{job.request.resolutionY}}
          </v-list-item-subtitle>
        </v-list-item>
        <v-list-item class="mb-n4">
          <v-list-item-title>Resolution Scale</v-list-item-title>
          <v-list-item-subtitle>
            {{job.request.resolutionScale}}%
          </v-list-item-subtitle>
        </v-list-item>
        <v-list-item class="mb-n4">
          <v-list-item-title>Samples</v-list-item-title>
          <v-list-item-subtitle>
            {{job.request.renderSamples}}
          </v-list-item-subtitle>
        </v-list-item>
        <v-list-item class="mb-n4">
          <v-list-item-title>Engine</v-list-item-title>
          <v-list-item-subtitle>
            {{job.request.renderEngine | formatRenderEngineName | toTitleCase }}
          </v-list-item-subtitle>
        </v-list-item>
        <v-list-item class="mb-n4">
          <v-list-item-title>Device</v-list-item-title>
          <v-list-item-subtitle>
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-icon
                  class="mr-2"
                  v-if="job.request.renderUseGpu"
                  dark
                  v-bind="attrs"
                  v-on="on"
                  color="green"
                >
                  mdi-expansion-card
                </v-icon>
                <v-icon
                  color="blue"
                  class="mr-2"
                  v-else
                  dark
                  v-bind="attrs"
                  v-on="on"
                >
                  mdi-memory
                </v-icon>
              </template>
              <span v-if="job.request.renderUseGpu">GPU</span>
              <span v-else>CPU</span>
            </v-tooltip>
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card-text>

    <v-card-text
      class="text-md-center"
    >
    </v-card-text>

    <v-card-actions>
      <v-chip v-if="job.status === 'RENDERING'"
        class="ma-2"
        color="blue"
        text-color="white"
      >
        {{titleize(job.status)}}
      </v-chip>
      <v-chip v-else-if="job.status !== undefined && job.status === 'ERROR'"
        class="ma-2"
        color="error"
        text-color="white"
      >
        {{titleize(job.status)}}
      </v-chip>
      <v-chip v-else-if="job.status !== undefined && job.status === 'FINISHED'"
        class="ma-2"
        color="green"
        text-color="white"
      >
        {{titleize(job.status)}}
      </v-chip>
      <v-chip v-else
        class="ma-2"
        color="grey"
        text-color="white"
      >
      <span v-if="job.status !== undefined">{{titleize(job.status)}}</span>
      <span v-else>Queued</span>
      </v-chip>

      <v-spacer />

      <v-speed-dial
        v-model="showActions"
      >
        <template v-slot:activator>
          <v-btn
            fab
            x-small
            v-model="showActions"
            color="grey darken-3"
          >
            <v-icon v-if="showActions">
              mdi-close
            </v-icon>
            <v-icon v-else>
              mdi-dots-vertical
            </v-icon>
          </v-btn>
        </template>
        <v-btn
            fab
            small
            color="grey"
            @click="showLogDialog = true"
        >
          <v-icon>mdi-text-box-outline</v-icon>
        </v-btn>
        <v-btn
            fab
            small
            color="primary lighten-1"
        >
          <v-icon>mdi-content-copy</v-icon>
        </v-btn>
        <v-btn
            fab
            small
            color="primary lighten-1"
        >
          <v-icon>mdi-restart</v-icon>
        </v-btn>

        <v-btn
            fab
            small
            color="red lighten-1"
            @click="showDeleteDialog = true"
        >
          <v-icon>mdi-trash-can</v-icon>
        </v-btn>
      </v-speed-dial>

      <v-dialog
        scrollable
        v-model="showLogDialog"
        width="550"
      >
        <v-card>
          <v-card-title class="grey darken-1">
            Job Logs: {{job.request.name}}
          </v-card-title>

          <v-card-text class="mt-5 text-center">
            <div v-html="jobLog"></div>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              @click="showLogDialog = false"
            >
              Close
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <v-dialog
        v-model="showDeleteDialog"
        width="500"
      >
        <v-card>
          <v-card-title class="grey">
            Delete Job
          </v-card-title>

          <v-card-text>
            <div class="text-center">
            Are you sure you want to delete {{job.request.name}}?
            </div>
            <div class="text-center">
              (Job: {{job.id}})
            </div>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="red"
              @click="deleteJob(job.id)"
            >
              Delete
            </v-btn>
            <v-btn
              @click="showDeleteDialog = false"
            >
              Cancel
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

    </v-card-actions>
  </v-card>
</template>

<script>
import axios from "axios";

export default {
  components: {},
  props: ['job'],
  data: () => ({
    loading: false,
    jobImage: '',
    showActions: false,
    showDeleteDialog: false,
    showLogDialog: false,
    cacheKey: '',
    jobStatus: '',
    jobLog: '',
    frames: 0,
    rendering: true,
    color: 'error',
  }),
  methods: {
    titleize(s) {
      return s.charAt(0).toUpperCase() + s.slice(1).toLowerCase();
    },
    deleteJob(id) {
      axios.delete(this.$apiHost + '/api/v1/jobs/' + id)
        .then( resp => {
          this.$parent.loadJobs()
        })
        .catch(err => {
          if (err.response != null) {
            this.text = err.response.data
          } else {
            this.text = err
          }
          this.snackbar = true
        })
        this.showDeleteDialog = false
    },
    loadJob() {
      if (this.job.status === 'FINISHED' || this.job.status === 'ERROR') {
        this.rendering = false
      }
      if (this.job.status === 'ERROR') {
        return
      }
      axios.get(this.$apiHost + '/api/v1/jobs/' + this.job.id + '/latest-render/' + this.job.request.renderStartFrame + '?cache=' + this.cacheKey)
        .then(resp => {
          this.jobImage = resp.data.url
        })
    },
    loadJobLog() {
      if (this.job.status !== 'FINISHED' && this.job.status !== 'ERROR') {
        return
      }
      axios.get(this.$apiHost + '/api/v1/jobs/' + this.job.id + '/log')
        .then(resp => {
          // reset
          this.jobLog = resp.data.replaceAll('\n', '<br />')
        })
        .catch(err => {
          var msg
          if (err.response != null) {
            msg = err.response.data
          } else {
            msg = err
          }
          this.jobLog = 'Error loading job log: ' + msg
        })
    },
  },
  created () {
    this.frames = parseInt(this.job.request.renderEndFrame) - parseInt(this.job.request.renderStartFrame) + 1
    this.loadJob()
  },
  updated() {
    this.loadJob()
  },
  watch: {
    showLogDialog(v) {
      if (v) {
        this.loadJobLog()
      }
    }
  }
}
</script>
