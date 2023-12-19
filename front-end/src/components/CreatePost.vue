<template>
    <div>
      <main>
        <form @submit.prevent="createPost">
          <div class="form-group">
            <label for="title">Title:</label>
            <input type="text" id="title" v-model="post.title" required />
            <label for="content">Content:</label>
            <textarea id="content" v-model="post.content" required></textarea>
            <label for="privacy">Privacy:</label>
            <select id="privacy" v-model="post.privacy">
              <option value="public">Public</option>
              <option value="private">Private</option>
              <option value="select">Select friends</option>
            </select>
            <label for="image">Add Image:</label>
            <input type="file" id="image" @change="handleImageUpload" />
          </div>
          <div class="image-preview" v-if="post.image">
            <img :src="post.image" alt="Image Preview" />
          </div>
          <button type="submit">Post</button>
        </form>
      </main>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        post: {
          title: '',
          content: '',
          privacy: 'public',
          image: null,
        },
      };
    },
    methods: {
      createPost() {
        console.log('Title:', this.post.title);
        console.log('Content:', this.post.content);
        this.resetForm();
      },
      handleImageUpload(event) {
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = () => {
            this.post.image = reader.result;
          };
          reader.readAsDataURL(file);
        }
      },
      resetForm() {
        this.post.title = '';
        this.post.content = '';
        this.post.privacy = 'public';
        this.post.image = null;
      },
    },
  };
  </script>
