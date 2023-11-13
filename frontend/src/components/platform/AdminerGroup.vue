<!-- eslint-disable vue/no-unused-vars -->
<template>
    <div>
        <!-- 面包屑导航 -->
        <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>系统管理</el-breadcrumb-item>
            <el-breadcrumb-item>管理组列表</el-breadcrumb-item>
        </el-breadcrumb>

        <!-- 卡片试图 -->
        <el-card>
            <el-row :gutter="10">
                <el-col :span="4">
                    <el-button type="primary" @click="saveDialogVisible = true">添加管理组</el-button>
                </el-col>
            </el-row>

            <!-- 管理组列表 -->
            <el-table :data="adminerGroupList" border stripe>
                <el-table-column label="ID" prop="group_id" width="50px"></el-table-column>
                <el-table-column label="组名" prop="name"></el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-tooltip class="item" e ffect="dark" content="编辑" placement="top" :enterable="false">
                            <el-button type="primary" icon="el-icon-edit" size="mini" @click="showEditDialog(scope.row.group_id)"></el-button>
                        </el-tooltip>
                        <el-tooltip class="item" effect="dark" content="添加管理员" placement="top" :enterable="false">
                            <el-button type="warning" icon="el-icon-plus" size="mini" @click="addAdminer(scope.row.group_id)"></el-button>
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

        <!-- 添加管理组弹窗 -->
        <el-dialog
        title="添加管理组"
        :visible.sync="saveDialogVisible"
        width="50%"
        @close="saveAdminerGroupFormClosed">
        <!-- 管理组信息主体区 -->
        <el-form
        ref="saveAdminerGroupFormRef"
        :model="saveAdminerGroupForm"
        :rules="saveAdminerGroupFormRules"
        label-width="80px">
            <!-- 隐藏的组ID -->
            <input type="hidden" v-model="saveAdminerGroupForm.group_id">
            <el-form-item label="组名" prop="name">
                <el-input v-model="saveAdminerGroupForm.name"></el-input>
            </el-form-item>
            <el-form-item label="备注">
                <el-input v-model="saveAdminerGroupForm.remark"></el-input>
            </el-form-item>
            <!-- 权限选择区 -->
            <el-form-item label="权限">
                <el-tree
                ref="adminerGroupTree"
                :data="authTree"
                :props="defaultProps"
                :default-checked-keys="checkedIds"
                node-key="auth_id"
                show-checkbox
                :default-expand-all="false"
                @check-change="handleCheckChange"
                >
                <template slot-scope="{ node, data }">
                <div :class="{ 'is-leaf': !data.children || data.children.length === 0 }">
                    {{ node.label }}
                </div>
                </template>
                </el-tree>
            </el-form-item>
        </el-form>

        <!-- 底部区 -->
        <span slot="footer" class="dialog-footer">
            <el-button @click="saveDialogVisible = false">取 消</el-button>
            <el-button type="primary" :disabled="isSaving" @click="saveAdminerGroupFormSub">确 定</el-button>
        </span>
        </el-dialog>
    </div>
</template>

<script>

