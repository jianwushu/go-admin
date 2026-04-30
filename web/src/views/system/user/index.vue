<template>
  <div class="p-4">
    <!-- 搜索区域 -->
    <el-card shadow="never" class="mb-4">
      <el-form :model="queryParams" inline>
        <el-form-item :label="t('user.username')">
          <el-input
            v-model="queryParams.username"
            :placeholder="t('user.usernamePlaceholder')"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="t('user.status')">
          <el-select v-model="queryParams.status" :placeholder="t('common.status')" clearable style="width: 120px">
            <el-option :label="t('user.enabled')" :value="1" />
            <el-option :label="t('user.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon class="mr-1"><Search /></el-icon>
            {{ t('common.search') }}
          </el-button>
          <el-button @click="handleReset">
            <el-icon class="mr-1"><Refresh /></el-icon>
            {{ t('common.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-medium">{{ t('user.title') }}</span>
          <div>
            <el-button type="primary" v-permission="'system:user:add'" @click="handleAdd">
              <el-icon class="mr-1"><Plus /></el-icon>
              {{ t('common.add') }}
            </el-button>
          </div>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" :label="t('user.username')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="nickname" :label="t('user.nickname')" min-width="120" show-overflow-tooltip />
        <el-table-column prop="email" :label="t('user.email')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="phone" :label="t('user.phone')" min-width="130" show-overflow-tooltip />
        <el-table-column :label="t('user.roles')" min-width="150">
          <template #default="{ row }">
            <el-tag
              v-for="role in row.roles"
              :key="role.id"
              size="small"
              class="mr-1 mb-1"
            >
              {{ role.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('user.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              v-permission="'system:user:edit'"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="t('common.createTime')" width="170" align="center" />
        <el-table-column :label="t('common.operation')" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link v-permission="'system:user:edit'" @click="handleEdit(row)">
              <el-icon class="mr-1"><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button type="warning" link v-permission="'system:user:resetPwd'" @click="handleResetPassword(row)">
              <el-icon class="mr-1"><Key /></el-icon>
              {{ t('user.resetPassword') }}
            </el-button>
            <el-button type="danger" link v-permission="'system:user:remove'" @click="handleDelete(row)">
              <el-icon class="mr-1"><Delete /></el-icon>
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 用户表单弹窗 -->
    <UserForm
      v-model:visible="formVisible"
      :data="currentUser"
      :is-edit="isEdit"
      @success="fetchData"
    />

    <!-- 重置密码弹窗 -->
    <el-dialog
      v-model="resetPwdVisible"
      :title="t('user.resetPassword')"
      width="400px"
      destroy-on-close
    >
      <el-form ref="resetPwdFormRef" :model="resetPwdForm" :rules="resetPwdRules" label-width="100px">
        <el-form-item :label="t('user.newPassword')" prop="password">
          <el-input
            v-model="resetPwdForm.password"
            type="password"
            :placeholder="t('user.passwordPlaceholder')"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetPwdVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="resetPwdLoading" @click="submitResetPassword">
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, Plus, Edit, Delete, Key } from '@element-plus/icons-vue'
import { getUserList, deleteUser, changeUserStatus, resetUserPassword } from '@/api/system'
import type { UserItem, UserListParams } from '@/types/api'
import UserForm from './components/UserForm.vue'

defineOptions({ name: 'UserManagement' })

const { t } = useI18n()

// 查询参数
const queryParams = reactive<UserListParams>({
  page: 1,
  pageSize: 10,
  username: undefined,
  status: undefined,
})

// 表格数据
const loading = ref(false)
const tableData = ref<UserItem[]>([])
const total = ref(0)

// 表单相关
const formVisible = ref(false)
const isEdit = ref(false)
const currentUser = ref<UserItem | null>(null)

// 重置密码相关
const resetPwdVisible = ref(false)
const resetPwdLoading = ref(false)
const resetPwdFormRef = ref<FormInstance>()
const resetPwdUserId = ref(0)
const resetPwdForm = reactive({
  password: '',
})
const resetPwdRules: FormRules = {
  password: [
    { required: true, message: () => t('user.passwordRequired'), trigger: 'blur' },
    { min: 6, message: () => t('user.passwordMin'), trigger: 'blur' },
  ],
}

/** 获取数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data: res } = await getUserList(queryParams)
    tableData.value = res.data || []
    total.value = res.total || 0
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    loading.value = false
  }
}

/** 搜索 */
function handleSearch() {
  queryParams.page = 1
  fetchData()
}

/** 重置 */
function handleReset() {
  queryParams.username = undefined
  queryParams.status = undefined
  queryParams.page = 1
  fetchData()
}

/** 分页大小变化 */
function handleSizeChange() {
  queryParams.page = 1
  fetchData()
}

/** 页码变化 */
function handleCurrentChange() {
  fetchData()
}

/** 新增 */
function handleAdd() {
  currentUser.value = null
  isEdit.value = false
  formVisible.value = true
}

/** 编辑 */
function handleEdit(row: UserItem) {
  currentUser.value = { ...row }
  isEdit.value = true
  formVisible.value = true
}

/** 删除 */
async function handleDelete(row: UserItem) {
  try {
    await ElMessageBox.confirm(t('common.confirmDelete'), t('common.tip'), {
      type: 'warning',
    })
    await deleteUser(row.id)
    ElMessage.success(t('common.success'))
    fetchData()
  } catch {
    // 取消或失败
  }
}

/** 状态变更 */
async function handleStatusChange(row: UserItem) {
  try {
    await changeUserStatus(row.id, row.status)
    ElMessage.success(t('common.success'))
  } catch {
    // 回滚状态
    row.status = row.status === 1 ? 0 : 1
  }
}

/** 重置密码 */
function handleResetPassword(row: UserItem) {
  resetPwdUserId.value = row.id
  resetPwdForm.password = ''
  resetPwdVisible.value = true
}

/** 提交重置密码 */
async function submitResetPassword() {
  const valid = await resetPwdFormRef.value?.validate().catch(() => false)
  if (!valid) return

  resetPwdLoading.value = true
  try {
    await resetUserPassword(resetPwdUserId.value, resetPwdForm.password)
    ElMessage.success(t('user.resetPasswordSuccess'))
    resetPwdVisible.value = false
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    resetPwdLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
