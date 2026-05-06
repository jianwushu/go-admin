<template>
  <div class="p-4">
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-medium">{{ t('systemConfig.title') }}</span>
          <el-button type="primary" v-permission="'system:config:edit'" :loading="saving" @click="handleSave">
            <el-icon class="mr-1"><Check /></el-icon>
            {{ t('systemConfig.save') }}
          </el-button>
        </div>
      </template>

      <el-form
        v-loading="loading"
        ref="formRef"
        :model="formData"
        label-width="140px"
        style="max-width: 700px"
      >
        <!-- 系统标题 -->
        <el-form-item :label="t('systemConfig.systemTitle')">
          <el-input
            v-model="formData.system_title"
            :placeholder="t('systemConfig.systemTitle')"
            maxlength="128"
            show-word-limit
          />
        </el-form-item>

        <!-- 版权信息 -->
        <el-form-item :label="t('systemConfig.systemCopyright')">
          <el-input
            v-model="formData.system_copyright"
            :placeholder="t('systemConfig.systemCopyright')"
            maxlength="512"
            show-word-limit
          />
        </el-form-item>

        <!-- 系统Logo -->
        <el-form-item :label="t('systemConfig.systemLogo')">
          <div class="flex items-start gap-4">
            <div v-if="formData.system_logo" class="relative">
              <el-image
                :src="formData.system_logo"
                fit="contain"
                style="width: 120px; height: 60px"
                :preview-src-list="[formData.system_logo]"
              />
              <el-button
                type="danger"
                size="small"
                circle
                class="absolute -top-2 -right-2"
                @click="formData.system_logo = ''"
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
            <el-upload
              :show-file-list="false"
              :before-upload="beforeLogoUpload"
              :http-request="handleLogoUpload"
              accept=".png,.jpg,.jpeg,.gif,.svg,.ico"
            >
              <el-button type="primary" plain>
                <el-icon class="mr-1"><Upload /></el-icon>
                {{ t('systemConfig.uploadLogo') }}
              </el-button>
            </el-upload>
          </div>
          <div class="text-xs text-gray-400 mt-1">{{ t('systemConfig.uploadHint') }}</div>
        </el-form-item>

        <!-- Favicon -->
        <el-form-item :label="t('systemConfig.systemFavicon')">
          <div class="flex items-start gap-4">
            <div v-if="formData.system_favicon" class="relative">
              <el-image
                :src="formData.system_favicon"
                fit="contain"
                style="width: 32px; height: 32px"
                :preview-src-list="[formData.system_favicon]"
              />
              <el-button
                type="danger"
                size="small"
                circle
                class="absolute -top-2 -right-2"
                @click="formData.system_favicon = ''"
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
            <el-upload
              :show-file-list="false"
              :before-upload="beforeLogoUpload"
              :http-request="handleFaviconUpload"
              accept=".png,.jpg,.jpeg,.gif,.svg,.ico"
            >
              <el-button type="primary" plain>
                <el-icon class="mr-1"><Upload /></el-icon>
                {{ t('systemConfig.uploadLogo') }}
              </el-button>
            </el-upload>
          </div>
          <div class="text-xs text-gray-400 mt-1">{{ t('systemConfig.uploadHint') }}</div>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Check, Close, Upload } from '@element-plus/icons-vue'
import { getAllSystemConfig, batchUpdateSystemConfig, uploadSystemLogo } from '@/api/system'
import type { SystemConfigItem } from '@/types/api'

defineOptions({ name: 'SystemConfig' })

const { t } = useI18n()

const loading = ref(false)
const saving = ref(false)

// 配置项ID映射
const configIdMap = reactive<Record<string, number>>({
  system_title: 0,
  system_copyright: 0,
  system_logo: 0,
  system_favicon: 0,
})

// 表单数据
const formData = reactive({
  system_title: '',
  system_copyright: '',
  system_logo: '',
  system_favicon: '',
})

/** 获取配置数据 */
async function fetchData() {
  loading.value = true
  try {
    const { data } = await getAllSystemConfig()
    if (data.code === 0 && data.data) {
      const configs = data.data as SystemConfigItem[]
      for (const item of configs) {
        if (item.configKey in formData) {
          ;(formData as Record<string, string>)[item.configKey] = item.configValue
          configIdMap[item.configKey] = item.id
        }
      }
    }
  } catch (error) {
    console.error('获取系统配置失败', error)
  } finally {
    loading.value = false
  }
}

/** 保存配置 */
async function handleSave() {
  saving.value = true
  try {
    const list = Object.keys(formData).map((key) => ({
      id: configIdMap[key],
      configValue: (formData as Record<string, string>)[key],
    }))
    const { data } = await batchUpdateSystemConfig(list)
    if (data.code === 0) {
      ElMessage.success(t('systemConfig.saveSuccess'))
    } else {
      ElMessage.error(data.msg)
    }
  } catch (error) {
    console.error('保存配置失败', error)
  } finally {
    saving.value = false
  }
}

/** 上传前校验 */
function beforeLogoUpload(file: File) {
  const allowedTypes = ['image/png', 'image/jpeg', 'image/gif', 'image/svg+xml', 'image/x-icon']
  if (!allowedTypes.includes(file.type) && !file.name.endsWith('.ico')) {
    ElMessage.error('不支持的文件类型')
    return false
  }
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.error('文件大小不能超过2MB')
    return false
  }
  return true
}

/** 上传Logo */
async function handleLogoUpload(options: { file: File }) {
  try {
    const { data } = await uploadSystemLogo(options.file)
    if (data.code === 0 && data.data) {
      formData.system_logo = data.data
      ElMessage.success(t('systemConfig.uploadSuccess'))
    } else {
      ElMessage.error(data.msg)
    }
  } catch (error) {
    console.error('上传失败', error)
  }
}

/** 上传Favicon */
async function handleFaviconUpload(options: { file: File }) {
  try {
    const { data } = await uploadSystemLogo(options.file)
    if (data.code === 0 && data.data) {
      formData.system_favicon = data.data
      ElMessage.success(t('systemConfig.uploadSuccess'))
    } else {
      ElMessage.error(data.msg)
    }
  } catch (error) {
    console.error('上传失败', error)
  }
}

onMounted(() => {
  fetchData()
})
</script>
