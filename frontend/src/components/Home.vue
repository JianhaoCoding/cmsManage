<!-- eslint-disable vue/no-parsing-error -->
<template>
  <el-container class="home-container">
    <!-- 头部区域 -->
    <el-header>
        <div>
            <img src="" alt="">
            <span>CMS管理系统</span>
        </div>

        <!-- 用户信息和下拉菜单 -->
        <el-dropdown @command="handleDropdownCommand">
          <span class="el-dropdown-link">
            <img :src="adminerInfo.avatar" class="user-avatar">
            {{ adminerInfo.name }}
            <i class="el-icon-arrow-down el-icon--right"></i>
          </span>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item command="settings">设置</el-dropdown-item>
            <el-dropdown-item command="logout">退出</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>

    </el-header>
    <el-container>
        <el-aside :width="isCollapse ? '64px' : '200px'">
          <div class="toggle-button" @click="toggleCollapse">|||</div>
            <!-- 侧边栏菜单区 -->
            <el-menu
            background-color="#333744" text-color="#fff" active-text-color="#409EFF" :collapse="isCollapse" :collapse-transition="isCollapseTransition" router :default-active="activePath">
                <!-- 一级菜单 -->
                <el-submenu :index="item.auth_id + ''" v-for="item in menuList" :key="item.auth_id">
                    <!-- 一级菜单模版区域 -->
                    <template slot="title">
                        <!-- 图标 -->
                        <i :class="iconsObj[item.auth_id]"></i>
                        <!-- 文本 -->
                        <span>{{item.auth_name}}</span>
                    </template>

                    <!-- 二级菜单 -->
                    <template v-for="subItem in item.children">
                        <el-menu-item
                        :index="subItem.path"
                        v-if="subItem.is_menu_show === 2"
                        :key="subItem.auth_id"
                        @click="saveNavState(subItem.path)"
                        >
                        <!-- 图标 -->
                        <i class="el-icon-menu"></i>
                            <!-- 文本 -->
                            <span>{{subItem.auth_name}}</span>
                        </el-menu-item>
                    </template>
                </el-submenu>

            </el-menu>
        </el-aside>
        <el-main>
          <!-- 路由占位符 -->
          <router-view></router-view>
        </el-main>
    </el-container>
  </el-container>
</template>

<script>

export default {
  data() {
    return {
      // 左侧菜单数据
      menuList: [],
      iconsObj: {
        24: 'el-icon-user',
        1: 'el-icon-s-tools'
      },
      // 是否折叠
      isCollapse: false,
      // 是否开启折叠动画
      isCollapseTransition: false,
      // 被激活的链接地址
      activePath: '',
      adminerInfo: {
        name: '',
        avatar: ''
      }
    }
  },
  created() {
    // Home组件创建的时候直接请求左侧导航
    this.getAdminerInfo()
    this.getMenuList()
    this.activePath = window.sessionStorage.getItem('activePath')
  },
  methods: {
    loginout () {
      // 退出登录
      localStorage.removeItem('userInfo')
      localStorage.removeItem('token')
      this.$router.push('/login')
    },
    async getAdminerInfo () {
      const { data: res } = await this.$http.get('adminer')
      // 管理员信息获取失败
      if (res.code !== 200) {
        if (res.code === 401) {
          this.$message.error('登录状态失效，请重新登录')
          this.$router.push('/login')
          return
        }
        this.$message.error('获取管理员信息失败！原因：' + res.message)
        return
      }

      // 管理员信息获取成功
      this.adminerInfo.name = res.data.username
      this.adminerInfo.avatar = this.generateAvatar(res.data.username)
    },
    async getMenuList () {
      const { data: res } = await this.$http.get('menurules')
      // 列表获取失败
      if (res.code !== 200) {
        if (res.code === 401) {
          this.$message.error('登录状态失效，请重新登录')
          this.$router.push('/login')
          return
        }
        this.$message.error('获取菜单列表失败！失败原因：' + res.message)
        return
      }

      // 列表获取成功
      this.menuList = res.data
    },
    // 点击按钮，切换菜单的折叠与展开
    toggleCollapse () {
      this.isCollapse = !this.isCollapse
    },
    // 保存链接激活状态
    saveNavState(activePath) {
      this.activePath = activePath
      window.sessionStorage.setItem('activePath', activePath)
    },
    handleDropdownCommand(command) {
      if (command === 'logout') {
        // 调用已有的退出登录方法
        this.loginout()
      } else if (command === 'settings') {
        // 处理设置逻辑
        // 例如跳转到设置页面
        // this.$router.push('/settings');
      }
    },
    generateAvatar(name) {
      // 获取用户名的最后一个字符
      const char = name.charAt(name.length - 1)
      const canvas = document.createElement('canvas')
      const size = 100
      canvas.width = size
      canvas.height = size

      const context = canvas.getContext('2d')
      context.fillStyle = '#f0f0f0'
      context.fillRect(0, 0, size, size)

      context.fillStyle = '#333'
      context.font = '48px sans-serif'
      context.textAlign = 'center'
      context.textBaseline = 'middle'
      context.fillText(char, size / 2, size / 2)

      return canvas.toDataURL()
    }
  }
}
</script>

<style lang="less" scoped>
.el-header{
    background-color: #373d41;
    display: flex;
    justify-content: space-between;
    padding-left: 0;
    align-items: center;
    color: #fff;
    font-size: 20px;
    > div {
        display: flex;
        align-items: center;
        span {
            margin-left: 15px;
        }
    }
}
.el-aside{
    background-color: #333744;
    .el-menu{
        border-right: none; //防止el-menu的边框对不齐
    }
}
.el-main{
    background-color: #eaedf1;
}
.home-container{
    height: 100%;
}

.iconfont {
  margin-right: 10px;
}

.toggle-button{
  background-color: #4a5064;
  font-size: 10px;
  line-height: 24px;
  color: #fff;
  text-align: center;
  letter-spacing: 0.2em;
  cursor: pointer;
}

.user-avatar {
  width: 30px; /* 调整为所需大小 */
  height: 30px;
  border-radius: 50%;
  margin-right: 10px;
}

.el-dropdown-link {
  display: flex; /* 使用 Flexbox 布局 */
  align-items: center; /* 垂直居中对齐 */
}
</style>
