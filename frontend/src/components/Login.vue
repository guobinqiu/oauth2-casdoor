<template>
  <div class="login">
    <h1>Login</h1>
    <button @click="login">Login with OAuth</button>
  </div>
</template>

<script>
export default {
  name: 'Login',
  methods: {
    login() {
      window.location.href = 'http://localhost:9000/login';
    },
  },
};
</script>

<style scoped>
.login {
  text-align: center;
  margin-top: 100px;
}
</style>
