<template>
  <div class="app">
    <div class="login">
      <div class="title">
        <span>DnsServer</span>
      </div>
      <div class="msg">
        <span v-text="msg"></span>
      </div>
      <div class="username">
        <el-input @focus="msg=''" @keydown.enter="login" v-model="username" placeholder="请输入用户名" class="w-50 m-2">
          <template #prefix>
            <el-icon class="el-input__icon">
              <User/>
            </el-icon>
          </template>
        </el-input>
      </div>

      <div class="password">
        <el-input  @focus="msg=''"  @keydown.enter="login" v-model="password" type="password" placeholder="请输入密码" class="w-50 m-2">
          <template #prefix>
            <el-icon class="el-input__icon">
              <Lock/>
            </el-icon>
          </template>
        </el-input>
      </div>

      <div class="submit">
        <el-button @click="login" type="primary">登录</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import {Lock, User} from '@element-plus/icons-vue'

export default {
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'LoginPage',
  data: () => {
    return {
      "username": "",
      "password": "",
      "msg": ""
    }
  },
  components: {
    User, Lock
  },
  methods: {
    login() {
      if (!this.username) {
        this.msg = "用户名不能为空"
        return
      }
      if (!this.password) {
        this.msg = "密码不能为空"
        return
      }
      this.msg = ""
      let loading = this.$loading({
        text:"登录中...",
        target:".app"
      })
      this.$api.login(this.username,this.password).then(res => {
        let data = res.data;
        if(!data.state){
          this.msg = data.msg
        }else {
          this.$message.success("登录成功,1秒后跳转主页")
          this.$store.setToken(data.data)
          setTimeout(() => {
            this.$parent.currentTemplate = "IndexPage"
          },1000)
        }
        loading.close()
      }).catch(err => {
        console.log(err)
        loading.close()
      })
    }
  },
  mounted() {

  }
}
</script>
<style scoped>

.app{
  position: absolute;
  background: url("../assets/background.jpg");
  background-size: 100% 100%;
  width: 100%;
  height: 100%;
}

.submit {
  display: flex;
  height: 35px;
  justify-content: left;
  font-size: 12px;
  width: 80%;
  margin-left: 10%;
  margin-top: 30px;
}

.submit > button {
  width: 100%;
}

.username {
  display: flex;
  height: 35px;
  justify-content: left;
  font-size: 12px;
  width: 80%;
  margin-left: 10%;
  margin-top: 5px;
}

.password {
  margin-top: 20px;
  height: 35px;
  display: flex;
  justify-content: left;
  font-size: 12px;
  width: 80%;
  margin-left: 10%;
}

.msg {
  display: flex;
  justify-content: center;
  font-size: 12px;
  color: #C55323FF;
  width: 80%;
  margin-left: 10%;
  padding-bottom: 5px;
  height: 12px;
}

.title {
  display: flex;
  justify-content: center;
  padding-top: 30px;
  padding-bottom: 15px;
  font-size: 18px;
  font-weight: bold;
  color: #0E645AFF;
  box-sizing: border-box;
}

.login {
  position: absolute;
  width: 330px;
  height: 280px;
  left: 50%;
  margin-left: -180px;
  top: 50%;
  margin-top: -180px;
  border-radius: 5px;
  box-shadow: 3px 3px 5px 1px #115E43D8;
  background: rgba(255, 255, 255, 0.38);
}
</style>

