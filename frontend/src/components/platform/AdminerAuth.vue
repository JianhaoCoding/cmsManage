<!-- eslint-disable vue/no-unused-vars -->
<template>
    <div>
        <!-- 面包屑导航 -->
        <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>系统管理</el-breadcrumb-item>
            <el-breadcrumb-item>权限列表</el-breadcrumb-item>
        </el-breadcrumb>

        <!-- 卡片试图 -->
        <el-card>
           <!-- 搜索与添加区域 -->
            <el-row :gutter="10">
                <el-col :span="4">
                    <el-input placeholder="请输入权限名" v-model="searchInfo.auth_name" clearable @clear="clearName"></el-input>
                </el-col>
                <el-col :span="4">
                    <el-input placeholder="请输入权限" v-model="searchInfo.auth_con_act" clearable></el-input>
                </el-col>
                <el-col :span="4">
                    <el-select v-model="searchInfo.auth_group_id" clearable placeholder="请选择权限组" @change="handleGroupChange">
                        <el-option
                        v-for="item in groupOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                </el-col>
                <el-col :span="6">
                    <el-select v-model="searchInfo.status" clearable placeholder="请选择状态" @change="handleStatusChange">
                        <el-option
                        v-for="item in authStatusOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                    <el-button icon="el-icon-search" @click="getAdminerAuthList"></el-button>
                </el-col>
                <el-col :span="4">
                    <el-button type="primary" @click="showAuthDialog">添加权限</el-button>
                </el-col>
            </el-row>

            <!-- 权限列表 -->
            <el-table :data="adminerAuthList" border stripe>
                <el-table-column label="ID" prop="auth_id" width="50px" align="center"></el-table-column>
                <el-table-column label="权限名称" prop="auth_name" width="250px"></el-table-column>
                <el-table-column label="权限" prop="auth_con_act" width="250px"></el-table-column>
                <el-table-column label="备注" prop="remark" width="200px"></el-table-column>
                <el-table-column label="菜单展示" width="100px">
                    <template slot-scope="scope">
                        <el-switch
                        v-model="scope.row.is_menu_show"
                        :active-value="2"
                        :inactive-value="1"
                        @change="handleSwitchMenuChange(scope.row)"
                        ></el-switch>
                    </template>
                </el-table-column>
                <el-table-column label="状态" width="100px">
                    <template slot-scope="scope">
                        <el-switch
                        :value="getSwitchStatus(scope.row)"
                        :active-value="2"
                        :inactive-value="1"
                        @change="handleSwitchChange(scope.row)"
                        ></el-switch>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="180px">
                    <template slot-scope="scope">
                        <el-tooltip class="item" e ffect="dark" content="编辑" placement="top" :enterable="false">
                            <el-button type="primary" icon="el-icon-edit" size="mini" @click="showEditDialog(scope.row.auth_id)"></el-button>
                        </el-tooltip>
                        <el-tooltip class="item" effect="dark" content="删除" placement="top" :enterable="false">
                            <el-button type="danger" icon="el-icon-delete" size="mini" @click="deleteAuth(scope.row)"></el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
             </el-table>

             <!-- 分页 -->
             <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="searchInfo.page"
            :page-sizes="[1, 2, 5, 10, 20]"
            :page-size="searchInfo.page_size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total">
            </el-pagination>
        </el-card>

        <!-- 添加编辑权限弹窗 -->
        <el-dialog
        :title="dialogTitle"
        :visible.sync="addDialogVisible"
        width="50%"
        @close="addAdminerAuthFormClosed">
        <!-- 管理员信息主体区 -->
        <el-form
        ref="addAdminerAuthFormRef"
        :model="addAdminerAuthForm"
        :rules="addAdminerAuthFormRules"
        label-width="80px">
            <input type="hidden" v-model="addAdminerAuthForm.auth_id">
            <el-form-item label="权限名称" prop="auth_name">
                <el-input v-model="addAdminerAuthForm.auth_name"></el-input>
            </el-form-item>
            <el-form-item label="权限">
                <el-input v-model="addAdminerAuthForm.auth_con_act"></el-input>
            </el-form-item>
            <!-- 权限组选择 -->
            <el-form-item label="权限组" prop="auth_group_id">
            <el-select v-model="addAdminerAuthForm.auth_group_id" placeholder="请选择用户组" clearable>
                <el-option
                v-for="item in groupOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
                ></el-option>
            </el-select>
            </el-form-item>
            <!-- 父级权限 -->
            <el-form-item label="父级权限">
            <el-select v-model="addAdminerAuthForm.prent_id" placeholder="请选择父级权限" clearable>
                <el-option
                v-for="item in prentOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
                ></el-option>
            </el-select>
            </el-form-item>
            <!-- 状态选择 -->
            <el-form-item label="状态">
            <el-switch
                v-model="addAdminerAuthForm.status"
                :active-value="2"
                :inactive-value="1"
                active-color="#13ce66"
                inactive-color="#ff4949"
            ></el-switch>
            </el-form-item>
            <!-- 是否菜单展示 -->
            <el-form-item label="菜单展示">
            <el-switch
                v-model="addAdminerAuthForm.is_menu_show"
                :active-value="2"
                :inactive-value="1"
                active-color="#13ce66"
                inactive-color="#ff4949"
            ></el-switch>
            </el-form-item>
            <el-form-item label="备注">
                <el-input v-model="addAdminerAuthForm.remark"></el-input>
            </el-form-item>
        </el-form>
        <!-- 底部区 -->
        <span slot="footer" class="dialog-footer">
            <el-button @click="addDialogVisible = false">取 消</el-button>
            <el-button type="primary" :disabled="isSaving" @click="addAdminerAuthFormSub">确 定</el-button>
        </span>
        </el-dialog>
    </div>
