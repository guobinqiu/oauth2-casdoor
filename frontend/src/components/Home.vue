<template>
  <div class="home">
    <h1>Welcome, {{ name }}!</h1>
    <button @click="logout">Logout</button>
  </div>
</template>

<script>
export default {
  name: 'Home',
  data() {
    return {
      name: '',
      token: '',
    };
  },
  created() {
    this.token = window.sessionStorage.getItem('token');
    this.getUserInfo();
  },
  methods: {
    getUserInfo() {
      this.$http.get('/userinfo', {
        headers: {
          Authorization: `Bearer ${this.token}`,
        },
      })
      .then(response => {
        console.log(response);
        this.name = response.data.name;
      })
      .catch(error => {
        console.log('Error fetching user info:', error);
      });
    },
    logout() {
      window.sessionStorage.clear()
      this.$router.push('/');
    },
  },
};
</script>

<style scoped>
.home {
  text-align: center;
  margin-top: 100px;
}
</style>
