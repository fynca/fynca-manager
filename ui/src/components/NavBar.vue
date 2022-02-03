<template>
  <div>
    <v-snackbar
      dark
      :top="true"
      v-model="showInfo"
      :timeout="5000"
      content-class="text-center"
    >
      {{ infoText }}
    </v-snackbar>

    <v-snackbar
      dark
      :top="true"
      v-model="showError"
      color="red lighten-3"
      :timeout="30000"
      content-class="text-center"
    >
      <span>{{ errorText }}</span>
      <template v-slot:action="{ attrs }">
        <v-btn
          color="black lighten-1"
          text
          v-bind="attrs"
          @click="showError = false"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </template>
    </v-snackbar>

    <v-app-bar
        app
        dense
        elevation="1"
    >
    <v-toolbar-title><router-link tag="span" class="link" :to="{name: 'main'}">Fynca</router-link></v-toolbar-title>
    <v-btn
      small
      outlined
      class="ml-8"
      color="grey"
      to="/jobs"
    >
      <span class="white--text">Jobs</span>
    </v-btn>

    <v-spacer></v-spacer>
    <template v-if="appConfig.authenticator !== 'none' && userConfig.account.admin">
    <v-menu
      bottom
      >
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          dark
          small
          class="ml-1"
          v-bind="attrs"
          v-on="on"
        >
        <span class="d-none d-sm-flex">System</span>
        <v-icon>arrow_drop_down</v-icon> <v-icon>mdi-cog</v-icon>
        </v-btn>
      </template>

      <!-- TODO: admin menu -->
      <v-list dense>
        <v-list-item
          dense
          to="/workers"
        >
          <v-list-item-title>Workers</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
    </template>


    <v-dialog
      scrollable
      id="queueRenderDialog"
      v-model="showQueueRenderDialog"
      width="700"
      @keydown.esc="showQueueRenderDialog = false"
    >
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          class="ml-1"
          color="primary"
          v-bind="attrs"
          v-on="on"
          small
        >
          <span class="d-none d-sm-flex">Render</span>
          <v-icon
              right
              dark
              >mdi-plus-circle-outline
          </v-icon>
        </v-btn>
      </template>

      <v-card>
        <v-card-title class="text-h5">
          Queue Render Job
        </v-card-title>

        <v-card-text>
          <v-form
            ref="queueJobForm"
            v-model="queueFormValid"
            lazy-validation
          >
            <v-text-field
              v-model="name"
              :counter="128"
              :rules="nameRules"
              label="Name"
              hint="Name of project to render"
              required
            ></v-text-field>

            <v-row>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="renderStartFrame"
                  label="Start Frame"
                  :rules="renderStartFrameRules"
                  required
                ></v-text-field>
              </v-col>

              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="renderEndFrame"
                  :rules="renderEndFrameRules"
                  label="End Frame"
                  hint="Match start frame for single frame render"
                ></v-text-field>
              </v-col>
            </v-row>

            <v-row>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="renderResolutionX"
                  :rules="renderResolutionXRules"
                  label="Resolution (X)"
                  required
                ></v-text-field>
              </v-col>

              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="renderResolutionY"
                  :rules="renderResolutionXRules"
                  label="Resolution (Y)"
                  required
                ></v-text-field>
              </v-col>
            </v-row>

            <v-slider
              v-model="renderSamples"
              class="my-6"
              label="Samples"
              :thumb-size="32"
              thumb-label="always"
              hint="Number of samples at which to render project"
              max="2048"
              step="10"
              min="10"
            ></v-slider>

            <v-slider
              v-model="resolutionScale"
              class="my-6"
              label="Scale"
              :thumb-size="28"
              thumb-label="always"
              hint="Scale in percentage of resolution to render project"
              max="100"
              min="1"
            ></v-slider>

            <v-checkbox
              v-model="renderUseGPU"
              label="Use GPU"
              hint="Enable GPU (if available) for rendering"
            ></v-checkbox>

            <v-slider
              v-model="renderSlices"
              class="my-6"
              label="Slices"
              :thumb-size="28"
              thumb-label="always"
              max="20"
              :rules="renderSlicesRules"
              min="0"
              step="1"
              hint="Split job among workers for sliced rendering (zero to disable)"
            ></v-slider>

            <v-file-input
              accept="application/x-blender,*.blend,*.zip"
              chips
              show-size
              truncate-length="24"
              ref="formFile"
              v-model="projectFile"
              :rules="projectFileRules"
              required
            ></v-file-input>

          </v-form>

          <v-overlay
            :absolute="true"
            :value="showUploadingOverlay"
          >
            <v-progress-circular
              size="80"
              indeterminate
              color="primary"
            ></v-progress-circular>
          </v-overlay>
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            @click="queueRenderJob()"
          >
            Submit
          </v-btn>
          <v-btn
            color="secondary"
            @click="showQueueRenderDialog = false"
          >
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <template v-if="appConfig.authenticator !== 'none'">
    <v-menu
      bottom
      >
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          dark
          small
          class="ml-1"
          v-bind="attrs"
          v-on="on"
        >
        <span v-model="username" class="d-none d-sm-flex">{{username}}</span>
        <v-icon>arrow_drop_down</v-icon>
        <v-icon v-if="userConfig.account.email === undefined || userConfig.account.email === ''">account_circle</v-icon>
        <v-avatar
          v-else
          size="24"
          >
          <img
            :src="userConfig.account.email | gravatar"
          >
        </v-avatar>
        </v-btn>
      </template>

      <!-- TODO: admin menu -->
      <v-list dense>
        <v-list-item
          dense
          to="/profile"
        >
          <v-list-item-title>Profile</v-list-item-title>
        </v-list-item>

        <v-divider />

        <v-list-item
          dense
          @click="logout()"
        >
          <v-list-item-title>Logout</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
    </template>

    </v-app-bar>
  </div>
