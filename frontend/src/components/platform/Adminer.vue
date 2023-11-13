<!-- eslint-disable vue/no-unused-vars -->
<template>
    <div>
        <!-- 面包屑导航 -->
        <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>系统管理</el-breadcrumb-item>
            <el-breadcrumb-item>管理员列表</el-breadcrumb-item>
        </el-breadcrumb>

        <!-- 卡片试图 -->
        <el-card>
           <!-- 搜索与添加区域 -->
            <el-row :gutter="10">
                <el-col :span="4">
                    <el-input placeholder="请输入用户名" v-model="searchInfo.username" clearable @clear="clearName"></el-input>
                </el-col>
                <el-col :span="4">
                    <el-input placeholder="请输入手机号" v-model="searchInfo.mobile" clearable></el-input>
                </el-col>
                <el-col :span="6">
                    <el-select v-model="searchInfo.group_id" clearable placeholder="请选择" @change="handleGroupChange">
                        <el-option
                        v-for="item in groupOptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                    <el-button icon="el-icon-search" @click="getAdminerList"></el-button>
                </el-col>
                <el-col :span="4">
                    <el-button type="primary" @click="addDialogVisible = true">添加管理员</el-button>
                </el-col>
            </el-row>

            <!-- 用户列表 -->
            <el-table :data="adminerList" border stripe>
                <el-table-column label="ID" prop="adminer_id" width="50px"></el-table-column>
                <el-table-column label="用户名" prop="username" width="100px"></el-table-column>
                <el-table-column label="昵称" prop="nickname" width="120px"></el-table-column>
                <el-table-column label="手机号" prop="mobile" width="120px"></el-table-column>
                <el-table-column label="Email" prop="email" width="200px"></el-table-column>
                <el-table-column label="用户组" prop="Group.name" width="100px"></el-table-column>
                <el-table-column label="状态">
                    <template slot-scope="scope">
                        <el-switch
                        v-model="scope.row.status"
                        :active-value="2"
                        :inactive-value="1"
                        @change="handleSwitchChange(scope.row)"
                        ></el-switch>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="180px">
                    <template slot-scope="scope">
                        <el-tooltip class="item" e ffect="dark" content="编辑" placement="top" :enterable="false">
                            <el-button type="primary" icon="el-icon-edit" size="mini" @click="showEditDialog(scope.row.adminer_id)"></el-button>
                        </el-tooltip>
                        <el-tooltip class="item" effect="dark" content="删除" placement="top" :enterable="false">
                            <el-button type="danger" icon="el-icon-delete" size="mini" @click="deleteAdminer(scope.row)"></el-button>
                        </el-tooltip>
                        <el-tooltip class="item" effect="dark" content="重置密码" placement="top" :enterable="false">
                            <el-button type="warning" icon="el-icon-refresh-left" size="mini" @click="resetPass(scope.row.adminer_id)"></el-button>
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

        <!-- 添加管理员弹窗 -->
        <el-dialog
        title="添加管理员"
        :visible.sync="addDialogVisible"
        width="50%"
        @close="addAdminerFormClosed">
        <!-- 管理员信息主体区 -->
        <el-form
        ref="addAdminerFormRef"
        :model="addAdminerForm"
        :rules="addAdminerFormRules"
        label-width="80px">
            <el-form-item label="用户名" prop="username">
                <el-input v-model="addAdminerForm.username"></el-input>
            </el-form-item>
            <el-form-item label="姓名" prop="nickname">
                <el-input v-model="addAdminerForm.nickname"></el-input>
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
                <el-input v-model="addAdminerForm.email"></el-input>
            </el-form-item>
            <el-form-item label="手机号" prop="mobile">
                <el-input v-model="addAdminerForm.mobile"></el-input>
            </el-form-item>
            <el-form-item label="备注">
                <el-input v-model="addAdminerForm.remark"></el-input>
            </el-form-item>
            <!-- 用户组选择 -->
            <el-form-item label="用户组" prop="group_id">
            <el-select v-model="addAdminerForm.group_id" placeholder="请选择用户组" clearable>
                <el-option
                v-for="item in groupOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
                ></el-option>
            </el-select>
            </el-form-item>
            <!-- 状态选择 -->
            <el-form-item label="状态" prop="status">
            <el-switch
                v-model="addAdminerForm.status"
                :active-value="2"
                :inactive-value="1"
                active-color="#13ce66"
                inactive-color="#ff4949"
            ></el-switch>
            </el-form-item>
        </el-form>
        <!-- 密码复制展示区 -->
        <div v-if="password" class="dialog-footer">
        <el-input
            ref="passwordRef"
            v-model="password"
            :suffix-icon="'el-icon-document-copy'"
            @click="copyPassword"
            readonly
        ></el-input>
        <el-button @click="copyPassword">复制密码</el-button>
        </div>
        <!-- 底部区 -->
        <span slot="footer" class="dialog-footer">
            <el-button @click="addDialogVisible = false">取 消</el-button>
            <el-button type="primary" :disabled="isSaving" @click="addAdminerFormSub">确 定</el-button>
        </span>
        </el-dialog>

        <!-- 修改管理员 -->
        <el-dialog title="编辑管理员" :visible.sync="editDialogVisible" width="50%" @close="editDialogClosed">
            <!-- 管理员信息主体区 -->
            <el-form
            ref="editAdminerFormRef"
            :model="editAdminerForm"
            :rules="editAdminerFormRules"
            label-width="80px">
                <el-form-item label="用户名">
                    <el-input v-model="editAdminerForm.username" disabled></el-input>
                </el-form-item>
                <el-form-item label="姓名" prop="nickname">
                    <el-input v-model="editAdminerForm.nickname"></el-input>
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                    <el-input v-model="editAdminerForm.email"></el-input>
                </el-form-item>
                <el-form-item label="手机号" prop="mobile">
                    <el-input v-model="editAdminerForm.mobile"></el-input>
                </el-form-item>
                <el-form-item label="备注">
                    <el-input v-model="editAdminerForm.remark"></el-input>
                </el-form-item>
                <!-- 用户组选择 -->
                <el-form-item label="用户组" prop="group_id">
                <el-select v-model="editAdminerForm.group_id" placeholder="请选择用户组" clearable>
                    <el-option
                    v-for="item in groupOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                    ></el-option>
                </el-select>
                </el-form-item>
                <!-- 状态选择 -->
                <el-form-item label="状态" prop="status">
                <el-switch
                    v-model="editAdminerForm.status"
                    :active-value="2"
                    :inactive-value="1"
                    active-color="#13ce66"
                    inactive-color="#ff4949"
                ></el-switch>
                </el-form-item>
            </el-form>

            <span slot="footer" class="dialog-footer">
                <el-button @click="editDialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="editAdminerInfo">确 定</el-button>
            </span>
        </el-dialog>

        <!-- 新密码弹窗 -->
        <el-dialog
            title="新密码"
            :visible.sync="resetPasswordDialogVisible"
            width="30%"
            @close="copyPassClosed"
        >
            <div style="text-align: center;">
            <el-input
                v-model="newPassword"
                ref="newPassword"
                readonly
            ></el-input>
            <el-button
                type="primary"
                icon="el-icon-copy-document"
                @click="copyNewPassword"
            >复制密码</el-button>
            </div>
        </el-dialog>
    </div>
