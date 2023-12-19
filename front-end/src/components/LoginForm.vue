<template>
    <div class="login-form">
      <h2>Login</h2>
      <form @submit.prevent="login">
        <div class="form-group">
          <label for="email">Email:</label>
          <input type="text" id="email" v-model="formData.email" required />
        </div>
        <div class="form-group">
          <label for="password">Password:</label>
          <input type="password" id="password" v-model="formData.password" required />
        </div>
        <span class="error-message" v-if="loginError">{{ loginError }}</span>
        <p>Don't have an account? <router-link to="/register">Register here!</router-link></p>
        <button type="submit">Login</button>
      </form>
    </div>
  </template>

<script>
import axios from 'axios';
export default {
  data() {
    return {
      formData: {
        email: '',
        password: '',
      },
      loginError: '',
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('http://localhost:8080/api/login', this.formData, {withCredentials: true});
        if (response.status === 200) {
          console.log('Successful login:', response.data.message);
          const id = response.data.id;
          sessionStorage.setItem('id', id);
          this.$router.push('/');
        }
      } catch (error) {
        console.error('Login failed:', error.response.data);
        this.loginError = error.response.data;
        }
    },
  },
};
</script>