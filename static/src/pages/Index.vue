<template>
  <div class="app">
    <el-row class="header">
      <el-col :span="3">
        <div class="header_left">
          <span>DnsServer</span>
        </div>
      </el-col>
      <el-col :span="21">
        <div class="header_right">
          <el-menu mode="horizontal" :ellipsis="false">
            <el-menu-item index="1">菜单1</el-menu-item>
            <el-menu-item index="2">菜单2</el-menu-item>
            <el-sub-menu index="Person">
              <template #title>
                <el-avatar :size="30" :src="avatar"/>
              </template>
              <el-menu-item @click="toRoute('PersonCenter')" index="PersonCenter">个人中心</el-menu-item>
              <el-menu-item @click="logout" index="Logout">退出登录</el-menu-item>
            </el-sub-menu>
          </el-menu>
        </div>
      </el-col>
    </el-row>
    <el-row class="body">
      <el-radio-group v-model="menuIsCollapse" size="small" style="position: absolute;margin-bottom: 10px;z-index: 999">
        <el-radio-button size="small" :label="false">展开</el-radio-button>
        <el-radio-button size="small" :label="true">收缩</el-radio-button>
      </el-radio-group>
      <el-col class="col" :span="menuIsCollapse ? 0.5 : 3">
        <el-menu
            :default-active="menuActive"
            :collapse="menuIsCollapse"
            style="border: none"
            active-text-color="#ffd04b"
            background-color="#545c64"
            text-color="#fff">
          <el-sub-menu index="1">
            <template #title>
              <el-icon>
                <Location/>
              </el-icon>
              <span>域名管理</span>
            </template>
            <el-menu-item @click="toRoute('DomainForward')" index="DomainForward">
              <template #title>
                <el-icon>
                  <Location/>
                </el-icon>
                <span>正向解析</span>
              </template>
            </el-menu-item>
            <el-menu-item @click="toRoute('DomainReverse')" index="DomainReverse">
              <template #title>
                <el-icon>
                  <Location/>
                </el-icon>
                <span>反向解析</span>
              </template>
            </el-menu-item>
          </el-sub-menu>

          <el-sub-menu index="2">
            <template #title>
              <el-icon>
                <Setting/>
              </el-icon>
              <span>设置</span>
            </template>
            <el-menu-item @click="toRoute('PersonCenter')" index="PersonCenter">
              <template #title>
                <el-icon>
                  <Location/>
                </el-icon>
                <span>个人中心</span>
              </template>
            </el-menu-item>
            <el-menu-item @click="toRoute('Other')" index="Other">
              <template #title>
                <el-icon>
                  <Location/>
                </el-icon>
                <span>其他</span>
              </template>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-col>
      <el-main class="main">
        <router-view/>
      </el-main>
    </el-row>
  </div>
</template>
<script>
import {Location, Setting} from '@element-plus/icons-vue'
import avatar from '@/assets/avatar.jpg'

export default {
  name: "IndexPage",
  data: () => {
    return {
      menuIsCollapse: false,
      menuActive: "Home",
      avatar
    }
  },
  components: {
    Location, Setting
  },
  methods: {
    logout() {
      this.$api.logout().then(res => {
        let data = res.data;
        if (!data.state){
          this.$message.error(data.msg)
        }else {
          this.$message.success(data.msg)
          this.$store.clearToken()
          this.$store.del("routeName")
          this.$parent.currentTemplate = "LoginPage"
        }
      }).catch(err => {
        console.log(err)
        this.$message.error("退出登录失败")
      })
    },
    toRoute(name) {
      this.$router.push({
        name: name
      })
      this.menuActive = name
      this.$store.set("routeName", name)
    }
  },
  mounted() {
    let routeName = this.$store.get("routeName");
    if (routeName) {
      this.menuActive = routeName
    }
  }
}
</script>

<style scoped>
.header_left {
  width: 100%;
  height: 50px;
  display: flex;
  justify-content: center;
  font-size: 25px;
  line-height: 45px;
  font-weight: bold;
  color: #237539FF;
  background: url("../assets/banner.jpg");
  opacity: 0.7;
  background-size: 100% 100%;
}

.el-sub-menu .el-menu-item {
  min-width: 0 !important;
}

.header_right {
  width: 100%;
  height: 50px;
  justify-content: right;
  display: flex;
  background: black;
  border-left: 1px solid white;
  box-sizing: border-box;
}

.app {
  width: 100vm;
  height: 100vh;
}

.header {
  height: 50px;
  width: 100%;
}

.body {
  height: calc(100% - 50px);
  width: 100%;
}

.main {
  box-sizing: border-box;
  background: gray;
  border-left: 1px solid white;
}

.col {
  padding-top: 25px;
  height: 100%;
  background: #545c64;
}
</style>