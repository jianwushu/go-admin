import type { App, Directive, DirectiveBinding } from 'vue'
import { usePermissionStore } from '@/store/modules/permission'

/**
 * v-permission 按钮权限指令
 *
 * 用法：
 *   单个权限：  v-permission="'system:user:add'"
 *   多个权限（满足其一）：v-permission="['system:user:add', 'system:user:edit']"
 *
 * 如果用户没有指定权限，元素将从 DOM 中移除。
 */
const permissionDirective: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string | string[]>) {
    const permissionStore = usePermissionStore()
    const requiredPermissions = binding.value

    if (!requiredPermissions) {
      console.warn('[v-permission] 缺少权限标识参数')
      return
    }

    // 超级管理员拥有所有权限
    if (permissionStore.buttons.includes('*:*:*')) {
      return
    }

    let hasPermission = false

    if (Array.isArray(requiredPermissions)) {
      // 多个权限，满足其一即可
      hasPermission = requiredPermissions.some((perm) =>
        permissionStore.buttons.includes(perm)
      )
    } else {
      // 单个权限
      hasPermission = permissionStore.buttons.includes(requiredPermissions)
    }

    if (!hasPermission) {
      // 移除元素
      el.parentNode?.removeChild(el)
    }
  },
}

/**
 * 注册权限指令
 */
export function setupPermissionDirective(app: App) {
  app.directive('permission', permissionDirective)
}

export default permissionDirective