export default {
  data() {
    return {
      searchInfo: {
        page_size: 10,
        page: 1
      },
      // 管理组列表
      adminerGroupList: [],
      total: 0,
      // 添加管理员组表单数据/编辑管理员组表单数据
      saveAdminerGroupForm: {
        group_id: 0,
        name: '',
        remark: ''
      },
      // 添加管理员表单验证规则
      saveAdminerGroupFormRules: {
        name: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 10, message: '用户名长度在3-10个字符', trigger: 'blur' }
        ]
      },
      // 控制添加管理员的显示与隐藏
      saveDialogVisible: false,
      saveAdminerGroupFormRef: {},
      // 这里将会是你从后端获取的权限数据
      authTree: [],
      copyAuthTree: [],
      defaultProps: {
        children: 'children',
        label: 'auth_name',
        isLeaf: 'leaf'
      },
      // 用于存储选中的权限ID
      selectedAuthIds: new Set(),
      // 用于存储默认展开的节点的ID
      defaultExpandedKeys: [],
      isSaving: false,
      listLoading: true,
      checkedIds: []
    }
  },
  created() {
    this.getAdminerGroupList()
  },
  methods: {
    async getAdminerGroupList() {
      const { data: res } = await this.$http.get('adminer/group/list', { params: this.searchInfo })

      if (res.code !== 200) {
        return this.$message.error('获取管理组列表失败, 失败原因：' + res.message)
      }

      this.adminerGroupList = res.data.list
      this.total = res.data.total
      this.authTree = res.data.auth_tree.sort((a, b) => a.auth_group_id - b.auth_group_id)
      this.copyAuthTree = this.authTree
    },
    // 监听pageSize改变事件
    handleSizeChange(newSize) {
      this.searchInfo.page_size = newSize
      this.getAdminerGroupList()
    },
    // 监听页码改变事件
    handleCurrentChange(newPage) {
      this.searchInfo.page = newPage
      this.getAdminerGroupList()
    },
    // 监听添加管理员表单的关闭事件
    saveAdminerGroupFormClosed() {
      this.$refs.saveAdminerGroupFormRef.resetFields()
      this.isSaving = false
      // 清空表单和选中的权限ID
      this.selectedAuthIds.clear()
      this.saveAdminerGroupForm = {
        group_id: null,
        name: '',
        remark: ''
      }
      if (this.$refs.adminerGroupTree) {
        this.$refs.adminerGroupTree.setCheckedKeys([])
      }
    },
    // 添加管理员表单的提交事件
    saveAdminerGroupFormSub() {
      this.isSaving = true

      const payload = {
        group_id: this.saveAdminerGroupForm.group_id,
        name: this.saveAdminerGroupForm.name,
        remark: this.saveAdminerGroupForm.remark,
        auth_ids: this.$refs.adminerGroupTree.getCheckedKeys().filter(key => key !== undefined)
      }

      this.$refs.saveAdminerGroupFormRef.validate(async valid => {
        if (!valid) {
          this.isSaving = false
          return
        }

        try {
          // 发起添加管理员请求
          const { data: res } = await this.$http.post('/adminer/group/save', payload)
          if (res.code !== 200) {
            return this.$message.error(res.message)
          } else {
            this.$message.success(res.message)
            this.saveDialogVisible = false
            setTimeout(() => {
              this.getAdminerGroupList()
            }, 1000)
          }
          this.isSaving = false
        } catch (error) {
          this.$message.error('保存失败')
          this.isSaving = false
        }
      })
    },
    // 展示编辑用户的对话框
    async showEditDialog(groupId) {
      const { data: res } = await this.$http.get(`/adminer/group/show/${groupId}`)
      if (res.code !== 200) {
        return this.$message.error(res.message)
      }

      this.saveAdminerGroupForm.group_id = res.data.group_id
      this.saveAdminerGroupForm.name = res.data.name
      this.saveAdminerGroupForm.remark = res.data.remark
      this.checkedIds = this.getCheckedAuthIds(res.data.auth_tree)
      this.saveDialogVisible = true
    },
    handleCheckChange(data, checked, indeterminate) {
      this.recursiveToggleCheck(data, checked)
    },
    recursiveToggleCheck(node, checked) {
      // 如果节点被选中，添加它的ID到selectedAuthIds
      if (checked) {
        this.selectedAuthIds.add(node.auth_id)
      } else {
        this.selectedAuthIds.delete(node.auth_id)
      }
      // 对子节点做同样的操作
      if (node.children && node.children.length) {
        node.children.forEach((child) => {
          this.recursiveToggleCheck(child, checked)
        })
      }
    },
    handleNodeClick(data, node) {
      // 如果是第一级节点，切换展开/折叠状态
      if (node.level === 1) {
        node.expanded ? node.collapse() : node.expand()
      }
    },
    renderContent(h, { node, data, store }) {
    // 如果节点没有子节点，则禁用复选框
      if (!data.children || data.children.length === 0) {
        return h('span', [
          h('span', node.label)
        ])
      } else {
        return h('span', [
          h('el-checkbox', {
            props: {
              checked: node.checked,
              indeterminate: node.indeterminate
            },
            on: {
              change: () => { store.commit('check', { data: node.data, checked: !node.checked, indeterminate: node.indeterminate }) }
            }
          }, node.label)
        ])
      }
    },
    addAdminer(groupId) {
      this.$router.push({ path: '/adminer/list', query: { openAddDialog: 'true', groupId: groupId } })
    },
    getCheckedAuthIds(authTree) {
      const checkedIds = []
      const traverseNodes = (nodes) => {
        for (const node of nodes) {
          if (node.is_check === 2) {
            checkedIds.push(node.auth_id)
          }
          if (node.children && node.children.length) {
            traverseNodes(node.children)
          }
        }
      }

      traverseNodes(authTree)
      return checkedIds
    },
    addGroupDialogInit() {
      this.saveDialogVisible = true
      this.checkedIds = []
      this.authTree = this.copyAuthTree
    }
  }
}

</script>

<style lang="less" scoped>
.dialog-footer {
  text-align: right;
}

.el-tree .is-leaf .el-checkbox {
  display: none;
}
</style>
