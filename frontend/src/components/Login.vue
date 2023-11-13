<template>
    <div class="login_container">
        <div class="login_box">
            <div class="avatar_box">
                <img src="../assets/logo.png" alt="">
            </div>
            <!-- 登录表单区域 -->
            <el-form ref="loginFormRef" label-width="0px" class="login_form" :model="loginForm" :rules="loginFormRules">
                <!-- 用户名 -->
                <el-form-item prop="username">
                    <el-input
                        placeholder="请输入用户名"
                        prefix-icon="iconfont el-icon-user"
                        v-model="loginForm.username"
                        @keyup.enter.native="login"
                    ></el-input>
                </el-form-item>
                <!-- 密码 -->
                <el-form-item prop="password">
                    <el-input
                        placeholder="请输入密码"
                        prefix-icon="iconfont el-icon-lock"
                        v-model="loginForm.password"
                        @keyup.enter.native="login"
                        show-password
                    ></el-input>
                </el-form-item>
                <!-- 提交按钮 -->
                <el-form-item class="btns">
                    <el-button type="primary" @click="login">登录</el-button>
                    <el-button type="info" @click="resetLoginForm">重置</el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script>
export default {
  data () {
    return {
      loginForm: {
        username: '',
        password: ''
      },
      loginFormRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    resetLoginForm () {
      this.$refs.loginFormRef.resetFields()
    },
    login () {
      this.$refs.loginFormRef.validate(async valid => {
        // eslint-disable-next-line no-useless-return
        if (!valid) return
        await this.$http.post('login', this.loginForm)
          .then(response => {
            // eslint-disable-next-line eqeqeq
            if (response.data.code != 200) {
              return this.$message.error(response.data.message)
            } else {
              const requestToken = 'Bearer ' + response.data.data.token
              this.$setLocalStorage('token', requestToken, 86400000)
              //   this.$setLocalStorage('userInfo', JSON.stringify(response.data.data), 86400)
              this.$router.push('/home')
              return this.$message.success(response.data.message)
            }
          })
          .catch(error => {
            return this.$message.error(error)
          })
      })
    }
  }
}
</script>

<style lang="less" scoped>
.login_container {
    background-color: #2b4b6b;
    height: 100%;
}

.login_box {
    width: 450px;
    height: 300px;
    background-color: #fff;
    border-radius: 3px;
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    .avatar_box{
        height: 130px;
        width: 130px;
        border: 1px solid #eee;
        border-radius: 50%;
        padding: 10px;
        box-shadow: 0 0 10px #ddd;
        position: absolute;
        left: 50%;
        transform: translate(-50%, -50%);
        background-color: #fff;
        img{
            width: 100%;
            height: 100%;
            border-radius: 50%;
            background-color: #eee;
        }
    }
}

.btns {
    display: flex;
    justify-content: flex-end;
}

.login_form {
    position: absolute;
    bottom: 0;
    width: 100%;
    padding: 0 20px;
    box-sizing: border-box;
}
</style>