</template>

<script>

export default {
  data() {
    return {
      // 查询表单数据
      searchInfo: {
        username: '',
        mobile: '',
        group_id: null,
        status: 0,
        page: 1,
        page_size: 1
      },
      // 管理员组选项
      groupOptions: [],
      groupOptionDefValue: null,
      // 管理员状态选项
      adminerStatusOptions: [],
      // 状态开关选项数据
      switchStatusInfo: {
        adminer_id: 0,
        status: 1
      },
      adminerStatusDefValue: null,
      adminerList: [],
      total: 0,
      listLoading: true,
      // 添加管理员表单数据
      addAdminerForm: {
        username: '',
        nickname: '',
        email: '',
        mobile: '',
        group_id: null,
        status: 2,
        remark: ''
      },
      // 添加管理员表单验证规则
      addAdminerFormRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 10, message: '用户名长度在3-10个字符', trigger: 'blur' }
        ],
        nickname: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
          { min: 2, max: 10, message: '姓名长度在2-10个字符', trigger: 'blur' }
        ],
        email: [
          { validator: this.validateEmail, trigger: 'blur' }
        ],
        mobile: [
          { validator: this.validateMobile, trigger: 'blur' }
        ],
        group_id: [
          { required: true, message: '请选择用户组', trigger: 'change' }
        ]
      },
      // 控制添加管理员的显示与隐藏
      addDialogVisible: false,
      addAdminerFormRef: {},
      // 编辑管理员表单数据
      editAdminerForm: {},
      // 添加管理员表单验证规则
      editAdminerFormRules: {
        nickname: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
          { min: 2, max: 10, message: '姓名长度在2-10个字符', trigger: 'blur' }
        ],
        email: [
          { validator: this.validateEmail, trigger: 'blur' }
        ],
        mobile: [
          { validator: this.validateMobile, trigger: 'blur' }
        ],
        group_id: [
          { required: true, message: '请选择用户组', trigger: 'change' }
        ]
      },
      editAdminerFormRef: {},
      isSaving: false,
      password: '',
      resetPasswordDialogVisible: false,
      newPassword: '',
      // 控制编辑管理员对话框的显示与隐藏
      editDialogVisible: false,
      isEditSaving: false
    }
  },
  created() {
    this.getAdminerList()
    // 当组件加载时检查查询参数
    if (this.$route.query.openAddDialog === 'true') {
      this.addDialogVisible = true
      this.addAdminerForm.group_id = parseInt(this.$route.query.groupId)
    }
  },
  methods: {
    async getAdminerList() {
      // eslint-disable-next-line no-unused-vars
      const { data: res } = await this.$http.get('adminer/list', {
        params: this.searchInfo
      })

      if (res.code !== 200) {
        return this.$message.error('获取管理员列表失败, 失败原因：' + res.message)
      }

      this.adminerList = res.data.adminerList
      this.total = res.data.adminerTotal
      this.groupOptions = res.data.groupOptions
      this.adminerStatusOptions = res.data.statusOptions
    },
    async handleSwitchChange(row) {
      const adminerId = row.adminer_id
      this.switchStatusInfo.adminer_id = adminerId
      this.switchStatusInfo.status = row.status
      const { data: res } = await this.$http.post('adminer/switchstatus', this.switchStatusInfo)
      if (res.code !== 200) {
        row.status = this.switchStatusInfo.status === 1 ? 2 : 1
        return this.$message.error('更新失败, 失败原因：' + res.message)
      }
      const statusStr = row.status === 1 ? '禁用' : '开启'
      this.$message.success(statusStr + '成功')
    },
    // 监听pageSize改变事件
    handleSizeChange(newSize) {
      this.searchInfo.page_size = newSize
      this.getAdminerList()
    },
    // 监听页码改变事件
    handleCurrentChange(newPage) {
      this.searchInfo.page = newPage
      this.getAdminerList()
    },
    handleGroupChange(newGroupId) {
      this.searchInfo.group_id = newGroupId
    },
    clearName() {
      this.searchInfo.name = ''
      this.getAdminerList()
    },
    // 验证邮箱的自定义方法
    validateEmail(rule, value, callback) {
      // 简单的邮箱正则表达式
      const regEmail = /^([a-zA-Z0-9_.-])+@([a-zA-Z0-9-])+(\.[a-zA-Z0-9]{2,4})+$/
      if (regEmail.test(value)) {
        callback()
      } else {
        callback(new Error('无效的邮箱地址'))
      }
    },
    validatePhone(rule, value, callback) {
      // 中国大陆手机号码的正则表达式
      const regPhone = /^1[3-9]\d{9}$/
      if (value === '' || value === undefined || value === null) {
        callback(new Error('请输入手机号码'))
      } else if (!regPhone.test(value)) {
        callback(new Error('手机号码格式不正确'))
      } else {
        callback()
      }
    },
    // 监听添加管理员表单的关闭事件
    addAdminerFormClosed() {
      this.$refs.addAdminerFormRef.resetFields()
      this.password = ''
      this.isSaving = false
    },
    // 添加管理员表单的提交事件
    addAdminerFormSub() {
      this.isSaving = true
      this.$refs.addAdminerFormRef.validate(async valid => {
        if (!valid) {
          this.isSaving = false
          return
        }

        try {
          // 发起添加管理员请求
          const { data: res } = await this.$http.post('/adminer/add', this.addAdminerForm)
          if (res.code !== 200) {
            return this.$message.error(res.message)
            // eslint-disable-next-line no-unreachable
            this.isSaving = false
          } else {
            this.password = res.data.password
            this.$message.success(res.message)
            setTimeout(() => {
              this.getAdminerList()
            }, 1500)
          }
        } catch (error) {
          this.$message.error('保存失败')
        }
      })
    },
    copyPassword() {
      if (this.password) {
        const input = this.$refs.passwordRef.$el.querySelector('input')
        input.select()
        document.execCommand('copy')
        this.$message.success('密码已复制到剪贴板')
      }
    },
    deleteAdminer(row) {
      // 发送请求删除数据
      this.$http.delete(`/adminer/delete/${row.adminer_id}`).then(response => {
        if (response.data.code === 200) {
          // 请求成功后，找到要删除的项的索引
          const index = this.adminerList.findIndex(item => item.adminer_id === row.adminer_id)
          if (index !== -1) {
            // 从数组中删除该项
            this.adminerList.splice(index, 1)
            this.$message.success('删除成功')
            setTimeout(() => {
              this.getAdminerList()
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
    resetPass(adminerId) {
      // 发送请求删除数据
      this.$http.get('/adminer/resetpass', { params: { adminer_id: adminerId } }).then(response => {
        if (response.data.code === 200) {
          // 请求成功后，弹窗展示新密码并且提供复制功能
          this.$message.success(response.data.message)
          this.newPassword = response.data.data.password
          this.resetPasswordDialogVisible = true
        } else {
          // 处理错误情况
          this.$message.error(response.data.message)
        }
      }).catch(error => {
        // 处理请求错误情况
        console.error('重置请求发生错误', error)
        this.$message.error('重置失败')
      })
    },
    copyNewPassword() {
      const textarea = document.createElement('textarea')
      textarea.value = this.newPassword
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
      this.$message.success('密码已复制到剪贴板')
    },
    copyPassClosed() {
      this.resetPasswordDialogVisible = false
      this.newPassword = ''
    },
    // 展示编辑用户的对话框
    async showEditDialog(adminerId) {
      const { data: res } = await this.$http.get(`/adminer/show/${adminerId}`)
      if (res.code !== 200) {
        return this.$message.error(res.message)
      }

      this.editAdminerForm = res.data
      this.editDialogVisible = true
    },
    // 监听编辑管理员的对话关闭事件
    editDialogClosed() {
      this.$refs.editAdminerFormRef.resetFields()
    },
    editAdminerInfo() {
      this.$refs.editAdminerFormRef.validate(async valid => {
        if (!valid) return

        const { data: res } = await this.$http.post('/adminer/edit', this.editAdminerForm)
        if (res.code !== 200) {
          return this.$message.error(res.message)
        }

        this.editDialogVisible = false
        this.getAdminerList()
        this.$message.success(res.message)
      })
    }
  }
}

</script>

<style lang="less" scoped>
</style>
