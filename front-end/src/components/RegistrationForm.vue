<template>
    <div class="registration-form">
      <h2>Registration</h2>
      <form @submit.prevent="register">
        <label for="email">Email:</label>
        <input type="text" id="email" v-model="formData.email" required />
        <span class="error-message" v-if="emailError">{{ emailError }}</span>
        <label for="password">Password:</label>
        <input type="password" id="password" v-model="formData.password" required />
        <label for="firstName">First Name:</label>
        <input type="text" id="firstName" v-model="formData.firstName" required />
        <label for="lastName">Last Name:</label>
        <input type="text" id="lastName" v-model="formData.lastName" required />
        <label for="phone">Phone:</label>
        <input type="number" id="phone" v-model="formData.phone" required />
        <label for="dob">Date of Birth:</label>
        <input type="date" id="dob" v-model="formData.dateOfBirth" required />
        <p></p>
        <label for="gender">Gender:</label>
        <select name="gender" id="gender" v-model="formData.gender" required >
        <option value="Male">Male</option>
        <option value="Female">Female</option>
        </select>
        <label for="country">Country:</label>
        <input type="text" id="country" v-model="formData.country" required />
        <label for="privacy">Privacy:</label>
        <select name="privacy"  id="privacy" v-model="formData.privacy" required >
        <option value="1">Public</option>
        <option value="0">Private</option>
        </select>
        <label for="avatar">Avatar/Image (Optional):</label>
        <input type="file" id="avatar" @change="handleAvatarUpload" accept="image/jpeg,image/png" />
        <label for="username">Username (Optional):</label>
        <input type="text" id="username" v-model="formData.username" />
        <label for="aboutMe">About Me (Optional):</label>
        <textarea id="aboutMe" v-model="formData.aboutMe"></textarea>
      <button type="submit">Register</button>
      </form>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  import LoginForm from './LoginForm.vue';
  export default {
    data() {
      return {
        formData: {
          email: '',
          password: '',
          firstName: '',
          lastName: '',
          phone: null,
          dateOfBirth: '',
          gender: '',
          country: '',
          privacy: '',
          username: '',
          aboutMe: '',
          avatar: '',
        },
        emailError: '',
      };
    },
    methods: {
      async register() {
        if (!this.isValidEmail(this.formData.email)) {
          this.emailError = 'Please enter a valid email address.';
          return; 
      }
      const dataToSend = new FormData();
          dataToSend.append('email', this.formData.email);
          dataToSend.append('password', this.formData.password);
          dataToSend.append('firstName', this.formData.firstName);
          dataToSend.append('lastName', this.formData.lastName);
          dataToSend.append('phone',this.formData.phone);
          dataToSend.append('dateOfBirth', this.formData.dateOfBirth);
          dataToSend.append('gender', this.formData.gender);
          dataToSend.append('country', this.formData.country);
          dataToSend.append('privacy', this.formData.privacy);
          dataToSend.append('username', this.formData.username);
          dataToSend.append('aboutMe', this.formData.aboutMe);
          if (this.formData.avatar) {
            dataToSend.append('avatar', this.formData.avatar);
          }
      try {
        console.log(dataToSend);
        const response = await axios.post('http://localhost:8080/api/register', dataToSend, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });
        
        if (response.status === 200) {
          console.log('Registration successful:', response.data);
          await LoginForm.methods.login.call(this);
          this.$router.push('/');
        }
      } catch (error) {
        console.error('Registration failed:', error.response.data);
        this.emailError = error.response.data;
      }
      },
      handleAvatarUpload(event) {
        const file = event.target.files[0];
        const allowedTypes = ['image/jpeg', 'image/png'];
        if (file) {
          if (allowedTypes.includes(file.type)) {
           this.formData.avatar = file;
        } else {
          alert('Please upload a JPEG or PNG image.');
          event.target.value = null;
          }
        } 
      },
      isValidEmail(email) {
      const emailPattern = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
      return emailPattern.test(email);
      },
    },
  };
  </script>
  
