<template>
    <div v-if="userData" class="profile-container">
      <div v-if="isPublicProfile">
        <div class="profile-header">
          <img class="avatar" src="userData.avatar" alt="User Avatar" /> 
          <h1 class="username">{{ userData.username }}</h1>
          <div class="user-details">
          <p class="full-name">{{ userData.firstname }} {{ userData.lastname }}</p>
          <p>{{ userData.gender }} | {{ userData.birthday }} | {{ userData.country }}</p>
          <p class="about">{{ userData.aboutme }}</p>
          </div>
          <p>Followers: <span class="clickable" @click="showModal('followers')">{{ followers ? followers.length : 0 }}</span></p>
          <p>Following: <span class="clickable" @click="showModal('following')">{{ following ? following.length : 0 }}</span></p>
          </div>
          <div v-if="modals.followers" class="modal"> 
            <div class="modal-content">
              <span class="close" @click="closeModal('followers')">&times;</span>
              <h3>Followers</h3>
              <ul>
              <li v-for="follower in this.followers" :key="follower.id">
                <router-link :to="'/profile/' + follower.first_user">
                  <span>{{ follower.firstname }} {{ follower.lastname }}</span>
              </router-link>
              </li>
            </ul>
            </div>
          </div>
          <div v-if="modals.following" class="modal">
            <div class="modal-content">
              <span class="close" @click="closeModal('following')">&times;</span>
              <h3>Following</h3>
              <ul>
                <li v-for="following in this.following" :key="following.id">
                  <router-link :to="'/profile/' + following.second_user">
                    <span>{{ following.firstname }} {{ following.lastname }}</span>
                  </router-link>
                </li>
              </ul>
            </div>
          </div>
          <div class="user-posts">
            <div class="post">
              <p>Post content</p>
            </div>

          </div>
        </div>
      <div v-else>
        <div class="profile-header">
          <img class="avatar" src="userData.avatar" alt="User Avatar" /> 
          <h1 class="username">{{ userData.username }}</h1>
          <div class="user-details">
          <p class="full-name">{{ userData.firstname }} {{ userData.lastname }}</p>
          </div>
        </div>
        <h2>This profile is private</h2>
        <button @click="sendFollowRequest(profileId)">Send Follow Request</button>
        <span class="error-message" v-if="followError">{{ followError }}</span>

        
      </div>
    </div>
    <div v-else>
      <h2>Profile doesn't exist</h2>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  export default {
    data() {
      return {
        userData: null,
        isPublicProfile: null,
        profileId: null,
        followError: '',
        following: null,
        followers: null,
        modals: {
          followers: false,
          following: false
        },
      };
    },
    mounted() {
      const id = this.$route.params.id;
      this.profileId = id;
      this.fetchUserData(id);
      //this.fetchUserAvatar(userData.avatar)
    },
    //BROKEN
    beforeRouteLeave(to, from, next) {
    // Check if modal is open, then close it before leaving the route
    if (this.modals.followers || this.modals.following) {
/*       this.closeModal('followers');
      this.closeModal('following');  
      this.modals.followers = false;
      this.modals.following = false;*/
    }
    next();
    },
    methods: {
      async fetchUserData(id) {
        try {
          const response = await fetch(`http://localhost:8080/api/profile?id=${id}`, {
            method: 'GET',
            credentials: 'include',
          });
          const data = await response.json();
          console.log(data);
          //if (data.message == "NOT_FOLLOWING") {
            if (Object.keys(data).length == 2) {
            this.isPublicProfile = false;
            this.userData = data.profile;
          } else if (data.message == "USER DOES NOT EXIST") {
            //change this
            //this.$router.push('/error-page');         
          } else {
          this.userData = data.profile;
          this.following = data.following;
          this.followers = data.followers;
          if (data.profile.privacy == 1) {
            this.isPublicProfile = true;
          } else {
            this.isPublicProfile = false;
          }
          await this.fetchUserAvatar(this.userData.avatar);
        }
        
        } catch (error) {
          console.error('Error fetching user data:', error);
        }
      },
      //BROKEN
      async fetchUserAvatar(avatarDir) {
        try {
          const avatarUrl = `http://localhost:8080/api/images/${avatarDir}`; 
          console.log(avatarUrl);
          const response = await fetch(avatarUrl);
          if (!response.ok) {
            throw new Error('Network response was not ok.');
          }
          this.userData.avatar = avatarUrl; 
          console.log(this.userData.avatar)

        } catch (error) {
          console.error('There was a problem fetching the user avatar:', error);
        }
      },
    async sendFollowRequest(id) {
      try {
        const response = await axios.post('http://localhost:8080/api/submitContent',{
          Type: 'followRequest',
          target_user_id: parseInt(id,10),
        }, {
          withCredentials: true
        });
        if (response.status === 200) {
          console.log("Follow request sent successfully.")
        }
      } catch (error) {
        console.error("Failed to send follow request:", error);
        this.followError = error.response.data;
        }
    },
    showModal(modalType) {
      this.modals[modalType] = true;
    },
    closeModal(modalType) {
      this.modals[modalType] = false;
    },
    },

    watch: {
      //fetches the new profile data if going to new profile
      $route(to, from) {
      if (to.params.id !== from.params.id) {
        this.fetchUserData(to.params.id);
      }
      }  
    },
/*     computed: {
      absoluteAvatarURL() {
      // Assuming userData.avatar is the relative path from the backend
      // Concatenate the backend base URL with the relative path of the avatar
      console.log(`http://localhost:8080/${this.userData.avatar}`);
      return `http://localhost:8080/${this.userData.avatar}`;
      }
    } */
    
  };
  </script>


<style>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.profile-header {
  text-align: center;
  margin-bottom: 20px;
}

.avatar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
}

.username {
  font-size: 24px;
  margin-bottom: 5px;
}

.user-details {
  margin-bottom: 15px;
}

.full-name {
  font-weight: bold;
}

.followers-following {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
}

.user-posts {
  border-top: 1px solid #ccc;
  padding-top: 20px;
}

.post {
  margin-bottom: 15px;
  border-bottom: 1px solid #ccc;
  padding-bottom: 15px;
}

.modal {
  position: fixed;
  z-index: 1;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  overflow: auto;
  background-color: rgba(0, 0, 0, 0.5);
}
.modal-content {
  background-color: #fefefe;
  margin: 15% auto;
  padding: 20px;
  border: 1px solid #888;
  width: 15%;
}
.close {
  color: #aaa;
  float: right;
  font-size: 28px;
  font-weight: bold;
}
.close:hover,
.close:focus {
  color: black;
  text-decoration: none;
  cursor: pointer;
}

ul {
  padding: 0;
}

li {
  width: 100%; 
  padding: 5px;
  text-align: left;
}

li:hover {
  background-color: #f0f0f0;
}

.clickable {
  cursor: pointer;
  text-decoration: underline;
  color: blue;
}
</style> 