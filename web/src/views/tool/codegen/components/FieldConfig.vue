<template>
  <div class="field-config">
    <el-table :data="fields" border stripe style="width: 100%" max-height="500">
      <el-table-column prop="columnName" :label="t('codegen.columnName')" width="140" fixed show-overflow-tooltip />
      <el-table-column prop="columnType" :label="t('codegen.columnType')" width="120" show-overflow-tooltip />
      <el-table-column prop="comment" :label="t('codegen.comment')" width="140" show-overflow-tooltip>
        <template #default="{ row }">
          <el-input v-model="row.comment" size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="goField" :label="t('codegen.goField')" width="140">
        <template #default="{ row }">
          <el-input v-model="row.goField" size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="goType" :label="t('codegen.goType')" width="100">
        <template #default="{ row }">
          <el-select v-model="row.goType" size="small" style="width: 100%">
            <el-option label="string" value="string" />
            <el-option label="int" value="int" />
            <el-option label="int64" value="int64" />
            <el-option label="float64" value="float64" />
            <el-option label="bool" value="bool" />
            <el-option label="time.Time" value="time.Time" />
            <el-option label="[]byte" value="[]byte" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column prop="tsType" :label="t('codegen.tsType')" width="100">
        <template #default="{ row }">
          <el-select v-model="row.tsType" size="small" style="width: 100%">
            <el-option label="string" value="string" />
            <el-option label="number" value="number" />
            <el-option label="boolean" value="boolean" />
            <el-option label="string[]" value="string[]" />
            <el-option label="Date" value="Date" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column prop="label" :label="t('codegen.label')" width="120">
        <template #default="{ row }">
          <el-input v-model="row.label" size="small" />
        </template>
      </el-table-column>
      <el-table-column prop="htmlType" :label="t('codegen.htmlType')" width="120">
        <template #default="{ row }">
          <el-select v-model="row.htmlType" size="small" style="width: 100%">
            <el-option :label="t('codegen.inputText')" value="input" />
            <el-option :label="t('codegen.textarea')" value="textarea" />
            <el-option :label="t('codegen.select')" value="select" />
            <el-option :label="t('codegen.radio')" value="radio" />
            <el-option :label="t('codegen.checkbox')" value="checkbox" />
            <el-option :label="t('codegen.date')" value="date" />
            <el-option :label="t('codegen.datetime')" value="datetime" />
            <el-option :label="t('codegen.upload')" value="upload" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column prop="queryType" :label="t('codegen.queryType')" width="120">
        <template #default="{ row }">
          <el-select v-model="row.queryType" size="small" style="width: 100%" :disabled="!row.isQuery">
            <el-option :label="t('codegen.eq')" value="eq" />
            <el-option :label="t('codegen.ne')" value="ne" />
            <el-option :label="t('codegen.gt')" value="gt" />
            <el-option :label="t('codegen.ge')" value="ge" />
            <el-option :label="t('codegen.lt')" value="lt" />
            <el-option :label="t('codegen.le')" value="le" />
            <el-option :label="t('codegen.like')" value="like" />
            <el-option :label="t('codegen.between')" value="between" />
          </el-select>
        </template>
      </el-table-column>
      <el-table-column :label="t('codegen.isList')" width="80" align="center">
        <template #default="{ row }">
          <el-checkbox v-model="row.isList" />
        </template>
      </el-table-column>
      <el-table-column :label="t('codegen.isQuery')" width="80" align="center">
        <template #default="{ row }">
          <el-checkbox v-model="row.isQuery" />
        </template>
      </el-table-column>
      <el-table-column :label="t('codegen.isRequired')" width="80" align="center">
        <template #default="{ row }">
          <el-checkbox v-model="row.isRequired" />
        </template>
      </el-table-column>
      <el-table-column :label="t('codegen.isEdit')" width="80" align="center">
        <template #default="{ row }">
          <el-checkbox v-model="row.isEdit" />
        </template>
      </el-table-column>
      <el-table-column :label="t('codegen.isPk')" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.isPk" type="danger" size="small">PK</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="dictType" :label="t('codegen.dictType')" width="130">
        <template #default="{ row }">
          <el-input v-model="row.dictType" size="small" :disabled="row.htmlType !== 'select' && row.htmlType !== 'radio' && row.htmlType !== 'checkbox'" />
        </template>
      </el-table-column>
      <el-table-column prop="sort" :label="t('codegen.sort')" width="80" align="center">
        <template #default="{ row }">
          <el-input-number v-model="row.sort" size="small" :min="0" :max="999" controls-position="right" style="width: 60px" />
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { ColumnConfig } from '@/api/codegen'

const { t } = useI18n()

defineProps<{
  fields: ColumnConfig[]
}>()
</script>
