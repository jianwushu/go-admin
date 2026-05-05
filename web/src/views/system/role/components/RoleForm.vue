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

      <el-form-item v-if="!isAdminRole" :label="t('role.dataScope')" prop="dataScope">
        <el-select v-model="formData.dataScope" :placeholder="t('role.dataScopeRequired')" style="width: 100%">
          <el-option :label="t('role.dataScopeAll')" :value="1" />
          <el-option :label="t('role.dataScopeDept')" :value="2" />
          <el-option :label="t('role.dataScopeDeptAndChild')" :value="3" />
          <el-option :label="t('role.dataScopeSelf')" :value="4" />
          <el-option :label="t('role.dataScopeCustom')" :value="5" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="!isAdminRole" :label="t('role.menus')" prop="menuIds">
        <el-cascader
          ref="menuCascaderRef"
          v-model="formData.menuIds"
          :options="menuTreeData"
          :props="menuCascaderProps"
          :placeholder="t('role.selectMenus')"
          collapse-tags
          collapse-tags-tooltip
          clearable
          filterable
          v-loading="menuTreeLoading"
          style="width: 100%"
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
        </el-cascader>
      </el-form-item>

      <el-form-item v-if="!isAdminRole && formData.dataScope === 5" :label="t('role.depts')" prop="deptIds">
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
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { CascaderInstance } from 'element-plus'
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
const menuCascaderRef = ref<CascaderInstance>()
const submitLoading = ref(false)

/** 是否为 admin 角色（code='admin'），admin 角色无需维护菜单和数据权限关联 */
const isAdminRole = computed(() => formData.code === 'admin')

/** 菜单树数据 */
const menuTreeData = ref<MenuItem[]>([])
const menuTreeLoading = ref(false)

/** 级联选择器配置 */
const menuCascaderProps = {
  value: 'id',
  label: 'name',
  children: 'children',
  multiple: true,
  checkStrictly: false,
  emitPath: false,
}

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

/**
 * 根据选中的叶子节点 ID，递归收集所有需要提交的节点 ID（包括父节点）
 * 级联选择器 checkStrictly=false 时 modelValue 只包含叶子节点，
 * 但后端需要所有选中节点（含父节点），因此需要遍历树结构推导
 */
function collectAllCheckedIds(leafIds: number[], tree: MenuItem[]): number[] {
  const resultSet = new Set<number>()
  const leafIdSet = new Set(leafIds)

  function walk(nodes: MenuItem[]): boolean {
    let allChildrenChecked = true
    let hasAnyChildChecked = false

    for (const node of nodes) {
      const children = node.children || []
      let nodeChecked = false

      if (children.length > 0) {
        // 有子节点，递归处理
        const childrenResult = walk(children)
        if (childrenResult) {
          // 子节点中有被选中的
          hasAnyChildChecked = true
          // 检查是否所有子节点都被选中
          const allChecked = children.every(c => resultSet.has(c.id))
          if (allChecked) {
            resultSet.add(node.id)
          } else {
            allChildrenChecked = false
          }
        } else {
          allChildrenChecked = false
        }
      } else {
        // 叶子节点，检查是否在选中列表中
        nodeChecked = leafIdSet.has(node.id)
        if (nodeChecked) {
          resultSet.add(node.id)
          hasAnyChildChecked = true
        } else {
          allChildrenChecked = false
        }
      }
    }

    return hasAnyChildChecked
  }

  walk(tree)
  return Array.from(resultSet)
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

  // 级联选择器只返回叶子节点 ID，需要推导出所有选中节点（含父节点）
  const allMenuIds = collectAllCheckedIds(formData.menuIds, menuTreeData.value)

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
