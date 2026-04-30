<template>
  <el-dialog
    :model-value="visible"
    :title="isEdit ? t('common.edit') + t('role.title') : t('common.add') + t('role.title')"
    width="650px"
    destroy-on-close
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="110px"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('role.name')" prop="name">
            <el-input
              v-model="formData.name"
              :placeholder="t('role.namePlaceholder')"
              maxlength="64"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('role.code')" prop="code">
            <el-input
              v-model="formData.code"
              :placeholder="t('role.codePlaceholder')"
              :disabled="isEdit"
              maxlength="64"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="t('role.sort')" prop="sort">
            <el-input-number
              v-model="formData.sort"
              :min="0"
              controls-position="right"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('role.status')" prop="status">
            <el-radio-group v-model="formData.status">
              <el-radio :value="1">{{ t('role.enabled') }}</el-radio>
              <el-radio :value="0">{{ t('role.disabled') }}</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item :label="t('role.dataScope')" prop="dataScope">
        <el-select v-model="formData.dataScope" :placeholder="t('role.dataScopeRequired')" style="width: 100%">
          <el-option :label="t('role.dataScopeAll')" :value="1" />
          <el-option :label="t('role.dataScopeDept')" :value="2" />
          <el-option :label="t('role.dataScopeDeptAndChild')" :value="3" />
          <el-option :label="t('role.dataScopeSelf')" :value="4" />
          <el-option :label="t('role.dataScopeCustom')" :value="5" />
        </el-select>
      </el-form-item>

      <el-form-item :label="t('role.menus')" prop="menuIds">
        <el-tree
          ref="menuTreeRef"
          :data="menuTreeData"
          :props="{ label: 'name', children: 'children' }"
          show-checkbox
          node-key="id"
          :default-checked-keys="formData.menuIds"
          check-strictly
          v-loading="menuTreeLoading"
          style="width: 100%; max-height: 300px; overflow-y: auto; border: 1px solid #dcdfe6; border-radius: 4px; padding: 8px"
        >
          <template #default="{ node, data }">
            <div class="flex items-center gap-2">
              <!-- 根据菜单类型显示不同图标 -->
              <el-icon v-if="data.type === 0" class="text-blue-500">
                <Folder />
              </el-icon>
              <el-icon v-else-if="data.type === 1" class="text-green-500">
                <Document />
              </el-icon>
              <el-icon v-else-if="data.type === 2" class="text-orange-500">
                <Operation />
              </el-icon>
              <!-- 节点名称 -->
              <span>{{ node.label }}</span>
              <!-- 按钮类型显示权限标识 -->
              <el-tag v-if="data.type === 2 && data.perms" size="small" type="info" class="ml-2">
                {{ data.perms }}
              </el-tag>
            </div>
          </template>
        </el-tree>
      </el-form-item>

      <el-form-item v-if="formData.dataScope === 5" :label="t('role.depts')" prop="deptIds">
        <el-select
          v-model="formData.deptIds"
          multiple
          :placeholder="t('role.selectDepts')"
          style="width: 100%"
        >
          <!-- 部门选项将在后续部门管理模块实现后动态加载 -->
        </el-select>
      </el-form-item>

      <el-form-item :label="t('role.remark')" prop="remark">
        <el-input
          v-model="formData.remark"
          type="textarea"
          :placeholder="t('role.remarkPlaceholder')"
          :rows="3"
          maxlength="512"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('update:visible', false)">
        {{ t('common.cancel') }}
      </el-button>
      <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
        {{ t('common.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { ElTree } from 'element-plus'
import { Folder, Document, Operation } from '@element-plus/icons-vue'
import { createRole, updateRole, getMenuTree } from '@/api/system'
import type { RoleInfo, MenuItem } from '@/types/api'

const props = defineProps<{
  visible: boolean
  data: RoleInfo | null
  isEdit: boolean
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

const { t } = useI18n()

const formRef = ref<FormInstance>()
const menuTreeRef = ref<InstanceType<typeof ElTree>>()
const submitLoading = ref(false)

/** 菜单树数据 */
const menuTreeData = ref<MenuItem[]>([])
const menuTreeLoading = ref(false)

/** 获取菜单树数据 */
async function fetchMenuTree() {
  menuTreeLoading.value = true
  try {
    const { data: res } = await getMenuTree()
    menuTreeData.value = res.data || []
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    menuTreeLoading.value = false
  }
}

/** 组件挂载时获取菜单树 */
onMounted(() => {
  fetchMenuTree()
})

/** 表单数据 */
const formData = reactive({
  id: 0,
  name: '',
  code: '',
  dataScope: 1,
  sort: 0,
  status: 1,
  remark: '',
  menuIds: [] as number[],
  deptIds: [] as number[],
})

/** 表单校验规则 */
const formRules: FormRules = {
  name: [
    { required: true, message: () => t('role.nameRequired'), trigger: 'blur' },
  ],
  code: [
    { required: true, message: () => t('role.codeRequired'), trigger: 'blur' },
  ],
  dataScope: [
    { required: true, message: () => t('role.dataScopeRequired'), trigger: 'change' },
  ],
}

/** 监听弹窗打开，初始化表单 */
watch(() => props.visible, (val) => {
  if (val) {
    if (props.isEdit && props.data) {
      // 编辑模式：填充数据
      formData.id = props.data.id
      formData.name = props.data.name
      formData.code = props.data.code
      formData.dataScope = props.data.dataScope
      formData.sort = props.data.sort
      formData.status = props.data.status
      formData.remark = props.data.remark
      formData.menuIds = props.data.menuIds || []
      formData.deptIds = props.data.deptIds || []
    } else {
      // 新增模式：重置表单
      formData.id = 0
      formData.name = ''
      formData.code = ''
      formData.dataScope = 1
      formData.sort = 0
      formData.status = 1
      formData.remark = ''
      formData.menuIds = []
      formData.deptIds = []
    }
  }
})

/** 提交表单 */
async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  // 获取选中的菜单ID
  const checkedMenuIds = menuTreeRef.value?.getCheckedKeys(false) as number[] || []
  const halfCheckedMenuIds = menuTreeRef.value?.getHalfCheckedKeys() as number[] || []
  const allMenuIds = [...checkedMenuIds, ...halfCheckedMenuIds]

  submitLoading.value = true
  try {
    const submitData = {
      ...formData,
      menuIds: allMenuIds,
    }

    if (props.isEdit) {
      await updateRole(submitData)
    } else {
      await createRole(submitData)
    }
    ElMessage.success(t('common.success'))
    emit('update:visible', false)
    emit('success')
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    submitLoading.value = false
  }
}
</script>