</template>

<script>

export default {
  data() {
    return {
      // 查询表单数据
      searchInfo: {
        auth_name: '',
        auth_con_act: '',
        auth_group_id: null,
        status: null,
        page: 1,
        page_size: 10
      },
      dialogTitle: '添加权限',
      // 管理员组选项
      groupOptions: [],
      groupOptionDefValue: null,
      // 父级权限选项
      prentOptions: [],
      prentOptionDefValue: null,
      // 管理员状态选项
      authStatusOptions: [],
      // 状态开关选项数据
      switchStatusInfo: {
        auth_id: 0,
        status: 1
      },
      authStatusDefValue: null,
      // 是否菜单选项
      isMenuOptions: [],
      // 状态开关选项数据
      switchMenuInfo: {
        auth_id: 0,
        is_menu_show: 1
      },
      isMenuShowDefValue: null,
      adminerAuthList: [],
      total: 0,
      listLoading: true,
      // 添加管理员表单数据
      addAdminerAuthForm: {
        auth_id: 0,
        auth_name: '',
        auth_con_act: '',
        auth_group_id: null,
        status: 2,
        remark: '',
        prent_id: null,
        is_menu_show: 1,
        sort: 0
      },
      // 添加管理员表单验证规则
      addAdminerAuthFormRules: {
        auth_name: [
          { required: true, message: '请输入权限名', trigger: 'blur' },
          { min: 3, max: 10, message: '权限名长度在3-10个字符', trigger: 'blur' }
        ],
        auth_group_id: [
          { required: true, message: '请选择权限组', trigger: 'change' }
        ]
      },
      // 控制添加管理员的显示与隐藏
      addDialogVisible: false,
      addAdminerAuthFormRef: {},
      isSaving: false
    }
  },
  created() {
    this.getAdminerAuthList()
  },
  methods: {
    async getAdminerAuthList() {
      const { data: res } = await this.$http.get('adminer/auth/list', {
        params: this.searchInfo
      })

      if (res.code !== 200) {
        return this.$message.error('获取权限列表失败, 失败原因：' + res.message)
      }

      this.total = res.data.total
      this.groupOptions = res.data.authGroupOptions
      this.prentOptions = res.data.prentAuthOptions
      this.isMenuOptions = res.data.isMenuOptions
      this.adminerAuthList = res.data.list
      this.authStatusOptions = res.data.statusOptions
    },
    async handleSwitchChange(row) {
      const authId = row.auth_id
      this.switchStatusInfo.auth_id = authId
      this.switchStatusInfo.status = row.status
      const { data: res } = await this.$http.post('adminer/auth/switchstatus', this.switchStatusInfo)
      if (res.code !== 200) {
        row.status = this.switchStatusInfo.status === 1 ? 2 : 1
        return this.$message.error('更新失败, 失败原因：' + res.message)
      }
      const statusStr = row.status === 1 ? '禁用' : '开启'
      this.$message.success(statusStr + '成功')
    },
    async handleSwitchMenuChange(row) {
      const authId = row.auth_id
      this.switchMenuInfo.auth_id = authId
      this.switchMenuInfo.is_menu_show = row.is_menu_show
      const { data: res } = await this.$http.post('adminer/auth/switchmenu', this.switchMenuInfo)
      if (res.code !== 200) {
        row.status = this.switchStatusInfo.status === 1 ? 2 : 1
        return this.$message.error('更新失败, 失败原因：' + res.message)
      }
      const statusStr = row.is_menu_show === 1 ? '隐藏' : '显示'
      this.$message.success(statusStr + '成功')
    },
    // 监听pageSize改变事件
    handleSizeChange(newSize) {
      this.searchInfo.page_size = newSize
      this.getAdminerAuthList()
    },
    // 监听页码改变事件
    handleCurrentChange(newPage) {
      this.searchInfo.page = newPage
      this.getAdminerAuthList()
    },
    handleGroupChange(newGroupId) {
      this.searchInfo.auth_group_id = newGroupId
      this.getAdminerAuthList()
    },
    handleStatusChange(newStatus) {
      this.searchInfo.status = newStatus
      this.getAdminerAuthList()
    },
    clearName() {
      this.searchInfo.auth_name = ''
      this.getAdminerAuthList()
    },
    // 监听添加权限表单的关闭事件
    addAdminerAuthFormClosed() {
      this.$refs.addAdminerAuthFormRef.resetFields()
      this.isSaving = false
    },
    // 添加权限表单的提交事件
    addAdminerAuthFormSub() {
      this.isSaving = true
      this.$refs.addAdminerAuthFormRef.validate(async valid => {
        if (!valid) {
          this.isSaving = false
          return
        }

        try {
          // 发起添加权限请求
          const { data: res } = await this.$http.post('/adminer/auth/save', this.addAdminerAuthForm)
          if (res.code !== 200) {
            return this.$message.error(res.message)
            // eslint-disable-next-line no-unreachable
            this.isSaving = false
          } else {
            this.$message.success(res.message)
            setTimeout(() => {
              this.getAdminerAuthList()
              this.addDialogVisible = false
            }, 1500)
          }
        } catch (error) {
          this.$message.error('保存失败')
        }
      })
    },
    deleteAuth(row) {
      // 发送请求删除数据
      this.$http.delete(`/adminer/auth/delete/${row.auth_id}`).then(response => {
        if (response.data.code === 200) {
          // 请求成功后，找到要删除的项的索引
          const index = this.adminerAuthList.findIndex(item => item.auth_id === row.auth_id)
          if (index !== -1) {
            // 从数组中删除该项
            this.adminerAuthList.splice(index, 1)
            this.$message.success('删除成功')
            setTimeout(() => {
              this.getAdminerAuthList()
            }, 1000)
          }
        } else {
          // 处理错误情况
          this.$message.error('删除失败')
        }
      }).catch(error => {
        // 处理请求错误情况
        console.error('删除请求发生错误', error)
        this.$message.error('删除失败')
      })
    },
    // 展示编辑用户的对话框
    async showEditDialog(authId) {
      const { data: res } = await this.$http.get(`/adminer/auth/show/${authId}`)
      if (res.code !== 200) {
        return this.$message.error(res.message)
      }

      this.addAdminerAuthForm = res.data
      this.addDialogVisible = true
    },
    getSwitchStatus(row) {
      return row.status === 99 ? 1 : row.status
    },
    showAuthDialog() {
      this.addDialogVisible = true
      this.$refs.addAdminerAuthFormRef.resetFields()
      this.addAdminerAuthForm.auth_con_act = ''
      this.addAdminerAuthForm.prent_id = null
    }
  }
}

</script>

<style lang="less" scoped>
</style>
