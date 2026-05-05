import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { TagView } from '@/types/store'

export const useTagsViewStore = defineStore('tagsView', () => {
  const visitedViews = ref<TagView[]>([])
  const cachedViews = ref<string[]>([])

  /** 添加访问视图 */
  function addVisitedView(view: TagView) {
    if (visitedViews.value.some((v) => v.path === view.path)) return
    const newView = Object.assign({}, view, {
      title: view.title || 'no-name',
    })
    // affix 标签始终插入到第一位，非 affix 标签追加到末尾
    if (view.meta?.affix) {
      visitedViews.value.unshift(newView)
    } else {
      visitedViews.value.push(newView)
    }
  }

  /** 添加缓存视图 */
  function addCachedView(view: TagView) {
    if (cachedViews.value.includes(view.name)) return
    if (!view.meta?.keepAlive) return
    cachedViews.value.push(view.name)
  }

  /** 删除访问视图 */
  function delVisitedView(view: TagView) {
    const index = visitedViews.value.findIndex((v) => v.path === view.path)
    if (index > -1) {
      visitedViews.value.splice(index, 1)
    }
  }

  /** 删除缓存视图 */
  function delCachedView(view: TagView) {
    const index = cachedViews.value.indexOf(view.name)
    if (index > -1) {
      cachedViews.value.splice(index, 1)
    }
  }

  /** 删除其他视图 */
  function delOthersViews(view: TagView) {
    visitedViews.value = visitedViews.value.filter(
      (v) => v.meta?.affix || v.path === view.path
    )
    cachedViews.value = cachedViews.value.filter((name) => name === view.name)
  }

  /** 删除所有视图 */
  function delAllViews() {
    visitedViews.value = visitedViews.value.filter((v) => v.meta?.affix)
    cachedViews.value = []
  }

  return {
    visitedViews,
    cachedViews,
    addVisitedView,
    addCachedView,
    delVisitedView,
    delCachedView,
    delOthersViews,
    delAllViews,
  }
})