</template>

<script>
import axios from "axios"
import md5 from "md5"

export default {
  props: ['appConfig', 'userConfig'],
  data: () => ({
    name: '',
    nameRules: [
      v => !!v || 'Name is required'
    ],
    renderSamples: 100,
    resolutionScale: 100,
    renderStartFrame: 1,
    renderStartFrameRules: [
      v => !!v || 'Start frame is required'
    ],
    renderEndFrame: 1,
    renderEndFrameRules: [
      v => !!v || 'End frame is required'
    ],
    renderResolutionX: 1920,
    renderResolutionXRules: [
      v => !!v || 'Resolution X is required'
    ],
    renderResolutionY: 1080,
    renderResolutionYRules: [
      v => !!v || 'Resolution Y is required'
    ],
    renderUseGPU: false,
    renderSlices: 0,
    renderSlicesRules: [
      v => (v == 0 || v > 1) || 'Render slices must be greater than 1'
    ],
    showInfo: false,
    infoText: '',
    showError: false,
    errorText: '',
    queueFormValid: true,
    showUploadingOverlay: false,
    projectFile: null,
    projectFileRules: [
      v => !!v || 'Project file is required'
    ],
    showQueueRenderDialog: false,
    color: 'error',
  }),
  computed: {
    currentPage() {
      return this.$route.name
    },
    username() {
      return localStorage.username
    },
  },
  created() {
  },
  mounted () {
    this.$root.$on('showInfo', (txt) => {
      this.infoText = txt
      this.showInfo = true
    })
    this.$root.$on('showError', (txt) => {
      this.errorText = txt
      this.showError = true
    })
  },
  methods: {
    resetQueueJobForm() {
	    this.$refs.queueJobForm.resetValidation()
      this.name = ''
      this.renderSamples = 100
      this.resolutionScale = 100
      this.renderStartFrame = 1
      this.renderEndFrame = 1
      this.renderResolutionX = 1920
      this.renderResolutionY = 1080
      this.renderUseGPU = false
      this.renderSlices = 0
      this.projectFile = null
      this.showUploadingOverlay = false
    },
    queueRenderJob() {
      if (!this.$refs.queueJobForm.validate()) {
        return
      }
      this.showUploadingOverlay = true
      // temporarily stop refresh
	    this.$root.$emit('cancelLoadJobs')
      let formData = new FormData()
      formData.append("name", this.name)
      formData.append("renderResolutionX", this.renderResolutionX)
      formData.append("renderResolutionY", this.renderResolutionY)
      formData.append("renderSamples", this.renderSamples)
      formData.append("resolutionScale", this.resolutionScale)
      formData.append("renderStartFrame", this.renderStartFrame)
      formData.append("renderEndFrame", this.renderEndFrame)
      formData.append("renderUseGPU", this.renderUseGPU)
      formData.append("renderSlices", this.renderSlices)
      formData.append("project", this.projectFile)
	    this.$root.$emit('showInfo', "Job uploading. Please wait...")
      axios.post(this.$apiHost+'/api/v1/jobs', formData)
        .then( res => {
          this.showQueueRenderDialog = false
          // refresh jobs
	        this.$root.$emit('enableLoadJobs')
          this.resetQueueJobForm()
      })
      .catch(err => {
        console.log(err)
	      this.$root.$emit('showError', 'Error uploading job: ' + err)
	      this.$root.$emit('enableLoadJobs')
        this.resetQueueJobForm()
      })
    },
    logout() {
      axios.get(this.$apiHost+'/auth/logout')
      localStorage.removeItem('token')
      window.location.href = '/'
    }
  }
}
</script>
