<template>
  <div class="load" v-if="!show" style="z-index: 100;position: absolute;width: 100%;height: 100%"></div>
  <component v-show="show" :is="currentTemplate"></component>
</template>

<script>
import LoginPage from "@/pages/Login";
import IndexPage from "@/pages/Index";

export default {
  name: 'App',
  data: () => {
    return {
      "currentTemplate": "LoginPage",
      "show": false
    }
  },
  components: {
    LoginPage,IndexPage
  },
  mounted() {
    let loading = this.$loading({
      fullscreen: true,
      text: "资源加载中...",
      target: ".load",
      lock: true
    })
    this.$api.checkLogin().then(res => {
      if (res.data.state) {
        this.currentTemplate = "IndexPage"
      }
      this.show = true
      loading.close()
    }).catch(() => {
      this.show = true
      loading.close()
    })
  }
}
</script>
<style scoped>
.load {
  width: 100%;
  height: 100%;
  position: absolute;
}
</style>

