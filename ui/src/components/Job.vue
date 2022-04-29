<template>
  <v-row class="pb-12">
    <v-col md="9">
        <v-card
          elevation="5"
        >
          <v-progress-linear
            v-if="rendering"
            stream
            :buffer-value="renderProgress+10"
            :value="renderProgress"
            color="primary"
          ></v-progress-linear>
          <v-card-text>
            <v-img
              v-bind:src="renderFrameSrc"
              height="650"
              contain
              max-height="650"
            >
            </v-img>

            <div class="text-center mt-2">
              <v-row justify="center">
                <v-col align="right">
                  <v-btn
                    @click="framePrev100"
                    v-if="frameJobs.length > 100"
                    class="mt-3"
                  >
                    <v-icon>mdi-rewind</v-icon>
                  </v-btn>

                  <v-btn
                    @click="framePrev10"
                    v-if="frameJobs.length > 10"
                    class="mt-3 ml-1"
                  >
                    <v-icon>mdi-rewind-10</v-icon>
                  </v-btn>

                  <v-btn
                    @click="framePrev"
                    v-if="frameJobs.length > 1"
                    class="mt-3 ml-1 mr-1"
                  >
                    <v-icon>mdi-skip-previous</v-icon>
                  </v-btn>
                </v-col>
                <v-col cols="1">
                <v-text-field
                  hide-details
                  :rules="frameRules"
                  label="Frame"
                  v-model="directFrame"
                  class="centered-input"
                  single-line
                  @change="updateDirectFrame"
                ></v-text-field>
                </v-col>
                <v-col align="left">
                  <v-btn
                    @click="frameNext"
                    v-if="frameJobs.length > 1"
                    class="mt-3 ml-1 mr-1"
                  >
                    <v-icon>mdi-skip-next</v-icon>
                  </v-btn>

                  <v-btn
                    @click="frameNext10"
                    v-if="frameJobs.length > 10"
                    class="mt-3"
                  >
                    <v-icon>mdi-fast-forward-10</v-icon>
                  </v-btn>

                  <v-btn
                    @click="frameNext100"
                    v-if="frameJobs.length > 100"
                    class="mt-3 ml-1"
                  >
                    <v-icon>mdi-fast-forward</v-icon>
                  </v-btn>
                </v-col>
              </v-row>
            </div>
            <v-fab-transition>
              <v-btn
                v-show="!this.rendering"
                @click="getCurrentFrameRenderSrc()"
                color="green"
                dark
                absolute
                bottom
                right
                fab
                target="_blank"
              >
                <v-icon>mdi-download</v-icon>
              </v-btn>
            </v-fab-transition>
          </v-card-text>
        </v-card>
    </v-col>

    <v-col md="3">
      <v-card
        elevation="5"
        >
        <v-list-item class="mb-3">
          <v-list-item-content>
            <v-list-item-title class="text-h5">
              {{jobTitle}}
              <v-chip
                v-if="job.status == 'RENDERING'"
                class="ma-1 mb-1 float-right"
                color="blue"
                text-color="white"
              >
                {{titleize(job.status)}}
              </v-chip>
              <v-chip
                v-else-if="job.status !== undefined && job.status === 'ERROR'"
                class="ma-1 mb-1 float-right"
                color="red"
                text-color="white"
              >
                {{titleize(job.status)}}
              </v-chip>
              <v-chip
                v-else-if="job.status === 'QUEUED'"
                class="ma-1 mb-1 float-right"
                color="grey"
                text-color="white"
              >
                {{titleize(job.status)}}
              </v-chip>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list class="transparent">
          <v-list-item>
            <v-list-item-title>Engine</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <span>{{renderEngine | formatRenderEngineName }}</span>
            </v-list-item-subtitle>
          </v-list-item>

          <v-list-item>
            <v-list-item-title>Device</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <v-icon
                    class="mr-2"
                    v-if="renderUseGPU"
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
                <span v-if="renderUseGPU">GPU</span>
                <span v-else>CPU</span>
              </v-tooltip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Frames</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                small
                class="ma-2"
                color="grey darken-1"
                text-color="white"
              >
                {{this.frameJobs.length}}
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Samples</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                small
                class="ma-2"
                color="grey darken-1"
                text-color="white"
              >
                {{renderSamples}}
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Resolution</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                small
                class="ma-2"
                color="grey darken-1"
                text-color="white"
              >
                {{resolutionX}} x {{resolutionY}}
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Resolution Scale</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                small
                class="ma-2"
                color="grey darken-1"
                text-color="white"
              >
                {{resolutionScale}}%
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Render Slices</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                small
                class="ma-2"
                color="grey darken-1"
                text-color="white"
              >
                {{renderSlices}}
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Created</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                small
                class="ma-2"
                color="grey darken-1"
                text-color="white"
              >
                {{jobCreated}}
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Queue Time</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                v-if="job.status === 'FINISHED'"
                small
                class="ma-2"
		color="grey darken-1"
                text-color="white"
              >
                {{queueTime}}
              </v-chip>
              <v-chip
                v-else
                small
                class="ma-2"
                color="primary"
                text-color="white"
              >
                In Progress
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Render Time</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                v-if="job.status === 'FINISHED'"
                small
                class="ma-2"
                :color="job.status === 'ERROR' ? 'red': 'green'"
                text-color="white"
              >
                {{jobRenderTime}}
              </v-chip>
              <v-chip
                v-else
                small
                class="ma-2"
                color="primary"
                text-color="white"
              >
                In Progress
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <v-list-item-title>Job Time</v-list-item-title>
            <v-list-item-subtitle class="text-right">
              <v-chip
                v-if="job.status === 'FINISHED'"
                small
                class="ma-2"
                :color="job.status === 'ERROR' ? 'red': 'green'"
                text-color="white"
              >
                {{jobTime}}
              </v-chip>
              <v-chip
                v-else
                small
                class="ma-2"
                color="primary"
                text-color="white"
              >
                In Progress
              </v-chip>
            </v-list-item-subtitle>
          </v-list-item>
        </v-list>

        <v-card-actions>
          <v-spacer />
            <template v-if="job.status === 'FINISHED' || job.status === 'ERROR'">
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                      dark
                      color="primary lighten-1"
                      class="ml-1 mr-1"
                    >
                    <v-icon
                      v-bind="attrs"
                      v-on="on"
                      dark
                      >
                        mdi-cog-play-outline
                    </v-icon>
                  </v-btn>
                </template>
                <span>Job Actions</span>
              </v-tooltip>
            </template>

            <template v-if="job.status === 'FINISHED' || job.status === 'ERROR'">
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                      dark
                      color="green lighten-2"
                      class="ml-1 mr-1"
                      @click="startGetJobArchive"
                    >
                    <v-icon
                      v-bind="attrs"
                      v-on="on"
                      dark
                      >
                        mdi-folder-zip
                    </v-icon>
                  </v-btn>
                </template>
                <span>Download Job Archive</span>
              </v-tooltip>
            </template>

            <template v-if="job.status === 'FINISHED' || job.status === 'ERROR'">
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                      dark
                      color="grey"
                      class="ml-1 mr-1"
                      @click="showRenderLogDialog = true"
                    >
                    <v-icon
                      v-bind="attrs"
                      v-on="on"
                      dark
                      >
                        mdi-text-box-outline
                    </v-icon>
                  </v-btn>
                </template>
                <span>View Render Logs</span>
              </v-tooltip>
            </template>

            <template v-if="job.status === 'FINISHED' || job.status === 'ERROR'">
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                      dark
                      color="red"
                      class="ml-1 mr-1"
                      @click="showDeleteJobDialog = true"
                    >
                    <v-icon
                      v-bind="attrs"
                      v-on="on"
                      dark
                      >
                        mdi-trash-can
                    </v-icon>
                  </v-btn>
                </template>
                <span>Delete Job</span>
              </v-tooltip>
            </template>
        </v-card-actions>
      </v-card>

    </v-col>

    <v-dialog
      scrollable
      v-model="showRenderLogDialog"
      width="75%"
    >
      <v-card>
        <v-card-title class="grey darken-1">
          Log: {{jobTitle}} (Frame: {{this.currentFrame}})
        </v-card-title>

        <v-card-text class="mt-5">
          <v-progress-linear
            v-if="renderLogLoading"
            indeterminate
            color="primary darken-2"
          ></v-progress-linear>

          <div v-if="renderSliceLogs.length === 0" v-html="renderLog"></div>
          <template v-else>
            <v-tabs>
              <v-tab v-for="(log, i) in renderSliceLogs" :key="i">Slice {{log.slice}}</v-tab>
              <v-tab-item v-for="(log, i) in renderSliceLogs" :key="i">
                <v-card flat>
                  <v-card-text v-html="formatRenderLog(log.log)"></v-card-text>
                </v-card>
              </v-tab-item>
            </v-tabs>
          </template>
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            @click="showRenderLogDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog
      v-model="showDeleteJobDialog"
      width="500"
    >
      <v-card>
        <v-card-title class="grey darken-1">
          Delete {{jobTitle}}
        </v-card-title>

        <v-card-text class="mt-5">
            <div class="text-center">
            Are you sure you want to delete {{jobTitle}}?
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
            @click="showDeleteJobDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </v-row>
