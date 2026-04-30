/** 全局类型声明 */

declare global {
  /** 任意对象 */
  type Recordable<T = unknown> = Record<string, T>

  /** 可空类型 */
  type Nullable<T> = T | null

  /** 可选类型 */
  type Optional<T> = T | undefined
}

export {}
