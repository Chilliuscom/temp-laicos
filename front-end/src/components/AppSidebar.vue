<template>
    <aside class="sidebar">
      <router-link to="/">Home</router-link>
      <router-link :to="`/profile/${sessionId}`">Profile</router-link>
    </aside>
  </template> 
  
  <script>
  export default {
    data() {
      return {
        sessionId: '', 
      };
    },
    created() {
      this.sessionId = sessionStorage.getItem('id');
    },
  };
  </script>