</template>

<style scoped>
  .centered-input >>> input {
    text-align: center
  }
</style>

<script>
import axios from "axios";

export default {
  props: ['id'],
  components: {},
  data: () => ({
    job: {},
    loading: true,
    jobTitle: '',
    renderEngine: '',
    renderFrame: '',
    renderStartFrame: -1,
    renderEndFrame: -1,
    renderSamples: '',
    resolutionX: '',
    resolutionY: '',
    resolutionScale: '',
    renderSlices: null,
    renderFrameSrc: '',
    queueTime: '',
    jobTime: '',
    renderTime: '',
    currentFrame: 1,
    carouselFrame: 0,
    directFrame: 1,
    snackbar: false,
    rendering: true,
    renderUseGPU: false,
    text: '',
    jobCreated: '',
    timer: '',
    frameJobs: [],
    frameRenderUrls: new Map(),
    renderProgress: 0,
    renderLog: '',
    renderLogLoading: false,
    renderSliceLogs: [],
    jobArchivePending: false,
    jobArchiveTimer: '',
    showRenderLogDialog: false,
    showDeleteJobDialog: false,
    frameRules: [
      v => (v != '') || 'Frame cannot be empty'
    ],
    maxFrames: 15,
    currentLoadFrameOffset: 0,
    color: 'error',
  }),
  methods: {
    humanTime(ms) {
      const delim = ''
      const showWith0 = value => (value < 10 ? `0${value}` : value)
      const hours = showWith0(Math.floor((ms / (1000 * 60 * 60)) % 60))
      const minutes = showWith0(Math.floor((ms / (1000 * 60)) % 60))
      const seconds = showWith0(Math.floor((ms / 1000) % 60))
      return `${parseInt(hours) ? `${hours}h${delim}` : ""}${minutes}m${delim}${seconds}s`
    },
    titleize(s) {
      return s.charAt(0).toUpperCase() + s.slice(1).toLowerCase();
    },
    loadJob() {
      axios.get(this.$apiHost + '/api/v1/jobs/'+this.id)
        .then(response => {
          this.job = response.data.job
          this.jobTitle = this.job.request.name
          this.renderSamples = this.job.request.renderSamples
          this.renderEngine = this.job.request.renderEngine
          this.renderStartFrame = this.job.request.renderStartFrame
          this.renderEndFrame = this.job.request.renderEndFrame
          this.renderUseGPU = this.job.request.renderUseGpu
          this.renderSlices = this.job.request.renderSlices
          this.resolutionX = this.job.request.resolutionX
          this.resolutionY = this.job.request.resolutionY
          this.resolutionScale = this.job.request.resolutionScale
          var created = new Date(this.job.createdAt)
          var started = new Date(this.job.startedAt)
          var finished = new Date(this.job.finishedAt)
          this.jobCreated = created.toLocaleString('en-US').replace(',', '')
          this.queueTime = this.humanTime(started - created)
          this.jobTime = this.humanTime(finished - started)
          this.loading = false
          this.cacheKey = this.job.finishedAt
          this.frameJobs = this.job.frameJobs
          if (this.job.status === 'FINISHED') {
            this.rendering = false
            this.cancelLoadJob()
          }
          // convert render time seconds
          if (this.job.duration !== null) {
            var renderTimeMs = parseInt(this.job.duration.replace('s', '')) * 1000
            this.jobRenderTime = this.humanTime(renderTimeMs)
          }
          this.currentFrame = parseInt(this.job.request.renderStartFrame)
          if (this.frameJobs.length > 100) {
            this.maxFrames = 5
          }
          this.directFrame = this.currentFrame
          // render progress
          if (this.job.frameJobs.length > 1) {
            var x
            var renderedFrames = 0
            for (x=0; x<this.job.frameJobs.length; x++) {
              if (this.job.frameJobs[x].status === 'FINISHED') {
                renderedFrames += 1
              }
            }
            this.renderProgress = (renderedFrames / this.job.frameJobs.length * 100)
          } else { // show slice progress
            var y
            var renderedSlices = 0
            var frameJob = this.job.frameJobs[0]
            for (y=0; y<frameJob.sliceJobs.length; y++) {
              if (frameJob.sliceJobs[y].status === 'FINISHED') {
                renderedSlices += 1
              }
            }
            this.renderProgress = (renderedSlices / frameJob.sliceJobs.length * 100)
          }
          this.loadFrameJobRenders(this.currentFrame)
        })
        .catch(err => {
          var msg
          if (err.response != null) {
            msg = err.response.data
          } else {
            msg = err
          }
          this.loading = false
          this.$root.$emit('showError', 'Error loading job: ' + msg)
        })
    },
    getCurrentFrameRenderSrc() {
      axios.get(this.$apiHost + '/api/v1/jobs/' + this.job.id + '/latest-render/' + this.currentFrame + '?ttl=3600s')
        .then(resp => {
          window.open(resp.data.url)
        })
    },
    getFrameJobRenderSrc(frame) {
      return this.frameRenderUrls.get(frame)
    },
    loadFrameJobRenders(v) {
      var targetFrame = parseInt(v)
      this.loading = true
      axios.get(this.$apiHost + '/api/v1/jobs/' + this.job.id + '/latest-render/'+ targetFrame + '?ttl=3600s')
      .then(resp => {
        this.renderFrameSrc = resp.data.url
      })
      .catch(err => {
        var msg
        if (err.response != null) {
          msg = err.response.data
        } else {
          msg = err
        }
        this.$root.$emit('showError', 'Error getting job latest render: ' + msg)
        this.loading = false
      })
    },
    loadRenderLog() {
      if (this.job.status !== 'FINISHED' && this.job.status !== 'ERROR') {
        return
      }
      this.renderLog = ''
      this.renderLogLoading = true
      axios.get(this.$apiHost + '/api/v1/renders/' + this.job.id + '/logs/' + this.currentFrame)
        .then(resp => {
          // reset
          this.renderLog = resp.data.renderLog.log.replaceAll('\n', '<br />')
          this.renderLogLoading = false
        })
        .catch(err => {
          var msg
          if (err.response != null) {
            msg = err.response.data
          } else {
            msg = err
          }
          this.renderLog = 'Error loading job log: ' + msg
          this.renderLogLoading = false
        })
    },
    formatRenderLog(v) {
      return v.replaceAll('\n', '<br />')
    },
    loadRenderSliceLogs() {
      if (this.job.status !== 'FINISHED' && this.job.status !== 'ERROR') {
        return
      }
      var idx = parseInt(this.currentFrame) - parseInt(this.job.request.renderStartFrame)
      if (this.job.frameJobs[idx].sliceJobs.length == 0) {
        return
      }
      var x
      let promises = []
      let renderSliceLogs = []
      for (x=0; x<this.job.frameJobs[idx].sliceJobs.length; x++) {
        var slice = x
        // calculate log frame by render start frame offset
        var logFrame = this.currentFrame
        promises.push(
          axios.get(this.$apiHost + '/api/v1/renders/' + this.job.id + '/logs/' + logFrame + '?slice=' + slice)
            .then(resp => {
              renderSliceLogs.push(resp.data.renderLog)
            })
            .catch(err => {
              var msg
              if (err.response != null) {
                msg = err.response.data
              } else {
                msg = err
              }
              this.renderSliceLogs[x] = 'Error loading job slice log: ' + msg
          }))
      }

      Promise.all(promises).then(() =>
        renderSliceLogs.sort((a,b) => (a.slice > b.slice) ? 1 : ((b.slice < a.slice) ? -1 : 0)),
      )
      this.renderSliceLogs = renderSliceLogs
    },
    deleteJob(id) {
      axios.delete(this.$apiHost + '/api/v1/jobs/' + id)
        // eslint-disable-next-line no-unused-vars
        .then( resp => {
          this.$router.push('/jobs')
        })
        .catch(err => {
          if (err.response != null) {
            this.text = err.response.data
          } else {
            this.text = err
          }
          this.snackbar = true
        })
        this.showDeleteJobDialog = false
    },
    updateCurrentFrame(v) {
      this.currentFrame = parseInt(v)
      this.directFrame = this.currentFrame
      this.currentFrameOffset = this.currentFrame
      this.loadFrameJobRenders(this.currentFrame)
      this.currentFrameUrl = this.getFrameJobRenderSrc(this.currentFrame)
    },
    updateCarouselFrame(v) {
      this.carouselFrame = v
      var offset = parseInt(this.job.request.renderStartFrame)
      var currentFrame = parseInt(v) + offset
      this.updateCurrentFrame(currentFrame)
      this.loadFrameJobRenders(currentFrame)
    },
    updateDirectFrame(v) {
      if (v < parseInt(this.job.request.renderStartFrame)) {
        v = parseInt(this.job.request.renderStartFrame)
      }
      if (v > parseInt(this.job.request.renderEndFrame)) {
        v = parseInt(this.job.request.renderEndFrame)
      }
      var offset = parseInt(this.job.request.renderStartFrame)
      this.updateCarouselFrame(v - offset)
      this.updateCurrentFrame(v)
      this.currentFrameUrl = this.getFrameJobRenderSrc(this.currentFrame)
    },
    framePrev() {
      var frame = (parseInt(this.currentFrame) - 1)
      if (frame < parseInt(this.job.request.renderStartFrame)) {
        frame = this.job.request.renderStartFrame
      }
      this.updateCurrentFrame(frame)
      this.loadFrameJobRenders(frame)
    },
    framePrev10() {
      var frame = (parseInt(this.currentFrame) - 10)
      if (frame < parseInt(this.job.request.renderStartFrame)) {
        frame = this.job.request.renderStartFrame
      }
      this.updateCurrentFrame(frame)
      this.loadFrameJobRenders(frame)
    },
    framePrev100() {
      var frame = (parseInt(this.currentFrame) - 100)
      if (frame < parseInt(this.job.request.renderStartFrame)) {
        frame = this.job.request.renderStartFrame
      }
      this.updateCurrentFrame(frame)
      this.loadFrameJobRenders(frame)
    },
    frameNext() {
      var frame = (parseInt(this.currentFrame) + 1)
      if (frame > parseInt(this.job.request.renderEndFrame)) {
        frame = this.job.request.renderEndFrame
      }
      this.updateCurrentFrame(frame)
    },
    frameNext10() {
      var frame = (parseInt(this.currentFrame) + 10)
      if (frame > parseInt(this.job.request.renderEndFrame)) {
        frame = this.job.request.renderEndFrame
      }
      this.updateCurrentFrame(frame)
      this.loadFrameJobRenders(frame)
    },
    frameNext100() {
      var frame = (parseInt(this.currentFrame) + 100)
      if (frame > parseInt(this.job.request.renderEndFrame)) {
        frame = this.job.request.renderEndFrame
      }
      this.updateCurrentFrame(frame)
      this.loadFrameJobRenders(frame)
    },
    startGetJobArchive() {
      this.getJobArchive()
      this.loading = true
      this.jobArchiveTimer = setInterval(this.getJobArchive, 5000)
    },
    getJobArchive() {
      if (this.job.status !== 'FINISHED') {
        return
      }
      this.loading = true
      axios.get(this.$apiHost + '/api/v1/jobs/' + this.job.id + '/archive')
        .then( resp => {
          var url = resp.data
          if (url === "") {
            if (!this.jobArchivePending) {
              this.$root.$emit('showInfo', 'Generating job archive. This can take a few minutes depending on project size. You can refresh and the archive will be generated in the background.')
              this.jobArchivePending = true
            }
            return
          }
          clearInterval(this.jobArchiveTimer)
          this.loading = false
          this.jobArchivePending = false
          window.location.href = url
        })
        .catch(err => {
          this.loading = false
          var msg
          if (err.response != null) {
            msg = err.response.data
          } else {
            msg = err
          }
          this.$root.$emit('showError', 'Error getting job archive: ' + msg)
        })
    },
    cancelLoadJob() {
      clearInterval(this.timer)
    }
  },
  created () {
    this.loadJob()
    this.timer = setInterval(this.loadJob, 10000)
  },
  beforeDestroy() {
    this.cancelLoadJob()
  },
  watch: {
    showRenderLogDialog(v) {
      if (v) {
        this.loadRenderLog()
        this.loadRenderSliceLogs()
      }
    },
  }
}
</script>
