<template>
  <div id="app">
    <header v-if="!$route.meta.hideHeader">
      <AppHeader />
    </header>
    <div class="container">
    <AppSidebar v-if="!$route.meta.hideSideBar"/>
    <main class="main-content">
      <router-view></router-view>
    </main>
    </div>
    <footer>
    </footer>
  </div>
</template>

<script>
import AppSidebar from './components/AppSidebar.vue';
import AppHeader from './components/AppHeader.vue'; 

export default {
  components: {
    AppSidebar,
    AppHeader,
  }
}; 
</script>


