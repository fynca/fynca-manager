<template>
  <v-row>
    <v-col class="d-none d-md-flex"/>
    <v-col
      sm="12"
      md="6"
      lg="6"
      xl="6"
    >
      <v-overlay
        :value="showUpdatingOverlay"
      >
        <v-progress-circular
          size="80"
          indeterminate
        ></v-progress-circular>
      </v-overlay>

      <template>
        <v-form
          ref="profileForm"
          v-model="profileFormValid"
          lazy-validation
        >
        <v-card
          class="pa-14"
          elevation="5"
          >
          <v-avatar
            size="128"
            rounded
            >
            <v-img
              :src="gravatar"
              lazy-src="/img/account.png"
            >
              <template v-slot:placeholder>
                <v-row
                  class="fill-height ma-0"
                  align="center"
                  justify="center"
                >
                  <v-progress-circular
                    indeterminate
                    color="grey lighten-5"
                  ></v-progress-circular>
                </v-row>
              </template>
            </v-img>
          </v-avatar>
          <v-card-title class="text-h2">{{username}}</v-card-title>
          <v-card-text>
            <v-row>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="firstName"
                  label="First Name"
                  :rules="firstNameRules"
                  required
                ></v-text-field>
              </v-col>

              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="lastName"
                  :value="lastName"
                  :rules="lastNameRules"
                  label="Last Name"
                  required
                ></v-text-field>
              </v-col>
            </v-row>

            <v-text-field
              v-model="email"
              :rules="emailRules"
              label="E-mail"
              required
            ></v-text-field>
          </v-card-text>

          <v-card-actions class="mt-5">
            <v-btn
              @click="showChangePasswordDialog = true"
            >
              Change Password
            </v-btn>

            <v-spacer></v-spacer>

            <v-btn
              color="green"
              @click="updateProfile"
            >
              Update
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-form>

      </template>
    </v-col>

    <v-col class="d-none d-md-flex"/>

    <v-dialog
      scrollable
      v-model="showChangePasswordDialog"
      width="400"
      @keydown.esc="showChangePasswordDialog = false"
    >
      <v-card>
        <v-card-title class="text-h5">
          Change Password
        </v-card-title>

        <v-card-text>
          <v-form
            ref="changePasswordForm"
            v-model="changePasswordFormValid"
            lazy-validation
          >
            <v-text-field
              type="password"
              v-model="currentPassword"
              :rules="currentPasswordRules"
              label="Current Password"
              required
            ></v-text-field>

            <v-text-field
              type="password"
              v-model="password"
              :rules="passwordRules"
              label="Password"
              required
            ></v-text-field>
            <v-text-field
              type="password"
              v-model="passwordConfirm"
              :rules="passwordConfirmRules"
              label="Password (confirm)"
              required
            ></v-text-field>

          </v-form>

          <v-overlay
            :absolute="true"
            :value="showChangePasswordOverlay"
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
            @click="changePassword()"
          >
            Submit
          </v-btn>
          <v-btn
            color="secondary"
            @click="closeChangePasswordDialog"
          >
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script>
import axios from 'axios'
import Vue from 'vue'

export default {
  props: [],
  data: () => ({
    username: '',
    gravatar: '',
    firstName: '',
    firstNameRules: [
      v => !!v || 'First name is required'
    ],
    lastName: '',
    lastNameRules: [
      v => !!v || 'Last name is required'
    ],
    email: '',
    emailRules: [
      v => !!v || 'Email is required'
    ],
    currentPassword: '',
    currentPasswordRules: [
      v => !!v || 'Password is required'
    ],
    password: '',
    passwordRules: [
      v => !!v || 'Password is required'
    ],
    passwordConfirm: '',
    passwordConfirmRules: [
      v => !!v || 'Password is required'
    ],
    profileFormValid: true,
    changePasswordFormValid: true,
    showChangePasswordOverlay: false,
    showChangePasswordDialog: false,
  }),
  methods: {
    changePassword() {
      if (!this.$refs.changePasswordForm.validate()) {
        return
      }
      this.showChangePasswordOverlay = true
      let formData = new FormData()
      formData.append("username", this.username)
      formData.append("password", this.password)
      formData.append("currentPassword", this.currentPassword)
      axios.post(this.$apiHost+'/api/v1/accounts/change-password', formData)
        .then( res => {
          this.showChangePasswordOverlay = false
          this.showChangePasswordDialog = false
          this.resetChangePasswordForm()
	        this.$root.$emit('showInfo', 'Password updated successfully')
      })
      .catch(err => {
        console.log(err)
        var msg
        if (err.response != null) {
          msg = err.response.data
        } else {
          msg = err
        }
        this.showChangePasswordOverlay = false
	      this.$root.$emit('showError', 'Error changing password: ' + msg)
        this.resetChangePasswordForm()
      })
    },
    updateProfile() {
      if (!this.$refs.profileForm.validate()) {
        return
      }
      this.showUpdatingOverlay = true
      let formData = new FormData()
      formData.append("username", this.username)
      formData.append("firstName", this.firstName)
      formData.append("lastName", this.lastName)
      formData.append("email", this.email)
      axios.post(this.$apiHost+'/api/v1/accounts/profile', formData)
        .then( res => {
	        this.$root.$emit('showInfo', 'Profile updated successfully')
          this.loadAccountInfo()
          this.showUpdatingOverlay = false
      })
      .catch(err => {
        console.log(err)
        var msg
        if (err.response != null) {
          msg = err.response.data
        } else {
          msg = err
        }
        this.showUpdatingOverlay = false
	      this.$root.$emit('showError', 'Error updating profile: ' + msg)
      })
    },
    closeChangePasswordDialog() {
      this.resetChangePasswordForm()
      this.showChangePasswordDialog = false
    },
    resetChangePasswordForm() {
	    this.$refs.changePasswordForm.resetValidation()
      this.password = ''
      this.passwordConfirm = ''
      this.showChangePasswordOverlay = false
    },
    loadAccountInfo() {
      this.showUpdatingOverlay = true
      axios.get(this.$apiHost + '/api/v1/accounts/profile')
        .then(resp => {
          var acct = resp.data.account
          this.username = acct.username
          this.firstName = acct.firstName
          this.lastName = acct.lastName
          this.email = acct.email
          if (acct.email !== '') {
            this.gravatar = Vue.filter('gravatar')(acct.email, 512)
          } else {
            this.gravatar = '/img/account.png'
          }
          this.showUpdatingOverlay = false
        })
        .catch(err => {
          var msg
          if (err.response != null) {
            msg = err.response.data
          } else {
            msg = err
          }
          this.showUpdatingOverlay = false
          console.log(err)
        });
    }
  },
  created() {
    this.loadAccountInfo()
  }
}
</script>
