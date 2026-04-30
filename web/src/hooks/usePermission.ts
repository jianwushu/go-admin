import { usePermissionStore } from '@/store/modules/permission'

/**
 * usePermission 按钮权限组合式函数
 *
 * 用法：
 *   const { hasPermission, hasAnyPermission, hasAllPermissions } = usePermission()
 *
 *   // 检查单个权限
 *   if (hasPermission('system:user:add')) { ... }
 *
 *   // 检查是否有任一权限
 *   if (hasAnyPermission(['system:user:add', 'system:user:edit'])) { ... }
 *
 *   // 检查是否拥有所有权限
 *   if (hasAllPermissions(['system:user:add', 'system:user:edit'])) { ... }
 */
export function usePermission() {
  const permissionStore = usePermissionStore()

  /**
   * 检查是否拥有指定权限
   * @param permission 权限标识，如 'system:user:add'
   * @returns 是否拥有该权限
   */
  function hasPermission(permission: string): boolean {
    // 超级管理员拥有所有权限
    if (permissionStore.buttons.includes('*:*:*')) {
      return true
    }
    return permissionStore.buttons.includes(permission)
  }

  /**
   * 检查是否拥有任一权限
   * @param permissions 权限标识数组
   * @returns 是否拥有其中任一权限
   */
  function hasAnyPermission(permissions: string[]): boolean {
    if (permissionStore.buttons.includes('*:*:*')) {
      return true
    }
    return permissions.some((perm) => permissionStore.buttons.includes(perm))
  }

  /**
   * 检查是否拥有所有权限
   * @param permissions 权限标识数组
   * @returns 是否拥有所有权限
   */
  function hasAllPermissions(permissions: string[]): boolean {
    if (permissionStore.buttons.includes('*:*:*')) {
      return true
    }
    return permissions.every((perm) => permissionStore.buttons.includes(perm))
  }

  /**
   * 获取当前用户所有按钮权限
   */
  function getPermissions(): string[] {
    return permissionStore.buttons
  }

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    getPermissions,
  }
}
