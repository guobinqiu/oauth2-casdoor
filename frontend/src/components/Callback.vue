<template>
  <div class="callback">
    <h1>Loading...</h1>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'Callback',
  created: async function () {
    const code = this.$route.query.code;
    if (code) {
      try {
        // 获取token
        const token = await this.getToken(code);

        // 存储访问令牌和刷新令牌
        window.sessionStorage.setItem('token', token.access_token);

        // 重定向到目标页面
        this.$router.push('/home');
      } catch (error) {
        console.log('Error processing login:', error);
        this.$router.push('/login');
      }
    } else {
      this.$router.push('/login');
    }
  },
  methods: {
    async getToken(code) {
      const response = await this.$http.post('/token', {
        code,
      }, {
        headers: {
          'Content-Type': 'application/json',
        },
      });
      return response.data;
    }
  },
}
</script>

<style scoped>
.callback {
  text-align: center;
  margin-top: 100px;
}
</style>