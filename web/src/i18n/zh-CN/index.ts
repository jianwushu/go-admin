import common from './common'
import login from './modules/login'
import layout from './modules/layout'
import system from './modules/system'
import menuModule from './modules/menu'
import dashboard from './modules/dashboard'
import monitor from './modules/monitor'
import operationLog from './modules/operationLog'
import loginLog from './modules/loginLog'
import codegen from './modules/codegen'

export default {
  ...common,
  ...login,
  ...layout,
  ...system,
  ...menuModule,
  ...dashboard,
  ...monitor,
  ...operationLog,
  ...loginLog,
  ...codegen,
}
