<template>
      <router-link to="/" class="sitename"><h2>BookFace</h2></router-link>
      <div class="right-elements">
      <img src="/icons/friends.png" alt="friends" class="icon"/>
      <img src="/icons/chat.png" alt="chat" class="icon"/>
      <img src="/icons/logout.png" alt="logout" class="icon" @click="logout"/>
      </div>
</template>

<script>
import axios from 'axios';
export default {
      methods: {
            async logout() {
            try {
                  const response = await axios.get('http://localhost:8080/api/logout', {
                        withCredentials: true, });
                  if (response.status === 200) {
                  
                  console.log('Successful logout:', response.data);
                  this.$router.push('/login');
                  }
            } catch (error) {
                  console.error('Logout failed:', error.response.data);
            }
            },
      },
}
</script>