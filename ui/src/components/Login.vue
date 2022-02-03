<template>
  <v-app id="inspire">
    <v-snackbar
      v-model="snackbar"
      :color="color"
      :multi-line=true
      :top=true
    >
      {{ text }}
      <v-btn
        dark
        text
        @click="snackbar = false"
      >
        Close
      </v-btn>
    </v-snackbar>

    <v-main>
      <v-container
        class="fill-height"
        fluid
      >
        <v-row
          align="center"
          justify="center"
        >
          <v-col
            cols="12"
            sm="8"
            md="4"
          >
            <v-card class="elevation-12">
              <v-card-text>
                <v-form
                  ref="form"
                  v-model="valid"
                  lazy-validation
                >
                  <v-text-field
                    v-model="username"
                    :rules="usernameRules"
                    label="Username"
                    prepend-icon="person"
                    type="text"
                    required
                  ></v-text-field>

                  <v-text-field
                    id="password"
                    v-model="password"
                    :rules="passwordRules"
                    label="Password"
                    prepend-icon="lock"
                    type="password"
                    required
                    @keydown="submit"
                  ></v-text-field>
                </v-form>
              </v-card-text>
              <v-card-actions>
                <div class="flex-grow-1"></div>
                <v-btn @click='login' color="primary">Login</v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
  import axios from "axios"

  export default {
    props: {},
    data: () => ({
      drawer: null,
      snackbar: false,
      text: '',
      valid: true,
      color: 'error',
      username: '',
      usernameRules: [
        v => !!v || 'Username is required',
      ],
      password: '',
      passwordRules: [
        v => !!v || 'Password is required',
      ],
    }),
    methods: {
      submit(v) {
        if (v.keyCode === 13) {
          this.login()
        }
      },
      login () {
        if (this.$refs.form.validate()) {
          // override config to ensure axios does not try to send an existing token
          let cfg = {
            headers: {
              "x-session-token": '',
            }
          }
          axios.post(this.$apiHost + '/auth/login', {
            username: this.username,
            password: this.password
          }, cfg)
          .then(response => {
            // store token
            var account = response.data.account
            var config = JSON.parse(atob(response.data.config))
            localStorage.username = account.username
            localStorage.token = config.token
            axios.defaults.headers.common['X-Session-Token'] = response.data.token
            window.location.href = '/'
          })
          .catch(err => {
            console.log(err)
            var msg
            if (err.response != null) {
              msg = err.response.data
            } else {
              msg = err
            }
            this.text = msg
            this.snackbar = true
          })
        }
      },
      reset () {
        this.$refs.form.reset()
      },
      resetValidation () {
        this.$refs.form.resetValidation()
      },
    },
  }
</script>
