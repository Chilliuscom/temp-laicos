import { createRouter, createWebHistory } from 'vue-router';

import LoginForm from './components/LoginForm.vue';
import RegistrationForm from './components/RegistrationForm.vue';
import HelloWorld from './components/HelloWorld.vue';
import CreatePost from './components/CreatePost.vue';
import ProfilePage from './components/ProfilePage.vue';
import isAuthenticated from './components/isAuthenticated'; 
import NotFound from './components/NotFound.vue';
import ErrorPage from './components/ErrorPage.vue';


const routes = [
  {
    path: '/',
    component: HelloWorld,
    name: 'HelloWorld', 
    meta: {
      requiresAuth: true,
    }
    },
  {
    path: '/login',
    component: LoginForm,
    name: 'Login', 
    meta: {
      hideHeader: true,
      hideSideBar: true,
      requiresAuth: false,
    }
  },
  {
    path: '/register',
    component: RegistrationForm,
    name: 'Register', 
    meta: {
      hideHeader: true,
      hideSideBar: true,
      requiresAuth: false,
    }
  },
  {
    path: '/post',
    component: CreatePost,
    name: 'createPost', 
    meta: {
      requiresAuth: true,
    }
  },
  {
    path: '/profile/:id',
    component: ProfilePage,
    name: 'profilePage',
    props: true,
    meta: {
      requiresAuth: true,
    }
  },
  {
    path: '/error-page',
    component: ErrorPage,
    meta: {
      hideHeader: true,
      hideSideBar: true,
    }
  },
  {
    path: '/:catchAll(.*)',
    component: NotFound,
    meta: {
      hideHeader: true,
      hideSideBar: true,
    }
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
  isAuthenticated()
    .then((authenticated) => {
      if (to.meta.requiresAuth && !authenticated) {
        next('/login');
      } else if ((to.name === 'Login' || to.name === 'Register') && authenticated) {
        next('/');
      }
      else {
        next();
      }
    })
    .catch((error) => {
      console.error('Error checking authentication:', error);
      next('/error-page'); 
    });
});


export default router;

